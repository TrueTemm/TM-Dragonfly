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

var Fortune fortune

type fortune struct{}

func (fortune) Name() string {
	return "Fortune"
}

func (fortune) MaxLevel() int {
	return 3
}

func (fortune) Cost(level int) (int, int) {
	minCost := 15 + (level-1)*9
	return minCost, minCost + 50 + level
}

func (fortune) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (fortune) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != SilkTouch
}

func (fortune) CompatibleWithItem(i world.Item) bool {
	t, ok := i.(item.Tool)
	return ok && (t.ToolType() == item.TypePickaxe || t.ToolType() == item.TypeShovel || t.ToolType() == item.TypeAxe || t.ToolType() == item.TypeHoe)
}
