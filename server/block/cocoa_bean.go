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
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type CocoaBean struct {
	transparent

	Facing cube.Direction

	Age int
}

func (c CocoaBean) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if c.Age == 2 {
		return item.BoneMealResultNone
	}
	c.Age++
	tx.SetBlock(pos, c, nil)
	return item.BoneMealResultSmall
}

func (c CocoaBean) HasLiquidDrops() bool {
	return true
}

func (c CocoaBean) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	var woodType WoodType
	switch b := tx.Block(pos.Side(c.Facing.Face())).(type) {
	case Log:
		woodType = b.Wood
	case Wood:
		woodType = b.Wood
	}
	if woodType != JungleWood() {
		breakBlock(c, pos, tx)
	}
}

func (c CocoaBean) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, c)
	if !used {
		return false
	}

	if face == cube.FaceUp || face == cube.FaceDown {
		return false
	}

	var woodType WoodType
	oppositePos := pos.Side(face.Opposite())
	if log, ok := tx.Block(oppositePos).(Log); ok {
		woodType = log.Wood
	} else if wood, ok := tx.Block(oppositePos).(Wood); ok {
		woodType = wood.Wood
	}
	if woodType == JungleWood() {
		c.Facing = face.Opposite().Direction()
		ctx.IgnoreBBox = true

		place(tx, pos, c, user, ctx)
		return placed(ctx)
	}

	return false
}

func (c CocoaBean) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if c.Age < 2 && r.IntN(5) == 0 {
		c.Age++
		tx.SetBlock(pos, c, nil)
	}
}

func (c CocoaBean) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, axeEffective, func(item.Tool, []item.Enchantment) []item.Stack {
		if c.Age == 2 {
			return []item.Stack{item.NewStack(c, rand.IntN(2)+2)}
		}
		return []item.Stack{item.NewStack(c, 1)}
	}).withBlastResistance(3)
}

func (CocoaBean) CompostChance() float64 {
	return 0.65
}

func (c CocoaBean) EncodeItem() (name string, meta int16) {
	return "minecraft:cocoa_beans", 0
}

func (c CocoaBean) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:cocoa", map[string]any{"age": int32(c.Age), "direction": int32(horizontalDirection(c.Facing))}
}

func (c CocoaBean) Model() world.BlockModel {
	return model.CocoaBean{Facing: c.Facing, Age: c.Age}
}

func allCocoaBeans() (cocoa []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		cocoa = append(cocoa, CocoaBean{Facing: i, Age: 0})
		cocoa = append(cocoa, CocoaBean{Facing: i, Age: 1})
		cocoa = append(cocoa, CocoaBean{Facing: i, Age: 2})
	}
	return
}
