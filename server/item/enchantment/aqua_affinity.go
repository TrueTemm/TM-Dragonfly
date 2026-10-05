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

var AquaAffinity aquaAffinity

type aquaAffinity struct{}

func (aquaAffinity) Name() string {
	return "Aqua Affinity"
}

func (aquaAffinity) MaxLevel() int {
	return 1
}

func (aquaAffinity) Cost(int) (int, int) {
	return 1, 41
}

func (aquaAffinity) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

func (aquaAffinity) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (aquaAffinity) CompatibleWithItem(i world.Item) bool {
	h, ok := i.(item.HelmetType)
	return ok && h.Helmet()
}
