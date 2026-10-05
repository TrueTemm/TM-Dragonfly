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

var ProjectileProtection projectileProtection

type projectileProtection struct{}

func (projectileProtection) Name() string {
	return "Projectile Protection"
}

func (projectileProtection) MaxLevel() int {
	return 4
}

func (projectileProtection) Cost(level int) (int, int) {
	minCost := 3 + (level-1)*6
	return minCost, minCost + 6
}

func (projectileProtection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityUncommon
}

func (projectileProtection) Modifier() float64 {
	return 0.08
}

func (projectileProtection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != BlastProtection && t != FireProtection && t != Protection
}

func (projectileProtection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
