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
	"github.com/df-mc/dragonfly/server/world/sound"
)

type PackedIce struct {
	solid
}

func (PackedIce) Instrument() sound.Instrument {
	return sound.Chimes()
}

func (p PackedIce) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, pickaxeEffective, silkTouchOnlyDrop(p))
}

func (p PackedIce) Friction() float64 {
	return 0.98
}

func (PackedIce) EncodeItem() (name string, meta int16) {
	return "minecraft:packed_ice", 0
}

func (PackedIce) EncodeBlock() (string, map[string]any) {
	return "minecraft:packed_ice", nil
}
