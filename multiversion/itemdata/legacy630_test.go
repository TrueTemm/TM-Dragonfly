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

package itemdata

import (
	"fmt"
	"testing"

	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

func TestLegacyItemVariantClosure(t *testing.T) {
	for _, proto := range []uint32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			checked := 0
			var unavailable []string
			for _, entry := range ItemsFor(proto) {
				if id, _ := ReverseTranslateStack(proto, int32(entry.RuntimeID), 0); id == 0 && entry.RuntimeID != 0 {
					unavailable = append(unavailable, entry.Name)
				}
				for meta := uint32(0); meta < 256; meta++ {
					native, aux := ReverseTranslateStack(proto, int32(entry.RuntimeID), meta)
					if native == 0 {
						continue
					}
					old, oldAux := TranslateStack(proto, native, aux)
					got, gotAux := ReverseTranslateStack(proto, old, oldAux)
					if native != got || aux != gotAux {
						t.Errorf("%s:%d -> %d:%d -> %d:%d -> %d:%d", entry.Name, meta, native, aux, old, oldAux, got, gotAux)
					}
					checked++
				}
			}
			t.Logf("checked %d legacy item variants", checked)
			t.Logf("legacy names without a native item: %v", unavailable)
		})
	}
}

func TestLegacyRegisteredItemVariants(t *testing.T) {
	for _, proto := range []uint32{589, 594, 598, 618, 622, 630} {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			supported, unavailable := 0, 0
			for _, it := range world.Items() {
				native, meta, ok := world.ItemRuntimeID(it)
				if !ok {
					t.Fatal("registered item has no runtime ID")
				}
				old, aux := TranslateStack(proto, native, uint32(uint16(meta)))
				if old == 0 {
					unavailable++
					continue
				}
				got, gotMeta := ReverseTranslateStack(proto, old, aux)
				if got != native || gotMeta != uint32(uint16(meta)) {
					name, _ := it.EncodeItem()
					t.Errorf("%s:%d -> %d:%d -> %d:%d", name, meta, old, aux, got, gotMeta)
				}
				supported++
			}
			t.Logf("%d supported registered variants and %d filtered variants", supported, unavailable)
		})
	}
}

func TestLegacy630ItemVariants(t *testing.T) {
	for _, c := range []struct {
		name string
		id   int32
		meta uint32
	}{
		{"grass_block", 2, 0},
		{"coarse_dirt", 3, 1},
		{"red_sand", 12, 1},
		{"stone_bricks", 98, 0},
		{"chiseled_stone_bricks", 98, 3},
		{"oak_slab", 158, 0},
		{"spruce_slab", 158, 1},
		{"short_grass", 31, 0},
	} {
		native, ok := world.ItemRuntimeIDByName("minecraft:" + c.name)
		if !ok {
			t.Fatalf("missing native item %s", c.name)
		}
		id, meta := TranslateStack(630, native, 0)
		if id != c.id || meta != c.meta {
			t.Errorf("%s: got %d:%d want %d:%d", c.name, id, meta, c.id, c.meta)
		}
		got, gotMeta := ReverseTranslateStack(630, c.id, c.meta)
		if got != native || gotMeta != 0 {
			t.Errorf("%s reverse: got %d:%d want %d:0", c.name, got, gotMeta, native)
		}
	}
	shield, _ := world.ItemRuntimeIDByName("minecraft:shield")
	id, meta := TranslateStack(630, shield, 12)
	native, damage := ReverseTranslateStack(630, id, meta)
	if native != shield || damage != 12 {
		t.Fatal("item damage lost")
	}
}

func TestOlderItemVariants(t *testing.T) {
	for _, proto := range []uint32{589, 594, 598, 618, 622} {
		for _, c := range []struct {
			name string
			id   int32
			meta uint32
		}{
			{"grass_block", 2, 0}, {"coarse_dirt", 3, 1}, {"red_sand", 12, 1}, {"oak_slab", 158, 0}, {"spruce_slab", 158, 1},
		} {
			native, ok := world.ItemRuntimeIDByName("minecraft:" + c.name)
			if !ok {
				t.Fatal(c.name)
			}
			id, meta := TranslateStack(proto, native, 0)
			if id != c.id || meta != c.meta {
				t.Errorf("protocol %d %s got %d:%d want %d:%d", proto, c.name, id, meta, c.id, c.meta)
			}
			got, damage := ReverseTranslateStack(proto, c.id, c.meta)
			if got != native || damage != 0 {
				t.Errorf("protocol %d %s failed reverse mapping", proto, c.name)
			}
		}
	}
}
