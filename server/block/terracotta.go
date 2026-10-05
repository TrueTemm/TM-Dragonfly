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

type Terracotta struct {
	solid
	bassDrum
}

func (Terracotta) SoilFor(block world.Block) bool {
	_, ok := block.(DeadBush)
	return ok
}

func (t Terracotta) BreakInfo() BreakInfo {
	return newBreakInfo(1.25, pickaxeHarvestable, pickaxeEffective, oneOf(t)).withBlastResistance(4.2)
}

func (Terracotta) EncodeItem() (name string, meta int16) {
	return "minecraft:hardened_clay", meta
}

func (Terracotta) EncodeBlock() (string, map[string]any) {
	return "minecraft:hardened_clay", nil
}
