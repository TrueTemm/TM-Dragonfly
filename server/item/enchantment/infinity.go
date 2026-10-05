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

var Infinity infinity

type infinity struct{}

func (infinity) Name() string {
	return "Infinity"
}

func (infinity) MaxLevel() int {
	return 1
}

func (infinity) Cost(int) (int, int) {
	return 20, 50
}

func (infinity) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

func (infinity) ConsumesArrows() bool {
	return false
}

func (infinity) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != Mending
}

func (infinity) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Bow)
	return ok
}
