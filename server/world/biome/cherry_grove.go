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

type CherryGrove struct{}

func (CherryGrove) Temperature() float64 {
	return 0.3
}

func (CherryGrove) Rainfall() float64 {
	return 0.8
}

func (CherryGrove) Depth() float64 {
	return 0.1
}

func (CherryGrove) Scale() float64 {
	return 0.2
}

func (CherryGrove) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (CherryGrove) Tags() []string {
	return []string{"mountains", "monster", "overworld", "cherry_grove", "bee_habitat"}
}

func (CherryGrove) String() string {
	return "cherry_grove"
}

func (CherryGrove) EncodeBiome() int {
	return 192
}
