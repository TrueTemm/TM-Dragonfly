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

type NetherWastes struct{}

func (NetherWastes) Temperature() float64 {
	return 2
}

func (NetherWastes) Rainfall() float64 {
	return 0
}

func (NetherWastes) Depth() float64 {
	return 0.1
}

func (NetherWastes) Scale() float64 {
	return 0.2
}

func (NetherWastes) WaterColour() color.RGBA {
	return color.RGBA{R: 0x90, G: 0x59, B: 0x57, A: 0xa5}
}

func (NetherWastes) Tags() []string {
	return []string{"nether", "nether_wastes", "spawn_endermen", "spawn_few_piglins", "spawn_ghast", "spawn_magma_cubes", "spawns_nether_mobs", "spawn_zombified_piglin", "spawns_warm_variant_farm_animals"}
}

func (NetherWastes) String() string {
	return "hell"
}

func (NetherWastes) EncodeBiome() int {
	return 8
}
