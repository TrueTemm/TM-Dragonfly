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
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"iter"
	"math"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type ProjectileBehaviourConfig struct {
	Owner *world.EntityHandle

	Gravity float64

	Drag float64

	Damage float64

	Potion potion.Potion

	KnockBackForceAddend float64

	KnockBackHeightAddend float64

	Particle world.Particle

	ParticleCount int

	Sound world.Sound

	Critical bool

	Hit func(e *Ent, tx *world.Tx, target trace.Result)

	SurviveBlockCollision bool

	BlockCollisionVelocityMultiplier float64

	DisablePickup bool

	PickupItem item.Stack

	CollisionPosition cube.Pos

	PiercingLevel int

	HitForce, HitHeight float64

	PassesThroughEntities bool
}

func (conf ProjectileBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf ProjectileBehaviourConfig) New() *ProjectileBehaviour {
	if conf.ParticleCount == 0 && conf.Particle != nil {
		conf.ParticleCount = 1
	}
	return &ProjectileBehaviour{
		BaseBehaviour: NewBaseBehaviour(),
		conf:          conf,
		collided:      conf.CollisionPosition != cube.Pos{},
		collisionPos:  conf.CollisionPosition,
		mc: &MovementComputer{
			Gravity:           conf.Gravity,
			Drag:              conf.Drag,
			DragBeforeGravity: true,
		},
	}
}

type ProjectileBehaviour struct {
	BaseBehaviour

	conf        ProjectileBehaviourConfig
	mc          *MovementComputer
	ageCollided int
	close       bool

	collisionPos cube.Pos
	collided     bool

	collidedEntities []*world.EntityHandle
	portalTravel     bool
}

func (lt *ProjectileBehaviour) Owner() *world.EntityHandle {
	return lt.conf.Owner
}

func (lt *ProjectileBehaviour) Explode(e *Ent, src world.ExplosionSource, impact float64) {
	e.data.Vel = e.Velocity().Add(e.Position().Sub(src.Position()).Normalize().Mul(impact))
}

func (lt *ProjectileBehaviour) Potion() potion.Potion {
	return lt.conf.Potion
}

func (lt *ProjectileBehaviour) Collided() bool {
	return lt.collided
}

func (lt *ProjectileBehaviour) Critical() bool {
	return lt.conf.Critical && !lt.collided
}

func (lt *ProjectileBehaviour) HandlePortalTravel(world.Dimension, world.Dimension) {
	lt.portalTravel = true
}

func (lt *ProjectileBehaviour) PortalTravel() bool {
	return lt.portalTravel
}

func (lt *ProjectileBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	if lt.close {
		_ = e.Close()
		return nil
	}

	if lt.collided && lt.tickAttached(e, tx) {
		if lt.ageCollided > 1200 {
			lt.close = true
		}
		return nil
	}
	vel := e.Velocity()
	m, result := lt.tickMovement(e, tx)
	e.data.Pos, e.data.Vel, e.data.Rot = m.pos, m.vel, m.rot

	lt.collisionPos, lt.collided, lt.ageCollided = cube.Pos{}, false, 0

	if result == nil {
		return m
	}

	for i := 0; i < lt.conf.ParticleCount; i++ {
		tx.AddParticle(result.Position(), lt.conf.Particle)
	}
	if lt.conf.Sound != nil {
		tx.PlaySound(result.Position(), lt.conf.Sound)
	}

	switch r := result.(type) {
	case trace.EntityResult:
		if lt.conf.Damage >= 0 {
			lt.hitEntity(r.Entity(), e, vel)
		}
		if DamageableEntity(r.Entity()) {
			lt.collidedEntities = append(lt.collidedEntities, r.Entity().H())
		}
	case trace.BlockResult:
		bpos := r.BlockPosition()
		if h, ok := tx.Block(bpos).(block.ProjectileHitter); ok {
			h.ProjectileHit(bpos, tx, e, r.Face())
		}
		if lt.selfHit(e, tx, r, vel) {

			lt.close = true
			return m
		}
		if lt.conf.SurviveBlockCollision {
			lt.hitBlockSurviving(e, r, m, tx)
			return m
		}
		lt.close = true
	}
	if lt.conf.Hit != nil {
		lt.conf.Hit(e, tx, result)
	}

	if len(lt.collidedEntities) > lt.conf.PiercingLevel {
		lt.close = true
	}
	return m
}

var ProjectileKnockback = struct {
	Horizontal, Vertical           float64
	PunchHorizontal, PunchVertical float64
}{Horizontal: 0.62, Vertical: 0.42, PunchHorizontal: 0.45, PunchVertical: 0.08}

var SelfHitRadius = 3.0

func (lt *ProjectileBehaviour) knockBack(l Living, vel mgl64.Vec3) {
	punch := lt.conf.KnockBackForceAddend / enchantment.Punch.KnockBackMultiplier()
	force := ProjectileKnockback.Horizontal + ProjectileKnockback.PunchHorizontal*punch
	height := ProjectileKnockback.Vertical + ProjectileKnockback.PunchVertical*punch
	if lt.conf.HitForce != 0 || lt.conf.HitHeight != 0 {
		force, height = lt.conf.HitForce, lt.conf.HitHeight
	}
	horizontal := mgl64.Vec3{vel[0], 0, vel[2]}
	if horizontal.LenSqr() < 1e-8 {

		l.KnockBack(l.Position(), 0, height)
		return
	}
	l.KnockBack(l.Position().Sub(horizontal), force, height)
}

func (lt *ProjectileBehaviour) selfHit(e *Ent, tx *world.Tx, r trace.BlockResult, vel mgl64.Vec3) bool {
	if lt.conf.Damage <= 0 || lt.conf.Owner == nil {
		return false
	}
	owner, ok := lt.conf.Owner.Entity(tx)
	if !ok || slices.Contains(lt.collidedEntities, owner.H()) {
		return false
	}
	box := owner.H().Type().BBox(owner).Translate(owner.Position()).Grow(SelfHitRadius)
	if !box.Vec3Within(r.Position()) {
		return false
	}
	lt.hitEntity(owner, e, vel)
	lt.collidedEntities = append(lt.collidedEntities, owner.H())
	return true
}

func (lt *ProjectileBehaviour) tickAttached(e *Ent, tx *world.Tx) bool {
	boxes := tx.Block(lt.collisionPos).Model().BBox(lt.collisionPos, tx)
	box := e.H().Type().BBox(e).Translate(e.Position())

	for _, bb := range boxes {
		if box.IntersectsWith(bb.Translate(lt.collisionPos.Vec3()).Grow(0.05)) {
			if lt.ageCollided > 5 && !lt.conf.DisablePickup {
				lt.tryPickup(e, tx)
			}
			lt.ageCollided++
			return true
		}
	}
	return false
}

func (lt *ProjectileBehaviour) tryPickup(e *Ent, tx *world.Tx) {
	translated := e.H().Type().BBox(e).Translate(e.Position())
	grown := translated.GrowVec3(mgl64.Vec3{1, 0.5, 1})
	for other := range tx.EntitiesWithin(translated.Grow(2)) {
		if !other.H().Type().BBox(other).Translate(other.Position()).IntersectsWith(grown) {
			continue
		}
		collector, ok := other.(Collector)
		if !ok {
			continue
		}
		if n, ok := collector.Collect(lt.conf.PickupItem); !ok || n == 0 {

			continue
		}

		lt.close = true
		for _, viewer := range tx.Viewers(e.Position()) {
			viewer.ViewEntityAction(e, PickedUpAction{Collector: collector})
		}

		return
	}
}

func (lt *ProjectileBehaviour) hitBlockSurviving(e *Ent, r trace.BlockResult, m *Movement, tx *world.Tx) {

	eps := math.Sqrt(0.1 * (1 - lt.conf.BlockCollisionVelocityMultiplier))
	if mgl64.FloatEqualThreshold(e.Velocity().Len(), 0, eps) {
		e.SetVelocity(mgl64.Vec3{})
		lt.collisionPos, lt.collided = r.BlockPosition(), true

		for _, v := range tx.Viewers(m.pos) {
			v.ViewEntityTeleport(e, m.pos)
			v.ViewEntityAction(e, ArrowShakeAction{Duration: time.Millisecond * 350})
			v.ViewEntityState(e)
		}
		return
	}
}

func (lt *ProjectileBehaviour) hitEntity(victim world.Entity, e *Ent, vel mgl64.Vec3) {
	owner, _ := lt.conf.Owner.Entity(e.tx)
	src := ProjectileDamageSource{Projectile: e, Owner: owner}
	dmg := math.Ceil(lt.conf.Damage * vel.Len())
	if lt.conf.Critical {
		dmg += rand.Float64() * dmg / 2
	}

	if _, vulnerable, ok := HurtEntity(victim, dmg, src); ok && vulnerable {
		l, ok := victim.(Living)
		if !ok {
			return
		}
		lt.knockBack(l, vel)

		for _, eff := range lt.conf.Potion.Effects() {
			if lasting, ok := eff.Type().(effect.LastingType); ok {
				l.AddEffect(effect.New(lasting, eff.Level(), time.Duration(float64(eff.Duration())/8)))
				continue
			}
			l.AddEffect(eff)
		}
		if flammable, ok := l.(Flammable); ok && e.OnFireDuration() > 0 {
			flammable.SetOnFire(time.Second * 5)
		}
	}
}

func (lt *ProjectileBehaviour) tickMovement(e *Ent, tx *world.Tx) (*Movement, trace.Result) {
	pos, vel := e.Position(), e.Velocity()
	viewers := tx.Viewers(pos)

	velBefore := vel
	vel = lt.mc.applyHorizontalForces(tx, pos, lt.mc.applyVerticalForces(vel))
	rot := cube.Rotation{
		mgl64.RadToDeg(math.Atan2(vel[0], vel[2])),
		mgl64.RadToDeg(math.Atan2(vel[1], math.Hypot(vel[0], vel[2]))),
	}

	var (
		end = pos.Add(vel)
		hit trace.Result
		ok  bool
	)
	if !finite(end) {

		_ = e.Close()
		return &Movement{v: tx.Viewers(pos), e: e, pos: pos, vel: mgl64.Vec3{}, rot: rot}, nil
	}
	if !mgl64.FloatEqual(end.Sub(pos).LenSqr(), 0) {
		if hit, ok = trace.Perform(pos, end, tx, e.H().Type().BBox(e).Grow(1.0), lt.ignores(e)); ok {
			if _, ok := hit.(trace.BlockResult); ok {

				vel[1] = (vel[1] + lt.mc.Gravity) / (1 - lt.mc.Drag)
				x, y, z := vel.Mul(lt.conf.BlockCollisionVelocityMultiplier).Elem()

				mx, my, mz := hit.Face().Axis().Vec3().Mul(-2).Add(mgl64.Vec3{1, 1, 1}).Elem()

				vel = mgl64.Vec3{x * mx, y * my, z * mz}
			} else if lt.conf.PiercingLevel == 0 {
				vel = zeroVec3
			}
			end = hit.Position()
		}
	}
	return &Movement{v: viewers, e: e, pos: end, vel: vel, dpos: end.Sub(pos), dvel: vel.Sub(velBefore), rot: rot}, hit
}

func finite(v mgl64.Vec3) bool {
	return !math.IsNaN(v[0]) && !math.IsNaN(v[1]) && !math.IsNaN(v[2]) &&
		!math.IsInf(v[0], 0) && !math.IsInf(v[1], 0) && !math.IsInf(v[2], 0)
}

func (lt *ProjectileBehaviour) ignores(e *Ent) trace.EntityFilter {
	return func(seq iter.Seq[world.Entity]) iter.Seq[world.Entity] {
		return func(yield func(world.Entity) bool) {
			if lt.conf.PassesThroughEntities {
				return
			}
			for other := range seq {
				g, ok := other.(interface{ GameMode() world.GameMode })
				spectator := ok && !g.GameMode().HasCollision()
				itself := e.H() == other.H()
				damageable := DamageableEntity(other)
				owner := e.data.Age < time.Second/4 && lt.conf.Owner == other.H()
				collidedEntity := slices.Contains(lt.collidedEntities, other.H())
				if spectator || itself || !damageable || owner || collidedEntity {
					continue
				}
				if !yield(other) {
					return
				}
			}
		}
	}
}
