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
)

type DeadBush struct {
	empty
	replaceable
	transparent
	sourceWaterDisplacer
}

func (d DeadBush) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !supportsVegetation(d, tx.Block(pos.Side(cube.FaceDown))) {
		breakBlock(d, pos, tx)
	}
}

func (d DeadBush) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, d)
	if !used || !supportsVegetation(d, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(tx, pos, d, user, ctx)
	return placed(ctx)
}

func (d DeadBush) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (d DeadBush) HasLiquidDrops() bool {
	return true
}

func (d DeadBush) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 100, true)
}

func (d DeadBush) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if t.ToolType() == item.TypeShears {
			return []item.Stack{item.NewStack(d, 1)}
		}
		if amount := rand.IntN(3); amount != 0 {
			return []item.Stack{item.NewStack(item.Stick{}, amount)}
		}
		return nil
	})
}

func (d DeadBush) EncodeItem() (name string, meta int16) {
	return "minecraft:deadbush", 0
}

func (d DeadBush) EncodeBlock() (string, map[string]any) {
	return "minecraft:deadbush", nil
}
