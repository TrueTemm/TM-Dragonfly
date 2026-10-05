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
	"iter"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
)

type Tx struct {
	w        *World
	closed   bool
	deferred []scheduledTransaction
}

type Context struct {
	*Tx
	cancel bool
}

func newTx(w *World) *Tx {
	return &Tx{w: w}
}

func (tx *Tx) Event() *Context {
	return &Context{Tx: tx}
}

func (ctx *Context) Cancelled() bool { return ctx.cancel }

func (ctx *Context) Cancel() { ctx.cancel = true }

func (tx *Tx) Defer(f func(tx *Tx)) *Task {
	return tx.DeferErr(func(tx *Tx) error {
		f(tx)
		return nil
	})
}

func (tx *Tx) DeferErr(f func(tx *Tx) error) *Task {
	return tx.deferTask(f)
}

func (tx *Tx) Range() cube.Range {
	return tx.World().ra
}

func (tx *Tx) SetBlock(pos cube.Pos, b Block, opts *SetOpts) {
	tx.setBlock(pos, b, opts)
}

func (tx *Tx) SetBlockEntity(pos cube.Pos, b Block) {
	tx.setBlockEntity(pos, b)
}

func (tx *Tx) Block(pos cube.Pos) Block {
	return tx.block(pos)
}

func (tx *Tx) BlockLoaded(pos cube.Pos) (Block, bool) {
	return tx.World().blockLoaded(pos)
}

func (tx *Tx) BlocksWithin(pos cube.Pos, radius int, blocks ...Block) iter.Seq[cube.Pos] {
	return tx.World().blocksWithin(pos, radius, blocks...)
}

func (tx *Tx) Liquid(pos cube.Pos) (Liquid, bool) {
	return tx.liquid(pos)
}

func (tx *Tx) SetLiquid(pos cube.Pos, b Liquid) {
	tx.setLiquid(pos, b)
}

func (tx *Tx) BuildStructure(pos cube.Pos, s Structure) {
	tx.buildStructure(pos, s)
}

func (tx *Tx) ScheduleBlockUpdate(pos cube.Pos, b Block, delay time.Duration) {
	tx.World().scheduleBlockUpdate(pos, b, delay)
}

func (tx *Tx) HighestLightBlocker(x, z int) int {
	return tx.highestLightBlocker(x, z)
}

func (tx *Tx) HighestBlock(x, z int) int {
	return tx.highestBlock(x, z)
}

func (tx *Tx) Light(pos cube.Pos) uint8 {
	return tx.light(pos)
}

func (tx *Tx) SkyLight(pos cube.Pos) uint8 {
	return tx.skyLight(pos)
}

func (tx *Tx) SetBiome(pos cube.Pos, b Biome) {
	tx.setBiome(pos, b)
}

func (tx *Tx) Biome(pos cube.Pos) Biome {
	return tx.biome(pos)
}

func (tx *Tx) Temperature(pos cube.Pos) float64 {
	return tx.temperature(pos)
}

func (tx *Tx) RainingAt(pos cube.Pos) bool {
	return tx.rainingAt(pos)
}

func (tx *Tx) SnowingAt(pos cube.Pos) bool {
	return tx.snowingAt(pos)
}

func (tx *Tx) ThunderingAt(pos cube.Pos) bool {
	return tx.thunderingAt(pos)
}

func (tx *Tx) Raining() bool {
	return tx.World().raining()
}

func (tx *Tx) Thundering() bool {
	return tx.World().thundering()
}

func (tx *Tx) AddParticle(pos mgl64.Vec3, p Particle) {
	tx.World().addParticle(pos, p)
}

func (tx *Tx) PlayEntityAnimation(e Entity, a EntityAnimation) {
	for _, viewer := range tx.World().viewersOf(e.Position()) {
		viewer.ViewEntityAnimation(e, a)
	}
}

func (tx *Tx) PlaySound(pos mgl64.Vec3, s Sound) {
	tx.World().playSound(tx, pos, s)
}

func (tx *Tx) AddEntity(e *EntityHandle) Entity {
	return tx.World().addEntity(tx, e)
}

func (tx *Tx) AddEntityAt(e *EntityHandle, pos mgl64.Vec3) Entity {
	return tx.World().addEntityAt(tx, e, pos)
}

func (tx *Tx) RemoveEntity(e Entity) *EntityHandle {
	return tx.World().removeEntity(e, tx)
}

func (tx *Tx) EntitiesWithin(box cube.BBox) iter.Seq[Entity] {
	return tx.World().entitiesWithin(tx, box)
}

func (tx *Tx) Entities() iter.Seq[Entity] {
	return tx.World().allEntities(tx)
}

func (tx *Tx) Players() iter.Seq[Entity] {
	return tx.World().allPlayers(tx)
}

func (tx *Tx) Viewers(pos mgl64.Vec3) []Viewer {
	return tx.World().viewersOf(pos)
}

func (tx *Tx) Sleepers() iter.Seq[Sleeper] {
	ent := tx.Entities()
	return func(yield func(Sleeper) bool) {
		for e := range ent {
			if sleeper, ok := e.(Sleeper); ok {
				if !yield(sleeper) {
					return
				}
			}
		}
	}
}

func (tx *Tx) BroadcastSleepingIndicator() {
	sleepers := tx.Sleepers()

	var sleeping, allSleepers int
	for s := range sleepers {
		allSleepers++
		if _, ok := s.Sleeping(); ok {
			sleeping++
		}
	}

	for s := range sleepers {
		s.SendSleepingIndicator(sleeping, allSleepers)
	}
}

func (tx *Tx) BroadcastSleepingReminder(sleeper Sleeper) {
	sleepers := tx.Sleepers()

	var notSleeping int
	for s := range sleepers {
		if _, ok := s.Sleeping(); !ok {
			notSleeping++
		}
	}

	for s := range sleepers {
		if _, ok := s.Sleeping(); !ok {
			s.Messaget(chat.MessageSleeping, sleeper.Name(), notSleeping)
		}
	}
}

func (tx *Tx) deferTask(f func(tx *Tx) error) *Task {
	if tx.closed {
		panic("world.Tx: use of transaction after transaction finishes is not permitted")
	}
	task := newTask()
	tx.deferred = append(tx.deferred, scheduledTransaction{task: task, f: f})
	return task
}

func (tx *Tx) World() *World {
	if tx.closed {
		panic("world.Tx: use of transaction after transaction finishes is not permitted")
	}
	return tx.w
}

func (tx *Tx) CurrentTick() int64 {
	w := tx.World()
	w.set.Lock()
	defer w.set.Unlock()
	return w.set.CurrentTick
}

func (tx *Tx) close() {
	tx.closed = true
}

func (tx *Tx) runDeferred() {
	for len(tx.deferred) > 0 {
		deferred := tx.deferred
		tx.deferred = nil
		for _, st := range deferred {
			st.Run(tx.w)
		}
	}
}

type normalTransaction struct {
	c chan struct{}
	f func(tx *Tx)
}

func (ntx normalTransaction) Run(w *World) {
	tx := newTx(w)
	ntx.f(tx)
	tx.close()
	tx.runDeferred()
	close(ntx.c)
}

type weakTransaction struct {
	c     chan bool
	f     func(tx *Tx)
	valid func() bool
	cond  *sync.Cond
}

func (wtx weakTransaction) Run(w *World) {
	valid := wtx.valid == nil || wtx.valid()
	if valid {
		tx := newTx(w)
		wtx.f(tx)
		tx.close()
		tx.runDeferred()
	}

	wtx.cond.L.Lock()
	defer wtx.cond.L.Unlock()

	wtx.c <- valid
	wtx.cond.Broadcast()
}

func (wtx weakTransaction) fail() {
	wtx.cond.L.Lock()
	defer wtx.cond.L.Unlock()
	wtx.c <- false
	wtx.cond.Broadcast()
}
