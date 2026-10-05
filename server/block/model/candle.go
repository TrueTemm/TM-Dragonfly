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

package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Candle struct {
	Count int
}

func (c Candle) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	switch c.Count {
	case 2:
		return []cube.BBox{cube.Box(0.3125, 0, 0.4375, 0.6875, 0.375, 0.625)}
	case 3:
		return []cube.BBox{cube.Box(0.3125, 0, 0.375, 0.625, 0.375, 0.6875)}
	case 4:
		return []cube.BBox{cube.Box(0.3125, 0, 0.3125, 0.6875, 0.375, 0.625)}
	default:
		return []cube.BBox{cube.Box(0.4375, 0, 0.4375, 0.5625, 0.375, 0.5625)}
	}
}

func (Candle) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
