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

type TaigaHills struct{}

func (TaigaHills) Temperature() float64 {
	return 0.25
}

func (TaigaHills) Rainfall() float64 {
	return 0.8
}

func (TaigaHills) Depth() float64 {
	return 0.45
}

func (TaigaHills) Scale() float64 {
	return 0.3
}

func (TaigaHills) WaterColour() color.RGBA {
	return color.RGBA{R: 0x23, G: 0x65, B: 0x83, A: 0xa5}
}

func (TaigaHills) Tags() []string {
	return []string{"animal", "hills", "monster", "overworld", "forest", "taiga", "spawns_cold_variant_farm_animals"}
}

func (TaigaHills) String() string {
	return "taiga_hills"
}

func (TaigaHills) EncodeBiome() int {
	return 19
}
