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

type InfestedStoneBricks struct {
	solid
	flute

	Type StoneBricksType
}

func (i InfestedStoneBricks) BreakInfo() BreakInfo {
	return newBreakInfo(0.75, pickaxeHarvestable, pickaxeEffective, silkTouchOnlyDrop(i)).withBlastResistance(0.75)
}

func (i InfestedStoneBricks) EncodeItem() (name string, meta int16) {
	return "minecraft:infested_" + i.Type.String(), 0
}

func (i InfestedStoneBricks) EncodeBlock() (string, map[string]any) {
	return "minecraft:infested_" + i.Type.String(), nil
}

func allInfestedStoneBricks() (s []world.Block) {
	for _, t := range StoneBricksTypes() {
		s = append(s, InfestedStoneBricks{Type: t})
	}
	return
}
