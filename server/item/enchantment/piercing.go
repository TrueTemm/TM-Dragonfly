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

var Piercing piercing

type piercing struct{}

func (p piercing) Name() string {
	return "Piercing"
}

func (p piercing) MaxLevel() int {
	return 4
}

func (p piercing) Cost(level int) (int, int) {
	return 1 + (level-1)*10, 50
}

func (p piercing) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

func (p piercing) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != Multishot
}

func (p piercing) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Crossbow)
	return ok
}

func (p piercing) Pierces() bool {
	return true
}
