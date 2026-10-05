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

type BambooJungleHills struct{}

func (BambooJungleHills) Temperature() float64 {
	return 0.95
}

func (BambooJungleHills) Rainfall() float64 {
	return 0.9
}

func (BambooJungleHills) Depth() float64 {
	return 0.45
}

func (BambooJungleHills) Scale() float64 {
	return 0.3
}

func (BambooJungleHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x1b, G: 0x9e, B: 0xd8, A: 0xa5}
}

func (BambooJungleHills) Tags() []string {
	return []string{"animal", "bamboo", "hills", "jungle", "monster", "overworld", "spawns_jungle_mobs", "spawns_warm_variant_farm_animals"}
}

func (BambooJungleHills) String() string {
	return "bamboo_jungle_hills"
}

func (BambooJungleHills) EncodeBiome() int {
	return 49
}
