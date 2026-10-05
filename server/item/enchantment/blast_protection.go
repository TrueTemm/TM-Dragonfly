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

var BlastProtection blastProtection

type blastProtection struct{}

func (blastProtection) Name() string {
	return "Blast Protection"
}

func (blastProtection) MaxLevel() int {
	return 4
}

func (blastProtection) Cost(level int) (int, int) {
	minCost := 5 + (level-1)*8
	return minCost, minCost + 8
}

func (blastProtection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (blastProtection) Modifier() float64 {
	return 0.08
}

func (blastProtection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != FireProtection && t != ProjectileProtection && t != Protection
}

func (blastProtection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
