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

type MushroomFieldShore struct{}

func (MushroomFieldShore) Temperature() float64 {
	return 0.9
}

func (MushroomFieldShore) Rainfall() float64 {
	return 1
}

func (MushroomFieldShore) Depth() float64 {
	return 0
}

func (MushroomFieldShore) Scale() float64 {
	return 0.025
}

func (MushroomFieldShore) WaterColour() color.RGBA {
	return color.RGBA{R: 0x81, G: 0x81, B: 0x93, A: 0xa5}
}

func (MushroomFieldShore) Tags() []string {
	return []string{"mooshroom_island", "overworld", "shore", "spawns_without_patrols"}
}

func (MushroomFieldShore) String() string {
	return "mushroom_island_shore"
}

func (MushroomFieldShore) EncodeBiome() int {
	return 15
}
