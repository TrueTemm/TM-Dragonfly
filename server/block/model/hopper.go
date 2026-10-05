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

type Hopper struct{}

func (h Hopper) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	bbox := []cube.BBox{full.ExtendTowards(cube.FaceUp, -0.375)}
	for _, f := range cube.HorizontalFaces() {
		bbox = append(bbox, full.ExtendTowards(f, -0.875))
	}
	return bbox
}

func (Hopper) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face == cube.FaceUp
}
