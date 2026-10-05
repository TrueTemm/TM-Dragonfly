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

package biome

import "image/color"

type DeepLukewarmOcean struct{}

func (DeepLukewarmOcean) Temperature() float64 {
	return 0.5
}

func (DeepLukewarmOcean) Rainfall() float64 {
	return 0.5
}

func (DeepLukewarmOcean) Depth() float64 {
	return -1.8
}

func (DeepLukewarmOcean) Scale() float64 {
	return 0.1
}

func (DeepLukewarmOcean) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0d, G: 0x96, B: 0xdb, A: 0xa5}
}

func (DeepLukewarmOcean) Tags() []string {
	return []string{"deep", "lukewarm", "monster", "ocean", "overworld", "spawns_warm_variant_farm_animals", "fast_fishing", "high_seas"}
}

func (DeepLukewarmOcean) String() string {
	return "deep_lukewarm_ocean"
}

func (DeepLukewarmOcean) EncodeBiome() int {
	return 43
}
