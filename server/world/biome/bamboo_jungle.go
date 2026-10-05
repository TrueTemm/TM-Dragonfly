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

type BambooJungle struct{}

func (BambooJungle) Temperature() float64 {
	return 0.95
}

func (BambooJungle) Rainfall() float64 {
	return 0.9
}

func (BambooJungle) Depth() float64 {
	return 0.1
}

func (BambooJungle) Scale() float64 {
	return 0.2
}

func (BambooJungle) WaterColour() color.RGBA {
	return color.RGBA{R: 0x14, G: 0xa2, B: 0xc5, A: 0xa5}
}

func (BambooJungle) Tags() []string {
	return []string{"animal", "bamboo", "jungle", "monster", "overworld", "spawns_jungle_mobs", "spawns_warm_variant_farm_animals"}
}

func (BambooJungle) String() string {
	return "bamboo_jungle"
}

func (BambooJungle) EncodeBiome() int {
	return 48
}
