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

type Fence struct {
	Wood bool
}

const (
	fenceHeight = 1.5
	fenceInset  = 0.375
)

func (f Fence) BBox(pos cube.Pos, s world.BlockSource) []cube.BBox {
	boxes := make([]cube.BBox, 0, 2)

	connectWest, connectEast := f.checkConnection(pos, cube.FaceWest, s), f.checkConnection(pos, cube.FaceEast, s)
	connectNorth, connectSouth := f.checkConnection(pos, cube.FaceNorth, s), f.checkConnection(pos, cube.FaceSouth, s)

	if connectWest || connectEast {
		sideBox := cube.Box(0, 0, 0, 1, fenceHeight, 1).Stretch(cube.Z, -fenceInset)
		if connectWest {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceEast, -fenceInset))
		}
		if connectEast {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceWest, -fenceInset))
		}
	}

	if connectNorth || connectSouth {
		sideBox := cube.Box(0, 0, 0, 1, fenceHeight, 1).Stretch(cube.X, -fenceInset)
		if connectNorth {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceSouth, -fenceInset))
		}
		if connectSouth {
			boxes = append(boxes, sideBox.ExtendTowards(cube.FaceNorth, -fenceInset))
		}
	}

	if len(boxes) == 0 {
		boxes = append(boxes, cube.Box(fenceInset, 0, fenceInset, 1-fenceInset, fenceHeight, 1-fenceInset))
	}

	return boxes
}

func (f Fence) FaceSolid(_ cube.Pos, face cube.Face, _ world.BlockSource) bool {
	return face == cube.FaceDown || face == cube.FaceUp
}

func (f Fence) checkConnection(pos cube.Pos, face cube.Face, src world.BlockSource) bool {
	sidePos := pos.Side(face)
	sideBlock := src.Block(sidePos)
	if fence, ok := sideBlock.Model().(Fence); ok && fence.Wood == f.Wood {
		return true
	}
	if sideBlock.Model().FaceSolid(sidePos, face.Opposite(), src) {
		return true
	}
	_, ok := sideBlock.Model().(FenceGate)
	return ok
}
