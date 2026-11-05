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
	"net"
	"testing"
	"time"
)

func TestGameUDPMuxRoutesPackets(t *testing.T) {
	mux, err := newGameUDPMux("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer mux.raknet.Close()
	defer mux.ice.Close()
	clientA, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer clientA.Close()
	clientB, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer clientB.Close()

	raknetPing := []byte{0x01, 0x02, 0x03}
	if _, err = clientA.WriteTo(raknetPing, mux.udp.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	stun := make([]byte, 20)
	stun[4], stun[5], stun[6], stun[7] = 0x21, 0x12, 0xa4, 0x42
	if _, err = clientB.WriteTo(stun, mux.udp.LocalAddr()); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 512)
	_ = mux.raknet.SetReadDeadline(time.Now().Add(time.Second))
	n, _, err := mux.raknet.ReadFrom(buf)
	if err != nil || !bytes.Equal(buf[:n], raknetPing) {
		t.Fatalf("RakNet packet = %x err = %v", buf[:n], err)
	}
	_ = mux.ice.SetReadDeadline(time.Now().Add(time.Second))
	n, _, err = mux.ice.ReadFrom(buf)
	if err != nil || !bytes.Equal(buf[:n], stun) {
		t.Fatalf("STUN packet = %x err = %v", buf[:n], err)
	}

	dtls := []byte{22, 0xfe, 0xfd, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if _, err = clientB.WriteTo(dtls, mux.udp.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	n, _, err = mux.ice.ReadFrom(buf)
	if err != nil || !bytes.Equal(buf[:n], dtls) {
		t.Fatalf("DTLS packet = %x err = %v", buf[:n], err)
	}
}

func TestGameUDPMuxReadDeadline(t *testing.T) {
	mux, err := newGameUDPMux("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer mux.raknet.Close()
	defer mux.ice.Close()
	_ = mux.ice.SetReadDeadline(time.Now().Add(15 * time.Millisecond))
	_, _, err = mux.ice.ReadFrom(make([]byte, 64))
	if e, ok := err.(net.Error); !ok || !e.Timeout() {
		t.Fatalf("expected timeout, got %v", err)
	}
}
