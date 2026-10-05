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

var Mending mending

type mending struct{}

func (mending) Name() string {
	return "Mending"
}

func (mending) MaxLevel() int {
	return 1
}

func (mending) Cost(level int) (int, int) {
	minCost := level * 25
	return minCost, minCost + 50
}

func (mending) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (mending) Treasure() bool {
	return true
}

func (mending) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != Infinity
}

func (mending) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Durable)
	return ok
}
