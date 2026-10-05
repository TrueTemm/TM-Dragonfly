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

var Protection protection

type protection struct{}

func (protection) Name() string {
	return "Protection"
}

func (protection) MaxLevel() int {
	return 4
}

func (protection) Cost(level int) (int, int) {
	minCost := 1 + (level-1)*11
	return minCost, minCost + 11
}

func (protection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

func (protection) Modifier() float64 {
	return 0.04
}

func (protection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != BlastProtection && t != FireProtection && t != ProjectileProtection
}

func (protection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
