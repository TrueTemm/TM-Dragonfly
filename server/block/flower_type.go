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

type FlowerType struct {
	flower
}

type flower uint8

func Dandelion() FlowerType {
	return FlowerType{0}
}

func Poppy() FlowerType {
	return FlowerType{1}
}

func BlueOrchid() FlowerType {
	return FlowerType{2}
}

func Allium() FlowerType {
	return FlowerType{3}
}

func AzureBluet() FlowerType {
	return FlowerType{4}
}

func RedTulip() FlowerType {
	return FlowerType{5}
}

func OrangeTulip() FlowerType {
	return FlowerType{6}
}

func WhiteTulip() FlowerType {
	return FlowerType{7}
}

func PinkTulip() FlowerType {
	return FlowerType{8}
}

func OxeyeDaisy() FlowerType {
	return FlowerType{9}
}

func Cornflower() FlowerType {
	return FlowerType{10}
}

func LilyOfTheValley() FlowerType {
	return FlowerType{11}
}

func WitherRose() FlowerType {
	return FlowerType{12}
}

func (f flower) Uint8() uint8 {
	return uint8(f)
}

func (f flower) Name() string {
	switch f {
	case 0:
		return "Dandelion"
	case 1:
		return "Poppy"
	case 2:
		return "Blue Orchid"
	case 3:
		return "Allium"
	case 4:
		return "Azure Bluet"
	case 5:
		return "Red Tulip"
	case 6:
		return "Orange Tulip"
	case 7:
		return "White Tulip"
	case 8:
		return "Pink Tulip"
	case 9:
		return "Oxeye Daisy"
	case 10:
		return "Cornflower"
	case 11:
		return "Lily of the Valley"
	case 12:
		return "Wither Rose"
	}
	panic("unknown flower type")
}

func (f flower) String() string {
	switch f {
	case 0:
		return "dandelion"
	case 1:
		return "poppy"
	case 2:
		return "blue_orchid"
	case 3:
		return "allium"
	case 4:
		return "azure_bluet"
	case 5:
		return "red_tulip"
	case 6:
		return "orange_tulip"
	case 7:
		return "white_tulip"
	case 8:
		return "pink_tulip"
	case 9:
		return "oxeye_daisy"
	case 10:
		return "cornflower"
	case 11:
		return "lily_of_the_valley"
	case 12:
		return "wither_rose"
	}
	panic("unknown flower type")
}

func FlowerTypes() []FlowerType {
	return []FlowerType{Dandelion(), Poppy(), BlueOrchid(), Allium(), AzureBluet(), RedTulip(), OrangeTulip(), WhiteTulip(), PinkTulip(), OxeyeDaisy(), Cornflower(), LilyOfTheValley(), WitherRose()}
}
