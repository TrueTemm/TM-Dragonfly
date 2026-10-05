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

type SwampHills struct{}

func (SwampHills) Temperature() float64 {
	return 0.8
}

func (SwampHills) Rainfall() float64 {
	return 0.5
}

func (SwampHills) Depth() float64 {
	return -0.1
}

func (SwampHills) Scale() float64 {
	return 0.3
}

func (SwampHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x61, G: 0x7b, B: 0x64, A: 0xa5}
}

func (SwampHills) Tags() []string {
	return []string{"animal", "monster", "mutated", "swamp", "overworld_generation", "spawns_slimes_on_surface", "slime", "swamp_water_huge_mushroom"}
}

func (SwampHills) String() string {
	return "swampland_mutated"
}

func (SwampHills) EncodeBiome() int {
	return 134
}
