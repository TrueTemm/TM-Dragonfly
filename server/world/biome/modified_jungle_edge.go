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

type ModifiedJungleEdge struct{}

func (ModifiedJungleEdge) Temperature() float64 {
	return 0.95
}

func (ModifiedJungleEdge) Rainfall() float64 {
	return 0.8
}

func (ModifiedJungleEdge) Depth() float64 {
	return 0.2
}

func (ModifiedJungleEdge) Scale() float64 {
	return 0.4
}

func (ModifiedJungleEdge) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0d, G: 0x8a, B: 0xe3, A: 0xa5}
}

func (ModifiedJungleEdge) Tags() []string {
	return []string{"animal", "edge", "jungle", "monster", "mutated", "overworld_generation", "spawns_jungle_mobs", "spawns_warm_variant_farm_animals"}
}

func (ModifiedJungleEdge) String() string {
	return "jungle_edge_mutated"
}

func (ModifiedJungleEdge) EncodeBiome() int {
	return 151
}
