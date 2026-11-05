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
	"reflect"
	"testing"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestLegacyVanillaComponentsAndFallingBlocks(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	for _, proto := range []int32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			p := Protocol{Version: proto}
			sg := p.ConvertFromLatest(&packet.StartGame{}, nil)[0].(*v766.StartGame)
			for _, entry := range sg.Items {
				if entry.ComponentBased {
					t.Fatalf("vanilla item %s requires an absent component definition", entry.Name)
				}
			}
			for _, pk := range []packet.Packet{&packet.JigsawStructureData{}, &packet.VoxelShapes{}, &packet.SetHud{}} {
				if len(p.ConvertFromLatest(pk, nil)) != 0 {
					t.Fatalf("unsupported %T sent to legacy client", pk)
				}
			}
			for rid := uint32(0); rid < uint32(world.DefaultBlockRegistry.BlockCount()); rid++ {
				pk := &packet.AddActor{EntityType: "minecraft:falling_block", EntityMetadata: protocol.EntityMetadata{protocol.EntityDataKeyVariant: int32(rid)}}
				out := p.ConvertFromLatest(pk, nil)[0].(*v786.AddActor)
				if out.EntityMetadata[protocol.EntityDataKeyVariant] != int32(blockpalette.HashFor(uint32(proto), rid)) || pk.EntityMetadata[protocol.EntityDataKeyVariant] != int32(rid) {
					t.Fatalf("falling block state %d was not translated independently", rid)
				}
			}
			mob := &packet.AddActor{EntityType: "minecraft:sheep", EntityMetadata: protocol.EntityMetadata{protocol.EntityDataKeyVariant: int32(15)}}
			if p.ConvertFromLatest(mob, nil)[0].(*v786.AddActor).EntityMetadata[protocol.EntityDataKeyVariant] != int32(15) {
				t.Fatal("mob variant treated as a block")
			}
		})
	}
}

func TestLegacyEveryItemWire(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	for _, proto := range []int32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			p := Protocol{Version: proto}
			checked := 0
			for _, entry := range itemdata.ItemsFor(uint32(proto)) {
				for _, meta := range []uint32{0, 1, 2, 3, 7, 8, 15, 16, 255, 32767, 65535} {
					id, aux := itemdata.ReverseTranslateStack(uint32(proto), int32(entry.RuntimeID), meta)
					if id == 0 {
						continue
					}
					stack := protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: id, MetadataValue: aux}, Count: 64, NBTData: map[string]any{"display": map[string]any{"Name": "Legacy audit"}}, CanBreak: []string{"minecraft:stone"}, CanBePlacedOn: []string{"minecraft:dirt"}, BlockingTick: 123}
					for _, instance := range []bool{false, true} {
						buf := new(bytes.Buffer)
						writer := p.NewWriter(buf, 0)
						var got protocol.ItemStack
						if instance {
							input := protocol.ItemInstance{Stack: stack, StackNetworkID: 7}
							writer.ItemInstance(&input)
							var output protocol.ItemInstance
							p.NewReader(buf, 0, true).ItemInstance(&output)
							if output.StackNetworkID != 7 {
								t.Fatal("stack network ID changed")
							}
							got = output.Stack
						} else {
							writer.Item(&stack)
							p.NewReader(buf, 0, true).Item(&got)
						}
						if got.NetworkID != id || got.MetadataValue != aux || got.Count != stack.Count || !reflect.DeepEqual(got.NBTData, stack.NBTData) || !reflect.DeepEqual(got.CanBreak, stack.CanBreak) || !reflect.DeepEqual(got.CanBePlacedOn, stack.CanBePlacedOn) || buf.Len() != 0 {
							t.Fatalf("%s:%d instance=%v lost item data", entry.Name, meta, instance)
						}
						if entry.Name == "minecraft:shield" && got.BlockingTick != 123 {
							t.Fatal("shield blocking tick lost")
						}
					}
					checked++
				}
			}
			t.Logf("checked %d item variants in stack and instance formats", checked)
		})
	}
}

func TestLegacyBlockInteractionWire(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	stone, _ := world.ItemRuntimeIDByName("minecraft:stone")
	rid, _ := world.DefaultBlockRegistry.StateToRuntimeID("minecraft:stone", nil)
	for _, proto := range []int32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			data, err := os.ReadFile(fmt.Sprintf("testdata/wire%d.json", proto))
			if err != nil {
				t.Fatal(err)
			}
			var fixtures map[string][]byte
			if err := json.Unmarshal(data, &fixtures); err != nil {
				t.Fatal(err)
			}
			p := Protocol{Version: proto}
			action := protocol.UseItemTransactionData{ActionType: 0, BlockPosition: protocol.BlockPos{1, -16, -2}, BlockFace: 1, HotBarSlot: 2,
				HeldItem: protocol.ItemInstance{StackNetworkID: 7, Stack: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: stone}, Count: 64, BlockRuntimeID: int32(rid), NBTData: map[string]any{"display": map[string]any{"Name": "Legacy audit"}}, CanBreak: []string{"minecraft:dirt"}}},
				Position: mgl32.Vec3{1, 65, -2}, ClickedPosition: mgl32.Vec3{0.5, 1, 0.5}, BlockRuntimeID: rid}
			buf := new(bytes.Buffer)
			p.NewWriter(buf, 0).PlayerInventoryAction(&action)
			if !bytes.Equal(buf.Bytes(), fixtures["player_inventory_action"]) || action.BlockRuntimeID != rid || action.HeldItem.Stack.NetworkID != stone {
				t.Fatal("block interaction differs from historical wire or changes its source")
			}
			var got protocol.UseItemTransactionData
			p.NewReader(buf, 0, true).PlayerInventoryAction(&got)
			if buf.Len() != 0 || got.BlockRuntimeID != rid || got.HeldItem.Stack.BlockRuntimeID != int32(rid) || got.HeldItem.Stack.NetworkID != stone || got.BlockPosition != action.BlockPosition {
				t.Fatal("block interaction did not recover native block and item IDs")
			}
			input := p.Packets(true)[packet.IDInventoryTransaction]()
			raw := bytes.NewBuffer(fixtures["use_item_transaction"])
			input.Marshal(p.NewReader(raw, 0, true))
			out := p.ConvertToLatest(input, nil)[0].(*packet.InventoryTransaction).TransactionData.(*protocol.UseItemTransactionData)
			if raw.Len() != 0 || out.BlockRuntimeID != rid || out.HeldItem.Stack.NetworkID != stone || out.BlockPosition != action.BlockPosition {
				t.Fatal("inventory transaction did not recover native IDs")
			}
			if proto < 630 {
				event := p.Packets(true)[53]()
				raw := bytes.NewBuffer(fixtures["crafting_event"])
				event.Marshal(p.NewReader(raw, 0, true))
				if raw.Len() != 0 || event.(*CraftingEvent).Input[0].Stack.NetworkID != stone || len(p.ConvertToLatest(event, nil)) != 0 {
					t.Fatal("deprecated crafting event was not decoded and ignored")
				}
				event.Marshal(p.NewWriter(raw, 0))
				if !bytes.Equal(raw.Bytes(), fixtures["crafting_event"]) {
					t.Fatal("crafting event differs from historical wire")
				}
			} else if p.Packets(true)[53] != nil {
				t.Fatal("removed crafting event accepted on 630")
			}
			for _, native := range []uint32{0, rid} {
				stack := action.HeldItem.Stack
				b, ok := world.DefaultBlockRegistry.BlockByRuntimeID(native)
				if !ok {
					t.Fatal("missing native block")
				}
				it, ok := b.(world.Item)
				if !ok {
					t.Fatal("native block has no item")
				}
				stack.NetworkID, _, _ = world.ItemRuntimeID(it)
				stack.BlockRuntimeID = int32(native)
				buf.Reset()
				p.NewWriter(buf, 0).Item(&stack)
				wire := bytes.NewBuffer(append([]byte(nil), buf.Bytes()...))
				var itemID, blockID int32
				var count uint16
				var meta uint32
				reader := protocol.NewReader(wire, 0, true)
				reader.Varint32(&itemID)
				reader.Uint16(&count)
				reader.Varuint32(&meta)
				reader.Varint32(&blockID)
				if uint32(blockID) != blockpalette.HashFor(uint32(proto), native) {
					t.Fatalf("item block state %d was not mapped", native)
				}
			}
		})
	}
}

func wireFixtures(t *testing.T) map[string][]byte {
	t.Helper()
	data, err := os.ReadFile("testdata/wire630.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string][]byte
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	return fixtures
}

func TestHistoricalWireFixtures(t *testing.T) {
	fixtures := wireFixtures(t)
	cases := map[string]packet.Packet{
		"level_chunk":         &packet.LevelChunk{Position: protocol.ChunkPos{-3, 5}, Dimension: 2, SubChunkCount: 2, RawPayload: []byte{9, 0, 0, 1, 2, 3}},
		"level_chunk_limited": &packet.LevelChunk{Position: protocol.ChunkPos{2, -8}, Dimension: 1, SubChunkCount: 0xfffffffe, SubChunkLimit: protocol.Option(int32(24))},
		"correction":          &packet.CorrectPlayerMovePrediction{Position: mgl32.Vec3{1, 65, -2}, Delta: mgl32.Vec3{0.25, -0.1, 1}, OnGround: true, Tick: 1234},
		"start_game":          &packet.StartGame{EntityUniqueID: 1, EntityRuntimeID: 1, PlayerPosition: mgl32.Vec3{0.5, 65, 0.5}, WorldSeed: 12345, WorldName: "Protocol 630", ServerAuthoritativeInventory: true, PropertyData: map[string]any{}},
		"player_list":         &packet.PlayerList{Entries: []protocol.PlayerListEntry{{EntityUniqueID: 1, Username: "Protocol630Test", Skin: protocol.Skin{ArmSize: protocol.ArmSizeWide, SkinColour: color.RGBA{}, Trusted: true}}}},
	}
	p := Protocol{}
	for name, pk := range cases {
		t.Run(name, func(t *testing.T) {
			converted := p.ConvertFromLatest(pk, nil)
			if len(converted) != 1 {
				t.Fatalf("converted to %d packets", len(converted))
			}
			buf := new(bytes.Buffer)
			converted[0].Marshal(p.NewWriter(buf, 0))
			if !bytes.Equal(buf.Bytes(), fixtures[name]) {
				t.Fatalf("wire differs from historical 630 codec: got %d bytes want %d", buf.Len(), len(fixtures[name]))
			}
		})
	}
}

func TestPlayerAuthInput630(t *testing.T) {
	p := Protocol{}
	raw := wireFixtures(t)["player_auth_input"]
	buf := bytes.NewBuffer(raw)
	pk := &PlayerAuthInput{}
	pk.Marshal(p.NewReader(buf, 0, true))
	if buf.Len() != 0 {
		t.Fatalf("%d unread bytes", buf.Len())
	}
	out := p.ConvertToLatest(pk, nil)[0].(*packet.PlayerAuthInput)
	actions, _ := out.BlockActions.Value()
	if out.Tick != 345 || out.Position != (mgl32.Vec3{1, 65, -2}) || len(actions) != 1 || actions[0].BlockPos != (protocol.BlockPos{1, 64, -2}) {
		t.Fatal("movement or block action lost")
	}
	encoded := new(bytes.Buffer)
	pk.Marshal(p.NewWriter(encoded, 0))
	if !bytes.Equal(encoded.Bytes(), raw) {
		t.Fatal("movement wire changed")
	}
	pk.InteractionModel = 1
	encoded.Reset()
	pk.Marshal(p.NewWriter(encoded, 0))
	if encoded.Bytes()[40] != 1 {
		t.Fatal("interaction model must use unsigned varint")
	}
}

func TestStartGameRegistries630(t *testing.T) {
	p := Protocol{}
	pk := p.ConvertFromLatest(&packet.StartGame{}, nil)[0].(*v766.StartGame)
	if pk.UseBlockNetworkIDHashes || pk.GameVersion != "1.20.50" || len(pk.Items) != 1461 {
		t.Fatal("incorrect start game registry")
	}
	if !blockpalette.Supported(630) || !itemdata.Supported(630) {
		t.Fatal("registry is not connected")
	}
	seen := map[int16]bool{}
	for _, it := range pk.Items {
		if seen[it.RuntimeID] {
			t.Fatalf("duplicate item ID %d", it.RuntimeID)
		}
		seen[it.RuntimeID] = true
	}
	for _, side := range []bool{false, true} {
		if _, ok := p.Packets(side)[v786.IDSetHud]; ok {
			t.Fatal("649-only packet accepted")
		}
	}
}

func TestCacheAndResourcePacks630(t *testing.T) {
	p := Protocol{}
	cache := &packet.ClientCacheStatus{Enabled: true}
	if p.ConvertToLatest(cache, nil)[0].(*packet.ClientCacheStatus).Enabled || !cache.Enabled {
		t.Fatal("cache policy or input packet changed")
	}
	custom := protocol.StackResourcePack{UUID: "custom-pack", Version: "2.0.0"}
	pk := &packet.ResourcePackStack{TexturePacks: []protocol.StackResourcePack{{UUID: "d34cfa4b-2ad1-453d-a0db-668b429a3ea0", Version: "1.26.40"}, custom}}
	out := p.ConvertFromLatest(pk, nil)[0].(*v786.ResourcePackStack)
	if out.BaseGameVersion != "1.20.50" || !reflect.DeepEqual(out.TexturePacks, []protocol.StackResourcePack{custom}) || len(pk.TexturePacks) != 2 {
		t.Fatal("resource pack stack changed incorrectly")
	}
}

func TestRegistriesNBT630(t *testing.T) {
	for _, data := range [][]byte{biomeDefinitions, entityIdentifiers} {
		var value map[string]any
		if err := nbt.UnmarshalEncoding(data, &value, nbt.NetworkLittleEndian); err != nil || len(value) == 0 {
			t.Fatalf("invalid registry: %v", err)
		}
	}
}

func TestBiomesPrecedeLegacyPlayerSpawn(t *testing.T) {
	for _, proto := range []int32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			p := Protocol{Version: proto}
			spawn := &packet.PlayStatus{Status: packet.PlayStatusPlayerSpawn}
			out := p.ConvertFromLatest(spawn, nil)
			if len(out) != 2 || out[1] != spawn {
				t.Fatal("player spawn must follow the full biome definitions")
			}
			biomes, ok := out[0].(*v786.BiomeDefinitionList)
			if !ok || !bytes.Equal(biomes.SerialisedBiomeDefinitions, registryFor(uint32(proto)).biomes) {
				t.Fatal("wrong biome registry before spawn")
			}
			var definitions map[string]map[string]any
			if err := nbt.UnmarshalEncoding(biomes.SerialisedBiomeDefinitions, &definitions, nbt.NetworkLittleEndian); err != nil {
				t.Fatal(err)
			}
			if len(definitions) < 60 {
				t.Fatalf("incomplete biome registry with %d entries", len(definitions))
			}
			data, err := os.ReadFile(fmt.Sprintf("testdata/wire%d.json", proto))
			if err != nil {
				t.Fatal(err)
			}
			var fixtures map[string][]byte
			if err := json.Unmarshal(data, &fixtures); err != nil {
				t.Fatal(err)
			}
			var requiredNames []string
			if err := json.Unmarshal(fixtures["required_biome_names"], &requiredNames); err != nil || len(requiredNames) < 60 {
				t.Fatal("historical biome defaults missing", err)
			}
			for _, name := range requiredNames {
				if _, ok := definitions[name]; !ok {
					t.Fatalf("required biome %s missing", name)
				}
			}
			t.Logf("%d biome definitions precede spawn and include all %d historical defaults", len(definitions), len(requiredNames))
			wire := new(bytes.Buffer)
			out[0].Marshal(p.NewWriter(wire, 0))
			if !bytes.Equal(wire.Bytes(), biomes.SerialisedBiomeDefinitions) {
				t.Fatal("biome packet must contain raw network NBT")
			}
			for _, status := range []int32{packet.PlayStatusLoginSuccess, packet.PlayStatusLoginFailedClient, packet.PlayStatusLoginFailedServer} {
				pk := &packet.PlayStatus{Status: status}
				if got := p.ConvertFromLatest(pk, nil); len(got) != 1 || got[0] != pk {
					t.Fatal("biome definitions added to a non spawn status")
				}
			}
		})
	}
}

func TestBlockAndItemTranslation630(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	stone, ok := world.DefaultBlockRegistry.StateToRuntimeID("minecraft:stone", nil)
	if !ok {
		t.Fatal("stone missing")
	}
	p := Protocol{}
	source := &packet.UpdateBlock{NewBlockRuntimeID: stone}
	out := p.ConvertFromLatest(source, nil)[0].(*v786.UpdateBlock)
	if out.NewBlockRuntimeID != 1774 || source.NewBlockRuntimeID != stone {
		t.Fatal("invalid block runtime ID")
	}
	id, ok := world.ItemRuntimeIDByName("minecraft:shield")
	if !ok {
		t.Fatal("shield missing")
	}
	item := protocol.ItemInstance{StackNetworkID: 4, Stack: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: id}, Count: 1, BlockingTick: 42, NBTData: map[string]any{"damage": int32(3)}}}
	buf := new(bytes.Buffer)
	p.NewWriter(buf, 0).ItemInstance(&item)
	var got protocol.ItemInstance
	p.NewReader(buf, 0, true).ItemInstance(&got)
	if got.Stack.NetworkID != id || got.Stack.BlockingTick != 42 || got.Stack.NBTData["damage"] != int32(3) || got.StackNetworkID != 4 || buf.Len() != 0 {
		t.Fatalf("invalid shield round trip: %+v", got)
	}
	var raw bytes.Buffer
	raw.Write([]byte{9, 1, 0, 1})
	_ = protocol.WriteVarint32(&raw, int32(stone))
	translated := blockpalette.RemapChunkBlockPalette(raw.Bytes(), 1, 630)
	want := bytes.NewBuffer([]byte{9, 1, 0, 1})
	_ = protocol.WriteVarint32(want, 1774)
	if !bytes.Equal(translated, want.Bytes()) {
		t.Fatal("chunk did not use 630 runtime IDs")
	}
	event := &packet.LevelEvent{EventType: packet.LevelEventParticlesCrackBlock, EventData: int32(stone) | (3 << 24)}
	ev := p.ConvertFromLatest(event, nil)[0].(*packet.LevelEvent)
	if ev.EventData != 1774|(3<<24) || event.EventData != int32(stone)|(3<<24) {
		t.Fatal("block face or event input lost")
	}
}
