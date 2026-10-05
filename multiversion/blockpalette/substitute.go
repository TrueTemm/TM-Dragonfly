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
	"regexp"
	"strings"
)

var substitutes = map[string]string{

	"minecraft:leaf_litter":       "minecraft:air",
	"minecraft:wildflowers":       "minecraft:air",
	"minecraft:cactus_flower":     "minecraft:air",
	"minecraft:pale_hanging_moss": "minecraft:air",
	"minecraft:resin_clump":       "minecraft:air",
	"minecraft:sulfur_spike":      "minecraft:air",
	"minecraft:dried_ghast":       "minecraft:air",

	"minecraft:bush":              "minecraft:short_grass",
	"minecraft:firefly_bush":      "minecraft:short_grass",
	"minecraft:short_dry_grass":   "minecraft:short_grass",
	"minecraft:tall_dry_grass":    "minecraft:tall_grass",
	"minecraft:open_eyeblossom":   "minecraft:oxeye_daisy",
	"minecraft:closed_eyeblossom": "minecraft:wither_rose",
	"minecraft:golden_dandelion":  "minecraft:dandelion",
	"minecraft:pale_moss_block":   "minecraft:moss_block",
	"minecraft:pale_moss_carpet":  "minecraft:moss_carpet",
	"minecraft:creaking_heart":    "minecraft:birch_log",
	"minecraft:mushroom_stem":     "minecraft:brown_mushroom_block",

	"minecraft:brown_mushroom_block": "minecraft:brown_mushroom_block",
	"minecraft:red_mushroom_block":   "minecraft:red_mushroom_block",

	"minecraft:straw_bed": "minecraft:bed",

	"minecraft:creeper_head":          "minecraft:skull",
	"minecraft:dragon_head":           "minecraft:skull",
	"minecraft:piglin_head":           "minecraft:skull",
	"minecraft:player_head":           "minecraft:skull",
	"minecraft:zombie_head":           "minecraft:skull",
	"minecraft:skeleton_skull":        "minecraft:skull",
	"minecraft:wither_skeleton_skull": "minecraft:skull",

	"minecraft:sulfur":                  "minecraft:yellow_terracotta",
	"minecraft:polished_sulfur":         "minecraft:yellow_terracotta",
	"minecraft:chiseled_sulfur":         "minecraft:yellow_terracotta",
	"minecraft:sulfur_bricks":           "minecraft:yellow_terracotta",
	"minecraft:potent_sulfur":           "minecraft:yellow_terracotta",
	"minecraft:cinnabar":                "minecraft:red_terracotta",
	"minecraft:polished_cinnabar":       "minecraft:red_terracotta",
	"minecraft:chiseled_cinnabar":       "minecraft:red_terracotta",
	"minecraft:cinnabar_bricks":         "minecraft:red_terracotta",
	"minecraft:resin_block":             "minecraft:orange_terracotta",
	"minecraft:resin_bricks":            "minecraft:orange_terracotta",
	"minecraft:chiseled_resin_bricks":   "minecraft:orange_terracotta",
	"minecraft:resin_brick_stairs":      "minecraft:smooth_red_sandstone_stairs",
	"minecraft:resin_brick_slab":        "minecraft:red_sandstone_slab",
	"minecraft:resin_brick_double_slab": "minecraft:red_sandstone_double_slab",
	"minecraft:resin_brick_wall":        "minecraft:red_sandstone_wall",
}

var (
	copperFamily = regexp.MustCompile(`^minecraft:(?:waxed_)?(?:exposed_|weathered_|oxidized_)?(copper_bars|copper_chain|copper_chest|copper_lantern|copper_torch|copper_golem_statue|lightning_rod)$`)
	copperTo     = map[string]string{
		"copper_bars": "minecraft:iron_bars", "copper_chain": "minecraft:iron_chain", "copper_chest": "minecraft:chest",
		"copper_lantern": "minecraft:lantern", "copper_torch": "minecraft:torch", "copper_golem_statue": "minecraft:air",
		"lightning_rod": "minecraft:lightning_rod",
	}

	stoneShape = regexp.MustCompile(`^minecraft:(?:polished_)?(sulfur|cinnabar)(?:_brick)?_(stairs|slab|double_slab|wall)$`)
	stoneTo    = map[string]string{"sulfur": "minecraft:sandstone_", "cinnabar": "minecraft:red_sandstone_"}

	poplarLeaves = regexp.MustCompile(`^minecraft:(?:red|yellow|orange)_poplar_leaves$`)
)

func substituteName(name string) string {
	if to, ok := substitutes[name]; ok {
		return to
	}
	if m := copperFamily.FindStringSubmatch(name); m != nil {
		return copperTo[m[1]]
	}
	if m := stoneShape.FindStringSubmatch(name); m != nil {
		return stoneTo[m[1]] + m[2]
	}
	if poplarLeaves.MatchString(name) {
		return "minecraft:birch_leaves"
	}

	mapped := false
	for _, fam := range []string{"pale_oak", "poplar"} {
		if strings.HasPrefix(name, "minecraft:"+fam+"_") {
			name, mapped = "minecraft:birch_"+strings.TrimPrefix(name, "minecraft:"+fam+"_"), true
			break
		}
		if strings.HasPrefix(name, "minecraft:stripped_"+fam+"_") {
			name, mapped = "minecraft:stripped_birch_"+strings.TrimPrefix(name, "minecraft:stripped_"+fam+"_"), true
			break
		}
	}

	if strings.HasSuffix(name, "_shelf") {
		return strings.TrimSuffix(name, "_shelf") + "_planks"
	}
	if mapped {
		return name
	}
	return ""
}

type oldPalette struct {
	byKey map[string]uint32
	first map[string]uint32
	props map[string]map[string]bool
}

func (p *statePalette) substituteHash(name string, props map[string]any, mapped []bool, old *oldPalette) uint32 {
	to := substituteName(name)
	if to == "" {
		return 0
	}
	try := func(pr map[string]any) (uint32, bool) {
		key := untypedKey(to, pr)
		if rid, ok := latestKeys[key]; ok && mapped[rid] {
			return p.toOld[rid], true
		}
		if h, ok := old.byKey[key]; ok {
			return h, true
		}
		return 0, false
	}
	if h, ok := try(props); ok {
		return h
	}

	keep := make(map[string]any, len(props))
	for k, v := range props {
		if latestPropKeys[to][k] || old.props[to][k] {
			keep[k] = v
		}
	}
	if len(keep) != len(props) {
		if h, ok := try(keep); ok {
			return h
		}
	}
	if rid, ok := latestFirst[to]; ok && mapped[rid] {
		return p.toOld[rid]
	}
	if h, ok := old.first[to]; ok {
		return h
	}
	return 0
}
