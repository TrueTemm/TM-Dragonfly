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

package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

var Efficiency efficiency

type efficiency struct{}

func (efficiency) Name() string {
	return "Efficiency"
}

func (efficiency) MaxLevel() int {
	return 5
}

func (efficiency) Cost(level int) (int, int) {
	minCost := 1 + 10*(level-1)
	return minCost, minCost + 50
}

func (efficiency) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

func (efficiency) Addend(level int) float64 {
	return float64(level*level + 1)
}

func (efficiency) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (efficiency) CompatibleWithItem(i world.Item) bool {
	t, ok := i.(item.Tool)
	return ok && (t.ToolType() != item.TypeSword && t.ToolType() != item.TypeNone)
}
