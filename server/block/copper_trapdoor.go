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
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type CopperTrapdoor struct {
	transparent
	bass
	sourceWaterDisplacer

	Oxidation OxidationType

	Waxed bool

	Facing cube.Direction

	Open bool

	Top bool
}

func (t CopperTrapdoor) Model() world.BlockModel {
	return model.Trapdoor{Facing: t.Facing, Top: t.Top, Open: t.Open}
}

func (t CopperTrapdoor) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, t)
	if !used {
		return false
	}
	t.Facing = user.Rotation().Direction().Opposite()
	t.Top = (clickPos.Y() > 0.5 && face != cube.FaceUp) || face == cube.FaceDown

	place(tx, pos, t, user, ctx)
	return placed(ctx)
}

func (t CopperTrapdoor) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if t.Waxed {
		return t, false
	}
	t.Waxed = true
	return t, true
}

func (t CopperTrapdoor) Strip() (world.Block, world.Sound, bool) {
	if t.Waxed {
		t.Waxed = false
		return t, sound.WaxRemoved{}, true
	} else if ot, ok := t.Oxidation.Decrease(); ok {
		t.Oxidation = ot
		return t, sound.CopperScraped{}, true
	}
	return t, nil, false
}

func (t CopperTrapdoor) CanOxidate() bool {
	return !t.Waxed
}

func (t CopperTrapdoor) OxidationLevel() OxidationType {
	return t.Oxidation
}

func (t CopperTrapdoor) WithOxidationLevel(o OxidationType) Oxidisable {
	t.Oxidation = o
	return t
}

func (t CopperTrapdoor) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, _ item.User, _ *item.UseContext) bool {
	t.Open = !t.Open
	tx.SetBlock(pos, t, nil)
	if t.Open {
		tx.PlaySound(pos.Vec3Centre(), sound.TrapdoorOpen{Block: t})
		return true
	}
	tx.PlaySound(pos.Vec3Centre(), sound.TrapdoorClose{Block: t})
	return true
}

func (t CopperTrapdoor) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, t)
}

func (t CopperTrapdoor) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(t)).withBlastResistance(6)
}

func (t CopperTrapdoor) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (t CopperTrapdoor) EncodeItem() (name string, meta int16) {
	return copperBlockName("copper_trapdoor", t.Oxidation, t.Waxed), 0
}

func (t CopperTrapdoor) EncodeBlock() (name string, properties map[string]any) {
	return copperBlockName("copper_trapdoor", t.Oxidation, t.Waxed), map[string]any{"direction": int32(math.Abs(float64(t.Facing) - 3)), "open_bit": t.Open, "upside_down_bit": t.Top}
}

func allCopperTrapdoors() (trapdoors []world.Block) {
	f := func(waxed bool) {
		for _, o := range OxidationTypes() {
			for i := cube.Direction(0); i <= 3; i++ {
				trapdoors = append(trapdoors, CopperTrapdoor{Oxidation: o, Waxed: waxed, Facing: i, Open: false, Top: false})
				trapdoors = append(trapdoors, CopperTrapdoor{Oxidation: o, Waxed: waxed, Facing: i, Open: false, Top: true})
				trapdoors = append(trapdoors, CopperTrapdoor{Oxidation: o, Waxed: waxed, Facing: i, Open: true, Top: true})
				trapdoors = append(trapdoors, CopperTrapdoor{Oxidation: o, Waxed: waxed, Facing: i, Open: true, Top: false})
			}
		}
	}
	f(false)
	f(true)
	return
}
