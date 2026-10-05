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

type FrozenRiver struct{}

func (FrozenRiver) Temperature() float64 {
	return 0
}

func (FrozenRiver) Rainfall() float64 {
	return 0.5
}

func (FrozenRiver) Depth() float64 {
	return -0.5
}

func (FrozenRiver) Scale() float64 {
	return 0
}

func (FrozenRiver) WaterColour() color.RGBA {
	return color.RGBA{R: 0x18, G: 0x53, B: 0x90, A: 0xa5}
}

func (FrozenRiver) Tags() []string {
	return []string{"frozen", "overworld", "river", "spawns_cold_variant_farm_animals", "spawns_cold_variant_frogs", "spawns_river_mobs", "spawns_snow_foxes", "spawns_white_rabbits"}
}

func (FrozenRiver) String() string {
	return "frozen_river"
}

func (FrozenRiver) EncodeBiome() int {
	return 11
}
