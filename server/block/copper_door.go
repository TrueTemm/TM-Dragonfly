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
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type CopperDoor struct {
	transparent
	bass
	sourceWaterDisplacer

	Oxidation OxidationType

	Waxed bool

	Facing cube.Direction

	Open bool

	Top bool

	Right bool
}

func (d CopperDoor) Strip() (world.Block, world.Sound, bool) {
	if d.Waxed {
		d.Waxed = false
		return d, sound.WaxRemoved{}, true
	} else if ot, ok := d.Oxidation.Decrease(); ok {
		d.Oxidation = ot
		return d, sound.CopperScraped{}, true
	}
	return d, nil, false
}

func (d CopperDoor) Model() world.BlockModel {
	return model.Door{Facing: d.Facing, Open: d.Open, Right: d.Right}
}

func (d CopperDoor) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if d.Waxed {
		return d, false
	}
	d.Waxed = true
	return d, true
}

func (d CopperDoor) CanOxidate() bool {
	return !d.Waxed
}

func (d CopperDoor) OxidationLevel() OxidationType {
	return d.Oxidation
}

func (d CopperDoor) WithOxidationLevel(o OxidationType) Oxidisable {
	d.Oxidation = o
	return d
}

func (d CopperDoor) NeighbourUpdateTick(pos, changedNeighbour cube.Pos, tx *world.Tx) {
	if pos == changedNeighbour {
		return
	}
	if d.Top {
		if b, ok := tx.Block(pos.Side(cube.FaceDown)).(CopperDoor); !ok {
			breakBlockNoDrops(d, pos, tx)
		} else if d.Oxidation != b.Oxidation || d.Waxed != b.Waxed {
			d.Oxidation = b.Oxidation
			d.Waxed = b.Waxed
			tx.SetBlock(pos, d, nil)
		}
	} else if solid := tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, tx); !solid {

		breakBlockNoDrops(d, pos, tx)
		dropItem(tx, item.NewStack(d, 1), pos.Vec3Centre())
	} else if b, ok := tx.Block(pos.Side(cube.FaceUp)).(CopperDoor); !ok {
		breakBlockNoDrops(d, pos, tx)
	} else if d.Oxidation != b.Oxidation || d.Waxed != b.Waxed {
		d.Oxidation = b.Oxidation
		d.Waxed = b.Waxed
		tx.SetBlock(pos, d, nil)
	}
}

func (d CopperDoor) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if face != cube.FaceUp {

		return false
	}
	below := pos
	pos = pos.Side(cube.FaceUp)
	if !replaceableWith(tx, pos, d) || !replaceableWith(tx, pos.Side(cube.FaceUp), d) {
		return false
	}
	if !tx.Block(below).Model().FaceSolid(below, cube.FaceUp, tx) {
		return false
	}
	d.Facing = user.Rotation().Direction()
	left := tx.Block(pos.Side(d.Facing.RotateLeft().Face()))
	right := tx.Block(pos.Side(d.Facing.RotateRight().Face()))
	if _, ok := left.Model().(model.Door); ok {
		d.Right = true
	}

	if diffuser, ok := right.(LightDiffuser); !ok || diffuser.LightDiffusionLevel() != 0 {
		if diffuser, ok := left.(LightDiffuser); ok && diffuser.LightDiffusionLevel() == 0 {
			d.Right = true
		}
	}

	ctx.IgnoreBBox = true
	place(tx, pos, d, user, ctx)
	place(tx, pos.Side(cube.FaceUp), CopperDoor{Oxidation: d.Oxidation, Waxed: d.Waxed, Facing: d.Facing, Top: true, Right: d.Right}, user, ctx)
	ctx.CountSub = 1
	return placed(ctx)
}

func (d CopperDoor) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, _ item.User, _ *item.UseContext) bool {
	d.Open = !d.Open
	tx.SetBlock(pos, d, nil)

	otherPos := pos.Side(cube.Face(boolByte(!d.Top)))
	other := tx.Block(otherPos)
	if door, ok := other.(CopperDoor); ok {
		door.Open = d.Open
		tx.SetBlock(otherPos, door, nil)
	}
	if d.Open {
		tx.PlaySound(pos.Vec3Centre(), sound.DoorOpen{Block: d})
		return true
	}
	tx.PlaySound(pos.Vec3Centre(), sound.DoorClose{Block: d})
	return true
}

func (d CopperDoor) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, d)
}

func (d CopperDoor) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(d)).withBlastResistance(6)
}

func (d CopperDoor) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (d CopperDoor) EncodeItem() (name string, meta int16) {
	return copperBlockName("copper_door", d.Oxidation, d.Waxed), 0
}

func (d CopperDoor) EncodeBlock() (name string, properties map[string]any) {
	return copperBlockName("copper_door", d.Oxidation, d.Waxed), map[string]any{"minecraft:cardinal_direction": d.Facing.RotateRight().String(), "door_hinge_bit": d.Right, "open_bit": d.Open, "upper_block_bit": d.Top}
}

func allCopperDoors() (doors []world.Block) {
	f := func(waxed bool) {
		for _, o := range OxidationTypes() {
			for i := cube.Direction(0); i <= 3; i++ {
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: false, Top: false, Right: false})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: false, Top: true, Right: false})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: true, Top: true, Right: false})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: true, Top: false, Right: false})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: false, Top: false, Right: true})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: false, Top: true, Right: true})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: true, Top: true, Right: true})
				doors = append(doors, CopperDoor{Oxidation: o, Waxed: waxed, Facing: i, Open: true, Top: false, Right: true})
			}
		}
	}
	f(false)
	f(true)
	return
}
