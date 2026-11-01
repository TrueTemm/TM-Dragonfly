package server

import (
	"encoding/binary"
	"net"
	"sync"
	"time"
)

type gameUDPMux struct {
	udp    net.PacketConn
	raknet *gamePacketConn
	ice    *gamePacketConn
	mu     sync.Mutex
	peers  map[string]time.Time
}

type gameDatagram struct {
	payload []byte
	address net.Addr
}

type gamePacketConn struct {
	mux           *gameUDPMux
	packets       chan gameDatagram
	closed        chan struct{}
	changed       chan struct{}
	closeOnce     sync.Once
	mu            sync.Mutex
	readDeadline  time.Time
	writeDeadline time.Time
}

func newGameUDPMux(address string) (*gameUDPMux, error) {
	udp, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, err
	}
	m := &gameUDPMux{udp: udp, peers: make(map[string]time.Time)}
	m.raknet = &gamePacketConn{mux: m, packets: make(chan gameDatagram, 4096), closed: make(chan struct{}), changed: make(chan struct{}, 1)}
	m.ice = &gamePacketConn{mux: m, packets: make(chan gameDatagram, 4096), closed: make(chan struct{}), changed: make(chan struct{}, 1)}
	go m.receive()
	return m, nil
}

func (m *gameUDPMux) receive() {
	buf := make([]byte, 65535)
	for {
		n, addr, err := m.udp.ReadFrom(buf)
		if err != nil {
			return
		}
		p := append([]byte(nil), buf[:n]...)
		now := time.Now()
		key := addr.String()
		m.mu.Lock()
		until := m.peers[key]
		stun := isSTUN(p)
		ice := stun || now.Before(until) || isDTLS(p)
		if ice {
			m.peers[key] = now.Add(2 * time.Minute)
		}
		if len(m.peers) > 4096 {
			for k, expiry := range m.peers {
				if now.After(expiry) {
					delete(m.peers, k)
				}
			}
		}
		m.mu.Unlock()
		dest := m.raknet
		if ice {
			dest = m.ice
		}
		select {
		case <-dest.closed:
		case dest.packets <- gameDatagram{payload: p, address: addr}:
		default:
		}
	}
}

func isSTUN(p []byte) bool {
	return len(p) >= 20 && p[0]&0xc0 == 0 && binary.BigEndian.Uint32(p[4:8]) == 0x2112a442
}
func isDTLS(p []byte) bool { return len(p) >= 13 && p[0] >= 20 && p[0] <= 63 && p[1] == 0xfe }

func (c *gamePacketConn) ReadFrom(p []byte) (int, net.Addr, error) {
	for {
		select {
		case <-c.closed:
			return 0, nil, net.ErrClosed
		default:
		}
		c.mu.Lock()
		deadline := c.readDeadline
		c.mu.Unlock()
		var clock <-chan time.Time
		var timer *time.Timer
		if !deadline.IsZero() {
			if time.Until(deadline) <= 0 {
				return 0, nil, gameDeadline{}
			}
			timer = time.NewTimer(time.Until(deadline))
			clock = timer.C
		}
		select {
		case <-c.closed:
			if timer != nil {
				timer.Stop()
			}
			return 0, nil, net.ErrClosed
		case d := <-c.packets:
			if timer != nil {
				timer.Stop()
			}
			return copy(p, d.payload), d.address, nil
		case <-c.changed:
			if timer != nil {
				timer.Stop()
			}
		case <-clock:
			return 0, nil, gameDeadline{}
		}
	}
}

func (c *gamePacketConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	select {
	case <-c.closed:
		return 0, net.ErrClosed
	default:
	}
	c.mu.Lock()
	deadline := c.writeDeadline
	c.mu.Unlock()
	if !deadline.IsZero() && time.Now().After(deadline) {
		return 0, gameDeadline{}
	}
	return c.mux.udp.WriteTo(p, addr)
}

func (c *gamePacketConn) Close() error {
	c.closeOnce.Do(func() {
		close(c.closed)
		m := c.mux
		select {
		case <-m.raknet.closed:
			select {
			case <-m.ice.closed:
				_ = m.udp.Close()
			default:
			}
		default:
		}
	})
	return nil
}

func (c *gamePacketConn) LocalAddr() net.Addr { return c.mux.udp.LocalAddr() }
func (c *gamePacketConn) SetDeadline(t time.Time) error {
	c.mu.Lock()
	c.readDeadline = t
	c.writeDeadline = t
	c.mu.Unlock()
	select {
	case c.changed <- struct{}{}:
	default:
	}
	return nil
}
func (c *gamePacketConn) SetReadDeadline(t time.Time) error {
	c.mu.Lock()
	c.readDeadline = t
	c.mu.Unlock()
	select {
	case c.changed <- struct{}{}:
	default:
	}
	return nil
}
func (c *gamePacketConn) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	c.writeDeadline = t
	c.mu.Unlock()
	return nil
}

type gameDeadline struct{}

func (gameDeadline) Error() string   { return "i/o timeout" }
func (gameDeadline) Timeout() bool   { return true }
func (gameDeadline) Temporary() bool { return true }

var _ net.PacketConn = (*gamePacketConn)(nil)
var _ net.Error = gameDeadline{}
