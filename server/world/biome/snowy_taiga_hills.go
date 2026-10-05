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

type SnowyTaigaHills struct{}

func (SnowyTaigaHills) Temperature() float64 {
	return -0.5
}

func (SnowyTaigaHills) Rainfall() float64 {
	return 0.4
}

func (SnowyTaigaHills) Depth() float64 {
	return 0.45
}

func (SnowyTaigaHills) Scale() float64 {
	return 0.3
}

func (SnowyTaigaHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x24, G: 0x5b, B: 0x78, A: 0xa5}
}

func (SnowyTaigaHills) Tags() []string {
	return []string{"animal", "cold", "forest", "hills", "monster", "overworld", "taiga", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_white_rabbits"}
}

func (SnowyTaigaHills) String() string {
	return "cold_taiga_hills"
}

func (SnowyTaigaHills) EncodeBiome() int {
	return 31
}
