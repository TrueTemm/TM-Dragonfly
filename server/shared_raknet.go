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
	return minecraft.RakNet{}.DialContext(ctx, address)
}
func (n sharedRakNetNetwork) PingContext(ctx context.Context, address string) ([]byte, error) {
	return minecraft.RakNet{}.PingContext(ctx, address)
}

var _ minecraft.Network = sharedRakNetNetwork{}
