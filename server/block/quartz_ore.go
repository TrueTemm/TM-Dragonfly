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

type NetherQuartzOre struct {
	solid
	bassDrum
}

func (q NetherQuartzOre) BreakInfo() BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, oreDrops(item.NetherQuartz{}, q)).withXPDropRange(0, 3)
}

func (NetherQuartzOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(item.NetherQuartz{}, 1), 0.2)
}

func (NetherQuartzOre) EncodeItem() (name string, meta int16) {
	return "minecraft:quartz_ore", 0
}

func (NetherQuartzOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:quartz_ore", nil
}
