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

var FireProtection fireProtection

type fireProtection struct{}

func (fireProtection) Name() string {
	return "Fire Protection"
}

func (fireProtection) MaxLevel() int {
	return 4
}

func (fireProtection) Cost(level int) (int, int) {
	minCost := 10 + (level-1)*8
	return minCost, minCost + 8
}

func (fireProtection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityUncommon
}

func (fireProtection) Modifier() float64 {
	return 0.08
}

func (fireProtection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != BlastProtection && t != ProjectileProtection && t != Protection
}

func (fireProtection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
