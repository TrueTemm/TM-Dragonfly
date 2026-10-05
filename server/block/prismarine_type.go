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

type PrismarineType struct {
	prismarine
}

type prismarine uint8

func NormalPrismarine() PrismarineType {
	return PrismarineType{0}
}

func DarkPrismarine() PrismarineType {
	return PrismarineType{1}
}

func BrickPrismarine() PrismarineType {
	return PrismarineType{2}
}

func (s prismarine) Uint8() uint8 {
	return uint8(s)
}

func (s prismarine) Name() string {
	switch s {
	case 0:
		return "Prismarine"
	case 1:
		return "Dark Prismarine"
	case 2:
		return "Prismarine Bricks"
	}
	panic("unknown prismarine type")
}

func (s prismarine) String() string {
	switch s {
	case 0:
		return "prismarine"
	case 1:
		return "dark_prismarine"
	case 2:
		return "prismarine_bricks"
	}
	panic("unknown prismarine type")
}

func PrismarineTypes() []PrismarineType {
	return []PrismarineType{NormalPrismarine(), DarkPrismarine(), BrickPrismarine()}
}
