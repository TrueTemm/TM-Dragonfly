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

import "github.com/df-mc/dragonfly/server/item"

type Snow struct {
	solid
}

func (s Snow) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, shovelEffective, silkTouchDrop(item.NewStack(item.Snowball{}, 4), item.NewStack(s, 1)))
}

func (Snow) EncodeItem() (name string, meta int16) {
	return "minecraft:snow", 0
}

func (Snow) EncodeBlock() (string, map[string]any) {
	return "minecraft:snow", nil
}
