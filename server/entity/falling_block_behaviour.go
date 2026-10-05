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
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"math/rand/v2"
)

type FallingBlockBehaviourConfig struct {
	Block world.Block

	Gravity float64

	Drag float64

	DistanceFallen float64
}

func (conf FallingBlockBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf FallingBlockBehaviourConfig) New() *FallingBlockBehaviour {
	behaviour := &FallingBlockBehaviour{block: conf.Block}
	behaviour.passive = PassiveBehaviourConfig{
		Gravity: conf.Gravity,
		Drag:    conf.Drag,
		Tick:    behaviour.tick,
	}.New()
	behaviour.passive.fallDistance = conf.DistanceFallen
	return behaviour
}

type FallingBlockBehaviour struct {
	passive *PassiveBehaviour
	block   world.Block
}

func (f *FallingBlockBehaviour) Block() world.Block {
	return f.block
}

func (f *FallingBlockBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	return f.passive.Tick(e, tx)
}

func (f *FallingBlockBehaviour) tick(e *Ent, tx *world.Tx) {
	pos := e.Position()
	bpos := cube.PosFromVec3(pos)
	if a, ok := f.block.(Solidifiable); (ok && a.Solidifies(bpos, tx)) || f.passive.mc.OnGround() {
		f.solidify(e, pos, tx)
	}
}

func (f *FallingBlockBehaviour) solidify(e *Ent, pos mgl64.Vec3, tx *world.Tx) {
	bpos := cube.PosFromVec3(pos)

	if d, ok := f.block.(damager); ok {
		f.damageEntities(e, d, pos, tx)
	}
	if l, ok := f.block.(landable); ok {
		l.Landed(tx, bpos)
	}
	f.passive.close = true

	if r, ok := tx.Block(bpos).(replaceable); ok && r.ReplaceableBy(f.block) {
		tx.SetBlock(bpos, f.block, nil)
	} else if i, ok := f.block.(world.Item); ok {
		opts := world.EntitySpawnOpts{Position: bpos.Vec3Middle()}
		tx.AddEntity(NewItem(opts, item.NewStack(i, 1)))
	}
}

func (f *FallingBlockBehaviour) damageEntities(e *Ent, d damager, pos mgl64.Vec3, tx *world.Tx) {
	damagePerBlock, maxDamage := d.Damage()
	dist := math.Ceil(f.passive.fallDistance - 1.0)
	if dist <= 0 {
		return
	}
	dmg := math.Min(math.Floor(dist*damagePerBlock), maxDamage)
	src := block.DamageSource{Block: f.block}

	for e := range filterLiving(tx.EntitiesWithin(e.H().Type().BBox(e).Translate(pos).Grow(0.05))) {
		e.(Living).Hurt(dmg, src)
	}
	if b, ok := f.block.(breakable); ok && dmg > 0.0 && rand.Float64() < (dist+1)*0.05 {
		f.block = b.Break()
	}
}

type Solidifiable interface {
	Solidifies(pos cube.Pos, tx *world.Tx) bool
}

type replaceable interface {
	ReplaceableBy(b world.Block) bool
}

type damager interface {
	Damage() (damagePerBlock, maxDamage float64)
}

type breakable interface {
	Break() world.Block
}

type landable interface {
	Landed(tx *world.Tx, pos cube.Pos)
}
