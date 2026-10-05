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

package cube

import "github.com/go-gl/mathgl/mgl64"

type Axis int

const (
	Y Axis = iota

	Z

	X
)

func (a Axis) String() string {
	switch a {
	case X:
		return "x"
	case Y:
		return "y"
	default:
		return "z"
	}
}

func (a Axis) RotateLeft() Axis {
	switch a {
	case X:
		return Z
	case Z:
		return X
	default:
		return 0
	}
}

func (a Axis) RotateRight() Axis {

	return a.RotateLeft()
}

func (a Axis) Faces() (negative, positive Face) {
	switch a {
	case X:
		return FaceWest, FaceEast
	case Y:
		return FaceDown, FaceUp
	default:
		return FaceNorth, FaceSouth
	}
}

func (a Axis) Vec3() mgl64.Vec3 {
	switch a {
	case X:
		return mgl64.Vec3{1, 0, 0}
	case Y:
		return mgl64.Vec3{0, 1, 0}
	default:
		return mgl64.Vec3{0, 0, 1}
	}
}

func Axes() []Axis {
	return []Axis{X, Y, Z}
}
