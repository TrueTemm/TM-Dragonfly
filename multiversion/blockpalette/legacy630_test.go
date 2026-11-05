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

package blockpalette

import (
	"bytes"
	"fmt"
	"testing"

	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestLegacyEveryStateAndChunkStorage(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	for proto, count := range map[uint32]uint32{589: 13141, 594: 13141, 598: 13141, 618: 13129, 622: 13121, 630: 14127} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			p := statePalettes[proto]
			p.build()
			if len(p.toLatest) != int(count) {
				t.Fatalf("mapped %d of %d old states", len(p.toLatest), count)
			}
			for old := uint32(0); old < count; old++ {
				native, ok := p.toLatest[old]
				if !ok || RuntimeIDFor(proto, HashFor(proto, native)) != native {
					t.Fatalf("state %d does not preserve its upgraded state", old)
				}
			}
			for _, bits := range []byte{0, 1, 2, 3, 4, 5, 6, 8, 16} {
				for _, version := range []byte{8, 9} {
					for base := 0; base < len(p.toOld); {
						size := min(1<<bits, 4096, len(p.toOld)-base)
						raw := new(bytes.Buffer)
						raw.Write([]byte{version, 2})
						if version == 9 {
							raw.WriteByte(0xfc)
						}
						for layer := 0; layer < 2; layer++ {
							raw.WriteByte(bits<<1 | 1)
							raw.Write(make([]byte, paletteWordCount(bits)*4))
							if bits > 0 {
								_ = protocol.WriteVarint32(raw, int32(size))
							}
							for rid := base; rid < base+size; rid++ {
								_ = protocol.WriteVarint32(raw, int32(rid))
							}
						}
						trailer := []byte{0xab, 0xcd, 0xef}
						before := append([]byte(nil), raw.Bytes()...)
						raw.Write(trailer)
						mapped := bytes.NewBuffer(RemapChunkBlockPalette(raw.Bytes(), 1, proto))
						header := 2
						if version == 9 {
							header++
						}
						if !bytes.Equal(mapped.Next(header), before[:header]) {
							t.Fatal("subchunk header changed")
						}
						for layer := 0; layer < 2; layer++ {
							if b, _ := mapped.ReadByte(); b != bits<<1|1 {
								t.Fatal("storage header changed")
							}
							mapped.Next(paletteWordCount(bits) * 4)
							n := int32(1)
							if bits > 0 {
								if err := protocol.Varint32(mapped, &n); err != nil || n != int32(size) {
									t.Fatal("palette size changed", err)
								}
							}
							for rid := base; rid < base+size; rid++ {
								var got int32
								if err := protocol.Varint32(mapped, &got); err != nil || uint32(got) != p.toOld[rid] || uint32(got) >= count {
									t.Fatalf("bits %d state %d has invalid wire ID %d: %v", bits, rid, got, err)
								}
							}
						}
						if !bytes.Equal(mapped.Bytes(), trailer) {
							t.Fatal("biome trailer changed")
						}
						base += size
					}
				}
			}
			t.Logf("checked %d old states and %d native states across all storage widths", count, len(p.toOld))
		})
	}
}

func TestRuntimeIDs630(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	cases := []struct {
		name  string
		props map[string]any
		id    uint32
	}{
		{"minecraft:air", nil, 11881},
		{"minecraft:stone", nil, 1774},
		{"minecraft:grass_block", nil, 12183},
		{"minecraft:dirt", nil, 9550},
		{"minecraft:coarse_dirt", nil, 9551},
		{"minecraft:info_update", nil, 668},
		{"minecraft:oak_stairs", map[string]any{"upside_down_bit": uint8(0), "weirdo_direction": int32(0)}, 717},
	}
	for _, c := range cases {
		rid, ok := world.DefaultBlockRegistry.StateToRuntimeID(c.name, c.props)
		if !ok {
			t.Fatalf("missing state %s", c.name)
		}
		if got := HashFor(630, rid); got != c.id {
			t.Errorf("%s: got %d want runtime ID %d", c.name, got, c.id)
		}
		if got := RuntimeIDFor(630, c.id); got != rid {
			t.Errorf("%s reverse: got %d want %d", c.name, got, rid)
		}
	}
	for rid := uint32(0); rid < uint32(world.DefaultBlockRegistry.BlockCount()); rid++ {
		if got := HashFor(630, rid); got >= 14127 {
			t.Fatalf("state %d leaked hash %d", rid, got)
		}
	}
	if got := HashFor(630, ^uint32(0)); got != 11881 {
		t.Fatalf("invalid state mapped to %d", got)
	}
	if got := RuntimeIDFor(630, 14127); got != world.DefaultBlockRegistry.AirRuntimeID() {
		t.Fatalf("invalid old runtime ID mapped to %d", got)
	}
}

func TestLegacyRuntimePaletteBounds(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	names := []string{"minecraft:stone", "minecraft:grass_block", "minecraft:air", "minecraft:dirt", "minecraft:info_update"}
	for _, c := range []struct {
		proto uint32
		count uint32
		ids   []uint32
	}{
		{589, 13141, []uint32{1743, 11372, 11076, 8918, 656}},
		{594, 13141, []uint32{1747, 11364, 11064, 8899, 657}},
		{618, 13129, []uint32{1757, 11347, 11045, 8903, 663}},
		{622, 13121, []uint32{1757, 11343, 11041, 8899, 663}},
	} {
		for i, name := range names {
			native, ok := world.DefaultBlockRegistry.StateToRuntimeID(name, nil)
			if !ok {
				t.Fatalf("missing native block %s", name)
			}
			old := HashFor(c.proto, native)
			if old != c.ids[i] || RuntimeIDFor(c.proto, old) != native {
				t.Errorf("protocol %d %s: got %d want %d", c.proto, name, old, c.ids[i])
			}
		}
		for native := uint32(0); native < uint32(world.DefaultBlockRegistry.BlockCount()); native++ {
			if old := HashFor(c.proto, native); old >= c.count {
				t.Fatalf("protocol %d native %d leaked hash %d", c.proto, native, old)
			}
		}
	}
}
