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

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Blackstone struct {
	solid
	bassDrum

	Type BlackstoneType
}

func (b Blackstone) BreakInfo() BreakInfo {
	drops := oneOf(b)
	hardness := 1.5

	switch b.Type {
	case GildedBlackstone():
		drops = func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
			if hasSilkTouch(enchantments) {
				return []item.Stack{item.NewStack(b, 1)}
			}
			nuggetChances := []float64{0.1, 1.0 / 7.0, 0.25, 1.0}
			if rand.Float64() < nuggetChances[min(fortuneLevel(enchantments), 3)] {
				return []item.Stack{item.NewStack(item.GoldNugget{}, rand.IntN(4)+2)}
			}
			return []item.Stack{item.NewStack(b, 1)}
		}
	case PolishedBlackstone():
		hardness = 2
	}

	return newBreakInfo(hardness, pickaxeHarvestable, pickaxeEffective, drops).withBlastResistance(6)
}

func (b Blackstone) EncodeItem() (name string, meta int16) {
	return "minecraft:" + b.Type.String(), 0
}

func (b Blackstone) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + b.Type.String(), nil
}

func allBlackstone() (s []world.Block) {
	for _, t := range BlackstoneTypes() {
		s = append(s, Blackstone{Type: t})
	}
	return
}
