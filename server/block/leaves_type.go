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

type LeavesType struct {
	leavesType
}

func OakLeaves() LeavesType {
	return LeavesType{0}
}

func SpruceLeaves() LeavesType {
	return LeavesType{1}
}

func BirchLeaves() LeavesType {
	return LeavesType{2}
}

func JungleLeaves() LeavesType {
	return LeavesType{3}
}

func AcaciaLeaves() LeavesType {
	return LeavesType{4}
}

func DarkOakLeaves() LeavesType {
	return LeavesType{5}
}

func MangroveLeaves() LeavesType {
	return LeavesType{6}
}

func CherryLeaves() LeavesType {
	return LeavesType{7}
}

func PaleOakLeaves() LeavesType {
	return LeavesType{8}
}

func AzaleaLeaves() LeavesType {
	return LeavesType{9}
}

func FloweringAzaleaLeaves() LeavesType {
	return LeavesType{10}
}

func LeavesTypes() []LeavesType {
	return []LeavesType{
		OakLeaves(),
		SpruceLeaves(),
		BirchLeaves(),
		JungleLeaves(),
		AcaciaLeaves(),
		DarkOakLeaves(),
		MangroveLeaves(),
		CherryLeaves(),
		PaleOakLeaves(),
		AzaleaLeaves(),
		FloweringAzaleaLeaves(),
	}
}

func WoodLeavesTypes() []LeavesType {
	return []LeavesType{
		OakLeaves(),
		SpruceLeaves(),
		BirchLeaves(),
		JungleLeaves(),
		AcaciaLeaves(),
		DarkOakLeaves(),
		MangroveLeaves(),
		CherryLeaves(),
		PaleOakLeaves(),
	}
}

type leavesType uint8

func (t leavesType) Uint8() uint8 {
	return uint8(t)
}

func (t leavesType) String() string {
	if wood, ok := t.Wood(); ok {
		return wood.String() + "_leaves"
	}
	switch t {
	case 9:
		return "azalea_leaves"
	case 10:
		return "azalea_leaves_flowered"
	}
	panic("unknown leaves type")
}

func (t leavesType) Wood() (WoodType, bool) {
	switch t {
	case 0:
		return OakWood(), true
	case 1:
		return SpruceWood(), true
	case 2:
		return BirchWood(), true
	case 3:
		return JungleWood(), true
	case 4:
		return AcaciaWood(), true
	case 5:
		return DarkOakWood(), true
	case 6:
		return MangroveWood(), true
	case 7:
		return CherryWood(), true
	case 8:
		return PaleOakWood(), true
	default:
		return WoodType{}, false
	}
}
