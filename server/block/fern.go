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
)

type Fern struct {
	replaceable
	transparent
	empty
}

func (g Fern) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 100, false)
}

func (g Fern) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, grassDrops(g))
}

func (g Fern) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	upper := DoubleTallGrass{Type: FernDoubleTallGrass(), UpperPart: true}
	if replaceableWith(tx, pos.Side(cube.FaceUp), upper) {
		tx.SetBlock(pos, DoubleTallGrass{Type: FernDoubleTallGrass()}, nil)
		tx.SetBlock(pos.Side(cube.FaceUp), upper, nil)
		return item.BoneMealResultSmall
	}
	return item.BoneMealResultNone
}

func (g Fern) CompostChance() float64 {
	return 0.3
}

func (g Fern) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !supportsVegetation(g, tx.Block(pos.Side(cube.FaceDown))) {
		breakBlock(g, pos, tx)
	}
}

func (g Fern) HasLiquidDrops() bool {
	return true
}

func (g Fern) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, g)
	if !used || !supportsVegetation(g, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(tx, pos, g, user, ctx)
	return placed(ctx)
}

func (g Fern) EncodeItem() (name string, meta int16) {
	return "minecraft:fern", 0
}

func (g Fern) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:fern", nil
}
