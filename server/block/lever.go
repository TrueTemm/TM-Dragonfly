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
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Lever struct {
	empty
	transparent
	flowingWaterDisplacer

	Powered bool

	Facing cube.Face

	Direction cube.Direction
}

func (l Lever) RedstonePower(cube.Pos, *world.Tx, cube.Face) int {
	if l.Powered {
		return 15
	}
	return 0
}

func (l Lever) RedstoneStrongPower(_ cube.Pos, _ *world.Tx, face cube.Face) int {
	if l.Powered && l.Facing.Opposite() == face {
		return 15
	}
	return 0
}

func (l Lever) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (l Lever) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	supportPos := pos.Side(l.Facing.Opposite())
	if !tx.Block(supportPos).Model().FaceSolid(supportPos, l.Facing, tx) {
		breakBlock(l, pos, tx)
	}
}

func (l Lever) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, l)
	if !used {
		return false
	}
	supportPos := pos.Side(face.Opposite())
	if !tx.Block(supportPos).Model().FaceSolid(supportPos, face, tx) {
		return false
	}

	l.Powered = false
	l.Facing = face
	l.Direction = cube.North
	if face.Axis() == cube.Y && user.Rotation().Direction().Face().Axis() == cube.X {
		l.Direction = cube.West
	}
	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

func (l Lever) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, _ item.User, _ *item.UseContext) bool {
	l.Powered = !l.Powered
	tx.SetBlock(pos, l, nil)
	if l.Powered {
		tx.PlaySound(pos.Vec3Centre(), sound.PowerOn{})
	} else {
		tx.PlaySound(pos.Vec3Centre(), sound.PowerOff{})
	}
	return true
}

func (l Lever) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, nothingEffective, oneOf(Lever{}))
}

func (l Lever) EncodeItem() (name string, meta int16) {
	return "minecraft:lever", 0
}

func (l Lever) EncodeBlock() (string, map[string]any) {
	direction := l.Facing.String()
	if l.Facing == cube.FaceDown || l.Facing == cube.FaceUp {
		axis := "east_west"
		if l.Direction == cube.North {
			axis = "north_south"
		}
		direction += "_" + axis
	}
	return "minecraft:lever", map[string]any{"open_bit": l.Powered, "lever_direction": direction}
}

func allLevers() (all []world.Block) {
	f := func(facing cube.Face, direction cube.Direction) {
		all = append(all, Lever{Facing: facing, Direction: direction})
		all = append(all, Lever{Facing: facing, Direction: direction, Powered: true})
	}
	for _, facing := range cube.Faces() {
		f(facing, cube.North)
		if facing == cube.FaceDown || facing == cube.FaceUp {
			f(facing, cube.West)
		}
	}
	return
}
