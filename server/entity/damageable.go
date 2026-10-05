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

package entity

import "github.com/df-mc/dragonfly/server/world"

type behaviourDamageable interface {
	Hurt(e *Ent, damage float64, src world.DamageSource) (n float64, vulnerable bool)
}

func HurtEntity(e world.Entity, damage float64, src world.DamageSource) (n float64, vulnerable, ok bool) {
	if l, ok := e.(Living); ok {
		n, vulnerable = l.Hurt(damage, src)
		return n, vulnerable, true
	}
	if ent, ok := e.(*Ent); ok {
		if d, ok := ent.Behaviour().(behaviourDamageable); ok {
			n, vulnerable = d.Hurt(ent, damage, src)
			return n, vulnerable, true
		}
	}
	return 0, false, false
}

func DamageableEntity(e world.Entity) bool {
	if _, ok := e.(Living); ok {
		return true
	}
	if ent, ok := e.(*Ent); ok {
		_, ok = ent.Behaviour().(behaviourDamageable)
		return ok
	}
	return false
}
