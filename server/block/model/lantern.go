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

type Lantern struct {
	Hanging bool
}

func (l Lantern) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if l.Hanging {
		return []cube.BBox{cube.Box(0.3125, 0.125, 0.3125, 0.6875, 0.625, 0.6875)}
	}
	return []cube.BBox{cube.Box(0.3125, 0, 0.3125, 0.6875, 0.5, 0.6875)}
}

func (l Lantern) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
