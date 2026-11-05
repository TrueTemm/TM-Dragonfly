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
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type projectRakNet struct{ log *slog.Logger }

func rakNetNetwork(log *slog.Logger) minecraft.Network {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	return projectRakNet{log: log}
}
func (n projectRakNet) DialContext(ctx context.Context, address string) (net.Conn, error) {
	return (raknet.Dialer{ErrorLog: n.log}).DialContext(ctx, address)
}
func (n projectRakNet) PingContext(ctx context.Context, address string) ([]byte, error) {
	return (raknet.Dialer{ErrorLog: n.log}).PingContext(ctx, address)
}
func (n projectRakNet) Listen(address string) (minecraft.NetworkListener, error) {
	return (raknet.ListenConfig{ErrorLog: n.log}).Listen(address)
}

type legacyRakNet struct{ minecraft.Network }

func (n legacyRakNet) Listen(address string) (minecraft.NetworkListener, error) {
	l, err := n.Network.Listen(address)
	if err != nil {
		return nil, err
	}
	w := &legacyRakNetListener{NetworkListener: l, ready: make(chan net.Conn), done: make(chan struct{})}
	go w.accept()
	return w, nil
}

type legacyRakNetListener struct {
	minecraft.NetworkListener
	ready chan net.Conn
	done  chan struct{}
	once  sync.Once
}

func (l *legacyRakNetListener) accept() {
	for {
		c, err := l.NetworkListener.Accept()
		if err != nil {
			_ = l.Close()
			return
		}
		go func() {
			_ = c.SetReadDeadline(time.Now().Add(10 * time.Second))
			first, err := readRakNetBatch(c)
			if err != nil {
				_ = c.Close()
				return
			}
			_ = c.SetReadDeadline(time.Time{})
			pks, err := unpackLegacyBatch(first)
			if err != nil || len(pks) != 1 || legacyPacketID(pks[0]) != packet.IDRequestNetworkSettings {
				_ = c.Close()
				return
			}
			buf := bytes.NewBuffer(pks[0])
			var header packet.Header
			_ = header.Read(buf)
			if buf.Len() != 4 {
				_ = c.Close()
				return
			}
			version := int32(binary.BigEndian.Uint32(buf.Bytes()))
			var prepared net.Conn = &rakNetReplay{Conn: c, first: first}
			if version < 649 {
				prepared = newLegacyRakNetConn(c, first)
			}
			select {
			case l.ready <- prepared:
			case <-l.done:
				_ = prepared.Close()
			}
		}()
	}
}

func (l *legacyRakNetListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.ready:
		return c, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *legacyRakNetListener) Close() (err error) {
	l.once.Do(func() { close(l.done); err = l.NetworkListener.Close() })
	return
}

type rakNetReplay struct {
	net.Conn
	first []byte
}

func (c *rakNetReplay) ReadPacket() ([]byte, error) {
	if c.first != nil {
		data := c.first
		c.first = nil
		return data, nil
	}
	return readRakNetBatch(c.Conn)
}

func (c *rakNetReplay) Latency() time.Duration {
	if lc, ok := c.Conn.(interface{ Latency() time.Duration }); ok {
		return lc.Latency()
	}
	return 0
}

func (c *rakNetReplay) Context() context.Context {
	if cc, ok := c.Conn.(interface{ Context() context.Context }); ok {
		return cc.Context()
	}
	return context.Background()
}

func readRakNetBatch(c net.Conn) ([]byte, error) {
	if pr, ok := c.(packet.PacketReader); ok {
		return pr.ReadPacket()
	}
	buf := make([]byte, 3*1024*1024)
	n, err := c.Read(buf)
	return buf[:n], err
}

type legacyBatchResult struct {
	data []byte
	err  error
}

type legacyRakNetConn struct {
	rakNetReplay
	reads      chan legacyBatchResult
	done       chan struct{}
	once       sync.Once
	writeMu    sync.Mutex
	compressed atomic.Bool
	clientKey  atomic.Pointer[ecdsa.PublicKey]
	crypto     atomic.Pointer[legacyPacketCrypto]
	handshake  chan struct{}
}

func newLegacyRakNetConn(c net.Conn, first []byte) *legacyRakNetConn {
	w := &legacyRakNetConn{rakNetReplay: rakNetReplay{Conn: c}, reads: make(chan legacyBatchResult, 8), done: make(chan struct{}), handshake: make(chan struct{})}
	go w.readBatches(first)
	return w
}

func (*legacyRakNetConn) DisableEncryption() bool { return true } // encryption is handled before batch conversion

func (c *legacyRakNetConn) Close() (err error) {
	c.once.Do(func() { close(c.done); err = c.Conn.Close() })
	return
}

func (c *legacyRakNetConn) ReadPacket() ([]byte, error) {
	select {
	case result := <-c.reads:
		return result.data, result.err
	case <-c.done:
		return nil, net.ErrClosed
	}
}

func (c *legacyRakNetConn) readBatches(first []byte) {
	data := first
	for {
		converted, err := c.incoming(data)
		if err != nil || len(converted) != 0 {
			select {
			case c.reads <- legacyBatchResult{converted, err}:
			case <-c.done:
				return
			}
		}
		if err != nil {
			return
		}
		data, err = readRakNetBatch(c.Conn)
		if err != nil {
			select {
			case c.reads <- legacyBatchResult{err: err}:
			case <-c.done:
			}
			return
		}
	}
}

func (c *legacyRakNetConn) incoming(data []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != 0xfe {
		return nil, fmt.Errorf("invalid legacy batch header")
	}
	if !c.compressed.Load() {
		return data, nil
	}
	body := data[1:]
	crypto := c.crypto.Load()
	if crypto != nil {
		var err error
		body, err = crypto.open(body)
		if err != nil {
			return nil, err
		}
	}
	body, err := packet.FlateCompression.Decompress(body, packet.DefaultMaxDecompressedLen)
	if err != nil {
		return nil, fmt.Errorf("legacy deflate: %w", err)
	}
	select {
	case <-c.handshake:
		return append([]byte{0xfe, 0xff}, body...), nil
	default:
	}
	pks, err := unpackLegacyBatch(append([]byte{0xfe}, body...))
	if err != nil {
		return nil, err
	}
	if len(pks) > 812 {
		return nil, fmt.Errorf("legacy batch exceeds packet limit")
	}
	filtered := pks[:0]
	for _, pk := range pks {
		switch legacyPacketID(pk) {
		case packet.IDLogin:
			buf := bytes.NewBuffer(pk)
			var header packet.Header
			_ = header.Read(buf)
			var request packet.Login
			if err := decodeLegacyPacket(buf, &request); err != nil {
				return nil, err
			}
			_, _, auth, err := login.Parse(request.ConnectionRequest, nil)
			if err != nil {
				return nil, fmt.Errorf("legacy login: %w", err)
			}
			c.clientKey.Store(auth.PublicKey)
		case packet.IDClientToServerHandshake:
			if crypto == nil || len(pk) != 1 {
				return nil, fmt.Errorf("invalid legacy encrypted handshake")
			}
			select {
			case <-c.handshake:
				return nil, fmt.Errorf("duplicate legacy encrypted handshake")
			default:
				close(c.handshake)
				_ = c.Conn.SetReadDeadline(time.Time{})
			}
			continue
		}
		if crypto != nil {
			select {
			case <-c.handshake:
			default:
				return nil, fmt.Errorf("legacy encrypted handshake required")
			}
		}
		filtered = append(filtered, pk)
	}
	if len(filtered) == 0 {
		return nil, nil
	}
	out := packLegacyBatch(filtered)
	return append([]byte{0xfe, 0xff}, out[1:]...), nil
}

func (c *legacyRakNetConn) Write(data []byte) (int, error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if !c.compressed.Load() {
		pks, err := unpackLegacyBatch(data)
		if err != nil {
			return 0, err
		}
		for i, pk := range pks {
			if legacyPacketID(pk) == packet.IDNetworkSettings {
				buf := bytes.NewBuffer(pk)
				var header packet.Header
				_ = header.Read(buf)
				var settings packet.NetworkSettings
				if err := decodeLegacyPacket(buf, &settings); err != nil {
					return 0, err
				}
				settings.CompressionThreshold = 1
				settings.CompressionAlgorithm = packet.CompressionAlgorithmFlate
				pks[i] = encodeLegacyPacket(&settings)
				c.compressed.Store(true)
			}
		}
		_, err = c.Conn.Write(packLegacyBatch(pks))
		return len(data), err
	}
	if len(data) < 2 || data[0] != 0xfe {
		return 0, fmt.Errorf("invalid upstream batch")
	}
	var body, compressed []byte
	var err error
	if data[1] == 0xff {
		body = data[2:]
		compressed, err = packet.FlateCompression.Compress(body)
	} else if data[1] == byte(packet.CompressionAlgorithmFlate) {
		compressed = data[2:]
		if c.crypto.Load() == nil && c.clientKey.Load() != nil {
			body, err = packet.FlateCompression.Decompress(compressed, packet.DefaultMaxDecompressedLen)
		}
	} else {
		err = fmt.Errorf("unexpected legacy compression %d", data[1])
	}
	if err != nil {
		return 0, err
	}
	if c.crypto.Load() == nil && c.clientKey.Load() != nil {
		pks, err := unpackLegacyBatch(append([]byte{0xfe}, body...))
		if err != nil {
			return 0, err
		}
		for _, pk := range pks {
			if legacyPacketID(pk) == packet.IDPlayStatus && len(pk) == 5 && int32(binary.BigEndian.Uint32(pk[1:])) == packet.PlayStatusLoginSuccess {
				if err := c.startEncryption(); err != nil {
					return 0, err
				}
				break
			}
		}
	}
	if crypto := c.crypto.Load(); crypto != nil {
		select {
		case <-c.handshake:
		case <-c.done:
			return 0, net.ErrClosed
		case <-time.After(10 * time.Second):
			_ = c.Close()
			return 0, fmt.Errorf("legacy encrypted handshake timed out")
		}
		compressed = crypto.seal(compressed)
	}
	_, err = c.Conn.Write(append([]byte{0xfe}, compressed...))
	return len(data), err
}

func (c *legacyRakNetConn) startEncryption() error {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	crypto, err := newLegacyPacketCrypto(key, c.clientKey.Load(), salt)
	if err != nil {
		return err
	}
	signer, err := jose.NewSigner(jose.SigningKey{Key: key, Algorithm: jose.ES384}, &jose.SignerOptions{ExtraHeaders: map[jose.HeaderKey]any{"x5u": login.MarshalPublicKey(&key.PublicKey)}})
	if err != nil {
		return err
	}
	token, err := jwt.Signed(signer).Claims(struct {
		Salt string `json:"salt"`
	}{base64.RawStdEncoding.EncodeToString(salt)}).Serialize()
	if err != nil {
		return err
	}
	batch := packLegacyBatch([][]byte{encodeLegacyPacket(&packet.ServerToClientHandshake{JWT: []byte(token)})})
	compressed, err := packet.FlateCompression.Compress(batch[1:])
	if err != nil {
		return err
	}
	c.crypto.Store(crypto)
	_ = c.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = c.Conn.Write(append([]byte{0xfe}, compressed...))
	return err
}

type legacyPacketCrypto struct {
	key     [32]byte
	send    cipher.Stream
	receive cipher.Stream
	sent    uint64
	read    uint64
}

func newLegacyPacketCrypto(key *ecdsa.PrivateKey, peer *ecdsa.PublicKey, salt []byte) (*legacyPacketCrypto, error) {
	private, err := key.ECDH()
	if err != nil {
		return nil, err
	}
	public, err := peer.ECDH()
	if err != nil {
		return nil, err
	}
	secret, err := private.ECDH(public)
	if err != nil {
		return nil, err
	}
	w := &legacyPacketCrypto{key: sha256.Sum256(append(append([]byte(nil), salt...), secret...))}
	block, err := aes.NewCipher(w.key[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	copy(iv, w.key[:12])
	iv[15] = 2
	w.send, w.receive = cipher.NewCTR(block, iv), cipher.NewCTR(block, iv)
	return w, nil
}

func (c *legacyPacketCrypto) checksum(counter uint64, data []byte) []byte {
	h := sha256.New()
	var sequence [8]byte
	binary.LittleEndian.PutUint64(sequence[:], counter)
	h.Write(sequence[:])
	h.Write(data)
	h.Write(c.key[:])
	return h.Sum(nil)[:8]
}

func (c *legacyPacketCrypto) seal(data []byte) []byte {
	out := append(bytes.Clone(data), c.checksum(c.sent, data)...)
	c.sent++
	c.send.XORKeyStream(out, out)
	return out
}

func (c *legacyPacketCrypto) open(data []byte) ([]byte, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("short encrypted legacy batch")
	}
	out := bytes.Clone(data)
	c.receive.XORKeyStream(out, out)
	body, sum := out[:len(out)-8], out[len(out)-8:]
	if subtle.ConstantTimeCompare(sum, c.checksum(c.read, body)) != 1 {
		return nil, fmt.Errorf("invalid encrypted legacy checksum")
	}
	c.read++
	return body, nil
}

func unpackLegacyBatch(data []byte) ([][]byte, error) {
	if len(data) == 0 || data[0] != 0xfe {
		return nil, fmt.Errorf("invalid batch header")
	}
	buf := bytes.NewReader(data[1:])
	var pks [][]byte
	for buf.Len() != 0 {
		n, err := binary.ReadUvarint(buf)
		if err != nil || n == 0 || n > uint64(buf.Len()) {
			return nil, fmt.Errorf("invalid legacy batch length")
		}
		pk := make([]byte, int(n))
		_, _ = io.ReadFull(buf, pk)
		pks = append(pks, pk)
	}
	return pks, nil
}

func packLegacyBatch(pks [][]byte) []byte {
	buf := bytes.NewBuffer([]byte{0xfe})
	var length [10]byte
	for _, pk := range pks {
		n := binary.PutUvarint(length[:], uint64(len(pk)))
		buf.Write(length[:n])
		buf.Write(pk)
	}
	return buf.Bytes()
}

func legacyPacketID(data []byte) uint32 {
	id, n := binary.Uvarint(data)
	if n <= 0 {
		return ^uint32(0)
	}
	return uint32(id) & 0x3ff
}

func encodeLegacyPacket(pk packet.Packet) []byte {
	buf := new(bytes.Buffer)
	_ = (&packet.Header{PacketID: pk.ID()}).Write(buf)
	pk.Marshal(protocol.NewWriter(buf, 0))
	return buf.Bytes()
}

func decodeLegacyPacket(buf *bytes.Buffer, pk packet.Packet) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("invalid legacy packet: %v", recovered)
		}
	}()
	pk.Marshal(protocol.NewReader(buf, 0, true))
	if buf.Len() != 0 {
		return fmt.Errorf("trailing legacy packet bytes")
	}
	return nil
}
