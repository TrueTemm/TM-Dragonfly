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

type SnowyTaigaMountains struct{}

func (SnowyTaigaMountains) Temperature() float64 {
	return -0.5
}

func (SnowyTaigaMountains) Rainfall() float64 {
	return 0.4
}

func (SnowyTaigaMountains) Depth() float64 {
	return 0.3
}

func (SnowyTaigaMountains) Scale() float64 {
	return 0.4
}

func (SnowyTaigaMountains) WaterColour() color.RGBA {
	return color.RGBA{R: 0x20, G: 0x5e, B: 0x83, A: 0xa5}
}

func (SnowyTaigaMountains) Tags() []string {
	return []string{"animal", "cold", "forest", "monster", "mutated", "taiga", "overworld_generation", "spawns_cold_variant_farm_animals"}
}

func (SnowyTaigaMountains) String() string {
	return "cold_taiga_mutated"
}

func (SnowyTaigaMountains) EncodeBiome() int {
	return 158
}
