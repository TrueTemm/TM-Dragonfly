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

type MangroveSwamp struct{}

func (MangroveSwamp) Temperature() float64 {
	return 0.8
}

func (MangroveSwamp) Rainfall() float64 {
	return 0.9
}

func (MangroveSwamp) Depth() float64 {
	return -0.2
}

func (MangroveSwamp) Scale() float64 {
	return 0.1
}

func (MangroveSwamp) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (MangroveSwamp) Tags() []string {
	return []string{"mangrove_swamp", "overworld", "monster", "spawns_slimes_on_surface", "spawns_warm_variant_farm_animals", "spawns_warm_variant_frogs", "slime"}
}

func (MangroveSwamp) String() string {
	return "mangrove_swamp"
}

func (MangroveSwamp) EncodeBiome() int {
	return 191
}
