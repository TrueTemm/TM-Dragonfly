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

type OreType struct {
	ore
}

func StoneOre() OreType {
	return OreType{0}
}

func DeepslateOre() OreType {
	return OreType{1}
}

func OreTypes() []OreType {
	return []OreType{StoneOre(), DeepslateOre()}
}

type ore uint8

func (o ore) Uint8() uint8 {
	return uint8(o)
}

func (o ore) Name() string {
	switch o {
	case 0:
		return "Stone"
	case 1:
		return "Deepslate"
	}
	panic("unknown ore type")
}

func (o ore) String() string {
	switch o {
	case 0:
		return "stone"
	case 1:
		return "deepslate"
	}
	panic("unknown ore type")
}

func (o ore) Prefix() string {
	switch o {
	case 0:
		return ""
	case 1:
		return "deepslate_"
	}
	panic("unknown ore type")
}

func (o ore) Hardness() float64 {
	switch o {
	case 0:
		return 3
	case 1:
		return 4.5
	}
	panic("unknown ore type")
}
