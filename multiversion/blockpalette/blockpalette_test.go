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

package blockpalette_test

import (
	"bytes"
	"testing"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestTranslateCommonBlocks(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	reg := world.DefaultBlockRegistry
	air := reg.AirRuntimeID()
	airMapped := blockpalette.Translate(786, air)

	names := map[string]uint32{}
	for rid := uint32(0); rid < uint32(reg.BlockCount()); rid++ {
		n, _, _ := reg.RuntimeIDToState(rid)
		if _, ok := names[n]; !ok {
			names[n] = rid
		}
	}
	for _, n := range []string{"minecraft:grass_block", "minecraft:dirt", "minecraft:bedrock", "minecraft:stone"} {
		rid, ok := names[n]
		if !ok {
			t.Fatalf("%s not registered in Dragonfly", n)
		}
		if got := blockpalette.Translate(786, rid); got == airMapped {
			t.Fatalf("%s translated to air (no 1.21.70 match) — expected a real block", n)
		}
	}
}

func TestRemapChunkPreservesStructure(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	reg := world.DefaultBlockRegistry
	rids := []uint32{reg.AirRuntimeID(), reg.AirRuntimeID() + 1, reg.AirRuntimeID() + 2}

	payload := new(bytes.Buffer)
	payload.WriteByte(8)
	payload.WriteByte(1)
	payload.WriteByte(byte(1<<1) | 1)
	payload.Write(make([]byte, 128*4))
	_ = protocol.WriteVarint32(payload, int32(len(rids)))
	for _, rid := range rids {
		_ = protocol.WriteVarint32(payload, int32(rid))
	}
	trailer := []byte{0xAB, 0xCD, 0xEF, 0x00}
	payload.Write(trailer)

	out := blockpalette.RemapChunkBlockPalette(payload.Bytes(), 1, 786)

	in := bytes.NewBuffer(out)
	if v, _ := in.ReadByte(); v != 8 {
		t.Fatalf("version changed: %d", v)
	}
	if v, _ := in.ReadByte(); v != 1 {
		t.Fatalf("storage count changed: %d", v)
	}
	if v, _ := in.ReadByte(); v != byte(1<<1)|1 {
		t.Fatalf("header changed: %d", v)
	}
	if idx := in.Next(128 * 4); !bytes.Equal(idx, make([]byte, 128*4)) {
		t.Fatalf("index words corrupted")
	}
	var count int32
	if err := protocol.Varint32(in, &count); err != nil || int(count) != len(rids) {
		t.Fatalf("palette count changed: %d (%v)", count, err)
	}
	for i := range rids {
		var v int32
		if err := protocol.Varint32(in, &v); err != nil {
			t.Fatalf("palette[%d]: %v", i, err)
		}

		want := blockpalette.HashOf(rids[i])
		if uint32(v) != want {
			t.Fatalf("palette[%d]: got %d want hash %d", i, uint32(v), want)
		}
	}
	if rest := in.Bytes(); !bytes.Equal(rest, trailer) {
		t.Fatalf("trailer corrupted: % x", rest)
	}
}
