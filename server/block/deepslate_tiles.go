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

type DeepslateTiles struct {
	solid
	bassDrum

	Cracked bool
}

func (d DeepslateTiles) BreakInfo() BreakInfo {
	return newBreakInfo(3.5, pickaxeHarvestable, pickaxeEffective, oneOf(d)).withBlastResistance(6)
}

func (d DeepslateTiles) SmeltInfo() item.SmeltInfo {
	if d.Cracked {
		return item.SmeltInfo{}
	}
	return newSmeltInfo(item.NewStack(DeepslateTiles{Cracked: true}, 1), 0.1)
}

func (d DeepslateTiles) EncodeItem() (name string, meta int16) {
	if d.Cracked {
		return "minecraft:cracked_deepslate_tiles", 0
	}
	return "minecraft:deepslate_tiles", 0
}

func (d DeepslateTiles) EncodeBlock() (string, map[string]any) {
	if d.Cracked {
		return "minecraft:cracked_deepslate_tiles", nil
	}
	return "minecraft:deepslate_tiles", nil
}
