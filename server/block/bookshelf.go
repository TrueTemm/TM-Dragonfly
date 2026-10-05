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
	"time"
)

type Bookshelf struct {
	solid
	bass
}

func (b Bookshelf) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, alwaysHarvestable, axeEffective, silkTouchDrop(item.NewStack(item.Book{}, 3), item.NewStack(b, 1)))
}

func (Bookshelf) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 20, true)
}

func (Bookshelf) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (Bookshelf) EncodeItem() (name string, meta int16) {
	return "minecraft:bookshelf", 0
}

func (Bookshelf) EncodeBlock() (string, map[string]any) {
	return "minecraft:bookshelf", nil
}
