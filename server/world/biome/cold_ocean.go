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

type ColdOcean struct{}

func (ColdOcean) Temperature() float64 {
	return 0.5
}

func (ColdOcean) Rainfall() float64 {
	return 0.5
}

func (ColdOcean) Depth() float64 {
	return -1
}

func (ColdOcean) Scale() float64 {
	return 0.1
}

func (ColdOcean) WaterColour() color.RGBA {
	return color.RGBA{R: 0x20, G: 0x80, B: 0xc9, A: 0xa5}
}

func (ColdOcean) Tags() []string {
	return []string{"cold", "monster", "ocean", "overworld", "spawns_cold_variant_farm_animals", "fast_fishing", "high_seas"}
}

func (ColdOcean) String() string {
	return "cold_ocean"
}

func (ColdOcean) EncodeBiome() int {
	return 44
}
