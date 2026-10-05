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

type LukewarmOcean struct{}

func (LukewarmOcean) Temperature() float64 {
	return 0.5
}

func (LukewarmOcean) Rainfall() float64 {
	return 0.5
}

func (LukewarmOcean) Depth() float64 {
	return -1
}

func (LukewarmOcean) Scale() float64 {
	return 0.1
}

func (LukewarmOcean) WaterColour() color.RGBA {
	return color.RGBA{R: 0x0d, G: 0x96, B: 0xdb, A: 0xa5}
}

func (LukewarmOcean) Tags() []string {
	return []string{"lukewarm", "monster", "ocean", "overworld", "spawns_warm_variant_farm_animals", "fast_fishing", "high_seas"}
}

func (LukewarmOcean) String() string {
	return "lukewarm_ocean"
}

func (LukewarmOcean) EncodeBiome() int {
	return 42
}
