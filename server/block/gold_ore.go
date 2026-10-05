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

type GoldOre struct {
	solid
	bassDrum

	Type OreType
}

func (g GoldOre) BreakInfo() BreakInfo {
	return newBreakInfo(g.Type.Hardness(), func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, oreDrops(item.RawGold{}, g)).withBlastResistance(3)
}

func (GoldOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(item.GoldIngot{}, 1), 1)
}

func (g GoldOre) EncodeItem() (name string, meta int16) {
	return "minecraft:" + g.Type.Prefix() + "gold_ore", 0
}

func (g GoldOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + g.Type.Prefix() + "gold_ore", nil
}
