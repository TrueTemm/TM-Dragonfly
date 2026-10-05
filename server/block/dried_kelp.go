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

type DriedKelp struct {
	solid
}

func (d DriedKelp) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, hoeEffective, oneOf(d)).withBlastResistance(2.5)
}

func (DriedKelp) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 5, false)
}

func (DriedKelp) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 200)
}

func (DriedKelp) CompostChance() float64 {
	return 0.5
}

func (DriedKelp) EncodeItem() (name string, meta int16) {
	return "minecraft:dried_kelp_block", 0
}

func (DriedKelp) EncodeBlock() (string, map[string]any) {
	return "minecraft:dried_kelp_block", nil
}
