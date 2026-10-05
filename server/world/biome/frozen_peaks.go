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

type FrozenPeaks struct{}

func (FrozenPeaks) Temperature() float64 {
	return -0.7
}

func (FrozenPeaks) Rainfall() float64 {
	return 0.9
}

func (FrozenPeaks) Depth() float64 {
	return 0.1
}

func (FrozenPeaks) Scale() float64 {
	return 0.2
}

func (FrozenPeaks) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (FrozenPeaks) Tags() []string {
	return []string{"mountains", "monster", "overworld", "frozen", "frozen_peaks", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_snow_foxes", "spawns_white_rabbits"}
}

func (FrozenPeaks) String() string {
	return "frozen_peaks"
}

func (FrozenPeaks) EncodeBiome() int {
	return 183
}
