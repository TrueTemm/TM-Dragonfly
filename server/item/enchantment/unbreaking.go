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
	"math/rand/v2"
)

var Unbreaking unbreaking

type unbreaking struct{}

func (unbreaking) Name() string {
	return "Unbreaking"
}

func (unbreaking) MaxLevel() int {
	return 3
}

func (unbreaking) Cost(level int) (int, int) {
	minCost := 5 + 8*(level-1)
	return minCost, minCost + 50
}

func (unbreaking) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityUncommon
}

func (unbreaking) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (unbreaking) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Durable)
	return ok
}

func (unbreaking) Reduce(it world.Item, level, amount int) int {
	after := amount
	_, ok := it.(item.Armour)
	for i := 0; i < amount; i++ {
		if (!ok || rand.Float64() >= 0.6) && rand.IntN(level+1) > 0 {
			after--
		}
	}
	return after
}
