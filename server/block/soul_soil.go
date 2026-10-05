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

type SoulSoil struct {
	solid
}

func (s SoulSoil) SoilFor(block world.Block) bool {
	_, ok := block.(NetherSprouts)
	return ok
}

func (s SoulSoil) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(s))
}

func (SoulSoil) EncodeItem() (name string, meta int16) {
	return "minecraft:soul_soil", 0
}

func (SoulSoil) EncodeBlock() (string, map[string]any) {
	return "minecraft:soul_soil", nil
}
