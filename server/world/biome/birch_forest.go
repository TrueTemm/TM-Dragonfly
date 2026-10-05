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

type BirchForest struct{}

func (BirchForest) Temperature() float64 {
	return 0.6
}

func (BirchForest) Rainfall() float64 {
	return 0.6
}

func (BirchForest) Depth() float64 {
	return 0.1
}

func (BirchForest) Scale() float64 {
	return 0.2
}

func (BirchForest) WaterColour() color.RGBA {
	return color.RGBA{R: 0x06, G: 0x77, B: 0xce, A: 0xa5}
}

func (BirchForest) Tags() []string {
	return []string{"animal", "birch", "forest", "monster", "overworld", "bee_habitat"}
}

func (BirchForest) String() string {
	return "birch_forest"
}

func (BirchForest) EncodeBiome() int {
	return 27
}
