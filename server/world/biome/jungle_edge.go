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

type JungleEdge struct{}

func (JungleEdge) Temperature() float64 {
	return 0.95
}

func (JungleEdge) Rainfall() float64 {
	return 0.8
}

func (JungleEdge) Depth() float64 {
	return 0.1
}

func (JungleEdge) Scale() float64 {
	return 0.2
}

func (JungleEdge) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0d, G: 0x8a, B: 0xe3, A: 0xa5}
}

func (JungleEdge) Tags() []string {
	return []string{"animal", "edge", "jungle", "monster", "overworld", "spawns_jungle_mobs", "spawns_warm_variant_farm_animals"}
}

func (JungleEdge) String() string {
	return "jungle_edge"
}

func (JungleEdge) EncodeBiome() int {
	return 23
}
