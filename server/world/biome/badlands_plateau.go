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

type BadlandsPlateau struct{}

func (BadlandsPlateau) Temperature() float64 {
	return 2
}

func (BadlandsPlateau) Rainfall() float64 {
	return 0
}

func (BadlandsPlateau) Depth() float64 {
	return 1.5
}

func (BadlandsPlateau) Scale() float64 {
	return 0.025
}

func (BadlandsPlateau) WaterColour() color.RGBA {
	return color.RGBA{R: 0x55, G: 0x80, B: 0x9e, A: 0xa5}
}

func (BadlandsPlateau) Tags() []string {
	return []string{"animal", "mesa", "monster", "overworld", "plateau", "rare", "spawns_mesa_mobs", "spawns_warm_variant_farm_animals", "surface_mineshaft"}
}

func (BadlandsPlateau) String() string {
	return "mesa_plateau"
}

func (BadlandsPlateau) EncodeBiome() int {
	return 39
}
