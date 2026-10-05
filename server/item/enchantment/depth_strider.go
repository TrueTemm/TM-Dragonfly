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

var DepthStrider depthStrider

type depthStrider struct{}

func (depthStrider) Name() string {
	return "Depth Strider"
}

func (depthStrider) MaxLevel() int {
	return 3
}

func (depthStrider) Cost(level int) (int, int) {
	minCost := level * 10
	return minCost, minCost + 15
}

func (depthStrider) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (depthStrider) CompatibleWithEnchantment(item.EnchantmentType) bool {

	return true
}

func (depthStrider) CompatibleWithItem(i world.Item) bool {
	b, ok := i.(item.BootsType)
	return ok && b.Boots()
}
