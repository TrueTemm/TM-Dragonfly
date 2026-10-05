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

type CrimsonForest struct{}

func (CrimsonForest) Temperature() float64 {
	return 2
}

func (CrimsonForest) Rainfall() float64 {
	return 0
}

func (CrimsonForest) Depth() float64 {
	return 0.1
}

func (CrimsonForest) Scale() float64 {
	return 0.2
}

func (CrimsonForest) WaterColour() color.RGBA {
	return color.RGBA{R: 0x90, G: 0x59, B: 0x57, A: 0xa5}
}

func (CrimsonForest) Tags() []string {
	return []string{"nether", "netherwart_forest", "crimson_forest", "spawn_few_zombified_piglins", "spawn_piglin", "spawns_warm_variant_farm_animals"}
}

func (CrimsonForest) String() string {
	return "crimson_forest"
}

func (CrimsonForest) EncodeBiome() int {
	return 179
}
