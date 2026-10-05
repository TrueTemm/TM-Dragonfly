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

var Respiration respiration

type respiration struct{}

func (respiration) Name() string {
	return "Respiration"
}

func (respiration) MaxLevel() int {
	return 3
}

func (respiration) Cost(level int) (int, int) {
	minCost := 10 * level
	return minCost, minCost + 30
}

func (respiration) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (respiration) Chance(level int) float64 {
	return float64(level) / float64(level+1)
}

func (respiration) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (respiration) CompatibleWithItem(i world.Item) bool {
	h, ok := i.(item.HelmetType)
	return ok && h.Helmet()
}
