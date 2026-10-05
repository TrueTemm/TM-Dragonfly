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

package block

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Gravel struct {
	gravityAffected
	solid
	snare
}

func (g Gravel) SoilFor(block world.Block) bool {
	switch block.(type) {
	case BambooSapling, Bamboo:
		return true
	}
	return false
}

func (g Gravel) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	g.fall(g, pos, tx)
}

func (g Gravel) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, shovelEffective, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(g, 1)}
		}
		flintChances := []float64{0.1, 1.0 / 7.0, 0.25, 1.0}
		flintChance := flintChances[min(fortuneLevel(enchantments), 3)]
		if rand.Float64() < flintChance {
			return []item.Stack{item.NewStack(item.Flint{}, 1)}
		}
		return []item.Stack{item.NewStack(g, 1)}
	})
}

func (Gravel) EncodeItem() (name string, meta int16) {
	return "minecraft:gravel", 0
}

func (Gravel) EncodeBlock() (string, map[string]any) {
	return "minecraft:gravel", nil
}
