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

type Forest struct{}

func (Forest) Temperature() float64 {
	return 0.7
}

func (Forest) Rainfall() float64 {
	return 0.8
}

func (Forest) Depth() float64 {
	return 0.1
}

func (Forest) Scale() float64 {
	return 0.2
}

func (Forest) WaterColour() color.RGBA {
	return color.RGBA{R: 0x1e, G: 0x97, B: 0xf2, A: 0xa5}
}

func (Forest) Tags() []string {
	return []string{"animal", "forest", "monster", "overworld", "bee_habitat"}
}

func (Forest) String() string {
	return "forest"
}

func (Forest) EncodeBiome() int {
	return 4
}
