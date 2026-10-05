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

type River struct{}

func (River) Temperature() float64 {
	return 0.5
}

func (River) Rainfall() float64 {
	return 0.5
}

func (River) Depth() float64 {
	return -0.5
}

func (River) Scale() float64 {
	return 0
}

func (River) WaterColour() color.RGBA {
	return color.RGBA{R: 0x00, G: 0x84, B: 0xff, A: 0xa5}
}

func (River) Tags() []string {
	return []string{"overworld", "spawns_more_frequent_drowned", "spawns_reduced_water_ambient_mobs", "spawns_river_mobs", "river"}
}

func (River) String() string {
	return "river"
}

func (River) EncodeBiome() int {
	return 7
}
