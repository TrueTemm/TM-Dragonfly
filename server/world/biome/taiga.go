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

type Taiga struct{}

func (Taiga) Temperature() float64 {
	return 0.25
}

func (Taiga) Rainfall() float64 {
	return 0.8
}

func (Taiga) Depth() float64 {
	return 0.1
}

func (Taiga) Scale() float64 {
	return 0.2
}

func (Taiga) WaterColour() color.RGBA {
	return color.RGBA{R: 0x28, G: 0x70, B: 0x82, A: 0xa5}
}

func (Taiga) Tags() []string {
	return []string{"animal", "forest", "monster", "overworld", "taiga", "has_structure_trail_ruins", "spawns_cold_variant_farm_animals"}
}

func (Taiga) String() string {
	return "taiga"
}

func (Taiga) EncodeBiome() int {
	return 5
}
