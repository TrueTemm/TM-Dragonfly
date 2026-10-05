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

type Wall struct {
	NorthConnection float64

	EastConnection float64

	SouthConnection float64

	WestConnection float64

	Post bool
}

func (w Wall) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	postHeight := 0.8125
	if w.Post {
		postHeight = 1
	}
	boxes := []cube.BBox{cube.Box(0.25, 0, 0.25, 0.75, postHeight, 0.75)}
	if w.NorthConnection > 0 {
		boxes = append(boxes, cube.Box(0.25, 0, 0, 0.75, w.NorthConnection, 0.25))
	}
	if w.EastConnection > 0 {
		boxes = append(boxes, cube.Box(0.75, 0, 0.25, 1, w.EastConnection, 0.75))
	}
	if w.SouthConnection > 0 {
		boxes = append(boxes, cube.Box(0.25, 0, 0.75, 0.75, w.SouthConnection, 1))
	}
	if w.WestConnection > 0 {
		boxes = append(boxes, cube.Box(0, 0, 0.25, 0.25, w.WestConnection, 0.75))
	}
	return boxes
}

func (w Wall) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face.Axis() == cube.Y
}
