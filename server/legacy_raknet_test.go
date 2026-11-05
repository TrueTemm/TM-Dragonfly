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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type legacyClientNetwork struct{ minecraft.Network }

func (n legacyClientNetwork) DialContext(ctx context.Context, address string) (net.Conn, error) {
	c, err := n.Network.DialContext(ctx, address)
	if err != nil {
		return nil, err
	}
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		_ = c.Close()
		return nil, err
	}
	return &legacyWireClient{rakNetReplay: rakNetReplay{Conn: c}, key: key}, nil
}

func dialTestClient(d minecraft.Dialer, address string, timeout time.Duration) (*minecraft.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var network minecraft.Network = rakNetNetwork(nil)
	if d.Protocol != nil && d.Protocol.ID() < 649 {
		network = legacyClientNetwork{Network: network}
	}
	return d.DialContextNetwork(ctx, network, address)
}

type legacyWireClient struct {
	rakNetReplay
	key        *ecdsa.PrivateKey
	crypto     atomic.Pointer[legacyPacketCrypto]
	compressed atomic.Bool
	writeMu    sync.Mutex
}

func (*legacyWireClient) DisableEncryption() bool { return true }

func (c *legacyWireClient) Write(data []byte) (int, error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if !c.compressed.Load() {
		return c.Conn.Write(data)
	}
	if len(data) < 2 {
		return 0, fmt.Errorf("short client batch")
	}
	var body []byte
	var err error
	if data[1] == 0xff {
		body = data[2:]
	} else {
		body, err = packet.FlateCompression.Decompress(data[2:], packet.DefaultMaxDecompressedLen)
	}
	if err != nil {
		return 0, err
	}
	pks, err := unpackLegacyBatch(append([]byte{0xfe}, body...))
	if err != nil {
		return 0, err
	}
	for i, pk := range pks {
		if legacyPacketID(pk) == packet.IDLogin {
			buf := bytes.NewBuffer(pk)
			var header packet.Header
			_ = header.Read(buf)
			var request packet.Login
			if err := decodeLegacyPacket(buf, &request); err != nil {
				return 0, err
			}
			id, client, _, err := login.Parse(request.ConnectionRequest, nil)
			if err != nil {
				return 0, err
			}
			request.ConnectionRequest = login.EncodeOffline(id, client, c.key, true)
			pks[i] = encodeLegacyPacket(&request)
		}
	}
	compressed, err := packet.FlateCompression.Compress(packLegacyBatch(pks)[1:])
	if err != nil {
		return 0, err
	}
	if crypto := c.crypto.Load(); crypto != nil {
		compressed = crypto.seal(compressed)
	}
	_, err = c.Conn.Write(append([]byte{0xfe}, compressed...))
	return len(data), err
}

func (c *legacyWireClient) ReadPacket() ([]byte, error) {
	for {
		data, err := readRakNetBatch(c.Conn)
		if err != nil {
			return nil, err
		}
		if !c.compressed.Load() {
			c.compressed.Store(true)
			return data, nil
		}
		if len(data) < 2 || data[0] != 0xfe {
			return nil, fmt.Errorf("invalid legacy server batch")
		}
		body := data[1:]
		if crypto := c.crypto.Load(); crypto != nil {
			body, err = crypto.open(body)
			if err != nil {
				return nil, err
			}
		}
		body, err = packet.FlateCompression.Decompress(body, packet.DefaultMaxDecompressedLen)
		if err != nil {
			return nil, err
		}
		pks, err := unpackLegacyBatch(append([]byte{0xfe}, body...))
		if err != nil {
			return nil, err
		}
		if len(pks) == 1 && legacyPacketID(pks[0]) == packet.IDServerToClientHandshake {
			buf := bytes.NewBuffer(pks[0])
			var header packet.Header
			_ = header.Read(buf)
			var handshake packet.ServerToClientHandshake
			if err := decodeLegacyPacket(buf, &handshake); err != nil {
				return nil, err
			}
			token, err := jwt.ParseSigned(string(handshake.JWT), []jose.SignatureAlgorithm{jose.ES384})
			if err != nil {
				return nil, err
			}
			keyString, _ := token.Headers[0].ExtraHeaders["x5u"].(string)
			public := new(ecdsa.PublicKey)
			if err := login.ParsePublicKey(keyString, public); err != nil {
				return nil, err
			}
			var claims struct {
				Salt string `json:"salt"`
			}
			if err := token.Claims(public, &claims); err != nil {
				return nil, err
			}
			salt, err := base64.RawStdEncoding.DecodeString(claims.Salt)
			if err != nil {
				return nil, err
			}
			crypto, err := newLegacyPacketCrypto(c.key, public, salt)
			if err != nil {
				return nil, err
			}
			c.writeMu.Lock()
			c.crypto.Store(crypto)
			response := packLegacyBatch([][]byte{encodeLegacyPacket(&packet.ClientToServerHandshake{})})
			compressed, compressErr := packet.FlateCompression.Compress(response[1:])
			if compressErr == nil {
				_, compressErr = c.Conn.Write(append([]byte{0xfe}, crypto.seal(compressed)...))
			}
			c.writeMu.Unlock()
			if compressErr != nil {
				return nil, compressErr
			}
			continue
		}
		return append([]byte{0xfe, 0xff}, body...), nil
	}
}

func TestLegacyEncryptedChecksum(t *testing.T) {
	a, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	b, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	sender, err := newLegacyPacketCrypto(a, &b.PublicKey, []byte("legacy checksum"))
	if err != nil {
		t.Fatal(err)
	}
	receiver, err := newLegacyPacketCrypto(b, &a.PublicKey, []byte("legacy checksum"))
	if err != nil {
		t.Fatal(err)
	}
	for _, plain := range [][]byte{[]byte("first"), []byte("second")} {
		decoded, err := receiver.open(sender.seal(plain))
		if err != nil || !bytes.Equal(decoded, plain) {
			t.Fatalf("encrypted stream: %x %v", decoded, err)
		}
	}
	corrupt := sender.seal([]byte("tampered"))
	corrupt[0] ^= 1
	if _, err := receiver.open(corrupt); err == nil {
		t.Fatal("tampered ciphertext accepted")
	}
}
