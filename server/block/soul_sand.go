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
	"github.com/df-mc/dragonfly/server/world/sound"
)

type SoulSand struct {
	solid
}

func (s SoulSand) SoilFor(block world.Block) bool {
	flower, ok := block.(Flower)
	return ok && flower.Type == WitherRose()
}

func (s SoulSand) Instrument() sound.Instrument {
	return sound.CowBell()
}

func (s SoulSand) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(s))
}

func (SoulSand) EncodeItem() (name string, meta int16) {
	return "minecraft:soul_sand", 0
}

func (SoulSand) EncodeBlock() (string, map[string]any) {
	return "minecraft:soul_sand", nil
}
