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

type BambooMosaic struct {
	solid
	bass
}

func (BambooMosaic) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(5, 20, true)
}

func (b BambooMosaic) BreakInfo() BreakInfo {
	return newBreakInfo(2, alwaysHarvestable, axeEffective, oneOf(b)).withBlastResistance(3)
}

func (BambooMosaic) RepairsWoodTools() bool {
	return true
}

func (BambooMosaic) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (BambooMosaic) EncodeItem() (name string, meta int16) {
	return "minecraft:bamboo_mosaic", 0
}

func (BambooMosaic) EncodeBlock() (string, map[string]any) {
	return "minecraft:bamboo_mosaic", nil
}
