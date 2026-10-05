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

type StoneBricks struct {
	solid
	bassDrum

	Type StoneBricksType
}

func (s StoneBricks) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(s)).withBlastResistance(6)
}

func (s StoneBricks) SmeltInfo() item.SmeltInfo {
	if s.Type == NormalStoneBricks() {
		return newSmeltInfo(item.NewStack(StoneBricks{Type: CrackedStoneBricks()}, 1), 0.1)
	}
	return item.SmeltInfo{}
}

func (s StoneBricks) EncodeItem() (name string, meta int16) {
	return "minecraft:" + s.Type.String(), 0
}

func (s StoneBricks) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + s.Type.String(), nil
}

func allStoneBricks() (s []world.Block) {
	for _, t := range StoneBricksTypes() {
		s = append(s, StoneBricks{Type: t})
	}
	return
}
