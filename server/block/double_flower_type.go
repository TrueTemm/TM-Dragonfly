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

type DoubleFlowerType struct {
	doubleFlower
}

type doubleFlower uint8

func Sunflower() DoubleFlowerType {
	return DoubleFlowerType{0}
}

func Lilac() DoubleFlowerType {
	return DoubleFlowerType{1}
}

func RoseBush() DoubleFlowerType {
	return DoubleFlowerType{4}
}

func Peony() DoubleFlowerType {
	return DoubleFlowerType{5}
}

func (d doubleFlower) Uint8() uint8 {
	return uint8(d)
}

func (d doubleFlower) Name() string {
	switch d {
	case 0:
		return "Sunflower"
	case 1:
		return "Lilac"
	case 4:
		return "Rose Bush"
	case 5:
		return "Peony"
	}
	panic("unknown double plant type")
}

func (d doubleFlower) String() string {
	switch d {
	case 0:
		return "sunflower"
	case 1:
		return "lilac"
	case 4:
		return "rose_bush"
	case 5:
		return "peony"
	}
	panic("unknown double plant type")
}

func DoubleFlowerTypes() []DoubleFlowerType {
	return []DoubleFlowerType{Sunflower(), Lilac(), RoseBush(), Peony()}
}
