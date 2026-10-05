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

type Trapdoor struct {
	Facing cube.Direction

	Open, Top bool
}

func (t Trapdoor) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if t.Open {
		return []cube.BBox{full.ExtendTowards(t.Facing.Face(), -0.8125)}
	} else if t.Top {
		return []cube.BBox{cube.Box(0, 0.8125, 0, 1, 1, 1)}
	}
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.1875, 1)}
}

func (t Trapdoor) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	if t.Open {
		return t.Facing.Face().Opposite() == face
	} else if t.Top {
		return face == cube.FaceUp
	}
	return face == cube.FaceDown
}
