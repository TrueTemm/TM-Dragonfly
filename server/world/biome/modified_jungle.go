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

type ModifiedJungle struct{}

func (ModifiedJungle) Temperature() float64 {
	return 0.95
}

func (ModifiedJungle) Rainfall() float64 {
	return 0.9
}

func (ModifiedJungle) Depth() float64 {
	return 0.2
}

func (ModifiedJungle) Scale() float64 {
	return 0.4
}

func (ModifiedJungle) WaterColour() color.RGBA {
	return color.RGBA{R: 0x1b, G: 0x9e, B: 0xd8, A: 0xa5}
}

func (ModifiedJungle) Tags() []string {
	return []string{"animal", "jungle", "monster", "mutated", "overworld_generation", "spawns_warm_variant_farm_animals"}
}

func (ModifiedJungle) String() string {
	return "jungle_mutated"
}

func (ModifiedJungle) EncodeBiome() int {
	return 149
}
