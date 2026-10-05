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

type CopperChain struct {
	transparent
	sourceWaterDisplacer

	Axis cube.Axis

	Oxidation OxidationType

	Waxed bool
}

func (CopperChain) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c CopperChain) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, c)
	if !used {
		return
	}
	c.Axis = face.Axis()

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (c CopperChain) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(c)).withBlastResistance(6)
}

func (c CopperChain) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if c.Waxed {
		return c, false
	}
	c.Waxed = true
	return c, true
}

func (c CopperChain) Strip() (world.Block, world.Sound, bool) {
	if c.Waxed {
		c.Waxed = false
		return c, sound.WaxRemoved{}, true
	} else if ot, ok := c.Oxidation.Decrease(); ok {
		c.Oxidation = ot
		return c, sound.CopperScraped{}, true
	}
	return c, nil, false
}

func (c CopperChain) CanOxidate() bool {
	return !c.Waxed
}

func (c CopperChain) OxidationLevel() OxidationType {
	return c.Oxidation
}

func (c CopperChain) WithOxidationLevel(o OxidationType) Oxidisable {
	c.Oxidation = o
	return c
}

func (c CopperChain) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, c)
}

func (c CopperChain) EncodeItem() (name string, meta int16) {
	return copperBlockName("copper_chain", c.Oxidation, c.Waxed), 0
}

func (c CopperChain) EncodeBlock() (name string, properties map[string]any) {
	return copperBlockName("copper_chain", c.Oxidation, c.Waxed), map[string]any{"pillar_axis": c.Axis.String()}
}

func (c CopperChain) Model() world.BlockModel {
	return model.Chain{Axis: c.Axis}
}

func allCopperChains() (chains []world.Block) {
	f := func(waxed bool) {
		for _, o := range OxidationTypes() {
			for _, axis := range cube.Axes() {
				chains = append(chains, CopperChain{Axis: axis, Oxidation: o, Waxed: waxed})
			}
		}
	}
	f(true)
	f(false)
	return
}
