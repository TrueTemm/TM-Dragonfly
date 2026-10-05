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

type OldGrowthPineTaiga struct{}

func (OldGrowthPineTaiga) Temperature() float64 {
	return 0.3
}

func (OldGrowthPineTaiga) Rainfall() float64 {
	return 0.8
}

func (OldGrowthPineTaiga) Depth() float64 {
	return 0.2
}

func (OldGrowthPineTaiga) Scale() float64 {
	return 0.2
}

func (OldGrowthPineTaiga) WaterColour() color.RGBA {
	return color.RGBA{R: 0x2d, G: 0x6d, B: 0x77, A: 0xa5}
}

func (OldGrowthPineTaiga) Tags() []string {
	return []string{"animal", "forest", "mega", "monster", "overworld", "rare", "taiga", "has_structure_trail_ruins", "spawns_cold_variant_farm_animals"}
}

func (OldGrowthPineTaiga) String() string {
	return "mega_taiga"
}

func (OldGrowthPineTaiga) EncodeBiome() int {
	return 32
}
