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

var QuickCharge quickCharge

type quickCharge struct{}

func (quickCharge) Name() string {
	return "Quick Charge"
}

func (quickCharge) MaxLevel() int {
	return 3
}

func (quickCharge) Cost(level int) (int, int) {
	minCost := 12 + (level-1)*20
	return minCost, 50
}

func (quickCharge) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityUncommon
}

func (quickCharge) ChargeDuration(level int) time.Duration {
	return time.Duration((1.25 - 0.25*float64(level)) * float64(time.Second))
}

func (quickCharge) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (quickCharge) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Crossbow)
	return ok
}
