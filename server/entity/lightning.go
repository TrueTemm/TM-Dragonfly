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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"time"
)

func NewLightning(opts world.EntitySpawnOpts) *world.EntityHandle {
	return NewLightningWithDamage(opts, 5, true, time.Second*8)
}

func NewLightningWithDamage(opts world.EntitySpawnOpts, dmg float64, blockFire bool, entityFireDuration time.Duration) *world.EntityHandle {
	conf := lightningConf
	conf.Tick = (&lightningState{
		Damage:             dmg,
		EntityFireDuration: entityFireDuration,
		BlockFire:          blockFire,
		state:              2,
		lifetime:           rand.IntN(4) + 1,
	}).tick
	return opts.New(LightningType, conf)
}

var lightningConf = StationaryBehaviourConfig{SpawnSounds: []world.Sound{sound.Explosion{}, sound.Thunder{}}, ExistenceDuration: time.Second}

type lightningState struct {
	Damage             float64
	EntityFireDuration time.Duration
	BlockFire          bool
	state, lifetime    int
}

func (s *lightningState) tick(e *Ent, tx *world.Tx) {
	pos := e.Position()

	if s.state--; s.state < 0 {
		if s.lifetime == 0 {
			_ = e.Close()
		} else if s.state < -rand.IntN(10) {
			s.lifetime--
			s.state = 1

			if s.BlockFire && tx.World().Difficulty().FireSpreadIncrease() >= 10 {
				s.spreadFire(tx, cube.PosFromVec3(pos))
			}
		}
	}
	if s.state > 0 {
		s.dealDamage(e, tx)
	}
}

func (s *lightningState) dealDamage(e *Ent, tx *world.Tx) {
	pos := e.Position()
	bb := e.H().Type().BBox(e).GrowVec3(mgl64.Vec3{3, 6, 3}).Translate(pos.Add(mgl64.Vec3{0, 3}))
	for e := range tx.EntitiesWithin(bb) {

		if l, ok := e.(Living); ok && l.Health() > 0 {
			if s.Damage > 0 {
				l.Hurt(s.Damage, LightningDamageSource{})
			}
			if f, ok := e.(Flammable); ok && f.OnFireDuration() < s.EntityFireDuration {
				f.SetOnFire(s.EntityFireDuration)
			}
		}
	}
}

func (s *lightningState) spreadFire(tx *world.Tx, pos cube.Pos) {
	s.fire().Start(tx, pos)
	for i := 0; i < 4; i++ {
		pos.Add(cube.Pos{rand.IntN(3) - 1, rand.IntN(3) - 1, rand.IntN(3) - 1})
		s.fire().Start(tx, pos)
	}
}

func (s *lightningState) fire() interface {
	Start(tx *world.Tx, pos cube.Pos)
} {
	return fire().(interface {
		Start(tx *world.Tx, pos cube.Pos)
	})
}

func fire() world.Block {
	f, ok := world.BlockByName("minecraft:fire", map[string]any{"age": int32(0)})
	if !ok {
		panic("could not find fire block")
	}
	return f
}

var LightningType lightningType

type lightningType struct{}

func (t lightningType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}
func (t lightningType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	data.Data = lightningConf.New()
}
func (t lightningType) EncodeNBT(*world.EntityData) map[string]any { return nil }
func (lightningType) EncodeEntity() string                         { return "minecraft:lightning_bolt" }
func (lightningType) BBox(world.Entity) cube.BBox                  { return cube.BBox{} }
