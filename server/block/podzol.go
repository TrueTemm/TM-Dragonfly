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

import "github.com/df-mc/dragonfly/server/world"

type Podzol struct {
	solid
}

func (p Podzol) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, NetherSprouts, DeadBush, SugarCane, BambooSapling, Bamboo:
		return true
	}
	return false
}

func (Podzol) Shovel() (world.Block, bool) {
	return DirtPath{}, true
}

func (p Podzol) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, silkTouchOneOf(Dirt{}, p))
}

func (Podzol) EncodeItem() (name string, meta int16) {
	return "minecraft:podzol", 0
}

func (Podzol) EncodeBlock() (string, map[string]any) {
	return "minecraft:podzol", nil
}
