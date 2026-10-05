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

type WindsweptSavanna struct{}

func (WindsweptSavanna) Temperature() float64 {
	return 2
}

func (WindsweptSavanna) Rainfall() float64 {
	return 0
}

func (WindsweptSavanna) Depth() float64 {
	return 0.363
}

func (WindsweptSavanna) Scale() float64 {
	return 1.225
}

func (WindsweptSavanna) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (WindsweptSavanna) Tags() []string {
	return []string{"animal", "monster", "mutated", "overworld", "savanna", "spawns_savanna_mobs", "spawns_warm_variant_farm_animals"}
}

func (WindsweptSavanna) String() string {
	return "savanna_mutated"
}

func (WindsweptSavanna) EncodeBiome() int {
	return 163
}
