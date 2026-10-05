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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
)

type Copper struct {
	solid
	bassDrum

	Type CopperType

	Oxidation OxidationType

	Waxed bool
}

func (c Copper) Strip() (world.Block, world.Sound, bool) {
	if c.Waxed {
		c.Waxed = false
		return c, sound.WaxRemoved{}, true
	} else if ot, ok := c.Oxidation.Decrease(); ok {
		c.Oxidation = ot
		return c, sound.CopperScraped{}, true
	}
	return c, nil, false
}

func (c Copper) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(c)).withBlastResistance(6)
}

func (c Copper) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	before := c.Waxed
	c.Waxed = true
	return c, !before
}

func (c Copper) CanOxidate() bool {
	return !c.Waxed
}

func (c Copper) OxidationLevel() OxidationType {
	return c.Oxidation
}

func (c Copper) WithOxidationLevel(o OxidationType) Oxidisable {
	c.Oxidation = o
	return c
}

func (c Copper) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, c)
}

func (c Copper) EncodeItem() (name string, meta int16) {
	if c.Type == NormalCopper() && c.Oxidation == UnoxidisedOxidation() && !c.Waxed {
		return "minecraft:copper_block", 0
	}
	name = "copper"
	if c.Type != NormalCopper() {
		name = c.Type.String() + "_" + name
	}
	if c.Oxidation != UnoxidisedOxidation() {
		name = c.Oxidation.String() + "_" + name
	}
	if c.Waxed {
		name = "waxed_" + name
	}
	return "minecraft:" + name, 0
}

func (c Copper) EncodeBlock() (string, map[string]any) {
	if c.Type == NormalCopper() && c.Oxidation == UnoxidisedOxidation() && !c.Waxed {
		return "minecraft:copper_block", nil
	}
	name := "copper"
	if c.Type != NormalCopper() {
		name = c.Type.String() + "_" + name
	}
	if c.Oxidation != UnoxidisedOxidation() {
		name = c.Oxidation.String() + "_" + name
	}
	if c.Waxed {
		name = "waxed_" + name
	}
	return "minecraft:" + name, nil
}

func allCopper() (c []world.Block) {
	f := func(waxed bool) {
		for _, t := range CopperTypes() {
			for _, o := range OxidationTypes() {
				c = append(c, Copper{Type: t, Oxidation: o, Waxed: waxed})
			}
		}
	}
	f(true)
	f(false)
	return
}
