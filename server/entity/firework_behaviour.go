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

import (
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"iter"
	"math"
	"time"
)

type FireworkBehaviourConfig struct {
	Firework item.Firework
	Owner    *world.EntityHandle

	ExistenceDuration time.Duration

	SidewaysVelocityMultiplier float64

	UpwardsAcceleration float64

	Attached bool
}

func (conf FireworkBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf FireworkBehaviourConfig) New() *FireworkBehaviour {
	b := &FireworkBehaviour{conf: conf}
	b.passive = PassiveBehaviourConfig{
		ExistenceDuration: conf.ExistenceDuration,
		Expire:            b.explode,
		Tick:              b.tick,
	}.New()
	return b
}

type FireworkBehaviour struct {
	conf    FireworkBehaviourConfig
	passive *PassiveBehaviour
}

func (f *FireworkBehaviour) PortalTravelComputer() *PortalTravelComputer {
	return f.passive.PortalTravelComputer()
}

func (f *FireworkBehaviour) Firework() item.Firework {
	return f.conf.Firework
}

func (f *FireworkBehaviour) Attached() bool {
	return f.conf.Attached
}

func (f *FireworkBehaviour) Owner() *world.EntityHandle {
	return f.conf.Owner
}

func (f *FireworkBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	return f.passive.Tick(e, tx)
}

func (f *FireworkBehaviour) tick(e *Ent, tx *world.Tx) {
	owner, ok := f.conf.Owner.Entity(tx)
	if f.conf.Attached && ok {

		e.data.Pos = owner.Position()
	} else {
		e.data.Vel[0] *= f.conf.SidewaysVelocityMultiplier
		e.data.Vel[1] += f.conf.UpwardsAcceleration
		e.data.Vel[2] *= f.conf.SidewaysVelocityMultiplier
	}
}

func (f *FireworkBehaviour) explode(e *Ent, tx *world.Tx) {
	owner, _ := f.conf.Owner.Entity(tx)
	pos, explosions := e.Position(), f.conf.Firework.Explosions

	for _, v := range tx.Viewers(pos) {
		v.ViewEntityAction(e, FireworkExplosionAction{})
	}
	for _, explosion := range explosions {
		if explosion.Shape == item.FireworkShapeHugeSphere() {
			tx.PlaySound(pos, sound.FireworkHugeBlast{})
		} else {
			tx.PlaySound(pos, sound.FireworkBlast{})
		}
		if explosion.Twinkle {
			tx.PlaySound(pos, sound.FireworkTwinkle{})
		}
	}

	if len(explosions) == 0 {
		return
	}

	force := float64(len(explosions)*2) + 5.0
	for victim := range filterLiving(tx.EntitiesWithin(e.H().Type().BBox(e).Translate(pos).Grow(5.25))) {
		tpos := victim.Position()
		dist := pos.Sub(tpos).Len()
		if dist > 5.0 {

			continue
		}
		dmg := force * math.Sqrt((5.0-dist)/5.0)
		src := ProjectileDamageSource{Owner: owner, Projectile: e}

		if pos == tpos {
			victim.(Living).Hurt(dmg, src)
			continue
		}
		if _, ok := trace.Perform(pos, tpos, tx, victim.H().Type().BBox(victim).Grow(0.3), nil); ok {
			victim.(Living).Hurt(dmg, src)
		}
	}
}

func filterLiving(seq iter.Seq[world.Entity]) iter.Seq[world.Entity] {
	return func(yield func(world.Entity) bool) {
		for e := range seq {
			if _, living := e.(Living); !living {
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}
