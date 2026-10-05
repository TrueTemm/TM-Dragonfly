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

var Sharpness sharpness

type sharpness struct{}

func (sharpness) Name() string {
	return "Sharpness"
}

func (sharpness) MaxLevel() int {
	return 5
}

func (sharpness) Cost(level int) (int, int) {
	minCost := 1 + (level-1)*11
	return minCost, minCost + 20
}

func (sharpness) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

func (sharpness) Addend(level int) float64 {
	return float64(level) * 1.25
}

func (sharpness) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (sharpness) CompatibleWithItem(i world.Item) bool {
	t, ok := i.(item.Tool)
	return ok && (t.ToolType() == item.TypeSword || t.ToolType() == item.TypeAxe)
}
