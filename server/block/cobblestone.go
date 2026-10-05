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

type Cobblestone struct {
	solid
	bassDrum

	Mossy bool
}

func (c Cobblestone) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(c)).withBlastResistance(6)
}

func (Cobblestone) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(Stone{}, 1), 0.1)
}

func (c Cobblestone) RepairsStoneTools() bool {
	return !c.Mossy
}

func (c Cobblestone) EncodeItem() (name string, meta int16) {
	if c.Mossy {
		return "minecraft:mossy_cobblestone", 0
	}
	return "minecraft:cobblestone", 0
}

func (c Cobblestone) EncodeBlock() (string, map[string]any) {
	if c.Mossy {
		return "minecraft:mossy_cobblestone", nil
	}
	return "minecraft:cobblestone", nil
}
