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

type Thin struct{}

const (
	thinHeight = 1
	thinInset  = 7.0 / 16.0
)

func (t Thin) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {
	boxes := make([]cube.BBox, 0, 2)

	connectWest, connectEast := t.checkConnection(pos, cube.FaceWest, s), t.checkConnection(pos, cube.FaceEast, s)
	if connectWest || connectEast {
		box := cube.Box(0, 0, 0, 1, thinHeight, 1).Stretch(cube.Z, -thinInset)
		if !connectWest {
			box = box.ExtendTowards(cube.FaceWest, -thinInset)
		} else if !connectEast {
			box = box.ExtendTowards(cube.FaceEast, -thinInset)
		}
		boxes = append(boxes, box)
	}

	connectNorth, connectSouth := t.checkConnection(pos, cube.FaceNorth, s), t.checkConnection(pos, cube.FaceSouth, s)
	if connectNorth || connectSouth {
		box := cube.Box(0, 0, 0, 1, thinHeight, 1).Stretch(cube.X, -thinInset)
		if !connectNorth {
			box = box.ExtendTowards(cube.FaceNorth, -thinInset)
		} else if !connectSouth {
			box = box.ExtendTowards(cube.FaceSouth, -thinInset)
		}
		boxes = append(boxes, box)
	}

	if len(boxes) == 0 {
		boxes = append(boxes, cube.Box(0, 0, 0, 1, thinHeight, 1).Stretch(cube.X, -thinInset).Stretch(cube.Z, -thinInset))
	}

	return boxes
}

func (t Thin) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face == cube.FaceDown
}

func (t Thin) checkConnection(pos cube.Pos, face cube.Face, s world.BlockSource) bool {
	sidePos := pos.Side(face)
	sideBlock := s.Block(sidePos)
	_, isThin := sideBlock.Model().(Thin)
	_, isWall := sideBlock.Model().(Wall)
	return isThin || isWall || sideBlock.Model().FaceSolid(sidePos, face.Opposite(), s)
}
