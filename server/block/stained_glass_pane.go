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

type StainedGlassPane struct {
	transparent
	thin
	clicksAndSticks
	sourceWaterDisplacer

	Colour item.Colour
}

func (p StainedGlassPane) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (p StainedGlassPane) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(p))
}

func (p StainedGlassPane) EncodeItem() (name string, meta int16) {
	return "minecraft:" + p.Colour.String() + "_stained_glass_pane", 0
}

func (p StainedGlassPane) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + p.Colour.String() + "_stained_glass_pane", nil
}

func allStainedGlassPane() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, StainedGlassPane{Colour: c})
	}
	return b
}
