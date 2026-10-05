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

type DeepFrozenOcean struct{}

func (DeepFrozenOcean) Temperature() float64 {
	return 0.5
}

func (DeepFrozenOcean) Rainfall() float64 {
	return 0.5
}

func (DeepFrozenOcean) Depth() float64 {
	return -1.8
}

func (DeepFrozenOcean) Scale() float64 {
	return 0.1
}

func (DeepFrozenOcean) WaterColour() color.RGBA {
	return color.RGBA{R: 0x25, G: 0x70, B: 0xb5, A: 0xa5}
}

func (DeepFrozenOcean) Tags() []string {
	return []string{"deep", "frozen", "monster", "ocean", "overworld", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_polar_bears_on_alternate_blocks", "fast_fishing", "high_seas"}
}

func (DeepFrozenOcean) String() string {
	return "deep_frozen_ocean"
}

func (DeepFrozenOcean) EncodeBiome() int {
	return 47
}
