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

type WindsweptForest struct{}

func (WindsweptForest) Temperature() float64 {
	return 0.2
}

func (WindsweptForest) Rainfall() float64 {
	return 0.3
}

func (WindsweptForest) Depth() float64 {
	return 1
}

func (WindsweptForest) Scale() float64 {
	return 0.5
}

func (WindsweptForest) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0e, G: 0x63, B: 0xab, A: 0xa5}
}

func (WindsweptForest) Tags() []string {
	return []string{"animal", "extreme_hills", "forest", "monster", "mountain", "overworld", "spawns_cold_variant_farm_animals"}
}

func (WindsweptForest) String() string {
	return "extreme_hills_plus_trees"
}

func (WindsweptForest) EncodeBiome() int {
	return 34
}
