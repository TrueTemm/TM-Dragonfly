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
	"math/rand/v2"
)

type Vines struct {
	replaceable
	transparent
	empty
	sourceWaterDisplacer

	NorthDirection bool

	EastDirection bool

	SouthDirection bool

	WestDirection bool
}

func (Vines) CompostChance() float64 {
	return 0.5
}

func (Vines) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (Vines) HasLiquidDrops() bool {
	return false
}

func (Vines) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(15, 100, true)
}

func (v Vines) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, func(t item.Tool) bool {
		return t.ToolType() == item.TypeShears
	}, axeEffective, oneOf(v))
}

func (Vines) EntityInside(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if fallEntity, ok := e.(fallDistanceEntity); ok {
		fallEntity.ResetFallDistance()
	}
}

func (v Vines) WithAttachment(direction cube.Direction, attached bool) Vines {
	switch direction {
	case cube.North:
		v.NorthDirection = attached
		return v
	case cube.East:
		v.EastDirection = attached
		return v
	case cube.South:
		v.SouthDirection = attached
		return v
	case cube.West:
		v.WestDirection = attached
		return v
	}
	panic("should never happen")
}

func (v Vines) Attachment(direction cube.Direction) bool {
	switch direction {
	case cube.North:
		return v.NorthDirection
	case cube.East:
		return v.EastDirection
	case cube.South:
		return v.SouthDirection
	case cube.West:
		return v.WestDirection
	}
	panic("should never happen")
}

func (v Vines) Attachments() (attachments []cube.Direction) {
	for _, d := range cube.Directions() {
		if v.Attachment(d) {
			attachments = append(attachments, d)
		}
	}
	return
}

func (v Vines) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if !v.canSpreadTo(tx, pos) || face.Axis() == cube.Y {
		return false
	}
	pos, face, used := firstReplaceable(tx, pos, face, v)
	if !used {
		return false
	}
	if _, ok := tx.Block(pos).(Vines); ok {

		return false
	}

	v = v.WithAttachment(face.Direction().Opposite(), true)

	place(tx, pos, v, user, ctx)
	return placed(ctx)
}

func (v Vines) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	above, updated := tx.Block(pos.Side(cube.FaceUp)), false
	for _, d := range v.Attachments() {
		if !v.canSpreadTo(tx, pos.Side(d.Face())) {
			if o, ok := above.(Vines); !ok || !o.Attachment(d) {

				v = v.WithAttachment(d, false)
				updated = true
			}
		}
	}
	if !updated {
		return
	}
	if len(v.Attachments()) == 0 {
		breakBlock(v, pos, tx)
		return
	}
	tx.SetBlock(pos, v, nil)
}

func (v Vines) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if r.Float64() > 0.25 {

		return
	}

	face := cube.Face(r.IntN(len(cube.Faces())))
	selectedPos := pos.Side(face)

	if face.Axis() != cube.Y && !v.Attachment(face.Direction()) {
		if !v.canSpread(tx, pos) {

			return
		}

		if _, ok := tx.Block(selectedPos).(Air); ok {
			rightRotatedFace := face.RotateRight()
			leftRotatedFace := face.RotateLeft()

			attachedOnRight := v.Attachment(rightRotatedFace.Direction())
			attachedOnLeft := v.Attachment(leftRotatedFace.Direction())

			rightSelectedPos := selectedPos.Side(rightRotatedFace)
			leftSelectedPos := selectedPos.Side(leftRotatedFace)

			if attachedOnRight && v.canSpreadTo(tx, rightSelectedPos) {
				tx.SetBlock(selectedPos, (Vines{}).WithAttachment(rightRotatedFace.Direction(), true), nil)
			} else if attachedOnLeft && v.canSpreadTo(tx, leftSelectedPos) {
				tx.SetBlock(selectedPos, (Vines{}).WithAttachment(leftRotatedFace.Direction(), true), nil)
			} else if _, ok = tx.Block(rightSelectedPos).(Air); ok && attachedOnRight && v.canSpreadTo(tx, pos.Side(rightRotatedFace)) {
				tx.SetBlock(rightSelectedPos, (Vines{}).WithAttachment(face.Opposite().Direction(), true), nil)
			} else if _, ok = tx.Block(leftSelectedPos).(Air); ok && attachedOnLeft && v.canSpreadTo(tx, pos.Side(leftRotatedFace)) {
				tx.SetBlock(leftSelectedPos, (Vines{}).WithAttachment(face.Opposite().Direction(), true), nil)
			}
		} else if v.canSpreadTo(tx, selectedPos) {

			tx.SetBlock(pos, v.WithAttachment(face.Direction(), true), nil)
		}
		return
	}

	if face == cube.FaceUp && selectedPos.OutOfBounds(tx.Range()) {

		if _, ok := tx.Block(selectedPos).(Air); ok {
			if !v.canSpread(tx, pos) {

				return
			}
			newVines := Vines{}
			for _, f := range cube.HorizontalFaces() {

				if r.IntN(2) == 0 && v.Attachment(f.Direction()) && v.canSpreadTo(tx, selectedPos.Side(f)) {
					newVines = newVines.WithAttachment(f.Direction(), true)
				}
			}
			if len(newVines.Attachments()) > 0 {
				tx.SetBlock(selectedPos, newVines, nil)
			}
			return
		}
	}

	selectedPos = pos.Side(cube.FaceDown)
	if selectedPos.OutOfBounds(tx.Range()) {
		return
	}
	newVines, vines := tx.Block(selectedPos).(Vines)
	if _, ok := tx.Block(selectedPos).(Air); !ok && !vines {

		return
	}
	var changed bool
	for _, f := range cube.HorizontalFaces() {

		if r.IntN(2) == 0 && v.Attachment(f.Direction()) && !newVines.Attachment(f.Direction()) {
			newVines, changed = newVines.WithAttachment(f.Direction(), true), true
		}
	}
	if changed {
		tx.SetBlock(selectedPos, newVines, nil)
	}
}

func (Vines) EncodeItem() (name string, meta int16) {
	return "minecraft:vine", 0
}

func (v Vines) EncodeBlock() (string, map[string]any) {
	var bits int
	for i, ok := range []bool{v.SouthDirection, v.WestDirection, v.NorthDirection, v.EastDirection} {
		if ok {
			bits |= 1 << i
		}
	}
	return "minecraft:vine", map[string]any{"vine_direction_bits": int32(bits)}
}

func (Vines) canSpreadTo(tx *world.Tx, pos cube.Pos) bool {
	switch tx.Block(pos).Model().(type) {
	case model.Solid, model.Leaves:
		return true
	}
	return false
}

func (v Vines) canSpread(tx *world.Tx, pos cube.Pos) bool {
	var count int
	for x := -4; x <= 4; x++ {
		for z := -4; z <= 4; z++ {
			for y := -1; y <= 1; y++ {
				if _, ok := tx.Block(pos.Add(cube.Pos{x, y, z})).(Vines); ok {
					count++

					if count >= 5 {
						return false
					}
				}
			}
		}
	}
	return true
}

func allVines() (b []world.Block) {
	for _, north := range []bool{true, false} {
		for _, east := range []bool{true, false} {
			for _, south := range []bool{true, false} {
				for _, west := range []bool{true, false} {
					b = append(b, Vines{
						NorthDirection: north,
						EastDirection:  east,
						SouthDirection: south,
						WestDirection:  west,
					})
				}
			}
		}
	}
	return
}
