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

type CopperTorch struct {
	transparent
	empty

	Facing cube.Face
}

func (t CopperTorch) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(t))
}

func (t CopperTorch) LightEmissionLevel() uint8 {
	return 14
}

func (t CopperTorch) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, t)
	if !used {
		return false
	}
	if face == cube.FaceDown {
		return false
	}
	if !tx.Block(pos.Side(face.Opposite())).Model().FaceSolid(pos.Side(face.Opposite()), face, tx) {
		found := false
		for _, i := range []cube.Face{cube.FaceSouth, cube.FaceWest, cube.FaceNorth, cube.FaceEast, cube.FaceDown} {
			if tx.Block(pos.Side(i)).Model().FaceSolid(pos.Side(i), i.Opposite(), tx) {
				found = true
				face = i.Opposite()
				break
			}
		}
		if !found {
			return false
		}
	}
	t.Facing = face.Opposite()

	place(tx, pos, t, user, ctx)
	return placed(ctx)
}

func (t CopperTorch) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !tx.Block(pos.Side(t.Facing)).Model().FaceSolid(pos.Side(t.Facing), t.Facing.Opposite(), tx) {
		breakBlock(t, pos, tx)
	}
}

func (t CopperTorch) HasLiquidDrops() bool {
	return true
}

func (t CopperTorch) EncodeItem() (name string, meta int16) {
	return "minecraft:copper_torch", 0
}

func (t CopperTorch) EncodeBlock() (name string, properties map[string]any) {
	var face string
	switch t.Facing {
	case cube.FaceDown:
		face = "top"
	case unknownFace:
		face = "unknown"
	default:
		face = t.Facing.String()
	}

	return "minecraft:copper_torch", map[string]any{"torch_facing_direction": face}
}

func allCopperTorches() (torch []world.Block) {
	for _, face := range cube.Faces() {
		if face == cube.FaceUp {
			face = unknownFace
		}

		torch = append(torch, CopperTorch{Facing: face})
	}
	return
}
