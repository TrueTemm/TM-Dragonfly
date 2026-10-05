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

type NetherSprouts struct {
	transparent
	replaceable
	empty
}

func (n NetherSprouts) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !supportsVegetation(n, tx.Block(pos.Side(cube.FaceDown))) {
		breakBlock(n, pos, tx)
	}
}

func (n NetherSprouts) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, n)
	if !used {
		return false
	}
	if !supportsVegetation(n, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(tx, pos, n, user, ctx)
	return placed(ctx)
}

func (n NetherSprouts) HasLiquidDrops() bool {
	return false
}

func (n NetherSprouts) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, true)
}

func (n NetherSprouts) BreakInfo() BreakInfo {
	return newBreakInfo(0, func(t item.Tool) bool {
		return t.ToolType() == item.TypeShears
	}, nothingEffective, oneOf(n))
}

func (NetherSprouts) CompostChance() float64 {
	return 0.5
}

func (n NetherSprouts) EncodeItem() (name string, meta int16) {
	return "minecraft:nether_sprouts", 0
}

func (n NetherSprouts) EncodeBlock() (string, map[string]any) {
	return "minecraft:nether_sprouts", nil
}
