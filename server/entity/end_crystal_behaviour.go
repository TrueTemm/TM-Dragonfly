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
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type endCrystalBehaviour struct {
	showBase      bool
	beamTarget    cube.Pos
	hasBeamTarget bool
	explosionSize float64
}

func (b endCrystalBehaviour) Apply(data *world.EntityData) {
	data.Data = b
}

func (endCrystalBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	if tx.World().Dimension() == world.End {
		pos := cube.PosFromVec3(e.Position())
		if _, air := tx.Block(pos.Side(cube.FaceDown)).(block.Air); !air {
			if _, air := tx.Block(pos).(block.Air); air {
				fire := block.Fire{}
				tx.SetBlock(pos, fire, nil)
				tx.ScheduleBlockUpdate(pos, fire, time.Duration(30+rand.IntN(10))*time.Second/20)
			}
		}
	}
	return nil
}

func (b endCrystalBehaviour) Explode(e *Ent, _ world.ExplosionSource, impact float64) {
	if impact <= 0 {
		return
	}
	explodeEndCrystal(e, b.explosionSize)
}

func (b endCrystalBehaviour) Hurt(e *Ent, damage float64, src world.DamageSource) (float64, bool) {
	damage = max(damage, 0)
	if _, ok := src.(VoidDamageSource); ok {
		_ = e.Close()
		return damage, true
	}
	explodeEndCrystal(e, b.explosionSize)
	return damage, true
}

func (endCrystalBehaviour) Immobile() bool {
	return true
}

func (b endCrystalBehaviour) ShowBase() bool {
	return b.showBase
}

func (b endCrystalBehaviour) BeamTarget() (cube.Pos, bool) {
	return b.beamTarget, b.hasBeamTarget
}

func explodeEndCrystal(e *Ent, explosionSize float64) {
	if _, ok := e.H().Entity(e.tx); !ok {
		return
	}
	_ = e.Close()
	block.ExplosionConfig{
		SuppressUnderwaterImpact: true,
	}.Explode(e.tx, world.EntityExplosionSource{
		Entity:        e,
		ExplosionSize: explosionSize,
	})
}
