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

type StainedTerracotta struct {
	solid
	bassDrum

	Colour item.Colour
}

func (t StainedTerracotta) SoilFor(block world.Block) bool {
	_, ok := block.(DeadBush)
	return ok
}

func (t StainedTerracotta) BreakInfo() BreakInfo {
	return newBreakInfo(1.25, pickaxeHarvestable, pickaxeEffective, oneOf(t))
}

func (t StainedTerracotta) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(GlazedTerracotta{Colour: t.Colour}, 1), 0.1)
}

func (t StainedTerracotta) EncodeItem() (name string, meta int16) {
	return "minecraft:" + t.Colour.String() + "_terracotta", 0
}

func (t StainedTerracotta) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + t.Colour.String() + "_terracotta", nil
}

func allStainedTerracotta() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, StainedTerracotta{Colour: c})
	}
	return b
}
