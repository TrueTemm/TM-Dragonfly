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

type CopperLantern struct {
	transparent
	sourceWaterDisplacer

	Hanging bool

	Oxidation OxidationType

	Waxed bool
}

func (c CopperLantern) Model() world.BlockModel {
	return model.Lantern{Hanging: c.Hanging}
}

func (c CopperLantern) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if c.Hanging {
		up := pos.Side(cube.FaceUp)
		if _, ok := tx.Block(up).(CopperChain); !ok && !tx.Block(up).Model().FaceSolid(up, cube.FaceDown, tx) {
			breakBlock(c, pos, tx)
		}
	} else {
		down := pos.Side(cube.FaceDown)
		if !tx.Block(down).Model().FaceSolid(down, cube.FaceUp, tx) {
			breakBlock(c, pos, tx)
		}
	}
}

func (CopperLantern) LightEmissionLevel() uint8 {
	return 15
}

func (c CopperLantern) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, c)
	if !used {
		return false
	}
	if face == cube.FaceDown {
		upPos := pos.Side(cube.FaceUp)
		if _, ok := tx.Block(upPos).(CopperChain); !ok && !tx.Block(upPos).Model().FaceSolid(upPos, cube.FaceDown, tx) {
			face = cube.FaceUp
		}
	}
	if face != cube.FaceDown {
		downPos := pos.Side(cube.FaceDown)
		if !tx.Block(downPos).Model().FaceSolid(downPos, cube.FaceUp, tx) {
			return false
		}
	}
	c.Hanging = face == cube.FaceDown

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (CopperLantern) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c CopperLantern) BreakInfo() BreakInfo {
	return newBreakInfo(3.5, pickaxeHarvestable, pickaxeEffective, oneOf(c))
}

func (c CopperLantern) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if c.Waxed {
		return c, false
	}
	c.Waxed = true
	return c, true
}

func (c CopperLantern) Strip() (world.Block, world.Sound, bool) {
	if c.Waxed {
		c.Waxed = false
		return c, sound.WaxRemoved{}, true
	} else if ot, ok := c.Oxidation.Decrease(); ok {
		c.Oxidation = ot
		return c, sound.CopperScraped{}, true
	}
	return c, nil, false
}

func (c CopperLantern) CanOxidate() bool {
	return !c.Waxed
}

func (c CopperLantern) OxidationLevel() OxidationType {
	return c.Oxidation
}

func (c CopperLantern) WithOxidationLevel(o OxidationType) Oxidisable {
	c.Oxidation = o
	return c
}

func (c CopperLantern) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, c)
}

func (c CopperLantern) EncodeItem() (name string, meta int16) {
	return copperBlockName("copper_lantern", c.Oxidation, c.Waxed), 0
}

func (c CopperLantern) EncodeBlock() (name string, properties map[string]any) {
	return copperBlockName("copper_lantern", c.Oxidation, c.Waxed), map[string]any{"hanging": c.Hanging}
}

func allCopperLanterns() (lanterns []world.Block) {
	f := func(waxed bool) {
		for _, o := range OxidationTypes() {
			lanterns = append(lanterns, CopperLantern{Hanging: false, Oxidation: o, Waxed: waxed})
			lanterns = append(lanterns, CopperLantern{Hanging: true, Oxidation: o, Waxed: waxed})
		}
	}
	f(true)
	f(false)
	return
}
