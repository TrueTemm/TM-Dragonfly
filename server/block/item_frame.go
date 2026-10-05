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
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
)

type ItemFrame struct {
	empty
	transparent
	sourceWaterDisplacer

	Facing cube.Face

	Item item.Stack

	Rotations int

	DropChance float64

	Glowing bool
}

func (i ItemFrame) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, ctx *item.UseContext) bool {
	if !i.Item.Empty() {

		i.Rotations = (i.Rotations + 1) % 8
		tx.PlaySound(pos.Vec3Centre(), sound.ItemFrameRotate{})
	} else if held, _ := u.HeldItems(); !held.Empty() {
		i.Item = held.Grow(-held.Count() + 1)

		ctx.SubtractFromCount(1)
		tx.PlaySound(pos.Vec3Centre(), sound.ItemAdd{})
	} else {
		return true
	}

	tx.SetBlock(pos, i, nil)
	return true
}

func (i ItemFrame) Punch(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User) {
	if i.Item.Empty() {
		return
	}

	if g, ok := u.(interface {
		GameMode() world.GameMode
	}); ok {
		if rand.Float64() <= i.DropChance && !g.GameMode().CreativeInventory() {
			dropItem(tx, i.Item, pos.Vec3Centre())
		}
	}
	i.Item, i.Rotations = item.Stack{}, 0
	tx.PlaySound(pos.Vec3Centre(), sound.ItemFrameRemove{})
	tx.SetBlock(pos, i, nil)
}

func (i ItemFrame) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, i)
	if !used {
		return false
	}
	if _, ok := tx.Block(pos.Side(face.Opposite())).Model().(model.Empty); ok {

		return false
	}
	i.Facing = face.Opposite()
	i.DropChance = 1.0

	place(tx, pos, i, user, ctx)
	return placed(ctx)
}

func (i ItemFrame) BreakInfo() BreakInfo {
	return newBreakInfo(0.25, alwaysHarvestable, nothingEffective, oneOf(ItemFrame{Glowing: i.Glowing})).withBreakHandler(func(pos cube.Pos, tx *world.Tx, _ item.User) {
		if !i.Item.Empty() {
			dropItem(tx, i.Item, pos.Vec3Centre())
		}
	})
}

func (i ItemFrame) EncodeItem() (name string, meta int16) {
	if i.Glowing {
		return "minecraft:glow_frame", 0
	}
	return "minecraft:frame", 0
}

func (i ItemFrame) EncodeBlock() (name string, properties map[string]any) {
	name = "minecraft:frame"
	if i.Glowing {
		name = "minecraft:glow_frame"
	}
	return name, map[string]any{
		"facing_direction":     int32(i.Facing.Opposite()),
		"item_frame_map_bit":   uint8(0),
		"item_frame_photo_bit": uint8(0),
	}
}

func (i ItemFrame) DecodeNBT(data map[string]any) any {
	i.DropChance = float64(nbtconv.Float32(data, "ItemDropChance"))
	i.Rotations = int(nbtconv.Uint8(data, "ItemRotation"))
	i.Item = item.MapNBT(data, "Item")
	return i
}

func (i ItemFrame) EncodeNBT() map[string]any {
	m := map[string]any{
		"ItemDropChance": float32(i.DropChance),
		"ItemRotation":   uint8(i.Rotations),
		"id":             "ItemFrame",
	}
	if i.Glowing {
		m["id"] = "GlowItemFrame"
	}
	if !i.Item.Empty() {
		m["Item"] = item.WriteNBT(i.Item, true)
	}
	return m
}

func (i ItemFrame) Pick() item.Stack {
	if i.Item.Empty() {
		return item.NewStack(ItemFrame{Glowing: i.Glowing}, 1)
	}
	return i.Item.Grow(-i.Item.Count() + 1)
}

func (ItemFrame) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (i ItemFrame) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Block(pos.Side(i.Facing)).Model().(model.Empty); ok {

		breakBlock(i, pos, tx)
	}
}

func allItemFrames() (frames []world.Block) {
	for _, f := range cube.Faces() {
		frames = append(frames, ItemFrame{Facing: f, Glowing: true})
		frames = append(frames, ItemFrame{Facing: f, Glowing: false})
	}
	return
}
