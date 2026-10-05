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

type EmeraldOre struct {
	solid
	bassDrum

	Type OreType
}

func (e EmeraldOre) BreakInfo() BreakInfo {
	i := newBreakInfo(e.Type.Hardness(), func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, oreDrops(item.Emerald{}, e)).withXPDropRange(3, 7)
	if e.Type == DeepslateOre() {
		i = i.withBlastResistance(3)
	}
	return i
}

func (EmeraldOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(item.Emerald{}, 1), 1)
}

func (e EmeraldOre) EncodeItem() (name string, meta int16) {
	return "minecraft:" + e.Type.Prefix() + "emerald_ore", 0
}

func (e EmeraldOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + e.Type.Prefix() + "emerald_ore", nil
}
