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

package item

type BannerPatternType struct {
	bannerPatternType
}

func CreeperBannerPattern() BannerPatternType {
	return BannerPatternType{0}
}

func SkullBannerPattern() BannerPatternType {
	return BannerPatternType{1}
}

func FlowerBannerPattern() BannerPatternType {
	return BannerPatternType{2}
}

func MojangBannerPattern() BannerPatternType {
	return BannerPatternType{3}
}

func FieldMasonedBannerPattern() BannerPatternType {
	return BannerPatternType{4}
}

func BordureIndentedBannerPattern() BannerPatternType {
	return BannerPatternType{5}
}

func PiglinBannerPattern() BannerPatternType {
	return BannerPatternType{6}
}

func GlobeBannerPattern() BannerPatternType {
	return BannerPatternType{7}
}

func FlowBannerPattern() BannerPatternType {
	return BannerPatternType{8}
}

func GusterBannerPattern() BannerPatternType {
	return BannerPatternType{9}
}

func BannerPatterns() []BannerPatternType {
	return []BannerPatternType{
		CreeperBannerPattern(),
		SkullBannerPattern(),
		FlowerBannerPattern(),
		MojangBannerPattern(),
		FieldMasonedBannerPattern(),
		BordureIndentedBannerPattern(),
		PiglinBannerPattern(),
		GlobeBannerPattern(),
		FlowBannerPattern(),
		GusterBannerPattern(),
	}
}

type bannerPatternType uint8

func (b bannerPatternType) Uint8() uint8 {
	return uint8(b)
}

func (b bannerPatternType) String() string {
	switch b {
	case 0:
		return "creeper"
	case 1:
		return "skull"
	case 2:
		return "flower"
	case 3:
		return "mojang"
	case 4:
		return "field_masoned"
	case 5:
		return "bordure_indented"
	case 6:
		return "piglin"
	case 7:
		return "globe"
	case 8:
		return "flow"
	case 9:
		return "guster"
	}
	panic("should never happen")
}
