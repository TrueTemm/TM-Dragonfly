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
	"github.com/df-mc/dragonfly/server/world"
)

type Dirt struct {
	solid

	Coarse bool
}

func (d Dirt) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, DeadBush:
		return !d.Coarse
	case Flower, DoubleFlower, NetherSprouts, PinkPetals, SugarCane, BambooSapling, Bamboo:
		return true
	}
	return false
}

func (d Dirt) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(d))
}

func (d Dirt) Till() (world.Block, bool) {
	if d.Coarse {
		return Dirt{Coarse: false}, true
	}
	return Farmland{}, true
}

func (d Dirt) Shovel() (world.Block, bool) {
	return DirtPath{}, true
}

func (d Dirt) EncodeItem() (name string, meta int16) {
	if d.Coarse {
		return "minecraft:coarse_dirt", 0
	}
	return "minecraft:dirt", 0
}

func (d Dirt) EncodeBlock() (string, map[string]any) {
	if d.Coarse {
		return "minecraft:coarse_dirt", nil
	}
	return "minecraft:dirt", nil
}

func supportsVegetation(vegetation, block world.Block) bool {
	soil, ok := block.(Soil)
	return ok && soil.SoilFor(vegetation)
}

type Soil interface {
	SoilFor(world.Block) bool
}
