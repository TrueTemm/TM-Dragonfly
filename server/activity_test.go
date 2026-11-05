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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestActivityPacketsOnRakNet(t *testing.T) {
	fixtures := map[int32]map[string][]byte{}
	for _, version := range []int32{630, 898} {
		data, err := os.ReadFile(fmt.Sprintf("../multiversion/conformance/testdata/activity%d.json", version))
		if err != nil {
			t.Fatal(err)
		}
		var values map[string][]byte
		if err := json.Unmarshal(data, &values); err != nil {
			t.Fatal(err)
		}
		fixtures[version] = values
	}
	cfg := (UserConfig{}).minecraftListenConfig(Config{Compression: packet.SnappyCompression, AuthDisabled: true, Allower: allower{}})
	for _, p := range cfg.AcceptedProtocols {
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			listener, err := cfg.ListenNetwork(legacyRakNet{Network: rakNetNetwork(nil)}, "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			pos := mgl32.Vec3{-16.5, 65.25, 3}
			source := []packet.Packet{
				&packet.LevelSoundEvent{SoundType: packet.SoundEventAttackStrong, Position: pos, ExtraData: -1, EntityType: "minecraft:player", BabyMob: true, DisableRelativeVolume: true, EntityUniqueID: -12345},
				&packet.LevelSoundEvent{SoundType: packet.SoundEventBreathe, Position: pos, ExtraData: -1, EntityType: "minecraft:player", BabyMob: true, DisableRelativeVolume: true, EntityUniqueID: -12345},
				&packet.Animate{ActionType: 1, EntityRuntimeID: 27, Data: 0.5, SwingSource: packet.AnimateSwingSourceAttack, Hand: 1},
			}
			done, release := make(chan error, 1), make(chan struct{})
			defer close(release)
			go func() {
				raw, err := listener.Accept()
				if err != nil {
					done <- err
					return
				}
				defer raw.Close()
				conn := raw.(*minecraft.Conn)
				if err := conn.StartGameTimeout(minecraft.GameData{WorldName: "Activity regression", EntityUniqueID: 1, EntityRuntimeID: 1}, 5*time.Second); err != nil {
					done <- err
					return
				}
				for _, pk := range source {
					if err := conn.WritePacket(pk); err != nil {
						done <- err
						return
					}
				}
				done <- conn.Flush()
				<-release
			}()
			type frame struct {
				id   uint32
				data []byte
			}
			received := make(chan frame, 8)
			dialer := minecraft.Dialer{Protocol: spawnProbeProtocol{Protocol: p}, EnableLegacyAuth: true, DisablePortFollowing: true, PacketFunc: func(h packet.Header, payload []byte, src, dst net.Addr) {
				if src.String() == listener.Addr().String() && (h.PacketID == packet.IDLevelSoundEvent || h.PacketID == packet.IDAnimate) {
					received <- frame{id: h.PacketID, data: append([]byte(nil), payload...)}
				}
			}}
			client, err := dialTestClient(dialer, listener.Addr().String(), 5*time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			if err := client.DoSpawnTimeout(5 * time.Second); err != nil {
				t.Fatal(err)
			}
			if err := <-done; err != nil {
				t.Fatal(err)
			}
			for i, name := range []string{"sound_attack", "sound_breathe", "animate_attack"} {
				var got frame
				select {
				case got = <-received:
				case <-time.After(5 * time.Second):
					t.Fatal("activity packet timeout")
				}
				want := fixtures[630][name]
				if i < 2 && p.ID() >= 786 {
					want = fixtures[898][name]
				}
				if i < 2 && p.ID() == 975 {
					want = append(append([]byte(nil), want...), 0)
				}
				if i == 2 && p.ID() == 859 {
					b := fixtures[898]["animate_swing"]
					want = append([]byte{2}, b[1:len(b)-1]...)
				}
				if i == 2 && p.ID() >= 898 {
					want = fixtures[898][name]
				}
				if (i < 2 && p.ID() >= 1001) || (i == 2 && p.ID() >= 2225) {
					buf := new(bytes.Buffer)
					source[i].Marshal(protocol.NewWriter(buf, 0))
					want = buf.Bytes()
				}
				if got.id != source[i].ID() || !bytes.Equal(got.data, want) {
					t.Fatalf("%s malformed: got %x want %x", name, got.data, want)
				}
			}
		})
	}
}
