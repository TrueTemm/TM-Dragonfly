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

type TaigaMountains struct{}

func (TaigaMountains) Temperature() float64 {
	return 0.25
}

func (TaigaMountains) Rainfall() float64 {
	return 0.8
}

func (TaigaMountains) Depth() float64 {
	return 0.2
}

func (TaigaMountains) Scale() float64 {
	return 0.4
}

func (TaigaMountains) WaterColour() color.RGBA {
	return color.RGBA{R: 0x1e, G: 0x6b, B: 0x82, A: 0xa5}
}

func (TaigaMountains) Tags() []string {
	return []string{"animal", "forest", "monster", "mutated", "taiga", "overworld_generation", "spawns_cold_variant_farm_animals"}
}

func (TaigaMountains) String() string {
	return "taiga_mutated"
}

func (TaigaMountains) EncodeBiome() int {
	return 133
}
