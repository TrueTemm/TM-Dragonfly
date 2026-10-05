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
	"math"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type ExperienceOrbBehaviourConfig struct {
	Gravity float64

	Drag float64

	ExistenceDuration time.Duration

	Experience int
}

func (conf ExperienceOrbBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf ExperienceOrbBehaviourConfig) New() *ExperienceOrbBehaviour {
	if conf.Experience == 0 {
		conf.Experience = 1
	}
	if conf.ExistenceDuration == 0 {
		conf.ExistenceDuration = time.Minute * 5
	}
	b := &ExperienceOrbBehaviour{conf: conf, lastSearch: time.Now()}
	b.passive = PassiveBehaviourConfig{
		Gravity:           conf.Gravity,
		Drag:              conf.Drag,
		ExistenceDuration: conf.ExistenceDuration,
		Tick:              b.tick,
	}.New()
	return b
}

type ExperienceOrbBehaviour struct {
	conf    ExperienceOrbBehaviourConfig
	passive *PassiveBehaviour

	lastSearch time.Time
	target     *world.EntityHandle
}

func (exp *ExperienceOrbBehaviour) PortalTravelComputer() *PortalTravelComputer {
	return exp.passive.PortalTravelComputer()
}

func (exp *ExperienceOrbBehaviour) Experience() int {
	return exp.conf.Experience
}

func (exp *ExperienceOrbBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	return exp.passive.Tick(e, tx)
}

var followBox = cube.Box(-8, -8, -8, 8, 8, 8)

func (exp *ExperienceOrbBehaviour) tick(e *Ent, tx *world.Tx) {
	targetEnt, ok := exp.target.Entity(tx)
	target, _ := targetEnt.(experienceCollector)

	pos := e.Position()
	hasTarget := ok && !target.Dead() && pos.Sub(target.Position()).Len() <= 8
	if !hasTarget && time.Since(exp.lastSearch) >= time.Second {
		exp.findTarget(tx, pos)
	} else if hasTarget {
		exp.moveToTarget(e, target)
	}
}

func (exp *ExperienceOrbBehaviour) findTarget(tx *world.Tx, pos mgl64.Vec3) {
	exp.target = nil
	for o := range tx.EntitiesWithin(followBox.Translate(pos)) {
		if ec, ok := o.(experienceCollector); ok && ec.CanCollectExperience() {
			exp.target = o.H()
			break
		}
	}
	exp.lastSearch = time.Now()
}

func (exp *ExperienceOrbBehaviour) moveToTarget(e *Ent, target experienceCollector) {
	pos, dst := e.Position(), target.Position()
	if o, ok := target.(Eyed); ok {
		dst[1] += o.EyeHeight() / 2
	}
	diff := dst.Sub(pos).Mul(0.125)
	if dist := diff.LenSqr(); dist < 1 {
		e.SetVelocity(e.Velocity().Add(diff.Normalize().Mul(0.2 * math.Pow(1-math.Sqrt(dist), 2))))
	}

	if e.H().Type().BBox(e).Translate(pos).IntersectsWith(target.H().Type().BBox(target).Translate(target.Position())) && target.CollectExperience(exp.conf.Experience) {
		_ = e.Close()
	}
}

type experienceCollector interface {
	Living

	CollectExperience(value int) bool

	CanCollectExperience() bool
}
