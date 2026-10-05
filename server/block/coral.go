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
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"time"
)

type Coral struct {
	empty
	transparent
	bassDrum
	sourceWaterDisplacer

	Type CoralType

	Dead bool
}

func (c Coral) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, c)
	if !used {
		return false
	}
	if !tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, tx) {
		return false
	}
	if liquid, ok := tx.Liquid(pos); ok {
		if water, ok := liquid.(Water); ok {
			if water.Depth != 8 {
				return false
			}
		}
	}

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (c Coral) HasLiquidDrops() bool {
	return false
}

func (c Coral) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c Coral) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, tx) {
		breakBlock(c, pos, tx)
		return
	} else if c.Dead {
		return
	}
	tx.ScheduleBlockUpdate(pos, c, time.Second*5/2)
}

func (c Coral) ScheduledTick(pos cube.Pos, tx *world.Tx, _ *rand.Rand) {
	adjacentWater := false
	pos.Neighbours(func(neighbour cube.Pos) {
		if liquid, ok := tx.Liquid(neighbour); ok {
			if _, ok := liquid.(Water); ok {
				adjacentWater = true
			}
		}
	}, tx.Range())
	if !adjacentWater {
		c.Dead = true
		tx.SetBlock(pos, c, nil)
	}
}

func (c Coral) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(c))
}

func (c Coral) EncodeBlock() (name string, properties map[string]any) {
	if c.Dead {
		return "minecraft:dead_" + c.Type.String() + "_coral", nil
	}
	return "minecraft:" + c.Type.String() + "_coral", nil
}

func (c Coral) EncodeItem() (name string, meta int16) {
	if c.Dead {
		return "minecraft:dead_" + c.Type.String() + "_coral", 0
	}
	return "minecraft:" + c.Type.String() + "_coral", 0
}

func allCoral() (c []world.Block) {
	f := func(dead bool) {
		for _, t := range CoralTypes() {
			c = append(c, Coral{Type: t, Dead: dead})
		}
	}
	f(true)
	f(false)
	return
}
