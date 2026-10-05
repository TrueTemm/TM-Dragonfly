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
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"image/color"
	"time"
)

type Sign struct {
	transparent
	empty
	bass
	sourceWaterDisplacer

	Wood WoodType

	Attach Attachment

	Waxed bool

	Front SignText

	Back SignText
}

type SignText struct {
	Text string

	BaseColour color.RGBA

	Glowing bool

	Owner string
}

func (s Sign) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (s Sign) MaxCount() int {
	return 16
}

func (s Sign) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, true)
}

func (s Sign) FuelInfo() item.FuelInfo {
	if !s.Wood.Flammable() {
		return item.FuelInfo{}
	}
	return newFuelInfo(time.Second * 10)
}

func (s Sign) EncodeItem() (name string, meta int16) {
	return "minecraft:" + s.Wood.String() + "_sign", 0
}

func (s Sign) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, axeEffective, oneOf(Sign{Wood: s.Wood}))
}

func (s Sign) Dye(pos cube.Pos, userPos mgl64.Vec3, c item.Colour) (world.Block, bool) {
	if s.EditingFrontSide(pos, userPos) {
		if s.Front.BaseColour == c.SignRGBA() {
			return s, false
		}
		s.Front.BaseColour = c.SignRGBA()
	} else {
		if s.Back.BaseColour == c.SignRGBA() {
			return s, false
		}
		s.Back.BaseColour = c.SignRGBA()
	}
	return s, true
}

func (s Sign) Ink(pos cube.Pos, userPos mgl64.Vec3, glowing bool) (world.Block, bool) {
	if s.EditingFrontSide(pos, userPos) {
		if s.Front.Glowing == glowing {
			return s, false
		}
		s.Front.Glowing = glowing
	} else {
		if s.Back.Glowing == glowing {
			return s, false
		}
		s.Back.Glowing = glowing
	}
	return s, true
}

func (s Sign) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if s.Waxed {
		return s, false
	}
	s.Waxed = true
	return s, true
}

func (s Sign) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if editor, ok := u.(SignEditor); ok && !s.Waxed {
		editor.OpenSign(pos, s.EditingFrontSide(pos, u.Position()))
	} else if s.Waxed {
		tx.PlaySound(pos.Vec3(), sound.WaxedSignFailedInteraction{})
	}
	return true
}

func (s Sign) EditingFrontSide(pos cube.Pos, userPos mgl64.Vec3) bool {
	return userPos.Sub(pos.Vec3Centre()).Dot(s.Attach.Rotation().Vec3()) > 0
}

type SignEditor interface {
	OpenSign(pos cube.Pos, frontSide bool)
}

func (s Sign) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, s)
	if !used || face == cube.FaceDown {
		return false
	}

	if face == cube.FaceUp {
		s.Attach = StandingAttachment(user.Rotation().Orientation().Opposite())
	} else {
		s.Attach = WallAttachment(face.Direction())
	}
	place(tx, pos, s, user, ctx)
	if editor, ok := user.(SignEditor); ok {
		editor.OpenSign(pos, true)
	}
	return placed(ctx)
}

func (s Sign) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if s.Attach.hanging {
		if _, ok := tx.Block(pos.Side(s.Attach.facing.Opposite().Face())).(Air); ok {
			breakBlock(s, pos, tx)
		}
	} else if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Air); ok {
		breakBlock(s, pos, tx)
	}
}

func (s Sign) EncodeBlock() (name string, properties map[string]any) {
	woodType := s.Wood.String() + "_"
	switch s.Wood {
	case OakWood():
		woodType = ""
	case DarkOakWood():
		woodType = "darkoak_"
	}
	if s.Attach.hanging {
		return "minecraft:" + woodType + "wall_sign", map[string]any{"facing_direction": int32(s.Attach.facing + 2)}
	}
	return "minecraft:" + woodType + "standing_sign", map[string]any{"ground_sign_direction": int32(s.Attach.o)}
}

func (s Sign) DecodeNBT(data map[string]any) any {
	if nbtconv.String(data, "Text") != "" {

		s.Front.Text = nbtconv.String(data, "Text")
		s.Front.BaseColour = nbtconv.RGBAFromInt32(nbtconv.Int32(data, "SignTextColor"))
		s.Front.Glowing = nbtconv.Bool(data, "IgnoreLighting") && nbtconv.Bool(data, "TextIgnoreLegacyBugResolved")
		return s
	}

	front, ok := data["FrontText"].(map[string]any)
	if ok {
		s.Front.BaseColour = nbtconv.RGBAFromInt32(nbtconv.Int32(front, "SignTextColor"))
		s.Front.Glowing = nbtconv.Bool(front, "IgnoreLighting")
		s.Front.Text = nbtconv.String(front, "Text")
		s.Front.Owner = nbtconv.String(front, "TextOwner")
	}

	s.Waxed = nbtconv.Bool(data, "IsWaxed")

	back, ok := data["BackText"].(map[string]any)
	if ok {
		s.Back.BaseColour = nbtconv.RGBAFromInt32(nbtconv.Int32(back, "SignTextColor"))
		s.Back.Glowing = nbtconv.Bool(back, "IgnoreLighting")
		s.Back.Text = nbtconv.String(back, "Text")
		s.Back.Owner = nbtconv.String(back, "TextOwner")
	}

	return s
}

func (s Sign) EncodeNBT() map[string]any {
	m := map[string]any{
		"id":      "Sign",
		"IsWaxed": boolByte(s.Waxed),
		"FrontText": map[string]any{
			"SignTextColor":  nbtconv.Int32FromRGBA(s.Front.BaseColour),
			"IgnoreLighting": boolByte(s.Front.Glowing),
			"Text":           s.Front.Text,
			"TextOwner":      s.Front.Owner,
		},
		"BackText": map[string]any{
			"SignTextColor":  nbtconv.Int32FromRGBA(s.Back.BaseColour),
			"IgnoreLighting": boolByte(s.Back.Glowing),
			"Text":           s.Back.Text,
			"TextOwner":      s.Back.Owner,
		},
	}
	return m
}

func allSigns() (signs []world.Block) {
	for _, w := range WoodTypes() {
		for _, d := range cube.Directions() {
			signs = append(signs, Sign{Wood: w, Attach: WallAttachment(d)})
		}
		for o := cube.Orientation(0); o <= 15; o++ {
			signs = append(signs, Sign{Wood: w, Attach: StandingAttachment(o)})
		}
	}
	return
}
