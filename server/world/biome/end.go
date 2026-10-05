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

type End struct{}

func (End) Temperature() float64 {
	return 0.5
}

func (End) Rainfall() float64 {
	return 0.5
}

func (End) Depth() float64 {
	return 0.1
}

func (End) Scale() float64 {
	return 0.2
}

func (End) WaterColour() color.RGBA {
	return color.RGBA{R: 0x62, G: 0x52, B: 0x9e, A: 0xa5}
}

func (End) Tags() []string {
	return []string{"the_end", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs"}
}

func (End) String() string {
	return "the_end"
}

func (End) EncodeBiome() int {
	return 9
}
