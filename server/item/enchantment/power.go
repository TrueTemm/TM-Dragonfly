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

var Power power

type power struct{}

func (power) Name() string {
	return "Power"
}

func (power) MaxLevel() int {
	return 5
}

func (power) Cost(level int) (int, int) {
	minCost := 1 + (level-1)*10
	return minCost, minCost + 15
}

func (power) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

func (power) PowerDamage(level int) float64 {
	return float64(level+1) * 0.5
}

func (power) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (power) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Bow)
	return ok
}
