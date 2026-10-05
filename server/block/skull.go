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
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Skull struct {
	transparent
	sourceWaterDisplacer

	Type SkullType

	Attach Attachment
}

func (Skull) Helmet() bool {
	return true
}

func (Skull) DefencePoints() float64 {
	return 0
}

func (Skull) Toughness() float64 {
	return 0
}

func (Skull) KnockBackResistance() float64 {
	return 0
}

func (s Skull) Model() world.BlockModel {
	return model.Skull{Direction: s.Attach.facing.Face(), Hanging: s.Attach.hanging}
}

func (s Skull) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, s)
	if !used || face == cube.FaceDown {
		return false
	}

	if face == cube.FaceUp {
		s.Attach = StandingAttachment(user.Rotation().Orientation())
	} else {
		s.Attach = WallAttachment(face.Direction())
	}
	place(tx, pos, s, user, ctx)
	return placed(ctx)
}

func (Skull) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (Skull) HasLiquidDrops() bool {
	return true
}

func (s Skull) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, nothingEffective, oneOf(Skull{Type: s.Type}))
}

func (s Skull) EncodeItem() (name string, meta int16) {
	return "minecraft:" + s.Type.String(), 0
}

func (s Skull) DecodeNBT(data map[string]interface{}) interface{} {
	if t := skull(nbtconv.Uint8(data, "SkullType")); t != 255 {

		s.Type = SkullType{t}
	}
	s.Attach.o = cube.OrientationFromYaw(float64(nbtconv.Float32(data, "Rotation")))
	return s
}

func (s Skull) EncodeNBT() map[string]interface{} {
	return map[string]interface{}{"id": "Skull", "SkullType": uint8(255), "Rotation": float32(s.Attach.o.Yaw())}
}

func (s Skull) EncodeBlock() (string, map[string]interface{}) {
	if s.Attach.hanging {
		if s.Attach.facing == unknownDirection {
			return "minecraft:" + s.Type.String(), map[string]interface{}{"facing_direction": int32(0)}
		}
		return "minecraft:" + s.Type.String(), map[string]interface{}{"facing_direction": int32(s.Attach.facing) + 2}
	}
	return "minecraft:" + s.Type.String(), map[string]interface{}{"facing_direction": int32(1)}
}

func allSkulls() (skulls []world.Block) {
	for _, t := range SkullTypes() {
		for _, d := range append(cube.Directions(), unknownDirection) {
			skulls = append(skulls, Skull{Type: t, Attach: WallAttachment(d)})
		}
		skulls = append(skulls, Skull{Type: t, Attach: StandingAttachment(0)})
	}
	return
}
