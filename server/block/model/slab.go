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

type Slab struct {
	Double, Top bool
}

func (s Slab) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	if s.Double {
		return []cube.BBox{full}
	}
	if s.Top {
		return []cube.BBox{cube.Box(0, 0.5, 0, 1, 1, 1)}
	}
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.5, 1)}
}

func (s Slab) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	if s.Double {
		return true
	} else if s.Top {
		return face == cube.FaceUp
	}
	return face == cube.FaceDown
}
