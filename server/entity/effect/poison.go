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

var Poison poison

type poison struct {
	nopLasting
}

func (poison) Apply(e world.Entity, eff Effect) {
	interval := max(25>>(eff.Level()-1), 1)
	if eff.Tick()%interval == 0 {
		if l, ok := e.(living); ok && l.Health() > 1 {
			l.Hurt(1, PoisonDamageSource{})
		}
	}
}

func (poison) RGBA() color.RGBA {
	return color.RGBA{R: 0x87, G: 0xa3, B: 0x63, A: 0xff}
}

type PoisonDamageSource struct {
	Fatal bool
}

func (PoisonDamageSource) ReducedByResistance() bool { return true }
func (PoisonDamageSource) ReducedByArmour() bool     { return false }
func (PoisonDamageSource) Fire() bool                { return false }
func (PoisonDamageSource) IgnoreTotem() bool         { return false }
