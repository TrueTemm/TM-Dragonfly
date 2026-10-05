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

type DoubleTallGrass struct {
	transparent
	replaceable
	empty

	UpperPart bool

	Type DoubleTallGrassType
}

func (d DoubleTallGrass) HasLiquidDrops() bool {
	return true
}

func (d DoubleTallGrass) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if d.UpperPart {
		if bottom, ok := tx.Block(pos.Side(cube.FaceDown)).(DoubleTallGrass); !ok || bottom.Type != d.Type || bottom.UpperPart {
			breakBlockNoDrops(d, pos, tx)
		}
	} else if upper, ok := tx.Block(pos.Side(cube.FaceUp)).(DoubleTallGrass); !ok || upper.Type != d.Type || !upper.UpperPart {
		breakBlockNoDrops(d, pos, tx)
	} else if !supportsVegetation(d, tx.Block(pos.Side(cube.FaceDown))) {
		breakBlock(d, pos, tx)
	}
}

func (d DoubleTallGrass) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, d)
	if !used || !replaceableWith(tx, pos.Side(cube.FaceUp), d) || !supportsVegetation(d, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(tx, pos, d, user, ctx)
	place(tx, pos.Side(cube.FaceUp), DoubleTallGrass{Type: d.Type, UpperPart: true}, user, ctx)
	return placed(ctx)
}

func (d DoubleTallGrass) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 100, true)
}

func (d DoubleTallGrass) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, grassDrops(d))
}

func (d DoubleTallGrass) CompostChance() float64 {
	if d.Type == FernDoubleTallGrass() {
		return 0.65
	}
	return 0.5
}

func (d DoubleTallGrass) EncodeItem() (name string, meta int16) {
	return "minecraft:" + d.Type.String(), 0
}

func (d DoubleTallGrass) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + d.Type.String(), map[string]any{"upper_block_bit": d.UpperPart}
}

func allDoubleTallGrass() (b []world.Block) {
	for _, g := range DoubleTallGrassTypes() {
		b = append(b, DoubleTallGrass{Type: g})
		b = append(b, DoubleTallGrass{Type: g, UpperPart: true})
	}
	return
}
