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

type SulfurCaves struct{}

func (SulfurCaves) Temperature() float64 {
	return 0.8
}

func (SulfurCaves) Rainfall() float64 {
	return 0.4
}

func (SulfurCaves) Depth() float64 {
	return 0.1
}

func (SulfurCaves) Scale() float64 {
	return 0.2
}

func (SulfurCaves) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (SulfurCaves) Tags() []string {
	return []string{"caves", "sulfur_caves", "overworld", "monster"}
}

func (SulfurCaves) String() string {
	return "sulfur_caves"
}

func (SulfurCaves) EncodeBiome() int {
	return 194
}
