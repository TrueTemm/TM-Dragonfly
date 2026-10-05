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

type Prismarine struct {
	solid
	bassDrum

	Type PrismarineType
}

func (p Prismarine) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(p)).withBlastResistance(6)
}

func (p Prismarine) EncodeItem() (id string, meta int16) {
	return "minecraft:" + p.Type.String(), 0
}

func (p Prismarine) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + p.Type.String(), nil
}

func allPrismarine() (c []world.Block) {
	for _, t := range PrismarineTypes() {
		c = append(c, Prismarine{Type: t})
	}
	return
}
