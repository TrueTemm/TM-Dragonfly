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
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"sync/atomic"
)

type enderChestOwner interface {
	ContainerOpener
	EnderChestInventory() *inventory.Inventory
}

type EnderChest struct {
	chest
	transparent
	bass
	sourceWaterDisplacer

	Facing cube.Direction

	viewers *atomic.Int64
}

func NewEnderChest() EnderChest {
	return EnderChest{viewers: &atomic.Int64{}}
}

func (EnderChest) ContainerSize() int { return 27 }

func (c EnderChest) BreakInfo() BreakInfo {
	return newBreakInfo(22.5, pickaxeHarvestable, pickaxeEffective, silkTouchDrop(item.NewStack(Obsidian{}, 8), item.NewStack(NewEnderChest(), 1))).withBlastResistance(600)
}

func (c EnderChest) LightEmissionLevel() uint8 {
	return 7
}

func (EnderChest) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c EnderChest) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, c)
	if !used {
		return
	}

	c = NewEnderChest()
	c.Facing = user.Rotation().Direction().Opposite()

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (c EnderChest) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(enderChestOwner); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (c EnderChest) AddViewer(tx *world.Tx, pos cube.Pos) {
	if c.viewers.Add(1) == 1 {
		c.open(tx, pos)
	}
}

func (c EnderChest) RemoveViewer(tx *world.Tx, pos cube.Pos) {
	if c.viewers.Load() == 0 {
		return
	}
	if c.viewers.Add(-1) == 0 {
		c.close(tx, pos)
	}
}

func (c EnderChest) open(tx *world.Tx, pos cube.Pos) {
	for _, v := range tx.Viewers(pos.Vec3()) {
		v.ViewBlockAction(pos, OpenAction{})
	}
	tx.PlaySound(pos.Vec3Centre(), sound.EnderChestOpen{})
}

func (c EnderChest) close(tx *world.Tx, pos cube.Pos) {
	for _, v := range tx.Viewers(pos.Vec3()) {
		v.ViewBlockAction(pos, CloseAction{})
	}
	tx.PlaySound(pos.Vec3Centre(), sound.EnderChestClose{})
}

func (c EnderChest) EncodeNBT() map[string]interface{} {
	return map[string]interface{}{"id": "EnderChest"}
}

func (c EnderChest) DecodeNBT(map[string]interface{}) interface{} {
	ec := NewEnderChest()
	ec.Facing = c.Facing
	return ec
}

func (EnderChest) EncodeItem() (name string, meta int16) {
	return "minecraft:ender_chest", 0
}

func (c EnderChest) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:ender_chest", map[string]any{"minecraft:cardinal_direction": c.Facing.String()}
}

func allEnderChests() (chests []world.Block) {
	for _, direction := range cube.Directions() {
		chests = append(chests, EnderChest{Facing: direction})
	}
	return
}
