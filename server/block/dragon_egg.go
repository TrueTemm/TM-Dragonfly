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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type DragonEgg struct {
	solid
	transparent
	gravityAffected
	sourceWaterDisplacer
}

func (d DragonEgg) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	d.fall(d, pos, tx)
}

func (d DragonEgg) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (d DragonEgg) LightEmissionLevel() uint8 {
	return 1
}

func (d DragonEgg) BreakInfo() BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, oneOf(d)).withBlastResistance(9)
}

func (DragonEgg) EncodeItem() (name string, meta int16) {
	return "minecraft:dragon_egg", 0
}

func (DragonEgg) EncodeBlock() (string, map[string]any) {
	return "minecraft:dragon_egg", nil
}
