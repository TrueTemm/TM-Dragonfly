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

type RedstoneBlock struct {
	solid
}

func (r RedstoneBlock) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(r)).withBlastResistance(6)
}

func (r RedstoneBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:redstone_block", 0
}

func (r RedstoneBlock) EncodeBlock() (string, map[string]any) {
	return "minecraft:redstone_block", nil
}

func (r RedstoneBlock) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, r)
	if !used {
		return false
	}
	place(tx, pos, r, user, ctx)
	return placed(ctx)
}

func (RedstoneBlock) RedstonePower(cube.Pos, *world.Tx, cube.Face) int {
	return 15
}

func (RedstoneBlock) RedstoneStrongPower(cube.Pos, *world.Tx, cube.Face) int {
	return 0
}

func (RedstoneBlock) RedstoneNonConductive() {}
