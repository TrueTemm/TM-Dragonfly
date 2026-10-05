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

type Swamp struct{}

func (Swamp) Temperature() float64 {
	return 0.8
}

func (Swamp) Rainfall() float64 {
	return 0.9
}

func (Swamp) Depth() float64 {
	return -0.2
}

func (Swamp) Scale() float64 {
	return 0.1
}

func (Swamp) WaterColour() color.RGBA {
	return color.RGBA{R: 0x61, G: 0x7b, B: 0x64, A: 0xa5}
}

func (Swamp) Tags() []string {
	return []string{"animal", "monster", "overworld", "swamp", "spawns_slimes_on_surface", "slime", "swamp_water_huge_mushroom"}
}

func (Swamp) String() string {
	return "swampland"
}

func (Swamp) EncodeBiome() int {
	return 6
}
