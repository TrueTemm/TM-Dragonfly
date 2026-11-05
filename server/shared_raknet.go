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
	"context"
	"net"
	"sync"

	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft"
)

type sharedGameUDP struct {
	once    sync.Once
	address string
	mux     *gameUDPMux
	err     error
}

func (s *sharedGameUDP) packet() (*gameUDPMux, error) {
	s.once.Do(func() { s.mux, s.err = newGameUDPMux(s.address) })
	return s.mux, s.err
}

type sharedRakNetNetwork struct{ conn net.PacketConn }

type reusedPacketListener struct{ conn net.PacketConn }

func (p reusedPacketListener) ListenPacket(string, string) (net.PacketConn, error) {
	return p.conn, nil
}

func (n sharedRakNetNetwork) Listen(address string) (minecraft.NetworkListener, error) {
	return raknet.ListenConfig{UpstreamPacketListener: reusedPacketListener{conn: n.conn}}.Listen(address)
}
func (n sharedRakNetNetwork) DialContext(ctx context.Context, address string) (net.Conn, error) {
	return rakNetNetwork(nil).DialContext(ctx, address)
}
func (n sharedRakNetNetwork) PingContext(ctx context.Context, address string) ([]byte, error) {
	return rakNetNetwork(nil).PingContext(ctx, address)
}

var _ minecraft.Network = sharedRakNetNetwork{}
