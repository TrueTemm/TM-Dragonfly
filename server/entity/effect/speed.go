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

var Speed speed

type speed struct {
	nopLasting
}

func (speed) Start(e world.Entity, lvl int) {
	speed := 1 + float64(lvl)*0.2
	if l, ok := e.(living); ok {
		l.SetSpeed(l.Speed() * speed)
	}
}

func (speed) End(e world.Entity, lvl int) {
	speed := 1 + float64(lvl)*0.2
	if l, ok := e.(living); ok {
		l.SetSpeed(l.Speed() / speed)
	}
}

func (speed) RGBA() color.RGBA {
	return color.RGBA{R: 0x33, G: 0xeb, B: 0xff, A: 0xff}
}
