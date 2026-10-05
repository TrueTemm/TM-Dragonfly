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

package block

type DoubleTallGrassType struct {
	doubleTallGrass
}

func NormalDoubleTallGrass() DoubleTallGrassType {
	return DoubleTallGrassType{0}
}

func FernDoubleTallGrass() DoubleTallGrassType {
	return DoubleTallGrassType{1}
}

func DoubleTallGrassTypes() []DoubleTallGrassType {
	return []DoubleTallGrassType{NormalDoubleTallGrass(), FernDoubleTallGrass()}
}

type doubleTallGrass uint8

func (t doubleTallGrass) Uint8() uint8 {
	return uint8(t)
}

func (t doubleTallGrass) Name() string {
	switch t {
	case 0:
		return "Tall Grass"
	case 1:
		return "Large Fern"
	}
	panic("unknown double tall grass type")
}

func (t doubleTallGrass) String() string {
	switch t {
	case 0:
		return "tall_grass"
	case 1:
		return "large_fern"
	}
	panic("unknown double tall grass type")
}
