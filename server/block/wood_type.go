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

type WoodType struct {
	wood
}

func OakWood() WoodType {
	return WoodType{0}
}

func SpruceWood() WoodType {
	return WoodType{1}
}

func BirchWood() WoodType {
	return WoodType{2}
}

func JungleWood() WoodType {
	return WoodType{3}
}

func AcaciaWood() WoodType {
	return WoodType{4}
}

func DarkOakWood() WoodType {
	return WoodType{5}
}

func CrimsonWood() WoodType {
	return WoodType{6}
}

func WarpedWood() WoodType {
	return WoodType{7}
}

func MangroveWood() WoodType {
	return WoodType{8}
}

func CherryWood() WoodType {
	return WoodType{9}
}

func PaleOakWood() WoodType {
	return WoodType{10}
}

func BambooWood() WoodType {
	return WoodType{11}
}

func WoodTypes() []WoodType {
	return []WoodType{OakWood(), SpruceWood(), BirchWood(), JungleWood(), AcaciaWood(), DarkOakWood(), CrimsonWood(), WarpedWood(), MangroveWood(), CherryWood(), PaleOakWood(), BambooWood()}
}

type wood uint8

func (w wood) Uint8() uint8 {
	return uint8(w)
}

func (w wood) Name() string {
	switch w {
	case 0:
		return "Oak Wood"
	case 1:
		return "Spruce Wood"
	case 2:
		return "Birch Wood"
	case 3:
		return "Jungle Wood"
	case 4:
		return "Acacia Wood"
	case 5:
		return "Dark Oak Wood"
	case 6:
		return "Crimson Wood"
	case 7:
		return "Warped Wood"
	case 8:
		return "Mangrove Wood"
	case 9:
		return "Cherry Wood"
	case 10:
		return "Pale Oak Wood"
	case 11:
		return "Bamboo Wood"
	}
	panic("unknown wood type")
}

func (w wood) String() string {
	switch w {
	case 0:
		return "oak"
	case 1:
		return "spruce"
	case 2:
		return "birch"
	case 3:
		return "jungle"
	case 4:
		return "acacia"
	case 5:
		return "dark_oak"
	case 6:
		return "crimson"
	case 7:
		return "warped"
	case 8:
		return "mangrove"
	case 9:
		return "cherry"
	case 10:
		return "pale_oak"
	case 11:
		return "bamboo"
	}
	panic("unknown wood type")
}

func (w wood) Flammable() bool {
	return w != CrimsonWood().wood && w != WarpedWood().wood
}

func (w WoodType) Leaves() (LeavesType, bool) {
	switch w {
	case OakWood():
		return OakLeaves(), true
	case SpruceWood():
		return SpruceLeaves(), true
	case BirchWood():
		return BirchLeaves(), true
	case JungleWood():
		return JungleLeaves(), true
	case AcaciaWood():
		return AcaciaLeaves(), true
	case DarkOakWood():
		return DarkOakLeaves(), true
	case MangroveWood():
		return MangroveLeaves(), true
	case CherryWood():
		return CherryLeaves(), true
	case PaleOakWood():
		return PaleOakLeaves(), true
	default:
		return LeavesType{}, false
	}
}
