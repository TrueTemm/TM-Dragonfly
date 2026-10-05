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

var Thorns thorns

type thorns struct{}

func (thorns) Name() string {
	return "Thorns"
}

func (thorns) MaxLevel() int {
	return 3
}

func (thorns) Cost(level int) (int, int) {
	minCost := 10 + 20*(level-1)
	return minCost, minCost + 50
}

func (thorns) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

func (thorns) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (thorns) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}

type ThornsDamageSource struct {
	Owner world.Entity
}

func (ThornsDamageSource) ReducedByResistance() bool { return true }
func (ThornsDamageSource) ReducedByArmour() bool     { return false }
func (ThornsDamageSource) Fire() bool                { return false }
func (ThornsDamageSource) IgnoreTotem() bool         { return false }
