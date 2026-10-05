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

type JaggedPeaks struct{}

func (JaggedPeaks) Temperature() float64 {
	return -0.7
}

func (JaggedPeaks) Rainfall() float64 {
	return 0.9
}

func (JaggedPeaks) Depth() float64 {
	return 0.1
}

func (JaggedPeaks) Scale() float64 {
	return 0.2
}

func (JaggedPeaks) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (JaggedPeaks) Tags() []string {
	return []string{"mountains", "monster", "overworld", "frozen", "jagged_peaks", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_snow_foxes", "spawns_white_rabbits"}
}

func (JaggedPeaks) String() string {
	return "jagged_peaks"
}

func (JaggedPeaks) EncodeBiome() int {
	return 182
}
