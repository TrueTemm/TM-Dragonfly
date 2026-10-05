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

var (
	bannerPatternsMap = map[string]BannerPatternType{}
	bannerPatternIDs  = map[BannerPatternType]string{}
)

func init() {
	registerBannerPattern("bo", BorderBannerPattern())
	registerBannerPattern("bri", BricksBannerPattern())
	registerBannerPattern("mc", CircleBannerPattern())
	registerBannerPattern("cre", CreeperBannerPattern())
	registerBannerPattern("cr", CrossBannerPattern())
	registerBannerPattern("cbo", CurlyBorderBannerPattern())
	registerBannerPattern("lud", DiagonalLeftBannerPattern())
	registerBannerPattern("rd", DiagonalRightBannerPattern())
	registerBannerPattern("ld", DiagonalUpLeftBannerPattern())
	registerBannerPattern("rud", DiagonalUpRightBannerPattern())
	registerBannerPattern("flo", FlowerBannerPattern())
	registerBannerPattern("gra", GradientBannerPattern())
	registerBannerPattern("gru", GradientUpBannerPattern())
	registerBannerPattern("hh", HalfHorizontalBannerPattern())
	registerBannerPattern("hhb", HalfHorizontalBottomBannerPattern())
	registerBannerPattern("vh", HalfVerticalBannerPattern())
	registerBannerPattern("vhr", HalfVerticalRightBannerPattern())
	registerBannerPattern("moj", MojangBannerPattern())
	registerBannerPattern("mr", RhombusBannerPattern())
	registerBannerPattern("sku", SkullBannerPattern())
	registerBannerPattern("ss", SmallStripesBannerPattern())
	registerBannerPattern("bl", SquareBottomLeftBannerPattern())
	registerBannerPattern("br", SquareBottomRightBannerPattern())
	registerBannerPattern("tl", SquareTopLeftBannerPattern())
	registerBannerPattern("tr", SquareTopRightBannerPattern())
	registerBannerPattern("sc", StraightCrossBannerPattern())
	registerBannerPattern("bs", StripeBottomBannerPattern())
	registerBannerPattern("cs", StripeCentreBannerPattern())
	registerBannerPattern("dls", StripeDownLeftBannerPattern())
	registerBannerPattern("drs", StripeDownRightBannerPattern())
	registerBannerPattern("ls", StripeLeftBannerPattern())
	registerBannerPattern("ms", StripeMiddleBannerPattern())
	registerBannerPattern("rs", StripeRightBannerPattern())
	registerBannerPattern("ts", StripeTopBannerPattern())
	registerBannerPattern("bt", TriangleBottomBannerPattern())
	registerBannerPattern("tt", TriangleTopBannerPattern())
	registerBannerPattern("bts", TrianglesBottomBannerPattern())
	registerBannerPattern("tts", TrianglesTopBannerPattern())
	registerBannerPattern("glb", GlobeBannerPattern())
	registerBannerPattern("pig", PiglinBannerPattern())
	registerBannerPattern("flw", FlowBannerPattern())
	registerBannerPattern("gus", GusterBannerPattern())
}

func registerBannerPattern(id string, pattern BannerPatternType) {
	bannerPatternsMap[id] = pattern
	bannerPatternIDs[pattern] = id
}

func BannerPatternByID(id string) (BannerPatternType, bool) {
	b, ok := bannerPatternsMap[id]
	return b, ok
}

func bannerPatternID(pattern BannerPatternType) string {
	id, ok := bannerPatternIDs[pattern]
	if !ok {
		panic("should never happen")
	}
	return id
}
