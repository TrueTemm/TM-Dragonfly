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
)

type Sand struct {
	gravityAffected
	solid
	snare

	Red bool
}

func (s Sand) SoilFor(block world.Block) bool {
	switch block.(type) {
	case Cactus, DeadBush, SugarCane, BambooSapling, Bamboo:
		return true
	}
	return false
}

func (s Sand) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	s.fall(s, pos, tx)
}

func (s Sand) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(s))
}

func (Sand) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(Glass{}, 1), 0.1)
}

func (s Sand) EncodeItem() (name string, meta int16) {
	if s.Red {
		return "minecraft:red_sand", 0
	}
	return "minecraft:sand", 0
}

func (s Sand) EncodeBlock() (string, map[string]any) {
	if s.Red {
		return "minecraft:red_sand", nil
	}
	return "minecraft:sand", nil
}
