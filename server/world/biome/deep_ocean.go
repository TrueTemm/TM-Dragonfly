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

type DeepOcean struct{}

func (DeepOcean) Temperature() float64 {
	return 0.5
}

func (DeepOcean) Rainfall() float64 {
	return 0.5
}

func (DeepOcean) Depth() float64 {
	return -1.8
}

func (DeepOcean) Scale() float64 {
	return 0.1
}

func (DeepOcean) WaterColour() color.RGBA {
	return color.RGBA{R: 0x17, G: 0x87, B: 0xd4, A: 0xa5}
}

func (DeepOcean) Tags() []string {
	return []string{"deep", "monster", "ocean", "overworld", "spawns_warm_variant_farm_animals", "fast_fishing", "high_seas"}
}

func (DeepOcean) String() string {
	return "deep_ocean"
}

func (DeepOcean) EncodeBiome() int {
	return 24
}
