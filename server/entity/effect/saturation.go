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

var Saturation saturation

type saturation struct {
	nopLasting
}

func (saturation) Apply(e world.Entity, eff Effect) {
	if i, ok := e.(interface {
		Saturate(food int, saturation float64)
	}); ok {
		i.Saturate(eff.Level(), 2*float64(eff.Level()))
	}
}

func (saturation) RGBA() color.RGBA {
	return color.RGBA{R: 0xf8, G: 0x24, B: 0x23, A: 0xff}
}
