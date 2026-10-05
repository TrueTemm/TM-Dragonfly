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

type OldGrowthSpruceTaiga struct{}

func (OldGrowthSpruceTaiga) Temperature() float64 {
	return 0.25
}

func (OldGrowthSpruceTaiga) Rainfall() float64 {
	return 0.8
}

func (OldGrowthSpruceTaiga) Depth() float64 {
	return 0.2
}

func (OldGrowthSpruceTaiga) Scale() float64 {
	return 0.2
}

func (OldGrowthSpruceTaiga) WaterColour() color.RGBA {
	return color.RGBA{R: 0x2d, G: 0x6d, B: 0x77, A: 0xa5}
}

func (OldGrowthSpruceTaiga) Tags() []string {
	return []string{"animal", "forest", "mega", "monster", "mutated", "overworld", "taiga", "has_structure_trail_ruins", "spawns_cold_variant_farm_animals"}
}

func (OldGrowthSpruceTaiga) String() string {
	return "redwood_taiga_mutated"
}

func (OldGrowthSpruceTaiga) EncodeBiome() int {
	return 160
}
