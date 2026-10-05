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

var InstantDamage instantDamage

type instantDamage struct{}

func (i instantDamage) Apply(e world.Entity, eff Effect) {
	base := 3 << eff.Level()
	if l, ok := e.(living); ok {
		l.Hurt(float64(base)*eff.potency, InstantDamageSource{})
	}
}

func (instantDamage) RGBA() color.RGBA {
	return color.RGBA{R: 0xa9, G: 0x65, B: 0x6a, A: 0xff}
}

type InstantDamageSource struct{}

func (InstantDamageSource) ReducedByArmour() bool     { return false }
func (InstantDamageSource) ReducedByResistance() bool { return true }
func (InstantDamageSource) Fire() bool                { return false }
func (InstantDamageSource) IgnoreTotem() bool         { return false }
