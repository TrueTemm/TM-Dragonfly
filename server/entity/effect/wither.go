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

var Wither wither

type wither struct {
	nopLasting
}

func (wither) Apply(e world.Entity, eff Effect) {
	interval := max(80>>eff.Level(), 1)
	if eff.Tick()%interval == 0 {
		if l, ok := e.(living); ok {
			l.Hurt(1, WitherDamageSource{})
		}
	}
}

func (wither) RGBA() color.RGBA {
	return color.RGBA{R: 0x73, G: 0x61, B: 0x56, A: 0xff}
}

type WitherDamageSource struct{}

func (WitherDamageSource) ReducedByResistance() bool { return true }
func (WitherDamageSource) ReducedByArmour() bool     { return false }
func (WitherDamageSource) Fire() bool                { return false }
func (WitherDamageSource) IgnoreTotem() bool         { return false }
