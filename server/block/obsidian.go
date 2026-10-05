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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Obsidian struct {
	solid
	bassDrum

	Crying bool
}

func (o Obsidian) LightEmissionLevel() uint8 {
	if o.Crying {
		return 10
	}
	return 0
}

func (o Obsidian) EncodeItem() (name string, meta int16) {
	if o.Crying {
		return "minecraft:crying_obsidian", 0
	}
	return "minecraft:obsidian", 0
}

func (o Obsidian) EncodeBlock() (string, map[string]any) {
	if o.Crying {
		return "minecraft:crying_obsidian", nil
	}
	return "minecraft:obsidian", nil
}

func (o Obsidian) Frame(dimension world.Dimension) bool {
	return dimension == world.Nether && !o.Crying
}

func (o Obsidian) SupportsEndCrystal() bool {
	return !o.Crying
}

func (o Obsidian) BreakInfo() BreakInfo {
	return newBreakInfo(35, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierDiamond.HarvestLevel
	}, pickaxeEffective, oneOf(o)).withBlastResistance(1200)
}
