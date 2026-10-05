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

var Flame flame

type flame struct{}

func (flame) Name() string {
	return "Flame"
}

func (flame) MaxLevel() int {
	return 1
}

func (flame) Cost(int) (int, int) {
	return 20, 50
}

func (flame) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (flame) BurnDuration() time.Duration {
	return time.Second * 5
}

func (flame) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (flame) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Bow)
	return ok
}
