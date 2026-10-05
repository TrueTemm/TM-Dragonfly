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

package block

import (
	"github.com/df-mc/dragonfly/server/item"
)

type NetherGoldOre struct {
	solid
}

func (n NetherGoldOre) BreakInfo() BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, multiOreDrops(item.GoldNugget{}, n, 2, 6)).withXPDropRange(0, 1)
}

func (NetherGoldOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(item.GoldIngot{}, 1), 1)
}

func (NetherGoldOre) EncodeItem() (name string, meta int16) {
	return "minecraft:nether_gold_ore", 0
}

func (NetherGoldOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:nether_gold_ore", nil
}
