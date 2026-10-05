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

type ShortGrass struct {
	replaceable
	transparent
	empty

	Double bool
}

func (g ShortGrass) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 100, false)
}

func (g ShortGrass) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, grassDrops(g))
}

func (g ShortGrass) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	upper := DoubleTallGrass{Type: NormalDoubleTallGrass(), UpperPart: true}
	if replaceableWith(tx, pos.Side(cube.FaceUp), upper) {
		tx.SetBlock(pos, DoubleTallGrass{Type: NormalDoubleTallGrass()}, nil)
		tx.SetBlock(pos.Side(cube.FaceUp), upper, nil)
		return item.BoneMealResultSmall
	}
	return item.BoneMealResultNone
}

func (g ShortGrass) CompostChance() float64 {
	return 0.65
}

func (g ShortGrass) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !supportsVegetation(g, tx.Block(pos.Side(cube.FaceDown))) {
		breakBlock(g, pos, tx)
	}
}

func (g ShortGrass) HasLiquidDrops() bool {
	return true
}

func (g ShortGrass) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, g)
	if !used || !supportsVegetation(g, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(tx, pos, g, user, ctx)
	return placed(ctx)
}

func (g ShortGrass) EncodeItem() (name string, meta int16) {
	return "minecraft:short_grass", 0
}

func (g ShortGrass) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:short_grass", nil
}
