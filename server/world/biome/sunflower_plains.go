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

type SunflowerPlains struct{}

func (SunflowerPlains) Temperature() float64 {
	return 0.8
}

func (SunflowerPlains) Rainfall() float64 {
	return 0.4
}

func (SunflowerPlains) Depth() float64 {
	return 0.125
}

func (SunflowerPlains) Scale() float64 {
	return 0.05
}

func (SunflowerPlains) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (SunflowerPlains) Tags() []string {
	return []string{"animal", "monster", "mutated", "overworld", "plains", "bee_habitat"}
}

func (SunflowerPlains) String() string {
	return "sunflower_plains"
}

func (SunflowerPlains) EncodeBiome() int {
	return 129
}
