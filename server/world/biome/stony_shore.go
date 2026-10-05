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

type StonyShore struct{}

func (StonyShore) Temperature() float64 {
	return 0.2
}

func (StonyShore) Rainfall() float64 {
	return 0.3
}

func (StonyShore) Depth() float64 {
	return 0.1
}

func (StonyShore) Scale() float64 {
	return 0.8
}

func (StonyShore) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0d, G: 0x67, B: 0xbb, A: 0xa5}
}

func (StonyShore) Tags() []string {
	return []string{"beach", "monster", "overworld", "stone"}
}

func (StonyShore) String() string {
	return "stone_beach"
}

func (StonyShore) EncodeBiome() int {
	return 25
}
