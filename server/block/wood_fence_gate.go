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
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type WoodFenceGate struct {
	transparent
	bass
	sourceWaterDisplacer

	Wood WoodType

	Facing cube.Direction

	Open bool

	Lowered bool
}

func (f WoodFenceGate) BreakInfo() BreakInfo {
	return newBreakInfo(2, alwaysHarvestable, axeEffective, oneOf(f)).withBlastResistance(3)
}

func (f WoodFenceGate) FlammabilityInfo() FlammabilityInfo {
	if !f.Wood.Flammable() {
		return newFlammabilityInfo(0, 0, false)
	}
	return newFlammabilityInfo(5, 20, true)
}

func (f WoodFenceGate) FuelInfo() item.FuelInfo {
	if !f.Wood.Flammable() {
		return item.FuelInfo{}
	}
	return newFuelInfo(time.Second * 15)
}

func (f WoodFenceGate) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, f)
	if !used {
		return false
	}
	f.Facing = user.Rotation().Direction()
	f.Lowered = f.shouldBeLowered(pos, tx)

	place(tx, pos, f, user, ctx)
	return placed(ctx)
}

func (f WoodFenceGate) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if f.shouldBeLowered(pos, tx) != f.Lowered {
		f.Lowered = !f.Lowered
		tx.SetBlock(pos, f, nil)
	}
}

func (f WoodFenceGate) shouldBeLowered(pos cube.Pos, tx *world.Tx) bool {
	leftSide := f.Facing.RotateLeft().Face()
	_, left := tx.Block(pos.Side(leftSide)).(Wall)
	_, right := tx.Block(pos.Side(leftSide.Opposite())).(Wall)
	return left || right
}

func (f WoodFenceGate) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	f.Open = !f.Open
	if f.Open && f.Facing.Opposite() == u.Rotation().Direction() {
		f.Facing = f.Facing.Opposite()
	}
	tx.SetBlock(pos, f, nil)
	if f.Open {
		tx.PlaySound(pos.Vec3Centre(), sound.FenceGateOpen{Block: f})
		return true
	}
	tx.PlaySound(pos.Vec3Centre(), sound.FenceGateClose{Block: f})
	return true
}

func (f WoodFenceGate) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (f WoodFenceGate) EncodeItem() (name string, meta int16) {
	if f.Wood == OakWood() {
		return "minecraft:fence_gate", 0
	}
	return "minecraft:" + f.Wood.String() + "_fence_gate", 0
}

func (f WoodFenceGate) EncodeBlock() (name string, properties map[string]any) {
	if f.Wood == OakWood() {
		return "minecraft:fence_gate", map[string]any{"minecraft:cardinal_direction": f.Facing.String(), "open_bit": f.Open, "in_wall_bit": f.Lowered}
	}
	return "minecraft:" + f.Wood.String() + "_fence_gate", map[string]any{"minecraft:cardinal_direction": f.Facing.String(), "open_bit": f.Open, "in_wall_bit": f.Lowered}
}

func (f WoodFenceGate) Model() world.BlockModel {
	return model.FenceGate{Facing: f.Facing, Open: f.Open}
}

func allFenceGates() (fenceGates []world.Block) {
	for _, w := range WoodTypes() {
		for i := cube.Direction(0); i <= 3; i++ {
			fenceGates = append(fenceGates, WoodFenceGate{Wood: w, Facing: i, Open: false, Lowered: false})
			fenceGates = append(fenceGates, WoodFenceGate{Wood: w, Facing: i, Open: false, Lowered: true})
			fenceGates = append(fenceGates, WoodFenceGate{Wood: w, Facing: i, Open: true, Lowered: true})
			fenceGates = append(fenceGates, WoodFenceGate{Wood: w, Facing: i, Open: true, Lowered: false})
		}
	}
	return
}
