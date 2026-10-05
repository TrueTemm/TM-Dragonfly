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
	"time"
)

type FletchingTable struct {
	solid
	bass
}

func (f FletchingTable) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, oneOf(f))
}

func (FletchingTable) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (FletchingTable) EncodeItem() (string, int16) {
	return "minecraft:fletching_table", 0
}

func (FletchingTable) EncodeBlock() (string, map[string]any) {
	return "minecraft:fletching_table", nil
}
