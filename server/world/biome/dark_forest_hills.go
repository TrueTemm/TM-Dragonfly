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

type DarkForestHills struct{}

func (DarkForestHills) Temperature() float64 {
	return 0.7
}

func (DarkForestHills) Rainfall() float64 {
	return 0.8
}

func (DarkForestHills) Depth() float64 {
	return 0.2
}

func (DarkForestHills) Scale() float64 {
	return 0.4
}

func (DarkForestHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x3b, G: 0x6c, B: 0xd1, A: 0xa5}
}

func (DarkForestHills) Tags() []string {
	return []string{"animal", "forest", "monster", "mutated", "roofed", "overworld_generation"}
}

func (DarkForestHills) String() string {
	return "roofed_forest_mutated"
}

func (DarkForestHills) EncodeBiome() int {
	return 157
}
