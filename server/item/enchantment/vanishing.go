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

var CurseOfVanishing curseOfVanishing

type curseOfVanishing struct{}

func (curseOfVanishing) Name() string {
	return "Curse of Vanishing"
}

func (curseOfVanishing) MaxLevel() int {
	return 1
}

func (curseOfVanishing) Cost(int) (int, int) {
	return 25, 50
}

func (curseOfVanishing) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

func (curseOfVanishing) CompatibleWithEnchantment(_ item.EnchantmentType) bool {
	return true
}

func (curseOfVanishing) CompatibleWithItem(i world.Item) bool {
	_, arm := i.(item.Armour)
	_, com := i.(item.Compass)
	_, dur := i.(item.Durable)
	_, rec := i.(item.RecoveryCompass)

	return arm || com || dur || rec
}

func (curseOfVanishing) Treasure() bool {
	return true
}

func (curseOfVanishing) Curse() bool {
	return true
}
