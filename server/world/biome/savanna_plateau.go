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

type SavannaPlateau struct{}

func (SavannaPlateau) Temperature() float64 {
	return 1
}

func (SavannaPlateau) Rainfall() float64 {
	return 0
}

func (SavannaPlateau) Depth() float64 {
	return 1.5
}

func (SavannaPlateau) Scale() float64 {
	return 0.025
}

func (SavannaPlateau) WaterColour() color.RGBA {
	return color.RGBA{R: 0x25, G: 0x90, B: 0xa8, A: 0xa5}
}

func (SavannaPlateau) Tags() []string {
	return []string{"animal", "monster", "overworld", "plateau", "savanna", "spawns_savanna_mobs", "spawns_warm_variant_farm_animals"}
}

func (SavannaPlateau) String() string {
	return "savanna_plateau"
}

func (SavannaPlateau) EncodeBiome() int {
	return 36
}
