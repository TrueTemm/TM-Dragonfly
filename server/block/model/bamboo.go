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

type Bamboo struct {
	Thick bool
}

func (b Bamboo) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {

	size := 0.5 + 2.0/16.0
	if b.Thick {
		size = 0.5 + 3.0/16.0
	}
	offset := randomOffset(pos, -0.25, 0.25, 16)
	return []cube.BBox{cube.Box(0.5, 0, 0.5, size, 1, size).Translate(offset)}
}

func (b Bamboo) FaceSolid(pos cube.Pos, face cube.Face, s world.BlockSource) bool {
	return false
}
