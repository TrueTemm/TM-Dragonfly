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

var Knockback knockback

type knockback struct{}

func (knockback) Name() string {
	return "Knockback"
}

func (knockback) MaxLevel() int {
	return 2
}

func (knockback) Cost(level int) (int, int) {
	minCost := 5 + (level-1)*20
	return minCost, minCost + 50
}

func (knockback) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityUncommon
}

func (knockback) Force(level int) float64 {
	return float64(level) / 2
}

func (knockback) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (knockback) CompatibleWithItem(i world.Item) bool {
	t, ok := i.(item.Tool)
	return ok && t.ToolType() == item.TypeSword
}
