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
	"sort"
	"testing"

	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

func TestDowngradeCoverage(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	protos := make([]int, 0, len(statePalettes))
	for p := range statePalettes {
		protos = append(protos, int(p))
	}
	sort.Ints(protos)
	for _, p := range protos {
		unmatchedOld, missing, _ := DowngradeStats(uint32(p))
		total := 0
		names := make([]string, 0, len(missing))
		for n, c := range missing {
			total += c
			names = append(names, n)
		}
		sort.Strings(names)
		t.Logf("%d: %d old states without a 1.26.45 equivalent; %d Dragonfly states (%d blocks) the client cannot spell", p, unmatchedOld, total, len(names))
		if unmatchedOld != 0 {
			t.Errorf("%d: %d states of the client palette upgrade to nothing Dragonfly knows — a schema is missing or misapplied", p, unmatchedOld)
		}

		for _, ancient := range []string{"minecraft:dirt", "minecraft:coarse_dirt", "minecraft:grass_block", "minecraft:stone", "minecraft:bedrock",
			"minecraft:oak_planks", "minecraft:oak_stairs", "minecraft:oak_leaves", "minecraft:sand", "minecraft:chest", "minecraft:oak_door",
			"minecraft:cobblestone_wall", "minecraft:tnt", "minecraft:sponge", "minecraft:anvil", "minecraft:purpur_block", "minecraft:mushroom_stem"} {
			if missing[ancient] > 0 && !(ancient == "minecraft:mushroom_stem" && p < 748) {
				t.Errorf("%d cannot spell %d states of %s", p, missing[ancient], ancient)
			}
		}
	}
}

func TestDowngradeDirt(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	reg := world.DefaultBlockRegistry
	var dirt, coarse uint32
	for rid := uint32(0); rid < uint32(reg.BlockCount()); rid++ {
		if name, _, ok := reg.RuntimeIDToState(rid); ok {
			switch name {
			case "minecraft:dirt":
				dirt = rid
			case "minecraft:coarse_dirt":
				coarse = rid
			}
		}
	}
	old := world.NetworkBlockHash("minecraft:dirt", map[string]any{"dirt_type": "normal"})
	oldCoarse := world.NetworkBlockHash("minecraft:dirt", map[string]any{"dirt_type": "coarse"})
	if got := HashFor(686, dirt); got != old {
		t.Fatalf("686 dirt hash = %d, want the dirt_type=normal spelling %d (latest %d)", got, old, HashOf(dirt))
	}
	if got := HashFor(685, coarse); got != oldCoarse {
		t.Fatalf("685 coarse dirt hash = %d, want %d", got, oldCoarse)
	}
	if HashFor(712, dirt) != HashOf(dirt) || HashFor(786, dirt) != HashOf(dirt) {
		t.Fatalf("dirt is flat from 1.21.20 on; 712/786 must hash it as latest does")
	}
	if got := RuntimeIDFor(686, old); got != dirt {
		t.Fatalf("686 reverse: %d -> rid %d, want dirt %d", old, got, dirt)
	}
	if got := RuntimeIDFor(686, HashOf(dirt)); got != dirt {
		t.Fatalf("686 reverse of the latest hash must still resolve (registry fallback): got %d", got)
	}
}

func TestDowngradePlaceholder(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	reg := world.DefaultBlockRegistry
	var pale, info uint32
	for rid := uint32(0); rid < uint32(reg.BlockCount()); rid++ {
		if name, _, ok := reg.RuntimeIDToState(rid); ok {
			switch name {
			case "minecraft:chalkboard":
				pale = rid
			case "minecraft:info_update":
				info = rid
			}
		}
	}
	placeholder := world.NetworkBlockHash("minecraft:info_update", nil)
	if got := HashFor(686, pale); got != placeholder {
		t.Fatalf("686 chalkboard = %d, want info_update %d", got, placeholder)
	}
	if got := RuntimeIDFor(686, placeholder); got != info {
		t.Fatalf("686 info_update maps back to rid %d, want %d", got, info)
	}
	if HashFor(786, pale) != HashOf(pale) {
		t.Fatalf("786 has the chalkboard; it must keep the real hash")
	}
	if HashFor(898, pale) != HashOf(pale) {
		t.Fatalf("898 borrows 844's palette; unknown states must keep the latest hash there")
	}
}

func TestDowngradeSubstitutes(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	reg := world.DefaultBlockRegistry
	rids := map[string]uint32{}
	for rid := uint32(0); rid < uint32(reg.BlockCount()); rid++ {
		if name, _, ok := reg.RuntimeIDToState(rid); ok {
			if _, seen := rids[name]; !seen {
				rids[name] = rid
			}
		}
	}
	p := statePalettes[686]
	p.build()
	for from, to := range map[string]string{
		"minecraft:leaf_litter": "minecraft:air", "minecraft:wildflowers": "minecraft:air",
		"minecraft:bush": "minecraft:short_grass", "minecraft:pale_oak_planks": "minecraft:birch_planks",
		"minecraft:oxidized_lightning_rod": "minecraft:lightning_rod", "minecraft:copper_lantern": "minecraft:lantern",
		"minecraft:waxed_weathered_copper_chain": "minecraft:iron_chain", "minecraft:poplar_shelf": "minecraft:birch_planks",
		"minecraft:sulfur_stairs": "minecraft:sandstone_stairs", "minecraft:pale_moss_carpet": "minecraft:moss_carpet",
	} {
		a, ok := rids[from]
		if !ok {
			t.Errorf("%s: not in the registry (rule is stale)", from)
			continue
		}
		b := rids[to]
		if to == "minecraft:air" {
			b = reg.AirRuntimeID()
		}
		if HashFor(686, a) != HashFor(686, b) {
			t.Errorf("686: %s -> %d, want %s (%d)", from, HashFor(686, a), to, HashFor(686, b))
		}
	}

	info := world.NetworkBlockHash("minecraft:info_update", nil)
	for name, rid := range rids {
		if substituteName(name) != "" && HashFor(686, rid) == info && name != "minecraft:info_update" {
			t.Errorf("686: %s has a substitute rule (%s) but still resolves to info_update", name, substituteName(name))
		}
	}
	names := make([]string, 0, len(p.unmatchedLatestNames))
	for n := range p.unmatchedLatestNames {
		names = append(names, n)
	}
	sort.Strings(names)
	t.Logf("686: %d blocks substituted, %d still info_update: %v", len(p.substituted), len(names), names)
}
