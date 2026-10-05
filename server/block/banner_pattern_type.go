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

import "github.com/df-mc/dragonfly/server/item"

type BannerPatternType struct {
	bannerPatternType
}

func BorderBannerPattern() BannerPatternType {
	return BannerPatternType{0}
}

func BricksBannerPattern() BannerPatternType {
	return BannerPatternType{1}
}

func CircleBannerPattern() BannerPatternType {
	return BannerPatternType{2}
}

func CreeperBannerPattern() BannerPatternType {
	return BannerPatternType{3}
}

func CrossBannerPattern() BannerPatternType {
	return BannerPatternType{4}
}

func CurlyBorderBannerPattern() BannerPatternType {
	return BannerPatternType{5}
}

func DiagonalLeftBannerPattern() BannerPatternType {
	return BannerPatternType{6}
}

func DiagonalRightBannerPattern() BannerPatternType {
	return BannerPatternType{7}
}

func DiagonalUpLeftBannerPattern() BannerPatternType {
	return BannerPatternType{8}
}

func DiagonalUpRightBannerPattern() BannerPatternType {
	return BannerPatternType{9}
}

func FlowerBannerPattern() BannerPatternType {
	return BannerPatternType{10}
}

func GradientBannerPattern() BannerPatternType {
	return BannerPatternType{11}
}

func GradientUpBannerPattern() BannerPatternType {
	return BannerPatternType{12}
}

func HalfHorizontalBannerPattern() BannerPatternType {
	return BannerPatternType{13}
}

func HalfHorizontalBottomBannerPattern() BannerPatternType {
	return BannerPatternType{14}
}

func HalfVerticalBannerPattern() BannerPatternType {
	return BannerPatternType{15}
}

func HalfVerticalRightBannerPattern() BannerPatternType {
	return BannerPatternType{16}
}

func MojangBannerPattern() BannerPatternType {
	return BannerPatternType{17}
}

func RhombusBannerPattern() BannerPatternType {
	return BannerPatternType{18}
}

func SkullBannerPattern() BannerPatternType {
	return BannerPatternType{19}
}

func SmallStripesBannerPattern() BannerPatternType {
	return BannerPatternType{20}
}

func SquareBottomLeftBannerPattern() BannerPatternType {
	return BannerPatternType{21}
}

func SquareBottomRightBannerPattern() BannerPatternType {
	return BannerPatternType{22}
}

func SquareTopLeftBannerPattern() BannerPatternType {
	return BannerPatternType{23}
}

func SquareTopRightBannerPattern() BannerPatternType {
	return BannerPatternType{24}
}

func StraightCrossBannerPattern() BannerPatternType {
	return BannerPatternType{25}
}

func StripeBottomBannerPattern() BannerPatternType {
	return BannerPatternType{26}
}

func StripeCentreBannerPattern() BannerPatternType {
	return BannerPatternType{27}
}

func StripeDownLeftBannerPattern() BannerPatternType {
	return BannerPatternType{28}
}

func StripeDownRightBannerPattern() BannerPatternType {
	return BannerPatternType{29}
}

func StripeLeftBannerPattern() BannerPatternType {
	return BannerPatternType{30}
}

func StripeMiddleBannerPattern() BannerPatternType {
	return BannerPatternType{31}
}

func StripeRightBannerPattern() BannerPatternType {
	return BannerPatternType{32}
}

func StripeTopBannerPattern() BannerPatternType {
	return BannerPatternType{33}
}

func TriangleBottomBannerPattern() BannerPatternType {
	return BannerPatternType{34}
}

func TriangleTopBannerPattern() BannerPatternType {
	return BannerPatternType{35}
}

func TrianglesBottomBannerPattern() BannerPatternType {
	return BannerPatternType{36}
}

func TrianglesTopBannerPattern() BannerPatternType {
	return BannerPatternType{37}
}

func GlobeBannerPattern() BannerPatternType {
	return BannerPatternType{38}
}

func PiglinBannerPattern() BannerPatternType {
	return BannerPatternType{39}
}

func FlowBannerPattern() BannerPatternType {
	return BannerPatternType{40}
}

func GusterBannerPattern() BannerPatternType {
	return BannerPatternType{41}
}

func BannerPatternTypes() []BannerPatternType {
	return []BannerPatternType{
		BorderBannerPattern(),
		BricksBannerPattern(),
		CircleBannerPattern(),
		CreeperBannerPattern(),
		CrossBannerPattern(),
		CurlyBorderBannerPattern(),
		DiagonalLeftBannerPattern(),
		DiagonalRightBannerPattern(),
		DiagonalUpLeftBannerPattern(),
		DiagonalUpRightBannerPattern(),
		FlowerBannerPattern(),
		GradientBannerPattern(),
		GradientUpBannerPattern(),
		HalfHorizontalBannerPattern(),
		HalfHorizontalBottomBannerPattern(),
		HalfVerticalBannerPattern(),
		HalfVerticalRightBannerPattern(),
		MojangBannerPattern(),
		RhombusBannerPattern(),
		SkullBannerPattern(),
		SmallStripesBannerPattern(),
		SquareBottomLeftBannerPattern(),
		SquareBottomRightBannerPattern(),
		SquareTopLeftBannerPattern(),
		SquareTopRightBannerPattern(),
		StraightCrossBannerPattern(),
		StripeBottomBannerPattern(),
		StripeCentreBannerPattern(),
		StripeDownLeftBannerPattern(),
		StripeDownRightBannerPattern(),
		StripeLeftBannerPattern(),
		StripeMiddleBannerPattern(),
		StripeRightBannerPattern(),
		StripeTopBannerPattern(),
		TriangleBottomBannerPattern(),
		TriangleTopBannerPattern(),
		TrianglesBottomBannerPattern(),
		TrianglesTopBannerPattern(),
		GlobeBannerPattern(),
		PiglinBannerPattern(),
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
		return "border"
	case 1:
		return "bricks"
	case 2:
		return "circle"
	case 3:
		return "creeper"
	case 4:
		return "cross"
	case 5:
		return "curly_border"
	case 6:
		return "diagonal_left"
	case 7:
		return "diagonal_right"
	case 8:
		return "diagonal_up_left"
	case 9:
		return "diagonal_up_right"
	case 10:
		return "flower"
	case 11:
		return "gradient"
	case 12:
		return "gradient_up"
	case 13:
		return "half_horizontal"
	case 14:
		return "half_horizontal_bottom"
	case 15:
		return "half_vertical"
	case 16:
		return "half_vertical_right"
	case 17:
		return "mojang"
	case 18:
		return "rhombus"
	case 19:
		return "skull"
	case 20:
		return "small_stripes"
	case 21:
		return "square_bottom_left"
	case 22:
		return "square_bottom_right"
	case 23:
		return "square_top_left"
	case 24:
		return "square_top_right"
	case 25:
		return "straight_cross"
	case 26:
		return "stripe_bottom"
	case 27:
		return "stripe_center"
	case 28:
		return "stripe_downleft"
	case 29:
		return "stripe_downright"
	case 30:
		return "stripe_left"
	case 31:
		return "stripe_middle"
	case 32:
		return "stripe_right"
	case 33:
		return "stripe_top"
	case 34:
		return "triangle_bottom"
	case 35:
		return "triangle_top"
	case 36:
		return "triangles_bottom"
	case 37:
		return "triangles_top"
	case 38:
		return "globe"
	case 39:
		return "piglin"
	case 40:
		return "flow"
	case 41:
		return "guster"
	}
	panic("should never happen")
}

func (b bannerPatternType) Item() (item.BannerPatternType, bool) {
	switch b {
	case 1:
		return item.FieldMasonedBannerPattern(), true
	case 3:
		return item.CreeperBannerPattern(), true
	case 5:
		return item.BordureIndentedBannerPattern(), true
	case 10:
		return item.FlowerBannerPattern(), true
	case 17:
		return item.MojangBannerPattern(), true
	case 19:
		return item.SkullBannerPattern(), true
	case 38:
		return item.GlobeBannerPattern(), true
	case 39:
		return item.PiglinBannerPattern(), true
	case 40:
		return item.FlowBannerPattern(), true
	case 41:
		return item.GusterBannerPattern(), true
	}
	return item.BannerPatternType{}, false
}
