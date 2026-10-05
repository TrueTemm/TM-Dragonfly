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

type WindsweptHills struct{}

func (WindsweptHills) Temperature() float64 {
	return 0.2
}

func (WindsweptHills) Rainfall() float64 {
	return 0.3
}

func (WindsweptHills) Depth() float64 {
	return 1
}

func (WindsweptHills) Scale() float64 {
	return 0.5
}

func (WindsweptHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x00, G: 0x7b, B: 0xf7, A: 0xa5}
}

func (WindsweptHills) Tags() []string {
	return []string{"animal", "extreme_hills", "monster", "overworld", "spawns_cold_variant_farm_animals"}
}

func (WindsweptHills) String() string {
	return "extreme_hills"
}

func (WindsweptHills) EncodeBiome() int {
	return 3
}
