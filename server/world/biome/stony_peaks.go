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

type StonyPeaks struct{}

func (StonyPeaks) Temperature() float64 {
	return 1
}

func (StonyPeaks) Rainfall() float64 {
	return 0.3
}

func (StonyPeaks) Depth() float64 {
	return 0.1
}

func (StonyPeaks) Scale() float64 {
	return 0.2
}

func (StonyPeaks) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (StonyPeaks) Tags() []string {
	return []string{"mountains", "monster", "overworld", "spawns_cold_variant_farm_animals"}
}

func (StonyPeaks) String() string {
	return "stony_peaks"
}

func (StonyPeaks) EncodeBiome() int {
	return 189
}
