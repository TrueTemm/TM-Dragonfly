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

type LushCaves struct{}

func (LushCaves) Temperature() float64 {
	return 0.9
}

func (LushCaves) Rainfall() float64 {
	return 0
}

func (LushCaves) Depth() float64 {
	return 0.1
}

func (LushCaves) Scale() float64 {
	return 0.2
}

func (LushCaves) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (LushCaves) Tags() []string {
	return []string{"caves", "lush_caves", "overworld", "monster", "spawns_tropical_fish_at_any_height"}
}

func (LushCaves) String() string {
	return "lush_caves"
}

func (LushCaves) EncodeBiome() int {
	return 187
}
