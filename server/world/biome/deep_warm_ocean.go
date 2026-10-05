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

type DeepWarmOcean struct{}

func (DeepWarmOcean) Temperature() float64 {
	return 0.5
}

func (DeepWarmOcean) Rainfall() float64 {
	return 0.5
}

func (DeepWarmOcean) Depth() float64 {
	return -1.8
}

func (DeepWarmOcean) Scale() float64 {
	return 0.1
}

func (DeepWarmOcean) WaterColour() color.RGBA {
	return color.RGBA{R: 0x02, G: 0xb0, B: 0xe5, A: 0xa5}
}

func (DeepWarmOcean) Tags() []string {
	return []string{"deep", "monster", "ocean", "overworld", "warm", "spawns_warm_variant_farm_animals", "fast_fishing", "high_seas"}
}

func (DeepWarmOcean) String() string {
	return "deep_warm_ocean"
}

func (DeepWarmOcean) EncodeBiome() int {
	return 41
}
