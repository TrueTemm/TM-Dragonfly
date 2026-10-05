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
	"context"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/portal"
	"github.com/go-gl/mathgl/mgl64"
)

type PortalTravelComputer struct {
	Instantaneous func(source, target world.Dimension) bool

	Teleport func(e Traveller, pos mgl64.Vec3)

	SpawnPoint func(tx *world.Tx) mgl64.Vec3

	Player bool

	CreatePortal bool

	Cooldown time.Duration

	mu             sync.Mutex
	start          time.Time
	cooldownUntil  time.Time
	inside         bool
	awaitingTravel bool
	travelling     bool
	timedOut       bool
	pending        *world.World
}

func NewPortalTravelComputer() *PortalTravelComputer {
	return &PortalTravelComputer{Instantaneous: func(world.Dimension, world.Dimension) bool { return true }, Cooldown: time.Second * 15}
}

const portalSearchRadius = 128

type portalTravelComputerProvider interface {
	PortalTravelComputer() *PortalTravelComputer
}

type Traveller interface {
	world.Entity

	Teleport(pos mgl64.Vec3)
}

type portalTravelHandler interface {
	HandlePortalTravel(source, destination world.Dimension)
}

func (t *PortalTravelComputer) EnterPortal(e Traveller, tx *world.Tx, target world.Dimension) {
	if destination := t.enterPortal(tx, target); destination != nil {
		t.travelQueued(e, tx, destination)
	}
}

func (t *PortalTravelComputer) queuePortalTravel(tx *world.Tx, target world.Dimension) {
	if destination := t.enterPortal(tx, target); destination != nil {
		t.mu.Lock()
		t.pending = destination
		t.mu.Unlock()
	}
}

func (t *PortalTravelComputer) enterPortal(tx *world.Tx, target world.Dimension) *world.World {
	source := tx.World()
	destination := source.PortalDestination(target)
	if destination == source {
		return nil
	}

	t.mu.Lock()
	t.inside = true
	if t.timedOut {

		t.mu.Unlock()
		return nil
	}
	if time.Now().Before(t.cooldownUntil) {
		t.mu.Unlock()
		return nil
	}
	travelNow := t.instantaneous(source.Dimension(), target) || (t.awaitingTravel && time.Since(t.start) >= time.Second*4)
	if !travelNow && !t.awaitingTravel {
		t.start, t.awaitingTravel = time.Now(), true
	}
	t.mu.Unlock()

	if travelNow {
		return destination
	}
	return nil
}

func (t *PortalTravelComputer) instantaneous(source, target world.Dimension) bool {
	return t.Instantaneous != nil && t.Instantaneous(source, target)
}

func (t *PortalTravelComputer) hasPendingPortalTravel() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.pending != nil
}

func (t *PortalTravelComputer) finishPendingPortalTravel(e Traveller, tx *world.Tx) bool {
	t.mu.Lock()
	destination := t.pending
	t.pending = nil
	t.mu.Unlock()

	if destination == nil {
		return false
	}
	t.travel(e, tx, destination)
	return true
}

func (t *PortalTravelComputer) StopPortalContact() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.inside {
		t.inside = false
		return
	}
	if t.travelling || t.pending != nil {
		return
	}
	t.timedOut, t.awaitingTravel = false, false
}

func (t *PortalTravelComputer) travel(e Traveller, tx *world.Tx, destination *world.World) {
	source := tx.World()
	if destination == nil || destination == source {
		return
	}

	sourceDim, destinationDim := source.Dimension(), destination.Dimension()
	origin := e.Position()
	pos := translatePortalPosition(cube.PosFromVec3(origin), sourceDim, destinationDim)

	t.mu.Lock()
	t.travelling, t.timedOut, t.awaitingTravel = true, true, false
	t.mu.Unlock()

	handle := tx.RemoveEntity(e)
	if handle == nil {
		t.mu.Lock()
		t.travelling, t.timedOut = false, false
		t.mu.Unlock()
		return
	}

	go t.transfer(handle, source, destination, origin, pos, sourceDim, destinationDim)
}

func (t *PortalTravelComputer) travelQueued(e Traveller, tx *world.Tx, destination *world.World) {
	source := tx.World()
	if destination == nil || destination == source {
		return
	}

	sourceDim, destinationDim := source.Dimension(), destination.Dimension()
	origin := e.Position()
	pos := translatePortalPosition(cube.PosFromVec3(origin), sourceDim, destinationDim)

	t.mu.Lock()
	t.travelling, t.timedOut, t.awaitingTravel = true, true, false
	t.mu.Unlock()

	h := e.H()
	tx.Defer(func(tx *world.Tx) {

		e, ok := h.Entity(tx)
		if !ok {
			t.mu.Lock()
			t.travelling, t.timedOut = false, false
			t.mu.Unlock()
			return
		}
		handle := tx.RemoveEntity(e)
		if handle == nil {
			t.mu.Lock()
			t.travelling, t.timedOut = false, false
			t.mu.Unlock()
			return
		}
		go t.transfer(handle, source, destination, origin, pos, sourceDim, destinationDim)
	})
}

func (t *PortalTravelComputer) transfer(handle *world.EntityHandle, source, destination *world.World, origin mgl64.Vec3, pos cube.Pos, sourceDim, destinationDim world.Dimension) {
	travelled, err := world.Call(context.Background(), destination, func(tx *world.Tx) (bool, error) {
		spawn, ok := t.destinationSpawn(tx, sourceDim, pos)
		if !ok {
			return false, nil
		}
		if e, ok := tx.AddEntityAt(handle, spawn).(Traveller); ok {
			t.finishTravel(e, spawn, sourceDim, destinationDim)
		}
		return true, nil
	})
	if err != nil {
		travelled = false
	}
	if !travelled {
		_, err = world.Call(context.Background(), source, func(tx *world.Tx) (struct{}, error) {
			tx.AddEntityAt(handle, origin)
			return struct{}{}, nil
		})
		if err != nil {
			_ = handle.Close()
		}
	}

	t.mu.Lock()
	t.travelling = false
	t.cooldownUntil = time.Now().Add(t.Cooldown)
	if !travelled {

		t.timedOut = false
	}
	t.mu.Unlock()
}

func (t *PortalTravelComputer) destinationSpawn(tx *world.Tx, sourceDim world.Dimension, pos cube.Pos) (mgl64.Vec3, bool) {
	if tx.World().Dimension() == world.End {
		portal.GenerateEndSpawnPlatform(tx)
		return portal.EndSpawnPosition(t.Player), true
	}
	if sourceDim == world.End && tx.World().Dimension() == world.Overworld {

		if t.SpawnPoint != nil {
			return t.SpawnPoint(tx), true
		}
		return tx.World().Spawn().Vec3Middle(), true
	}
	if !t.CreatePortal {
		n, ok := portal.FindNetherPortal(tx, pos, portalSearchRadius)
		if !ok {
			return mgl64.Vec3{}, false
		}
		return n.Spawn().Vec3Middle(), true
	}
	if n, ok := portal.FindOrCreateNetherPortal(tx, pos, portalSearchRadius); ok {
		return n.Spawn().Vec3Middle(), true
	}
	return mgl64.Vec3{}, false
}

func (t *PortalTravelComputer) finishTravel(e Traveller, pos mgl64.Vec3, source, destination world.Dimension) {
	handlePortalTravel(e, source, destination)
	if t.Teleport != nil {
		t.Teleport(e, pos)
		return
	}
	e.Teleport(pos)
}

func handlePortalTravel(e Traveller, source, destination world.Dimension) {
	if ent, ok := e.(*Ent); ok {
		if h, ok := ent.Behaviour().(portalTravelHandler); ok {
			h.HandlePortalTravel(source, destination)
		}
		return
	}
	if h, ok := e.(portalTravelHandler); ok {
		h.HandlePortalTravel(source, destination)
	}
}

func translatePortalPosition(pos cube.Pos, source, target world.Dimension) cube.Pos {
	switch source {
	case world.Overworld:
		pos[0], pos[2] = pos[0]>>3, pos[2]>>3
	case world.Nether:
		pos[0], pos[2] = pos[0]*8, pos[2]*8
	}
	r := target.Range()
	pos[1] = min(max(pos[1], r.Min()), r.Max())
	return pos
}
