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

type BirchForestHills struct{}

func (BirchForestHills) Temperature() float64 {
	return 0.6
}

func (BirchForestHills) Rainfall() float64 {
	return 0.6
}

func (BirchForestHills) Depth() float64 {
	return 0.45
}

func (BirchForestHills) Scale() float64 {
	return 0.3
}

func (BirchForestHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0a, G: 0x74, B: 0xc4, A: 0xa5}
}

func (BirchForestHills) Tags() []string {
	return []string{"animal", "birch", "forest", "hills", "monster", "overworld", "bee_habitat"}
}

func (BirchForestHills) String() string {
	return "birch_forest_hills"
}

func (BirchForestHills) EncodeBiome() int {
	return 28
}
