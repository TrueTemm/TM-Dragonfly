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
	"maps"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
)

type ticker struct {
	interval time.Duration
}

func (t ticker) tickLoop(w *World) {
	tc := time.NewTicker(t.interval)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			<-w.exec(t.tick)
		case <-w.closing:

			w.running.Done()
			return
		}
	}
}

func (w *World) AdvanceTick() {
	<-w.exec(ticker{}.tick)
}

func (t ticker) tick(tx *Tx) {
	viewers, loaders := tx.World().allViewers()
	w := tx.World()
	start := time.Now()
	defer func() { w.stats.record(start, time.Since(start), len(w.chunks), len(w.entities)) }()

	w.set.Lock()
	if s := w.set.Spawn; s[1] > tx.Range()[1] && w.Dimension() == Overworld {

		w.set.Spawn[1] = tx.highestObstructingBlock(s[0], s[2]) + 1
	}
	if len(viewers) == 0 && w.set.CurrentTick != 0 && !w.conf.Synchronous {

		w.set.Unlock()
		return
	}
	if w.advance {
		w.set.CurrentTick++
		if w.set.TimeCycle {
			w.set.Time++
		}
		if w.set.WeatherCycle {
			w.advanceWeather()
		}
	}

	rain, thunder, tick, tim, cycle := w.set.Raining, w.set.Thundering && w.set.Raining, w.set.CurrentTick, int(w.set.Time), w.set.TimeCycle

	tryAdvanceDay := false
	if tx.w.set.RequiredSleepTicks > 0 {
		tx.w.set.RequiredSleepTicks--
		tryAdvanceDay = tx.w.set.RequiredSleepTicks <= 0
	}

	w.set.Unlock()

	if tryAdvanceDay {
		t.tryAdvanceDay(tx, cycle)
	}

	if tick%20 == 0 {

		sendWeather := w.Dimension().WeatherCycle() && (rain || thunder || w.lastWeather != [2]bool{rain, thunder})
		w.lastWeather = [2]bool{rain, thunder}
		sendTime := w.Dimension().TimeCycle() && cycle
		if sendTime || sendWeather {
			for _, viewer := range viewers {
				if sendTime {
					viewer.ViewTime(tim)
				}
				if sendWeather {
					viewer.ViewWeather(rain, thunder)
				}
			}
		}
	}
	if thunder {
		w.tickLightning(tx)
	}

	t.tickEntities(tx, tick)
	w.scheduledUpdates.tick(tx, tick)
	t.tickBlocksRandomly(tx, loaders, tick)
	t.performNeighbourUpdates(tx)
	w.redstone.tick(tx, tick)
}

func (t ticker) performNeighbourUpdates(tx *Tx) {
	updates := slices.Clone(tx.World().neighbourUpdates)
	clear(tx.World().neighbourUpdates)
	tx.World().neighbourUpdates = tx.World().neighbourUpdates[:0]

	for _, update := range updates {
		pos, changedNeighbour := update.pos, update.neighbour
		if ticker, ok := tx.Block(pos).(NeighbourUpdateTicker); ok {
			ticker.NeighbourUpdateTick(pos, changedNeighbour, tx)
		}
		if liquid, ok := tx.additionalLiquid(pos); ok {
			if ticker, ok := liquid.(NeighbourUpdateTicker); ok {
				ticker.NeighbourUpdateTick(pos, changedNeighbour, tx)
			}
		}
	}
}

func (t ticker) tickBlocksRandomly(tx *Tx, loaders []*Loader, tick int64) {
	var (
		r             = int32(tx.World().tickRange())
		g             randUint4
		blockEntities []cube.Pos
		randomBlocks  []cube.Pos
	)
	if r == 0 {

		return
	}

	randomTicks := tx.World().conf.RandomTickSpeed > 0

	loaded := make([]ChunkPos, 0, len(loaders))
	if tx.World().conf.Synchronous {
		loaded = slices.Collect(maps.Keys(tx.World().chunks))
	} else {
		for _, loader := range loaders {
			loader.mu.RLock()
			pos := loader.pos
			loader.mu.RUnlock()

			loaded = append(loaded, pos)
		}
	}

	for pos, c := range tx.World().chunks {
		hasBlockEntities := len(c.BlockEntities) > 0
		if !randomTicks && !hasBlockEntities {
			continue
		}
		if !t.anyWithinDistance(pos, loaded, r) {

			continue
		}
		if hasBlockEntities {
			blockEntities = slices.AppendSeq(blockEntities, maps.Keys(c.BlockEntities))
		}
		if !randomTicks {
			continue
		}

		cx, cz := int(pos[0]<<4), int(pos[1]<<4)

		for j := 0; j < tx.World().conf.RandomTickSpeed; j++ {
			x, y, z := g.uint4(tx.World().r), g.uint4(tx.World().r), g.uint4(tx.World().r)

			for i, sub := range c.Sub() {
				if sub.Empty() {

					continue
				}

				if rid := sub.Layers()[0].At(x, y, z); tx.World().conf.Blocks.RandomTickBlock(rid) {
					subY := (i + (tx.Range().Min() >> 4)) << 4
					randomBlocks = append(randomBlocks, cube.Pos{cx + int(x), subY + int(y), cz + int(z)})

					x, y, z = g.uint4(tx.World().r), g.uint4(tx.World().r), g.uint4(tx.World().r)
				}
			}
		}
	}

	for _, pos := range randomBlocks {
		if rb, ok := tx.Block(pos).(RandomTicker); ok {
			rb.RandomTick(pos, tx, tx.World().r)
		}
	}
	for _, pos := range blockEntities {
		if tb, ok := tx.Block(pos).(TickerBlock); ok {
			tb.Tick(tick, pos, tx)
		}
	}
}

func (t ticker) anyWithinDistance(pos ChunkPos, loaded []ChunkPos, r int32) bool {
	for _, chunkPos := range loaded {
		xDiff, zDiff := chunkPos[0]-pos[0], chunkPos[1]-pos[1]
		if (xDiff*xDiff)+(zDiff*zDiff) <= r*r {

			return true
		}
	}
	return false
}

func (t ticker) tickEntities(tx *Tx, tick int64) {
	for handle, lastPos := range tx.World().entities {
		e := handle.mustEntity(tx)
		chunkPos := chunkPosFromVec3(handle.data.Pos)

		c, ok := tx.World().chunks[chunkPos]
		if !ok {
			continue
		}

		if lastPos != chunkPos {

			tx.World().entities[handle] = chunkPos
			c.Entities = append(c.Entities, handle)

			var viewers []Viewer

			if old, ok := tx.World().chunks[lastPos]; ok {
				old.Entities = sliceutil.DeleteVal(old.Entities, handle)
				viewers = old.viewers
			}

			for _, viewer := range viewers {
				if slices.Index(c.viewers, viewer) == -1 {

					viewer.HideEntity(e)
				}
			}
			for _, viewer := range c.viewers {
				if slices.Index(viewers, viewer) == -1 {

					showEntity(e, viewer)
				}
			}
		}

		if tx.World().conf.Synchronous || len(c.viewers) > 0 {
			if te, ok := e.(TickerEntity); ok {
				te.Tick(tx, tick)
			}
		} else if sv, ok := e.(SelfViewer); ok && sv.SelfViewing() {
			sv.Tick(tx, tick)
		}
	}
}

type randUint4 struct {
	x uint64
	n uint8
}

func (g *randUint4) uint4(r *rand.Rand) uint8 {
	if g.n == 0 {
		g.x = r.Uint64()
		g.n = 16
	}
	val := g.x & 0b1111

	g.x >>= 4
	g.n--
	return uint8(val)
}

type scheduledTickQueue struct {
	ticks         []scheduledTick
	furthestTicks map[scheduledTickIndex]int64
	currentTick   int64
}

type scheduledTick struct {
	pos   cube.Pos
	b     Block
	bhash uint64
	t     int64
}

type scheduledTickIndex struct {
	pos  cube.Pos
	hash uint64
}

func newScheduledTickQueue(tick int64) *scheduledTickQueue {
	return &scheduledTickQueue{furthestTicks: make(map[scheduledTickIndex]int64), currentTick: tick}
}

func (queue *scheduledTickQueue) tick(tx *Tx, tick int64) {
	queue.currentTick = tick

	w := tx.World()
	for _, t := range queue.ticks {
		if t.t > tick {
			continue
		}
		b := tx.Block(t.pos)
		if ticker, ok := b.(ScheduledTicker); ok && w.conf.Blocks.BlockHash(b) == t.bhash {
			ticker.ScheduledTick(t.pos, tx, w.r)
		} else if liquid, ok := tx.additionalLiquid(t.pos); ok && w.conf.Blocks.BlockHash(liquid) == t.bhash {
			if ticker, ok := liquid.(ScheduledTicker); ok {
				ticker.ScheduledTick(t.pos, tx, w.r)
			}
		}
	}

	queue.ticks = slices.DeleteFunc(queue.ticks, func(t scheduledTick) bool {
		return t.t <= tick
	})
	maps.DeleteFunc(queue.furthestTicks, func(index scheduledTickIndex, t int64) bool {
		return t <= tick
	})
}

func (queue *scheduledTickQueue) schedule(br BlockRegistry, pos cube.Pos, b Block, delay time.Duration) {
	resTick := queue.currentTick + int64(max(delay/(time.Second/20), 1))
	index := scheduledTickIndex{pos: pos, hash: br.BlockHash(b)}
	if t, ok := queue.furthestTicks[index]; ok && t >= resTick && t > queue.currentTick {
		return
	}
	queue.furthestTicks[index] = resTick
	queue.ticks = append(queue.ticks, scheduledTick{pos: pos, t: resTick, b: b, bhash: index.hash})
}

func (queue *scheduledTickQueue) fromChunk(pos ChunkPos) []scheduledTick {
	m := make([]scheduledTick, 0, 8)
	for _, t := range queue.ticks {
		if pos == chunkPosFromBlockPos(t.pos) {
			m = append(m, t)
		}
	}
	return m
}

func (queue *scheduledTickQueue) removeChunk(pos ChunkPos) {
	queue.ticks = slices.DeleteFunc(queue.ticks, func(tick scheduledTick) bool {
		return chunkPosFromBlockPos(tick.pos) == pos
	})
	maps.DeleteFunc(queue.furthestTicks, func(index scheduledTickIndex, _ int64) bool {
		return chunkPosFromBlockPos(index.pos) == pos
	})
}

func (queue *scheduledTickQueue) add(ticks []scheduledTick) {
	queue.ticks = append(queue.ticks, ticks...)
	for _, t := range ticks {
		index := scheduledTickIndex{pos: t.pos, hash: t.bhash}
		if existing, ok := queue.furthestTicks[index]; ok {
			queue.furthestTicks[index] = max(existing, t.t)
		} else {
			queue.furthestTicks[index] = t.t
		}
	}
}
