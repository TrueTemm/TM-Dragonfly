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

type EndPortalFrame struct {
	bassDrum

	Eye bool

	Facing cube.Direction
}

func (f EndPortalFrame) Model() world.BlockModel {
	return model.EndPortalFrame{}
}

func (EndPortalFrame) LightEmissionLevel() uint8 {
	return 1
}

func (EndPortalFrame) EncodeItem() (name string, meta int16) {
	return "minecraft:end_portal_frame", 0
}

func (f EndPortalFrame) EncodeBlock() (string, map[string]any) {
	return "minecraft:end_portal_frame", map[string]any{
		"end_portal_eye_bit":           f.Eye,
		"minecraft:cardinal_direction": f.Facing.String(),
	}
}

func (f EndPortalFrame) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, f)
	if !used {
		return false
	}
	f.Facing = user.Rotation().Direction().Opposite()
	f.Eye = false
	place(tx, pos, f, user, ctx)
	return placed(ctx)
}

func (f EndPortalFrame) EndPortalFrameState() (eye bool, facing cube.Direction) {
	return f.Eye, f.Facing
}

func (EndPortalFrame) EncodeNBT() map[string]any {
	return map[string]any{"id": "EndPortal"}
}

func (f EndPortalFrame) DecodeNBT(map[string]any) any {
	return f
}

func (f EndPortalFrame) InsertEndPortalEye() (world.Block, bool) {
	if f.Eye {
		return f, false
	}
	f.Eye = true
	return f, true
}

func allEndPortalFrames() []world.Block {
	frames := make([]world.Block, 0, len(cube.Directions())*2)
	for _, dir := range cube.Directions() {
		for _, eye := range []bool{false, true} {
			frames = append(frames, EndPortalFrame{Facing: dir, Eye: eye})
		}
	}
	return frames
}
