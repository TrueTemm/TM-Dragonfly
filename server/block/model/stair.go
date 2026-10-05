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
	"github.com/go-gl/mathgl/mgl64"
)

type Stair struct {
	Facing cube.Direction

	UpsideDown bool
}

func (s Stair) BBox(pos cube.Pos, bs world.BlockSource) []cube.BBox {
	b := []cube.BBox{cube.Box(0, 0, 0, 1, 0.5, 1)}
	if s.UpsideDown {
		b[0] = cube.Box(0, 0.5, 0, 1, 1, 1)
	}
	t := s.cornerType(pos, bs)

	face, oppositeFace := s.Facing.Face(), s.Facing.Opposite().Face()
	switch t {
	case noCorner, cornerRightInner, cornerLeftInner:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(face, 0.5).
			Stretch(s.Facing.RotateRight().Face().Axis(), 0.5))
	}
	switch t {
	case cornerRightOuter:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(face, 0.5).
			ExtendTowards(s.Facing.RotateLeft().Face(), 0.5))
	case cornerLeftOuter:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(face, 0.5).
			ExtendTowards(s.Facing.RotateRight().Face(), 0.5))
	case cornerRightInner:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(oppositeFace, 0.5).
			ExtendTowards(s.Facing.RotateRight().Face(), 0.5))
	case cornerLeftInner:
		b = append(b, cube.Box(0.5, 0.5, 0.5, 0.5, 1, 0.5).
			ExtendTowards(oppositeFace, 0.5).
			ExtendTowards(s.Facing.RotateLeft().Face(), 0.5))
	}
	if s.UpsideDown {
		for i := range b[1:] {
			b[i+1] = b[i+1].Translate(mgl64.Vec3{0, -0.5})
		}
	}
	return b
}

func (s Stair) FaceSolid(pos cube.Pos, face cube.Face, bs world.BlockSource) bool {

	if (face == cube.FaceUp && s.UpsideDown) || (face == cube.FaceDown && !s.UpsideDown) {
		return true
	}

	switch t := s.cornerType(pos, bs); t {
	case cornerRightOuter, cornerLeftOuter:

		return false
	case noCorner:

		return s.Facing.Face() == face
	case cornerRightInner:
		return face == s.Facing.RotateRight().Face() || face == s.Facing.Face()
	default:
		return face == s.Facing.RotateLeft().Face() || face == s.Facing.Face()
	}
}

const (
	noCorner = iota
	cornerRightInner
	cornerLeftInner
	cornerRightOuter
	cornerLeftOuter
)

func (s Stair) cornerType(pos cube.Pos, bs world.BlockSource) uint8 {
	rotatedFacing := s.Facing.RotateRight()
	if closedSide, ok := bs.Block(pos.Side(s.Facing.Face())).Model().(Stair); ok && closedSide.UpsideDown == s.UpsideDown {
		if closedSide.Facing == rotatedFacing {
			return cornerLeftOuter
		} else if closedSide.Facing == rotatedFacing.Opposite() {

			if side, ok := bs.Block(pos.Side(s.Facing.RotateRight().Face())).Model().(Stair); !ok || side.Facing != s.Facing || side.UpsideDown != s.UpsideDown {
				return cornerRightOuter
			}
			return noCorner
		}
	}
	if openSide, ok := bs.Block(pos.Side(s.Facing.Opposite().Face())).Model().(Stair); ok && openSide.UpsideDown == s.UpsideDown {
		if openSide.Facing == rotatedFacing {

			if side, ok := bs.Block(pos.Side(s.Facing.RotateRight().Face())).Model().(Stair); !ok || side.Facing != s.Facing || side.UpsideDown != s.UpsideDown {
				return cornerRightInner
			}
		} else if openSide.Facing == rotatedFacing.Opposite() {
			return cornerLeftInner
		}
	}
	return noCorner
}
