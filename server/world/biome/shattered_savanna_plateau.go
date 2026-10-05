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

type ShatteredSavannaPlateau struct{}

func (ShatteredSavannaPlateau) Temperature() float64 {
	return 1
}

func (ShatteredSavannaPlateau) Rainfall() float64 {
	return 0.5
}

func (ShatteredSavannaPlateau) Depth() float64 {
	return 1.05
}

func (ShatteredSavannaPlateau) Scale() float64 {
	return 1.212
}

func (ShatteredSavannaPlateau) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (ShatteredSavannaPlateau) Tags() []string {
	return []string{"animal", "monster", "mutated", "overworld", "plateau", "savanna", "spawns_savanna_mobs", "spawns_warm_variant_farm_animals"}
}

func (ShatteredSavannaPlateau) String() string {
	return "savanna_plateau_mutated"
}

func (ShatteredSavannaPlateau) EncodeBiome() int {
	return 164
}
