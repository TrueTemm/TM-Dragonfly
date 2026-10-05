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

type MushroomFields struct{}

func (MushroomFields) Temperature() float64 {
	return 0.9
}

func (MushroomFields) Rainfall() float64 {
	return 1
}

func (MushroomFields) Depth() float64 {
	return 0.2
}

func (MushroomFields) Scale() float64 {
	return 0.3
}

func (MushroomFields) WaterColour() color.RGBA {
	return color.RGBA{R: 0x8a, G: 0x89, B: 0x97, A: 0xa5}
}

func (MushroomFields) Tags() []string {
	return []string{"mooshroom_island", "overworld", "spawns_without_patrols"}
}

func (MushroomFields) String() string {
	return "mushroom_island"
}

func (MushroomFields) EncodeBiome() int {
	return 14
}
