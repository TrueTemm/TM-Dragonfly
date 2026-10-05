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

type Skull struct {
	Direction cube.Face

	Hanging bool
}

func (s Skull) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	box := cube.Box(0.25, 0, 0.25, 0.75, 0.5, 0.75)
	if !s.Hanging {
		return []cube.BBox{box}
	}
	return []cube.BBox{box.TranslateTowards(s.Direction.Opposite(), 0.25).TranslateTowards(cube.FaceUp, 0.25)}
}

func (Skull) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
