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

package effect

import (
	"github.com/df-mc/dragonfly/server/world"
	"image/color"
)

var Resistance resistance

type resistance struct {
	nopLasting
}

func (resistance) Multiplier(e world.DamageSource, lvl int) float64 {
	if !e.ReducedByResistance() {
		return 1
	}
	if v := 1 - 0.2*float64(lvl); v >= 0 {
		return v
	}
	return 0
}

func (resistance) RGBA() color.RGBA {
	return color.RGBA{R: 0x91, G: 0x46, B: 0xf0, A: 0xff}
}
