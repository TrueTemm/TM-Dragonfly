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

package v786

import (
	"testing"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestConvertFromLatestUpdateBlockTranslates(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	reg := world.DefaultBlockRegistry

	var grass uint32
	for rid := uint32(0); rid < uint32(reg.BlockCount()); rid++ {
		if n, _, _ := reg.RuntimeIDToState(rid); n == "minecraft:grass_block" {
			grass = rid
			break
		}
	}
	want := blockpalette.HashOf(grass)
	if want == 0 {
		t.Fatalf("grass produced no block hash")
	}

	out := Protocol{}.ConvertFromLatest(&packet.UpdateBlock{NewBlockRuntimeID: grass, Position: [3]int32{1, 100, 2}}, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 packet, got %d", len(out))
	}
	ub, ok := out[0].(*UpdateBlock)
	if !ok {
		t.Fatalf("expected *v786.UpdateBlock, got %T (conversion not wired)", out[0])
	}
	if ub.NewBlockRuntimeID != want {
		t.Fatalf("grass rid not translated: got %d want %d (raw df rid was %d)", ub.NewBlockRuntimeID, want, grass)
	}
	if ub.NewBlockRuntimeID == grass {
		t.Fatalf("NewBlockRuntimeID still the raw Dragonfly rid %d — client would render the wrong block", grass)
	}
}
