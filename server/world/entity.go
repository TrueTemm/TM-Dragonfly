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

package world

import (
	"encoding/binary"
	"io"
	"maps"
	"math"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type EntityType interface {
	Open(tx *Tx, handle *EntityHandle, data *EntityData) Entity

	EncodeEntity() string

	BBox(e Entity) cube.BBox

	DecodeNBT(m map[string]any, data *EntityData)

	EncodeNBT(data *EntityData) map[string]any
}

type EntityConfig interface {
	Apply(data *EntityData)
}

type EntityHandle struct {
	id uuid.UUID
	t  EntityType

	cond         *sync.Cond
	worldless    *atomic.Bool
	weakTxActive bool
	w            *World

	worldReady bool

	worldVersion atomic.Uint64

	closed       chan struct{}
	worldChanged chan struct{}
	closeOnce    sync.Once

	data EntityData
}

type EntitySpawnOpts struct {
	Position mgl64.Vec3

	Rotation cube.Rotation

	Velocity mgl64.Vec3

	ID uuid.UUID

	NameTag string
}

func (opts EntitySpawnOpts) New(t EntityType, conf EntityConfig) *EntityHandle {
	if opts.ID == uuid.Nil {

		opts.ID = uuid.New()
		clear(opts.ID[:8])
	}
	handle := &EntityHandle{
		id:        opts.ID,
		t:         t,
		cond:      sync.NewCond(&sync.Mutex{}),
		worldless: &atomic.Bool{},
		closed:    make(chan struct{}),
		data:      EntityData{AlwaysShowNameTag: true},
	}
	handle.worldless.Store(true)
	handle.data.Pos, handle.data.Rot, handle.data.Vel = opts.Position, opts.Rotation, opts.Velocity
	handle.data.Name = opts.NameTag
	conf.Apply(&handle.data)
	return handle
}

func NewEntity(t EntityType, conf EntityConfig) *EntityHandle {
	var opts EntitySpawnOpts
	return opts.New(t, conf)
}

func entityFromData(t EntityType, id int64, data map[string]any) *EntityHandle {
	handle := &EntityHandle{
		t:         t,
		cond:      sync.NewCond(&sync.Mutex{}),
		worldless: &atomic.Bool{},
		closed:    make(chan struct{}),
		data:      EntityData{AlwaysShowNameTag: true},
	}
	binary.LittleEndian.PutUint64(handle.id[8:], uint64(id))
	handle.decodeNBT(data)
	if !decodeSafely(t, data, &handle.data) {
		return nil
	}
	return handle
}

func decodeSafely(t EntityType, data map[string]any, into *EntityData) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	t.DecodeNBT(data, into)
	return true
}

func (e *EntityHandle) Type() EntityType {
	return e.t
}

func (e *EntityHandle) Entity(tx *Tx) (Entity, bool) {
	if e == nil || e.w != tx.World() {
		return nil, false
	}
	return e.t.Open(tx, e, &e.data), true
}

func (e *EntityHandle) mustEntity(tx *Tx) Entity {
	if ent, ok := e.Entity(tx); ok {
		return ent
	}
	panic("can't load entity with Tx of different world")
}

func (e *EntityHandle) UUID() uuid.UUID {
	return e.id
}

func (e *EntityHandle) Closed() bool {
	if e == nil {
		return true
	}
	select {
	case <-e.closed:
		return true
	default:
		return false
	}
}

func (e *EntityHandle) Close() error {
	e.closeOnce.Do(func() {
		e.setAndUnlockWorld(closeWorld)
		close(e.closed)
	})
	return nil
}

func cancelled(c <-chan struct{}) bool {
	if c == nil {
		return false
	}
	select {
	case <-c:
		return true
	default:
		return false
	}
}

func (e *EntityHandle) execWorld(f func(tx *Tx, e Entity), weak bool, cancel <-chan struct{}, allowedCloseWorld *World) bool {
	e.cond.L.Lock()
	for e.w == nil || (e.w != closeWorld && (!e.worldReady || (!weak && e.weakTxActive))) {
		if cancelled(cancel) {
			if weak {
				e.clearWeakTxActiveLocked()
			}
			e.cond.L.Unlock()
			return false
		}

		e.cond.Wait()
	}
	if cancelled(cancel) {
		if weak {
			e.clearWeakTxActiveLocked()
		}
		e.cond.L.Unlock()
		return false
	}

	e.worldless.Store(false)
	if e.w == closeWorld {

		if weak {
			e.clearWeakTxActiveLocked()
		}
		e.cond.L.Unlock()
		return false
	}
	if e.w.closed.Load() && !e.w.closeAcceptingEntityTasks.Load() && e.w != allowedCloseWorld {
		if weak {
			e.clearWeakTxActiveLocked()
		}
		e.cond.L.Unlock()
		return false
	}

	var ran atomic.Bool
	ret := e.weakExec(func(tx *Tx) {
		ent := e.mustEntity(tx)
		ran.Store(true)
		f(tx, ent)
	}, allowedCloseWorld)
	if !ret && e.w != nil && e.w != closeWorld && e.w.closed.Load() && !e.w.closeAcceptingEntityTasks.Load() && e.w != allowedCloseWorld {
		e.clearWeakTxActiveLocked()
		e.cond.L.Unlock()
		return false
	}
	e.cond.L.Unlock()

	if ran.Load() {
		return true
	}
	if !ret {

		return e.execWorld(f, true, cancel, allowedCloseWorld)
	}
	return true
}

func (e *EntityHandle) weakExec(f execFunc, allowedCloseWorld *World) bool {
	e.weakTxActive = true
	w, version := e.w, e.worldVersion.Load()

	c := w.weakExec(func() bool {
		return e.worldVersion.Load() == version && !e.worldless.Load()
	}, e.cond, f, w == allowedCloseWorld)
	for len(c) == 0 && e.w != closeWorld {

		e.cond.Wait()
	}
	if e.w != closeWorld && !<-c {

		return false
	}

	closed := e.w == closeWorld
	e.weakTxActive = false
	e.cond.Broadcast()
	return !closed
}

func (e *EntityHandle) clearWeakTxActiveLocked() {
	if e.weakTxActive {
		e.weakTxActive = false
		e.cond.Broadcast()
	}
}

var closeWorld = &World{}

func (e *EntityHandle) unsetAndLockWorld() {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()

	e.worldless.Store(true)
	e.w = nil
	e.worldReady = false
	e.worldVersion.Add(1)
	e.notifyWorldChangedLocked()
}

func (e *EntityHandle) setAndUnlockWorld(w *World) {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()

	if e.w != nil {
		panic("cannot add entity to new world before removing from old world")
	}
	e.w = w
	e.worldReady = false
	e.worldVersion.Add(1)
	e.notifyWorldChangedLocked()
}

func (e *EntityHandle) markWorldReady(w *World) {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()
	if e.w == w {
		e.worldReady = true
		e.cond.Broadcast()
	}
}

func (e *EntityHandle) notifyWorldChangedLocked() {
	if e.worldChanged != nil {
		close(e.worldChanged)
		e.worldChanged = nil
	}
	e.cond.Broadcast()
}

func (e *EntityHandle) setAndUnlockWorldAt(w *World, pos mgl64.Vec3) {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()

	if e.w != nil {
		panic("cannot add entity to new world before removing from old world")
	}
	e.data.Pos = pos
	e.w = w
	e.cond.Broadcast()
}

func (e *EntityHandle) decodeNBT(m map[string]any) {
	e.data.Pos = readVec3(m, "Pos")
	e.data.Vel = readVec3(m, "Motion")
	e.data.Rot = readRotation(m)
	e.data.Age = time.Duration(readInt16(m, "Age")) * (time.Second / 20)
	e.data.FireDuration = time.Duration(readInt16(m, "Fire")) * time.Second / 20
	e.data.Name, _ = m["NameTag"].(string)
}

func (e *EntityHandle) encodeNBT() map[string]any {
	return map[string]any{
		"Pos":     []float32{float32(e.data.Pos[0]), float32(e.data.Pos[1]), float32(e.data.Pos[2])},
		"Motion":  []float32{float32(e.data.Vel[0]), float32(e.data.Vel[1]), float32(e.data.Vel[2])},
		"Yaw":     float32(e.data.Rot[0]),
		"Pitch":   float32(e.data.Rot[1]),
		"Fire":    int16(e.data.FireDuration.Seconds() * 20),
		"Age":     int16(min(e.data.Age/(time.Second/20), math.MaxInt16)),
		"NameTag": e.data.Name,
	}
}

type EntityData struct {
	Pos, Vel          mgl64.Vec3
	Rot               cube.Rotation
	Name              string
	AlwaysShowNameTag bool
	FireDuration      time.Duration
	Age               time.Duration

	Data any
}

type Entity interface {
	io.Closer

	H() *EntityHandle

	Position() mgl64.Vec3

	Rotation() cube.Rotation
}

type TickerEntity interface {
	Entity

	Tick(tx *Tx, current int64)
}

type SelfViewer interface {
	TickerEntity

	SelfViewing() bool
}

type EntityAction interface {
	EntityAction()
}

type DamageSource interface {
	ReducedByArmour() bool

	ReducedByResistance() bool

	Fire() bool

	IgnoreTotem() bool
}

type HealingSource interface {
	HealingSource()
}

type EntityRegistry struct {
	conf EntityRegistryConfig
	ent  map[string]EntityType
}

type EntityRegistryConfig struct {
	Item               func(opts EntitySpawnOpts, it any) *EntityHandle
	FallingBlock       func(opts EntitySpawnOpts, bl Block) *EntityHandle
	TNT                func(opts EntitySpawnOpts, fuse time.Duration) *EntityHandle
	BottleOfEnchanting func(opts EntitySpawnOpts, owner Entity) *EntityHandle
	Arrow              func(opts EntitySpawnOpts, conf ArrowSpawnConfig) *EntityHandle
	Egg                func(opts EntitySpawnOpts, owner Entity) *EntityHandle
	EndCrystal         func(opts EntitySpawnOpts) *EntityHandle
	EnderPearl         func(opts EntitySpawnOpts, owner Entity) *EntityHandle
	FishingHook        func(opts EntitySpawnOpts, owner Entity) *EntityHandle
	Firework           func(opts EntitySpawnOpts, firework Item, owner Entity, sidewaysVelocityMultiplier, upwardsAcceleration float64, attached bool) *EntityHandle
	LingeringPotion    func(opts EntitySpawnOpts, t any, owner Entity) *EntityHandle
	Snowball           func(opts EntitySpawnOpts, owner Entity) *EntityHandle
	SplashPotion       func(opts EntitySpawnOpts, t any, owner Entity) *EntityHandle
	Lightning          func(opts EntitySpawnOpts) *EntityHandle
}

type ArrowSpawnConfig struct {
	Damage float64

	Owner Entity

	Critical bool

	DisablePickup bool

	ObtainArrowOnPickup bool

	PunchLevel int

	PiercingLevel int

	Tip any
}

func (conf EntityRegistryConfig) New(ent []EntityType) EntityRegistry {
	m := make(map[string]EntityType, len(ent))
	for _, e := range ent {
		name := e.EncodeEntity()
		if _, ok := m[name]; ok {
			panic("cannot register the same entity (" + name + ") twice")
		}
		m[name] = e
	}
	return EntityRegistry{conf: conf, ent: m}
}

func (reg EntityRegistry) Config() EntityRegistryConfig {
	return reg.conf
}

func (reg EntityRegistry) Lookup(name string) (EntityType, bool) {
	t, ok := reg.ent[name]
	return t, ok
}

func (reg EntityRegistry) Types() []EntityType {
	return slices.Collect(maps.Values(reg.ent))
}

func readVec3(x map[string]any, k string) mgl64.Vec3 {
	if i, ok := x[k].([]any); ok {
		if len(i) != 3 {
			return mgl64.Vec3{}
		}
		var v mgl64.Vec3
		for index, f := range i {
			f32, _ := f.(float32)
			v[index] = float64(f32)
		}
		return v
	} else if i, ok := x[k].([]float32); ok {
		if len(i) != 3 {
			return mgl64.Vec3{}
		}
		return mgl64.Vec3{float64(i[0]), float64(i[1]), float64(i[2])}
	}
	return mgl64.Vec3{}
}

func readFloat32(m map[string]any, k string) float32 {
	v, _ := m[k].(float32)
	return v
}

func readRotation(m map[string]any) cube.Rotation {
	return cube.Rotation{float64(readFloat32(m, "Yaw")), float64(readFloat32(m, "Pitch"))}
}

func readInt16(m map[string]any, k string) int16 {
	v, _ := m[k].(int16)
	return v
}
