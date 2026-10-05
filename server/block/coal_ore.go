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

import "github.com/df-mc/dragonfly/server/item"

type CoalOre struct {
	solid
	bassDrum

	Type OreType
}

func (c CoalOre) BreakInfo() BreakInfo {
	return newBreakInfo(c.Type.Hardness(), pickaxeHarvestable, pickaxeEffective, oreDrops(item.Coal{}, c)).withXPDropRange(0, 2).withBlastResistance(3)
}

func (CoalOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(item.Coal{}, 1), 0.1)
}

func (c CoalOre) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.Type.Prefix() + "coal_ore", 0
}

func (c CoalOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + c.Type.Prefix() + "coal_ore", nil

}
