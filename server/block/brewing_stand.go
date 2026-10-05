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
	"time"
)

type BrewingStand struct {
	sourceWaterDisplacer
	transparent
	*brewer

	LeftSlot bool

	MiddleSlot bool

	RightSlot bool
}

func NewBrewingStand() BrewingStand {
	return BrewingStand{brewer: newBrewer()}
}

func (BrewingStand) ContainerSize() int { return 5 }

func (b BrewingStand) Model() world.BlockModel {
	return model.BrewingStand{}
}

func (b BrewingStand) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (b BrewingStand) Tick(_ int64, pos cube.Pos, tx *world.Tx) {

	left, _ := b.inventory.Item(1)
	middle, _ := b.inventory.Item(2)
	right, _ := b.inventory.Item(3)

	displayLeft, displayMiddle, displayRight := b.LeftSlot, b.MiddleSlot, b.RightSlot
	b.LeftSlot, b.MiddleSlot, b.RightSlot = !left.Empty(), !middle.Empty(), !right.Empty()
	if b.LeftSlot != displayLeft || b.MiddleSlot != displayMiddle || b.RightSlot != displayRight {
		tx.SetBlock(pos, b, nil)
	}

	b.tickBrewing("brewing_stand", pos, tx)
}

func (b BrewingStand) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (b BrewingStand) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, b)
	if !used {
		return
	}

	b = NewBrewingStand()
	place(tx, pos, b, user, ctx)
	return placed(ctx)
}

func (b BrewingStand) EncodeNBT() map[string]any {
	if b.brewer == nil {

		b = NewBrewingStand()
	}
	duration := b.Duration()
	fuel, totalFuel := b.Fuel()
	return map[string]any{
		"id":         "BrewingStand",
		"Items":      nbtconv.InvToNBT(b.inventory),
		"CookTime":   int16(duration.Milliseconds() / 50),
		"FuelTotal":  int16(totalFuel),
		"FuelAmount": int16(fuel),
	}
}

func (b BrewingStand) DecodeNBT(data map[string]any) any {
	brew := time.Duration(nbtconv.Int16(data, "CookTime")) * time.Millisecond * 50

	fuel := int32(nbtconv.Int16(data, "FuelAmount"))
	maxFuel := int32(nbtconv.Int16(data, "FuelTotal"))

	b = NewBrewingStand()
	b.setDuration(brew)
	b.setFuel(fuel, maxFuel)
	nbtconv.InvFromNBT(b.inventory, nbtconv.Slice(data, "Items"))
	return b
}

func (b BrewingStand) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, pickaxeEffective, oneOf(BrewingStand{})).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		for _, i := range b.Inventory(tx, pos).Clear() {
			dropItem(tx, i, pos.Vec3Centre())
		}
	})
}

func (b BrewingStand) EncodeBlock() (string, map[string]any) {
	return "minecraft:brewing_stand", map[string]any{
		"brewing_stand_slot_a_bit": b.LeftSlot,
		"brewing_stand_slot_b_bit": b.MiddleSlot,
		"brewing_stand_slot_c_bit": b.RightSlot,
	}
}

func (b BrewingStand) EncodeItem() (name string, meta int16) {
	return "minecraft:brewing_stand", 0
}

func allBrewingStands() (stands []world.Block) {
	for _, left := range []bool{false, true} {
		for _, middle := range []bool{false, true} {
			for _, right := range []bool{false, true} {
				stands = append(stands, BrewingStand{LeftSlot: left, MiddleSlot: middle, RightSlot: right})
			}
		}
	}
	return
}
