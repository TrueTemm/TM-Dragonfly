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

type Beach struct{}

func (Beach) Temperature() float64 {
	return 0.8
}

func (Beach) Rainfall() float64 {
	return 0.4
}

func (Beach) Depth() float64 {
	return 0
}

func (Beach) Scale() float64 {
	return 0.025
}

func (Beach) WaterColour() color.RGBA {
	return color.RGBA{R: 0x15, G: 0x7c, B: 0xab, A: 0xa5}
}

func (Beach) Tags() []string {
	return []string{"beach", "monster", "overworld", "warm"}
}

func (Beach) String() string {
	return "beach"
}

func (Beach) EncodeBiome() int {
	return 16
}
