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

type StainedGlass struct {
	transparent
	solid
	clicksAndSticks

	Colour item.Colour
}

func (g StainedGlass) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(g))
}

func (g StainedGlass) EncodeItem() (name string, meta int16) {
	return "minecraft:" + g.Colour.String() + "_stained_glass", 0
}

func (g StainedGlass) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + g.Colour.String() + "_stained_glass", nil
}

func allStainedGlass() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, StainedGlass{Colour: c})
	}
	return b
}
