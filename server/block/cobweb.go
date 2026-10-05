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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Cobweb struct {
	empty
	transparent
}

func (Cobweb) Cobweb() {}

func (Cobweb) EntityInside(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if fallEntity, ok := e.(fallDistanceEntity); ok {
		fallEntity.ResetFallDistance()
	}
	if v, ok := e.(velocityEntity); ok {
		vel := v.Velocity()
		vel[0] *= 0.25
		vel[1] *= 0.05
		vel[2] *= 0.25
		v.SetVelocity(vel)
	}
}

func (c Cobweb) BreakInfo() BreakInfo {
	return newBreakInfo(
		4,
		alwaysHarvestable,
		func(t item.Tool) bool {
			return swordEffective(t) || shearsEffective(t)
		},
		func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
			if t.ToolType() == item.TypeShears {
				return oneOf(c)(t, enchantments)
			}
			if t.ToolType() == item.TypeSword {
				return oneOf(String{})(t, enchantments)
			}
			return nil
		},
	)
}

func (Cobweb) HasLiquidDrops() bool {
	return true
}

func (Cobweb) EncodeItem() (name string, meta int16) {
	return "minecraft:web", 0
}

func (Cobweb) EncodeBlock() (string, map[string]any) {
	return "minecraft:web", nil
}
