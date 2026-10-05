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

type Melon struct {
	solid
}

func (m Melon) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, axeEffective, discreteDrops(item.MelonSlice{}, m, 3, 7, 9))
}

func (Melon) CompostChance() float64 {
	return 0.65
}

func (Melon) EncodeItem() (name string, meta int16) {
	return "minecraft:melon_block", 0
}

func (Melon) EncodeBlock() (string, map[string]any) {
	return "minecraft:melon_block", nil
}
