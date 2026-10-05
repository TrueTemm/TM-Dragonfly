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
	"time"
)

var FireAspect fireAspect

type fireAspect struct{}

func (fireAspect) Name() string {
	return "Fire Aspect"
}

func (fireAspect) MaxLevel() int {
	return 2
}

func (fireAspect) Cost(level int) (int, int) {
	minCost := 10 + (level-1)*20
	return minCost, minCost + 50
}

func (fireAspect) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (fireAspect) Duration(level int) time.Duration {
	return time.Second * 4 * time.Duration(level)
}

func (fireAspect) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (fireAspect) CompatibleWithItem(i world.Item) bool {
	t, ok := i.(item.Tool)
	return ok && t.ToolType() == item.TypeSword
}
