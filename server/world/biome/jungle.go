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

type Jungle struct{}

func (Jungle) Temperature() float64 {
	return 0.95
}

func (Jungle) Rainfall() float64 {
	return 0.9
}

func (Jungle) Depth() float64 {
	return 0.1
}

func (Jungle) Scale() float64 {
	return 0.2
}

func (Jungle) WaterColour() color.RGBA {
	return color.RGBA{R: 0x14, G: 0xa2, B: 0xc5, A: 0xa5}
}

func (Jungle) Tags() []string {
	return []string{"animal", "has_structure_trail_ruins", "jungle", "monster", "overworld", "rare", "spawns_jungle_mobs", "spawns_warm_variant_farm_animals"}
}

func (Jungle) String() string {
	return "jungle"
}

func (Jungle) EncodeBiome() int {
	return 21
}
