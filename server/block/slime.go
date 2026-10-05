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

type Slime struct {
	solid
	transparent
}

func (Slime) EntityLand(_ cube.Pos, _ *world.Tx, e world.Entity, distance *float64) {
	if _, ok := e.(fallDistanceEntity); ok {
		*distance = 0
	}
	if s, ok := e.(interface{ Sneaking() bool }); ok && s.Sneaking() {
		return
	}
	if v, ok := e.(velocityEntity); ok {
		vel := v.Velocity()
		if vel[1] < 0 {
			vel[1] = -vel[1]
			v.SetVelocity(vel)
		}
	}
}

func (Slime) Friction() float64 {
	return 0.8
}

func (s Slime) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(s))
}

func (Slime) EncodeItem() (name string, meta int16) {
	return "minecraft:slime", 0
}

func (Slime) EncodeBlock() (string, map[string]any) {
	return "minecraft:slime", nil
}
