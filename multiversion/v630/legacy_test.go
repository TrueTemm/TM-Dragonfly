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

package v630

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestHistoricalLegacyWire(t *testing.T) {
	for _, id := range []int32{589, 594, 618, 622} {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			p := Protocol{Version: id}
			data, err := os.ReadFile(fmt.Sprintf("testdata/wire%d.json", id))
			if err != nil {
				t.Fatal(err)
			}
			var fixtures map[string][]byte
			if err := json.Unmarshal(data, &fixtures); err != nil {
				t.Fatal(err)
			}
			cases := map[string]packet.Packet{
				"level_chunk":         &packet.LevelChunk{Position: protocol.ChunkPos{-3, 5}, Dimension: 2, SubChunkCount: 2, RawPayload: []byte{9, 0, 0, 1, 2, 3}},
				"level_chunk_limited": &packet.LevelChunk{Position: protocol.ChunkPos{2, -8}, Dimension: 1, SubChunkCount: 0xfffffffe, SubChunkLimit: protocol.Option(int32(24))},
				"correction":          &packet.CorrectPlayerMovePrediction{Position: mgl32.Vec3{1, 65, -2}, Delta: mgl32.Vec3{0.25, -0.1, 1}, OnGround: true, Tick: 1234},
				"start_game":          &packet.StartGame{EntityUniqueID: 1, EntityRuntimeID: 1, PlayerPosition: mgl32.Vec3{0.5, 65, 0.5}, WorldSeed: 12345, WorldName: "Protocol 630", EditorWorldType: 1, ServerAuthoritativeInventory: true, PropertyData: map[string]any{}},
				"player_list":         &packet.PlayerList{Entries: []protocol.PlayerListEntry{{EntityUniqueID: 1, Username: "Protocol630Test", Skin: protocol.Skin{ArmSize: protocol.ArmSizeWide, SkinColour: color.RGBA{}, Trusted: true}}}},
				"resource_packs_info": &packet.ResourcePacksInfo{},
				"disconnect":          &packet.Disconnect{Message: "Legacy protocol test"},
				"commands":            &packet.AvailableCommands{Commands: []protocol.Command{{Name: "probe", Description: "sample", AliasesOffset: 0xffffffff, Overloads: []protocol.CommandOverload{{Parameters: []protocol.CommandParameter{{Name: "message", Type: protocol.CommandArgValid | protocol.CommandArgTypeString}}}}}}},
			}
			for name, pk := range cases {
				t.Run(name, func(t *testing.T) {
					out := p.ConvertFromLatest(pk, nil)
					if len(out) != 1 {
						t.Fatalf("converted into %d packets", len(out))
					}
					buf := new(bytes.Buffer)
					out[0].Marshal(p.NewWriter(buf, 0))
					expected := fixtures[name]
					if !bytes.Equal(buf.Bytes(), expected) {
						index := 0
						for index < min(buf.Len(), len(expected)) && buf.Bytes()[index] == expected[index] {
							index++
						}
						t.Fatalf("wire differs at %d: got %d bytes want %d", index, buf.Len(), len(expected))
					}
				})
			}
			raw := bytes.NewBuffer(fixtures["player_auth_input"])
			input := &PlayerAuthInput{}
			input.Marshal(p.NewReader(raw, 0, true))
			out := p.ConvertToLatest(input, nil)[0].(*packet.PlayerAuthInput)
			if raw.Len() != 0 || out.Tick != 345 {
				t.Fatal("legacy movement was not decoded")
			}
		})
	}
}

func TestHistoricalPacketPools(t *testing.T) {
	for _, id := range []int32{589, 594, 618, 622, 630} {
		data, err := os.ReadFile(fmt.Sprintf("testdata/wire%d.json", id))
		if err != nil {
			t.Fatal(err)
		}
		var fixtures map[string][]byte
		if err := json.Unmarshal(data, &fixtures); err != nil {
			t.Fatal(err)
		}
		p := Protocol{Version: id}
		historical := make(map[uint32]string)
		for _, name := range []string{"client_pool", "server_pool"} {
			var pool map[uint32]string
			if err := json.Unmarshal(fixtures[name], &pool); err != nil {
				t.Fatal(err)
			}
			for packetID, kind := range pool {
				historical[packetID] = kind
			}
		}
		pool := p.Packets(true)
		for packetID, constructor := range p.Packets(false) {
			pool[packetID] = constructor
		}
		for packetID, kind := range historical {
			if pool[packetID] == nil && packetID != 127 && packetID != 128 {
				t.Errorf("protocol %d missing %d %s", id, packetID, kind)
			}
		}
		for packetID := range pool {
			if _, ok := historical[packetID]; !ok {
				t.Errorf("protocol %d unexpected %d %T", id, packetID, pool[packetID]())
			}
		}
	}
}
