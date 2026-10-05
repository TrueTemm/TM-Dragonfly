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

type SnowyTaiga struct{}

func (SnowyTaiga) Temperature() float64 {
	return -0.5
}

func (SnowyTaiga) Rainfall() float64 {
	return 0.4
}

func (SnowyTaiga) Depth() float64 {
	return 0.2
}

func (SnowyTaiga) Scale() float64 {
	return 0.2
}

func (SnowyTaiga) WaterColour() color.RGBA {
	return color.RGBA{R: 0x20, G: 0x5e, B: 0x83, A: 0xa5}
}

func (SnowyTaiga) Tags() []string {
	return []string{"animal", "cold", "forest", "monster", "overworld", "taiga", "has_structure_trail_ruins", "spawns_cold_variant_farm_animals"}
}

func (SnowyTaiga) String() string {
	return "cold_taiga"
}

func (SnowyTaiga) EncodeBiome() int {
	return 30
}
