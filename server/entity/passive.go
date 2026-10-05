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

	"github.com/df-mc/dragonfly/server/world"
)

type PassiveBehaviourConfig struct {
	Gravity float64

	Drag float64

	ExistenceDuration time.Duration

	Expire func(e *Ent, tx *world.Tx)

	Tick func(e *Ent, tx *world.Tx)
}

func (conf PassiveBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf PassiveBehaviourConfig) New() *PassiveBehaviour {
	if conf.ExistenceDuration == 0 {
		conf.ExistenceDuration = math.MaxInt64
	}
	return &PassiveBehaviour{
		BaseBehaviour: NewBaseBehaviour(),
		conf:          conf,
		fuse:          conf.ExistenceDuration,
		mc: &MovementComputer{
			Gravity:           conf.Gravity,
			Drag:              conf.Drag,
			DragBeforeGravity: true,
		},
	}
}

type PassiveBehaviour struct {
	BaseBehaviour

	conf PassiveBehaviourConfig
	mc   *MovementComputer

	close        bool
	fallDistance float64
	fuse         time.Duration
}

func (p *PassiveBehaviour) Explode(e *Ent, src world.ExplosionSource, impact float64) {
	e.data.Vel = e.data.Vel.Add(e.data.Pos.Sub(src.Position()).Normalize().Mul(impact))
}

func (p *PassiveBehaviour) Fuse() time.Duration {
	if p.conf.Expire != nil {
		return p.fuse
	}
	return -1
}

func (p *PassiveBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	if p.close {
		_ = e.Close()
		return nil
	}

	m := p.mc.TickMovement(e, e.data.Pos, e.data.Vel, e.data.Rot, tx)
	e.data.Pos, e.data.Vel = m.pos, m.vel
	p.fallDistance = math.Max(p.fallDistance-m.dpos[1], 0)

	p.fuse = p.conf.ExistenceDuration - e.Age()

	if p.conf.Tick != nil {
		p.conf.Tick(e, tx)
	}

	if p.Fuse()%(time.Second/4) == 0 {
		for _, v := range tx.Viewers(m.pos) {
			v.ViewEntityState(e)
		}
	}

	if e.Age() > p.conf.ExistenceDuration {
		p.close = true
		if p.conf.Expire != nil {
			p.conf.Expire(e, tx)
		}
	}
	return m
}
