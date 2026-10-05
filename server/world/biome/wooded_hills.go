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

type WoodedHills struct{}

func (WoodedHills) Temperature() float64 {
	return 0.7
}

func (WoodedHills) Rainfall() float64 {
	return 0.8
}

func (WoodedHills) Depth() float64 {
	return 0.45
}

func (WoodedHills) Scale() float64 {
	return 0.3
}

func (WoodedHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x05, G: 0x6b, B: 0xd1, A: 0xa5}
}

func (WoodedHills) Tags() []string {
	return []string{"animal", "hills", "monster", "overworld", "forest", "bee_habitat"}
}

func (WoodedHills) String() string {
	return "forest_hills"
}

func (WoodedHills) EncodeBiome() int {
	return 18
}
