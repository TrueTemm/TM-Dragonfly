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
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type Lectern struct {
	bass
	sourceWaterDisplacer

	Facing cube.Direction

	Book item.Stack

	Page int
}

func (Lectern) Model() world.BlockModel {
	return model.Lectern{}
}

func (Lectern) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (Lectern) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (l Lectern) BreakInfo() BreakInfo {
	d := []item.Stack{item.NewStack(Lectern{}, 1)}
	if !l.Book.Empty() {
		d = append(d, l.Book)
	}
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, simpleDrops(d...))
}

func (l Lectern) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, l)
	if !used {
		return false
	}
	l.Facing = user.Rotation().Direction().Opposite()
	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

type readableBook interface {
	TotalPages() int

	Page(page int) (string, bool)
}

func (l Lectern) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, ctx *item.UseContext) bool {
	if !l.Book.Empty() {
		if opener, ok := u.(ContainerOpener); ok {
			opener.OpenBlockContainer(pos, tx)
			return true
		}
		return false
	}

	held, _ := u.HeldItems()
	if _, ok := held.Item().(readableBook); !ok {

		return false
	}

	l.Book, l.Page = held.Grow(-held.Count()+1), 0
	tx.SetBlock(pos, l, nil)

	tx.PlaySound(pos.Vec3Centre(), sound.LecternBookPlace{})
	ctx.SubtractFromCount(1)
	return true
}

func (l Lectern) Punch(pos cube.Pos, _ cube.Face, tx *world.Tx, _ item.User) {
	if l.Book.Empty() {

		return
	}

	dropItem(tx, l.Book, pos.Side(cube.FaceUp).Vec3Middle())

	l.Book = item.Stack{}
	tx.SetBlock(pos, l, nil)
	tx.PlaySound(pos.Vec3Centre(), sound.Attack{})
}

func (l Lectern) TurnPage(pos cube.Pos, tx *world.Tx, page int) error {
	if page == l.Page {

		return nil
	}
	if l.Book.Empty() {
		return fmt.Errorf("lectern at %v is empty", pos)
	}
	if r, ok := l.Book.Item().(readableBook); ok && (page >= r.TotalPages() || page < 0) {
		return fmt.Errorf("page number %d is out of bounds", page)
	}
	l.Page = page
	tx.SetBlock(pos, l, nil)
	return nil
}

func (l Lectern) EncodeNBT() map[string]any {
	m := map[string]any{
		"hasBook": boolByte(!l.Book.Empty()),
		"page":    int32(l.Page),
		"id":      "Lectern",
	}
	if r, ok := l.Book.Item().(readableBook); ok {
		m["book"] = item.WriteNBT(l.Book, true)
		m["totalPages"] = int32(r.TotalPages())
	}
	return m
}

func (l Lectern) DecodeNBT(m map[string]any) any {
	l.Page = int(nbtconv.Int32(m, "page"))
	l.Book = item.MapNBT(m, "book")
	return l
}

func (Lectern) EncodeItem() (name string, meta int16) {
	return "minecraft:lectern", 0
}

func (l Lectern) EncodeBlock() (string, map[string]any) {
	return "minecraft:lectern", map[string]any{
		"minecraft:cardinal_direction": l.Facing.String(),
		"powered_bit":                  uint8(0),
	}
}

func allLecterns() (lecterns []world.Block) {
	for _, f := range cube.Directions() {
		lecterns = append(lecterns, Lectern{Facing: f})
	}
	return
}
