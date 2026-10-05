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

type SnowyPlains struct{}

func (SnowyPlains) Temperature() float64 {
	return 0
}

func (SnowyPlains) Rainfall() float64 {
	return 0.5
}

func (SnowyPlains) Depth() float64 {
	return 0.125
}

func (SnowyPlains) Scale() float64 {
	return 0.05
}

func (SnowyPlains) WaterColour() color.RGBA {
	return color.RGBA{R: 0x14, G: 0x55, B: 0x9b, A: 0xa5}
}

func (SnowyPlains) Tags() []string {
	return []string{"frozen", "ice", "ice_plains", "overworld", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_snow_foxes", "spawns_white_rabbits"}
}

func (SnowyPlains) String() string {
	return "ice_plains"
}

func (SnowyPlains) EncodeBiome() int {
	return 12
}
