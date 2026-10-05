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

package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"math"
)

type AffectedDamageSource interface {
	world.DamageSource

	AffectedByEnchantment(e item.EnchantmentType) bool
}

type DamageModifier interface {
	Modifier() float64
}

func ProtectionFactor(src world.DamageSource, enchantments []item.Enchantment) float64 {
	f := 0.0
	for _, e := range enchantments {
		t := e.Type()
		modifier, ok := t.(DamageModifier)
		if !ok {
			continue
		}
		reduced := false
		if _, ok := t.(protection); ok && src.ReducedByResistance() {

			reduced = true
		} else if asrc, ok := src.(AffectedDamageSource); ok && asrc.AffectedByEnchantment(t) {
			reduced = true
		}

		if reduced {
			f += float64(e.Level()) * modifier.Modifier()
		}
	}
	return math.Min(f, 0.8)
}
