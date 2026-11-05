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
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/multiversion/v630"
	"github.com/df-mc/dragonfly/multiversion/v686"
	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type spawnProbeProtocol struct{ minecraft.Protocol }

func (p spawnProbeProtocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if _, ok := pk.(*v686.ResourcePacksInfo); ok {
		return []packet.Packet{&packet.ResourcePacksInfo{}}
	}
	return p.Protocol.ConvertToLatest(pk, conn)
}

func TestLegacySpawnPacketOrder(t *testing.T) {
	for _, id := range []int32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			conf := Config{Compression: packet.SnappyCompression, AuthDisabled: true, Allower: allower{}}
			cfg := (UserConfig{}).minecraftListenConfig(conf)
			listener, err := cfg.ListenNetwork(legacyRakNet{Network: rakNetNetwork(nil)}, "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			finished, release := make(chan error, 1), make(chan struct{})
			defer close(release)
			go func() {
				conn, err := listener.Accept()
				if err != nil {
					finished <- err
					return
				}
				defer conn.Close()
				data := minecraft.GameData{WorldName: "Legacy spawn regression", EntityUniqueID: 1, EntityRuntimeID: 1}
				finished <- conn.(*minecraft.Conn).StartGameTimeout(data, 5*time.Second)
				<-release
			}()
			type event struct {
				id    uint32
				count int
				err   error
			}
			events := make(chan event, 8)
			dialer := minecraft.Dialer{EnableLegacyAuth: true, Protocol: spawnProbeProtocol{Protocol: v630.Protocol{Version: id}}, DisablePortFollowing: true, PacketFunc: func(header packet.Header, payload []byte, src, dst net.Addr) {
				if src.String() != listener.Addr().String() {
					return
				}
				switch header.PacketID {
				case packet.IDBiomeDefinitionList:
					var definitions map[string]any
					err := nbt.UnmarshalEncoding(payload, &definitions, nbt.NetworkLittleEndian)
					events <- event{id: header.PacketID, count: len(definitions), err: err}
				case packet.IDPlayStatus:
					var status packet.PlayStatus
					status.Marshal(protocol.NewReader(bytes.NewBuffer(payload), 0, true))
					if status.Status == packet.PlayStatusPlayerSpawn {
						events <- event{id: header.PacketID}
					}
				}
			}}
			conn, err := dialTestClient(dialer, listener.Addr().String(), 5*time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			if err := conn.DoSpawnTimeout(5 * time.Second); err != nil {
				t.Fatal(err)
			}
			for i, want := range []uint32{packet.IDBiomeDefinitionList, packet.IDPlayStatus} {
				select {
				case got := <-events:
					if got.id != want || got.err != nil || (i == 0 && got.count != 86) {
						t.Fatalf("spawn packet %d: got %+v want ID %d", i, got, want)
					}
				case <-time.After(5 * time.Second):
					t.Fatal("spawn sequence timed out")
				}
			}
			if err := <-finished; err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLegacyCompressionNegotiation(t *testing.T) {
	conf := Config{Compression: packet.SnappyCompression, AuthDisabled: true, Allower: allower{}}
	cfg := (UserConfig{}).minecraftListenConfig(conf)
	listener, err := cfg.ListenNetwork(legacyRakNet{Network: rakNetNetwork(nil)}, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	for _, id := range []int32{589, 594, 598, 618, 622, 630, 649, 2225} {
		checkNetworkSettings(t, listener.Addr().String(), id)
	}
}

func checkNetworkSettings(t *testing.T, address string, id int32) {
	t.Helper()
	conn, err := raknet.DialTimeout(address, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := new(bytes.Buffer)
	header := packet.Header{PacketID: packet.IDRequestNetworkSettings}
	if err := header.Write(buf); err != nil {
		t.Fatal(err)
	}
	request := packet.RequestNetworkSettings{ClientProtocol: id}
	request.Marshal(protocol.NewWriter(buf, 0))
	if err := packet.NewEncoder(conn).Encode([][]byte{buf.Bytes()}); err != nil {
		t.Fatal(err)
	}
	packets, err := packet.NewDecoder(conn).Decode()
	if err != nil {
		t.Fatal(err)
	}
	if len(packets) != 1 {
		t.Fatalf("protocol %d: received %d packets", id, len(packets))
	}
	payload := bytes.NewBuffer(packets[0])
	if err := header.Read(payload); err != nil {
		t.Fatal(err)
	}
	if header.PacketID != packet.IDNetworkSettings {
		t.Fatalf("protocol %d: packet ID %d", id, header.PacketID)
	}
	var settings packet.NetworkSettings
	settings.Marshal(protocol.NewReader(payload, 0, true))
	wantAlgorithm, wantThreshold := uint16(packet.CompressionAlgorithmSnappy), uint16(256)
	if id < 649 {
		wantAlgorithm, wantThreshold = packet.CompressionAlgorithmFlate, 1
	}
	if payload.Len() != 0 || settings.CompressionAlgorithm != wantAlgorithm || settings.CompressionThreshold != wantThreshold {
		t.Errorf("protocol %d: got algorithm=%d threshold=%d trailing=%d want algorithm=%d threshold=%d", id, settings.CompressionAlgorithm, settings.CompressionThreshold, payload.Len(), wantAlgorithm, wantThreshold)
	}
}
