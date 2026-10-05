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

const (
	FaceDown Face = iota

	FaceUp

	FaceNorth

	FaceSouth

	FaceWest

	FaceEast
)

type Face int

func (f Face) Direction() Direction {
	return Direction(f - 2)
}

func (f Face) Opposite() Face {
	switch f {
	default:
		return FaceUp
	case FaceUp:
		return FaceDown
	case FaceNorth:
		return FaceSouth
	case FaceSouth:
		return FaceNorth
	case FaceWest:
		return FaceEast
	case FaceEast:
		return FaceWest
	}
}

func (f Face) Axis() Axis {
	switch f {
	default:
		return Y
	case FaceEast, FaceWest:
		return X
	case FaceNorth, FaceSouth:
		return Z
	}
}

func (f Face) RotateRight() Face {
	switch f {
	case FaceNorth:
		return FaceEast
	case FaceEast:
		return FaceSouth
	case FaceSouth:
		return FaceWest
	case FaceWest:
		return FaceNorth
	}
	return f
}

func (f Face) RotateLeft() Face {
	switch f {
	case FaceNorth:
		return FaceWest
	case FaceEast:
		return FaceNorth
	case FaceSouth:
		return FaceEast
	case FaceWest:
		return FaceSouth
	}
	return f
}

func (f Face) String() string {
	switch f {
	case FaceUp:
		return "up"
	case FaceDown:
		return "down"
	case FaceNorth:
		return "north"
	case FaceSouth:
		return "south"
	case FaceWest:
		return "west"
	case FaceEast:
		return "east"
	}
	panic("invalid face")
}

func Faces() []Face {
	return faces[:]
}

func HorizontalFaces() []Face {
	return hFaces[:]
}

var hFaces = [...]Face{FaceNorth, FaceEast, FaceSouth, FaceWest}

var faces = [...]Face{FaceDown, FaceUp, FaceNorth, FaceEast, FaceSouth, FaceWest}
