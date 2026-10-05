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

type WarpedForest struct{}

func (WarpedForest) Temperature() float64 {
	return 2
}

func (WarpedForest) Rainfall() float64 {
	return 0
}

func (WarpedForest) Depth() float64 {
	return 0.1
}

func (WarpedForest) Scale() float64 {
	return 0.2
}

func (WarpedForest) WaterColour() color.RGBA {
	return color.RGBA{R: 0x90, G: 0x59, B: 0x57, A: 0xa5}
}

func (WarpedForest) Tags() []string {
	return []string{"nether", "netherwart_forest", "warped_forest", "spawn_endermen", "spawns_warm_variant_farm_animals"}
}

func (WarpedForest) String() string {
	return "warped_forest"
}

func (WarpedForest) EncodeBiome() int {
	return 180
}
