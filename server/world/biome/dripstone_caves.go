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

type DripstoneCaves struct{}

func (DripstoneCaves) Temperature() float64 {
	return 0.2
}

func (DripstoneCaves) Rainfall() float64 {
	return 0
}

func (DripstoneCaves) Depth() float64 {
	return 0.1
}

func (DripstoneCaves) Scale() float64 {
	return 0.2
}

func (DripstoneCaves) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (DripstoneCaves) Tags() []string {
	return []string{"caves", "overworld", "dripstone_caves", "monster"}
}

func (DripstoneCaves) String() string {
	return "dripstone_caves"
}

func (DripstoneCaves) EncodeBiome() int {
	return 188
}
