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

type PaleGarden struct{}

func (PaleGarden) Temperature() float64 {
	return 0.7
}

func (PaleGarden) Rainfall() float64 {
	return 0.8
}

func (PaleGarden) Depth() float64 {
	return 0.1
}

func (PaleGarden) Scale() float64 {
	return 0.2
}

func (PaleGarden) WaterColour() color.RGBA {
	return color.RGBA{R: 0x60, G: 0xb7, B: 0xff, A: 0xa6}
}

func (PaleGarden) Tags() []string {
	return []string{"monster", "overworld", "pale_garden"}
}

func (PaleGarden) String() string {
	return "pale_garden"
}

func (PaleGarden) EncodeBiome() int {
	return 193
}
