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
	"github.com/df-mc/dragonfly/server/world"
)

type GlassPane struct {
	transparent
	thin
	clicksAndSticks
	sourceWaterDisplacer
}

func (p GlassPane) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (p GlassPane) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(p))
}

func (GlassPane) EncodeItem() (name string, meta int16) {
	return "minecraft:glass_pane", meta
}

func (GlassPane) EncodeBlock() (string, map[string]any) {
	return "minecraft:glass_pane", nil
}
