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

type CocoaBean struct {
	Facing cube.Direction

	Age int
}

func (c CocoaBean) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{full.
		Stretch(c.Facing.RotateRight().Face().Axis(), -(6-float64(c.Age))/16).
		ExtendTowards(cube.FaceUp, -0.25).
		ExtendTowards(cube.FaceDown, -((7-float64(c.Age)*2)/16)).
		ExtendTowards(c.Facing.Face(), -0.0625).
		ExtendTowards(c.Facing.Opposite().Face(), -((11 - float64(c.Age)*2) / 16))}
}

func (c CocoaBean) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
