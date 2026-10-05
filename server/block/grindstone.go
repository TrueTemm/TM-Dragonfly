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
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Grindstone struct {
	transparent

	Attach GrindstoneAttachment

	Facing cube.Direction
}

func (g Grindstone) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(g)).withBlastResistance(6)
}

func (g Grindstone) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (g Grindstone) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, g)
	if !used {
		return false
	}
	g.Facing = user.Rotation().Direction().Opposite()
	if face == cube.FaceDown {
		g.Attach = HangingGrindstoneAttachment()
	} else if face != cube.FaceUp {
		g.Attach = WallGrindstoneAttachment()
		g.Facing = face.Direction()
	}
	place(tx, pos, g, user, ctx)
	return placed(ctx)
}

func (g Grindstone) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	supportFace := g.Facing.Face().Opposite()
	if g.Attach == HangingGrindstoneAttachment() {
		supportFace = cube.FaceUp
	} else if g.Attach == StandingGrindstoneAttachment() {
		supportFace = cube.FaceDown
	}
	if _, ok := tx.Block(pos.Side(supportFace)).Model().(model.Empty); ok {

		breakBlockNoDrops(g, pos, tx)
		dropItem(tx, item.NewStack(g, 1), pos.Vec3Centre())
	}
}

func (g Grindstone) Model() world.BlockModel {
	axis := cube.Y
	if g.Attach == WallGrindstoneAttachment() {
		axis = g.Facing.Face().Axis()
	}
	return model.Grindstone{Axis: axis}
}

func (g Grindstone) EncodeBlock() (string, map[string]any) {
	return "minecraft:grindstone", map[string]any{
		"attachment": g.Attach.String(),
		"direction":  int32(horizontalDirection(g.Facing)),
	}
}

func (g Grindstone) EncodeItem() (name string, meta int16) {
	return "minecraft:grindstone", 0
}

func allGrindstones() (grindstones []world.Block) {
	for _, a := range GrindstoneAttachments() {
		for _, d := range cube.Directions() {
			grindstones = append(grindstones, Grindstone{Attach: a, Facing: d})
		}
	}
	return
}
