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

var Multishot multishot

type multishot struct{}

func (multishot) Name() string {
	return "Multishot"
}

func (multishot) MaxLevel() int {
	return 1
}

func (m multishot) Cost(level int) (int, int) {
	return 20, 50
}

func (multishot) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (multishot) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != Piercing
}

func (multishot) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Crossbow)
	return ok
}

func (multishot) MultipleProjectiles() bool {
	return true
}
