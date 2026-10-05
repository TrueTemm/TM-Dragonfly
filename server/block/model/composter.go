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
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Composter struct {
	Level int
}

func (c Composter) BBox(_ cube.Pos, _ world.BlockSource) []cube.BBox {
	compostHeight := math.Abs(math.Min(float64(c.Level), 7)*0.125 - 0.0625)
	return []cube.BBox{
		cube.Box(0, 0, 0, 1, 1, 0.125),
		cube.Box(0, 0, 0.875, 1, 1, 1),
		cube.Box(0.875, 0, 0, 1, 1, 1),
		cube.Box(0, 0, 0, 0.125, 1, 1),
		cube.Box(0.125, 0, 0.125, 0.875, 0.125+compostHeight, 0.875),
	}
}

func (Composter) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face != cube.FaceUp
}
