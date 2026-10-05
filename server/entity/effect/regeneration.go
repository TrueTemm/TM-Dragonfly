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

var Regeneration regeneration

type regeneration struct {
	nopLasting
}

func (regeneration) Apply(e world.Entity, eff Effect) {
	interval := max(50>>(eff.Level()-1), 1)
	if eff.Tick()%interval == 0 {
		if l, ok := e.(living); ok {
			l.Heal(1, RegenerationHealingSource{})
		}
	}
}

func (regeneration) RGBA() color.RGBA {
	return color.RGBA{R: 0xcd, G: 0x5c, B: 0xab, A: 0xff}
}

type RegenerationHealingSource struct{}

func (RegenerationHealingSource) HealingSource() {}
