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

type SnowyBeach struct{}

func (SnowyBeach) Temperature() float64 {
	return 0.05
}

func (SnowyBeach) Rainfall() float64 {
	return 0.3
}

func (SnowyBeach) Depth() float64 {
	return 0
}

func (SnowyBeach) Scale() float64 {
	return 0.025
}

func (SnowyBeach) WaterColour() color.RGBA {
	return color.RGBA{R: 0x14, G: 0x63, B: 0xa5, A: 0xa5}
}

func (SnowyBeach) Tags() []string {
	return []string{"beach", "cold", "monster", "overworld", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_snow_foxes", "spawns_white_rabbits"}
}

func (SnowyBeach) String() string {
	return "cold_beach"
}

func (SnowyBeach) EncodeBiome() int {
	return 26
}
