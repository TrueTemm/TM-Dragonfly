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
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Behaviour interface {
	Tick(e *Ent, tx *world.Tx) *Movement
}

type Ent struct {
	tx                *world.Tx
	handle            *world.EntityHandle
	data              *world.EntityData
	deferPortalTravel bool
	once              sync.Once
}

func Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) *Ent {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (e *Ent) H() *world.EntityHandle {
	return e.handle
}

func (e *Ent) Owner() *world.EntityHandle {
	if o, ok := e.Behaviour().(interface{ Owner() *world.EntityHandle }); ok {
		return o.Owner()
	}
	return nil
}

func (e *Ent) Behaviour() Behaviour {
	return e.data.Data.(Behaviour)
}

func (e *Ent) Explode(src world.ExplosionSource, impact float64) {
	if expl, ok := e.Behaviour().(interface {
		Explode(e *Ent, src world.ExplosionSource, impact float64)
	}); ok {
		expl.Explode(e, src, impact)
	}
}

func (e *Ent) Position() mgl64.Vec3 {
	return e.data.Pos
}

func (e *Ent) Velocity() mgl64.Vec3 {
	return e.data.Vel
}

func (e *Ent) SetVelocity(v mgl64.Vec3) {
	e.data.Vel = v
}

func (e *Ent) Teleport(pos mgl64.Vec3) {
	viewers := e.tx.Viewers(e.data.Pos)
	e.data.Pos = pos
	for _, v := range viewers {
		v.ViewEntityTeleport(e, pos)
	}
}

func (e *Ent) Rotation() cube.Rotation {
	return e.data.Rot
}

func (e *Ent) Age() time.Duration {
	return e.data.Age
}

func (e *Ent) OnFireDuration() time.Duration {
	return e.data.FireDuration
}

func (e *Ent) SetOnFire(duration time.Duration) {
	duration = max(duration, 0)
	stateChanged := (e.data.FireDuration > 0) != (duration > 0)

	e.data.FireDuration = duration
	if stateChanged {
		e.updateState()
	}
}

func (e *Ent) Extinguish() {
	e.SetOnFire(0)
}

func (e *Ent) NameTag() string {
	return e.data.Name
}

func (e *Ent) SetNameTag(s string) {
	e.data.Name = s
	e.updateState()
}

func (e *Ent) AlwaysShowNameTag() bool {
	return e.data.AlwaysShowNameTag
}

func (e *Ent) SetAlwaysShowNameTag(alwaysShow bool) {
	e.data.AlwaysShowNameTag = alwaysShow
	e.updateState()
}

func (e *Ent) updateState() {
	for _, v := range e.tx.Viewers(e.data.Pos) {
		v.ViewEntityState(e)
	}
}

func (e *Ent) Tick(tx *world.Tx, current int64) {
	e.deferPortalTravel = true
	defer func() {
		e.deferPortalTravel = false
	}()

	y := e.data.Pos[1]
	if y < float64(tx.Range()[0]) && current%10 == 0 {
		_ = e.Close()
		return
	}
	e.SetOnFire(e.OnFireDuration() - time.Second/20)

	m := e.Behaviour().Tick(e, tx)
	if e.finishPendingPortalTravel(tx) {
		return
	}
	if m != nil {
		m.Send()
	}
	if e.checkPortalInsiders() && e.finishPendingPortalTravel(tx) {
		return
	}
	e.stopPortalContact()
	e.data.Age += time.Second / 20
}

func (e *Ent) Close() error {
	e.once.Do(func() {
		e.tx.RemoveEntity(e)
		_ = e.handle.Close()
	})
	return nil
}

func (e *Ent) TravelThroughPortal(tx *world.Tx, target world.Dimension) {
	if tc := e.portalTravelComputer(); tc != nil {
		if e.deferPortalTravel {
			tc.queuePortalTravel(tx, target)
			return
		}
		tc.EnterPortal(e, tx, target)
	}
}

func (e *Ent) portalTravelComputer() *PortalTravelComputer {
	if b, ok := e.Behaviour().(portalTravelComputerProvider); ok {
		return b.PortalTravelComputer()
	}
	return nil
}

func (e *Ent) stopPortalContact() {
	if tc := e.portalTravelComputer(); tc != nil {
		tc.StopPortalContact()
	}
}

func (e *Ent) pendingPortalTravel() bool {
	if tc := e.portalTravelComputer(); tc != nil {
		return tc.hasPendingPortalTravel()
	}
	return false
}

func (e *Ent) finishPendingPortalTravel(tx *world.Tx) bool {
	if tc := e.portalTravelComputer(); tc != nil {
		return tc.finishPendingPortalTravel(e, tx)
	}
	return false
}

type portalBlock interface {
	Portal() world.Dimension
}

func (e *Ent) checkPortalInsiders() bool {
	box := e.H().Type().BBox(e).Translate(e.Position()).Grow(-0.0001)
	low, high := cube.PosFromVec3(box.Min()), cube.PosFromVec3(box.Max())

	for blockPos := range cube.Range3D(low, high) {
		if p, ok := e.tx.Block(blockPos).(portalBlock); ok {
			e.TravelThroughPortal(e.tx, p.Portal())
			if e.pendingPortalTravel() {
				return true
			}
		}
	}
	return false
}
