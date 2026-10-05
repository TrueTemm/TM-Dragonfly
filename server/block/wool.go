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
	"github.com/df-mc/dragonfly/server/world/sound"
)

type Wool struct {
	solid

	Colour item.Colour
}

func (w Wool) Instrument() sound.Instrument {
	return sound.Guitar()
}

func (w Wool) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 60, true)
}

func (w Wool) BreakInfo() BreakInfo {
	return newBreakInfo(0.8, alwaysHarvestable, shearsEffective, oneOf(w))
}

func (w Wool) EncodeItem() (name string, meta int16) {
	return "minecraft:" + w.Colour.String() + "_wool", 0
}

func (w Wool) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + w.Colour.String() + "_wool", nil
}

func allWool() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, Wool{Colour: c})
	}
	return b
}
