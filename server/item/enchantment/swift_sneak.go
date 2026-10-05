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

var SwiftSneak swiftSneak

type swiftSneak struct{}

func (swiftSneak) Name() string {
	return "Swift Sneak"
}

func (swiftSneak) MaxLevel() int {
	return 3
}

func (swiftSneak) Cost(level int) (int, int) {
	minCost := level * 25
	return minCost, minCost + 50
}

func (swiftSneak) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

func (swiftSneak) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (swiftSneak) Treasure() bool {
	return true
}

func (swiftSneak) CompatibleWithItem(i world.Item) bool {
	b, ok := i.(item.LeggingsType)
	return ok && b.Leggings()
}
