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

var Slowness slowness

type slowness struct {
	nopLasting
}

func (slowness) Start(e world.Entity, lvl int) {
	slowness := 1 - float64(lvl)*0.15
	if slowness <= 0 {
		slowness = 0.00001
	}
	if l, ok := e.(living); ok {
		l.SetSpeed(l.Speed() * slowness)
	}
}

func (slowness) End(e world.Entity, lvl int) {
	slowness := 1 - float64(lvl)*0.15
	if slowness <= 0 {
		slowness = 0.00001
	}
	if l, ok := e.(living); ok {
		l.SetSpeed(l.Speed() / slowness)
	}
}

func (slowness) RGBA() color.RGBA {
	return color.RGBA{R: 0x8b, G: 0xaf, B: 0xe0, A: 0xff}
}
