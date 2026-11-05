/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/df-mc/go-nethernet"
	"github.com/df-mc/go-nethernet/endpoint"
	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"
	"github.com/sandertv/gophertunnel/minecraft"
)

type netherNetListener struct {
	Listener
	signaling *http.Server
	ice       net.PacketConn
}

func (l *netherNetListener) Close() error {
	err := l.Listener.Close()
	_ = l.signaling.Close()
	_ = l.ice.Close()
	return err
}

func (uc UserConfig) netherNetListener(conf Config) (Listener, error) {
	return uc.netherNetListenerWithUDP(conf, nil)
}

func (uc UserConfig) netherNetListenerWithUDP(conf Config, shared *sharedGameUDP) (Listener, error) {
	settings := uc.Network.NetherNet
	iceAddress := settings.ICEAddress
	if iceAddress == "" {
		iceAddress = ":19133"
	}
	if shared != nil {
		iceAddress = uc.Network.Address
	}
	var udp net.PacketConn
	var err error
	if shared != nil {
		var mux *gameUDPMux
		mux, err = shared.packet()
		if err == nil {
			udp = mux.ice
		}
	} else {
		udp, err = net.ListenPacket("udp", iceAddress)
	}
	if err != nil {
		logNetherNetListenFailure(conf.Log, uc.Network.Address, iceAddress, err)
		return nil, fmt.Errorf("listen ICE UDP: %w", err)
	}

	se := webrtc.SettingEngine{}
	se.SetICEUDPMux(ice.NewUDPMuxDefault(ice.UDPMuxParams{UDPConn: udp}))
	if settings.PublicIP != "" {
		se.SetNAT1To1IPs([]string{settings.PublicIP}, webrtc.ICECandidateTypeHost)
	}

	handler := endpoint.NewHandler()
	transport := minecraft.NetherNet{
		Signaling: handler,
		ListenConfig: nethernet.ListenConfig{
			Log:               conf.Log,
			API:               webrtc.NewAPI(webrtc.WithSettingEngine(se)),
			DisableTrickleICE: true,
			ConnContext: func(parent context.Context, _ *nethernet.Conn) (context.Context, context.CancelFunc) {
				return context.WithTimeout(parent, 20*time.Second)
			},
		},
	}

	cfg := uc.minecraftListenConfig(conf)
	cfg.ErrorLog = conf.Log.With("net origin", "gophertunnel", "transport", "nethernet")
	minecraftListener, err := cfg.ListenNetwork(transport, handler.NetworkID())
	if err != nil {
		_ = udp.Close()
		logNetherNetListenFailure(conf.Log, uc.Network.Address, iceAddress, err)
		return nil, fmt.Errorf("listen NetherNet: %w", err)
	}

	tcp, err := net.Listen("tcp", uc.Network.Address)
	if err != nil {
		_ = minecraftListener.Close()
		_ = udp.Close()
		logNetherNetListenFailure(conf.Log, uc.Network.Address, iceAddress, err)
		return nil, fmt.Errorf("listen signaling TCP: %w", err)
	}

	var tlsConfig *tls.Config
	if settings.Certificate != "" || settings.PrivateKey != "" {
		cert, loadErr := tls.LoadX509KeyPair(settings.Certificate, settings.PrivateKey)
		if loadErr != nil {
			_ = tcp.Close()
			_ = minecraftListener.Close()
			_ = udp.Close()
			return nil, fmt.Errorf("load NetherNet TLS certificate: %w", loadErr)
		}
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{cert}}
	}

	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	go func() {
		if serveErr := server.Serve(mixedHTTPListener{Listener: tcp, TLS: tlsConfig}); serveErr != nil && serveErr != http.ErrServerClosed {
			conf.Log.Error("NetherNet signaling stopped", "error", serveErr)
		}
	}()

	conf.Log.Info("Listener running", "transport", "nethernet", "tcp", uc.Network.Address, "iceUDP", udp.LocalAddr(), "tls", tlsConfig != nil)
	if tlsConfig == nil {
		conf.Log.Warn("NetherNet HTTPS is unavailable", "hint", "configure NetherNet certificate and private key for HTTPS signaling")
	}
	return &netherNetListener{
		Listener:  listener{Listener: minecraftListener, log: conf.Log},
		signaling: server,
		ice:       udp,
	}, nil
}

type mixedHTTPListener struct {
	net.Listener
	TLS *tls.Config
}

func (l mixedHTTPListener) Accept() (net.Conn, error) {
	for {
		conn, err := l.Listener.Accept()
		if err != nil {
			return nil, err
		}
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buffer := bufio.NewReader(conn)
		first, err := buffer.Peek(1)
		_ = conn.SetReadDeadline(time.Time{})
		if err != nil {
			_ = conn.Close()
			continue
		}
		wrapped := &bufferedConnection{Conn: conn, buffer: buffer}
		if first[0] == 0x16 && l.TLS != nil {
			return tls.Server(wrapped, l.TLS), nil
		}
		return wrapped, nil
	}
}

type bufferedConnection struct {
	net.Conn
	buffer *bufio.Reader
}

func (c *bufferedConnection) Read(b []byte) (int, error) {
	return c.buffer.Read(b)
}
