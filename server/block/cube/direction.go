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

type Direction int

const (
	North Direction = iota

	South

	West

	East
)

func (d Direction) Face() Face {
	return Face(d + 2)
}

func (d Direction) Opposite() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case West:
		return East
	case East:
		return West
	}
	panic("invalid direction")
}

func (d Direction) RotateRight() Direction {
	switch d {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	}
	panic("invalid direction")
}

func (d Direction) RotateLeft() Direction {
	switch d {
	case North:
		return West
	case East:
		return North
	case South:
		return East
	case West:
		return South
	}
	panic("invalid direction")
}

func (d Direction) String() string {
	switch d {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	}
	panic("invalid direction")
}

var directions = [...]Direction{North, East, South, West}

func Directions() []Direction {
	return directions[:]
}
