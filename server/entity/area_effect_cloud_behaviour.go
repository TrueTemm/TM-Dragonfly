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
	"iter"
	"time"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type AreaEffectCloudBehaviourConfig struct {
	Potion potion.Potion

	Radius float64

	RadiusUseGrowth float64

	RadiusTickGrowth float64

	Duration time.Duration

	DurationUseGrowth time.Duration

	ReapplicationDelay time.Duration
}

func (conf AreaEffectCloudBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf AreaEffectCloudBehaviourConfig) New() *AreaEffectCloudBehaviour {
	if conf.Radius == 0 {
		conf.Radius = 3.0
	}
	if conf.Duration == 0 {
		conf.Duration = time.Second * 30
	}
	stationary := StationaryBehaviourConfig{ExistenceDuration: conf.Duration}
	return &AreaEffectCloudBehaviour{
		conf:       conf,
		stationary: stationary.New(),
		duration:   conf.Duration,
		radius:     conf.Radius,
		targets:    make(map[*world.EntityHandle]time.Duration),
	}
}

type AreaEffectCloudBehaviour struct {
	conf AreaEffectCloudBehaviourConfig

	stationary *StationaryBehaviour

	duration time.Duration
	radius   float64
	targets  map[*world.EntityHandle]time.Duration
}

func (a *AreaEffectCloudBehaviour) PortalTravelComputer() *PortalTravelComputer {
	return a.stationary.PortalTravelComputer()
}

func (a *AreaEffectCloudBehaviour) Radius() float64 {
	return a.radius
}

func (a *AreaEffectCloudBehaviour) Effects() []effect.Effect {
	return a.conf.Potion.Effects()
}

func (a *AreaEffectCloudBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	a.stationary.Tick(e, tx)
	if a.stationary.close || e.Age() < time.Second/2 {

		return nil
	}

	pos := e.Position()
	if a.subtractTickRadius() {
		for _, v := range tx.Viewers(pos) {
			v.ViewEntityState(e)
		}
	}

	if (e.Age()/(time.Second/20))%10 != 0 {

		return nil
	}

	for target, expiration := range a.targets {
		if e.Age() >= expiration {
			delete(a.targets, target)
		}
	}
	if a.applyEffects(pos, e, a.filter(tx.EntitiesWithin(e.H().Type().BBox(e).Translate(pos)))) {
		for _, v := range tx.Viewers(pos) {
			v.ViewEntityState(e)
		}
	}
	return nil
}

func (a *AreaEffectCloudBehaviour) filter(seq iter.Seq[world.Entity]) iter.Seq[world.Entity] {
	return func(yield func(world.Entity) bool) {
		for e := range seq {
			_, target := a.targets[e.H()]
			_, living := e.(Living)
			if !living || target {
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

func (a *AreaEffectCloudBehaviour) applyEffects(pos mgl64.Vec3, ent *Ent, entities iter.Seq[world.Entity]) bool {
	var update bool
	for e := range entities {
		delta := e.Position().Sub(pos)
		delta[1] = 0
		if delta.Len() <= a.radius {
			l := e.(Living)
			for _, eff := range a.Effects() {
				if lasting, ok := eff.Type().(effect.LastingType); ok {
					l.AddEffect(effect.New(lasting, eff.Level(), eff.Duration()/4))
					continue
				}
				l.AddEffect(eff)
			}

			a.targets[e.H()] = ent.Age() + a.conf.ReapplicationDelay
			a.subtractUseDuration()
			a.subtractUseRadius()

			update = true
		}
	}
	return update
}

func (a *AreaEffectCloudBehaviour) subtractTickRadius() bool {
	a.radius += a.conf.RadiusTickGrowth
	if a.radius < 0.5 {
		a.stationary.close = true
	}
	return a.conf.RadiusTickGrowth != 0
}

func (a *AreaEffectCloudBehaviour) subtractUseDuration() {
	a.duration += a.conf.DurationUseGrowth
	if a.duration <= 0 {
		a.stationary.close = true
	}
}

func (a *AreaEffectCloudBehaviour) subtractUseRadius() {
	a.radius += a.conf.RadiusUseGrowth
	if a.radius <= 0.5 {
		a.stationary.close = true
	}
}
