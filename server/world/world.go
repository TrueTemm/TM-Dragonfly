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
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"iter"
	"maps"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/goleveldb/leveldb"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type World struct {
	conf Config
	ra   cube.Range

	queue        chan transaction
	queueClosing chan struct{}
	queueing     sync.WaitGroup

	scheduleMu sync.Mutex
	scheduling sync.WaitGroup

	closed                    atomic.Bool
	closeAcceptingEntityTasks atomic.Bool

	advance bool

	o sync.Once

	set     *Settings
	handler atomic.Pointer[Handler]

	weather

	lastWeather [2]bool

	closeStarted chan struct{}
	closing      chan struct{}
	running      sync.WaitGroup

	chunks        map[ChunkPos]*Column
	chunkRequests map[ChunkPos]*chunkRequest
	chunkWorkers  *chunkWorkerPool

	entities map[*EntityHandle]ChunkPos

	r *rand.Rand

	scheduledUpdates *scheduledTickQueue
	redstone         *redstoneEngine
	neighbourUpdates []neighbourUpdate

	viewerMu sync.Mutex
	viewers  map[*Loader]Viewer

	stats tickStats
}

type transaction interface {
	Run(w *World)
}

func New() *World {
	var conf Config
	return conf.New()
}

func (w *World) ReadOnly() bool {
	return w.conf.ReadOnly
}

func (w *World) Name() string {
	w.set.Lock()
	defer w.set.Unlock()
	return w.set.Name
}

func (w *World) Dimension() Dimension {
	return w.conf.Dim
}

func (w *World) Range() cube.Range {
	return w.ra
}

func (w *World) BlockRegistry() BlockRegistry {
	return w.conf.Blocks
}

type execFunc func(tx *Tx)

func (w *World) exec(f execFunc) <-chan struct{} {
	c := make(chan struct{})
	ntx := normalTransaction{c: c, f: f}
	if w.conf.Synchronous {
		ntx.Run(w)
		return c
	}
	w.queue <- ntx
	return c
}

func (w *World) weakExec(valid func() bool, cond *sync.Cond, f execFunc, allowClosed bool) <-chan bool {
	c := make(chan bool, 1)
	if w.conf.Synchronous {
		run := valid == nil || valid()
		if run {

			cond.L.Unlock()
			tx := newTx(w)
			f(tx)
			tx.close()
			tx.runDeferred()
			cond.L.Lock()
		}
		c <- run
		return c
	}
	w.scheduleMu.Lock()
	if w.closed.Load() && !w.closeAcceptingEntityTasks.Load() && !allowClosed {
		w.scheduleMu.Unlock()
		c <- false
		return c
	}
	wtx := weakTransaction{c: c, f: f, valid: valid, cond: cond}
	select {
	case w.queue <- wtx:
		w.scheduleMu.Unlock()
	default:
		w.scheduling.Add(1)
		w.scheduleMu.Unlock()
		go func() {
			defer w.scheduling.Done()
			select {
			case w.queue <- wtx:
			case <-w.closing:
				wtx.fail()
			case <-w.queueClosing:
				wtx.fail()
			}
		}()
	}
	return c
}

func (w *World) handleTransactions() {
	for {
		select {
		case tx := <-w.queue:
			tx.Run(w)
		case <-w.queueClosing:
			w.queueing.Done()
			return
		}
	}
}

func (w *World) EntityRegistry() EntityRegistry {
	return w.conf.Entities
}

func (tx *Tx) block(pos cube.Pos) Block {
	return tx.World().blockInChunk(tx.chunk(chunkPosFromBlockPos(pos)), pos)
}

func (w *World) blockLoaded(pos cube.Pos) (Block, bool) {
	if pos.OutOfBounds(w.ra) {
		return w.conf.Blocks.Air(), false
	}
	c, ok := w.chunks[chunkPosFromBlockPos(pos)]
	if !ok {
		return w.conf.Blocks.Air(), false
	}
	rid := c.Block(uint8(pos[0]), int16(pos[1]), uint8(pos[2]), 0)
	if w.conf.Blocks.NBTBlock(rid) {
		if b, ok := c.BlockEntities[pos]; ok {
			return b, true
		}
	}
	return w.conf.Blocks.BlockByRuntimeIDOrAir(rid), true
}

func (w *World) blockInChunk(c *Column, pos cube.Pos) Block {
	if pos.OutOfBounds(w.ra) {

		return w.conf.Blocks.Air()
	}
	rid := c.Block(uint8(pos[0]), int16(pos[1]), uint8(pos[2]), 0)
	if w.conf.Blocks.NBTBlock(rid) {

		if b, ok := c.BlockEntities[pos]; ok {
			return b
		}

		nbtB := w.conf.Blocks.BlockByRuntimeIDOrAir(rid).(NBTer).DecodeNBT(map[string]any{}).(Block)
		c.BlockEntities[pos] = nbtB
		c.net = nil
		for _, v := range c.viewers {
			v.ViewBlockUpdate(pos, nbtB, 0)
		}
		return nbtB
	}
	return w.conf.Blocks.BlockByRuntimeIDOrAir(rid)
}

func (tx *Tx) biome(pos cube.Pos) Biome {
	if pos.OutOfBounds(tx.Range()) {

		return ocean()
	}
	id := int(tx.chunk(chunkPosFromBlockPos(pos)).Biome(uint8(pos[0]), int16(pos[1]), uint8(pos[2])))
	b, ok := BiomeByID(id)
	if !ok {
		tx.World().conf.Log.Error("biome not found by ID", "ID", id)
		return unknownBiome{id: id}
	}
	return b
}

func (w *World) HighestLightBlocker(x, z int) int {
	y, _ := Call(context.Background(), w, func(tx *Tx) (int, error) {
		return tx.highestLightBlocker(x, z), nil
	})
	return y
}

func (tx *Tx) highestLightBlocker(x, z int) int {
	return int(tx.chunk(ChunkPos{int32(x >> 4), int32(z >> 4)}).HighestLightBlocker(uint8(x), uint8(z)))
}

func (tx *Tx) highestBlock(x, z int) int {
	return int(tx.chunk(ChunkPos{int32(x >> 4), int32(z >> 4)}).HighestBlock(uint8(x), uint8(z)))
}

func (tx *Tx) highestObstructingBlock(x, z int) int {
	yHigh := tx.highestBlock(x, z)
	src := worldSource{tx: tx}
	for y := yHigh; y >= tx.Range()[0]; y-- {
		pos := cube.Pos{x, y, z}
		m := tx.block(pos).Model()
		if m.FaceSolid(pos, cube.FaceUp, src) || m.FaceSolid(pos, cube.FaceDown, src) {
			return y
		}
	}
	return tx.Range()[0]
}

type SetOpts struct {
	DisableBlockUpdates bool

	DisableLiquidDisplacement bool

	DisableRedstoneUpdates bool
}

func (tx *Tx) setBlock(pos cube.Pos, b Block, opts *SetOpts) {
	w := tx.World()
	if pos.OutOfBounds(w.Range()) {

		return
	}
	if opts == nil {
		opts = &SetOpts{}
	}

	x, y, z := uint8(pos[0]), int16(pos[1]), uint8(pos[2])
	c := tx.chunk(chunkPosFromBlockPos(pos))

	rid := w.conf.Blocks.BlockRuntimeID(b)
	redstoneAfterRelevant := isRedstoneRelevant(b)
	needOldBlock := !opts.DisableRedstoneUpdates || !redstoneAfterRelevant
	needOldRID := needOldBlock || (rid != w.conf.Blocks.AirRuntimeID() && !opts.DisableLiquidDisplacement)

	var oldRID uint32
	if needOldRID {
		oldRID = c.Block(x, y, z, 0)
	}
	var oldBlock Block
	if needOldBlock {
		oldBlock = w.conf.Blocks.BlockByRuntimeIDOrAir(oldRID)
		if w.conf.Blocks.NBTBlock(oldRID) {
			if blockEntity, ok := c.BlockEntities[pos]; ok {
				oldBlock = blockEntity
			}
		}
	}

	var before uint32
	if rid != w.conf.Blocks.AirRuntimeID() && !opts.DisableLiquidDisplacement {
		before = oldRID
	}

	c.markModified()
	c.SetBlock(x, y, z, 0, rid)
	if w.conf.Blocks.NBTBlock(rid) {
		c.BlockEntities[pos] = b
	} else {
		delete(c.BlockEntities, pos)
	}

	viewers := slices.Clone(c.viewers)

	if !opts.DisableLiquidDisplacement {
		var secondLayer Block

		airRID := w.conf.Blocks.AirRuntimeID()
		if rid == airRID {
			if li := c.Block(x, y, z, 1); li != airRID {
				c.SetBlock(x, y, z, 0, li)
				c.SetBlock(x, y, z, 1, airRID)
				secondLayer = w.conf.Blocks.Air()
				b = w.conf.Blocks.BlockByRuntimeIDOrAir(li)
			}
		} else if w.conf.Blocks.LiquidDisplacingBlock(rid) {
			if w.conf.Blocks.LiquidBlock(before) {
				l := w.conf.Blocks.BlockByRuntimeIDOrAir(before)
				if b.(LiquidDisplacer).CanDisplace(l.(Liquid)) {
					c.SetBlock(x, y, z, 1, before)
					secondLayer = l
				}
			}
		} else if li := c.Block(x, y, z, 1); li != airRID {
			c.SetBlock(x, y, z, 1, airRID)
			secondLayer = w.conf.Blocks.Air()
		}

		if secondLayer != nil {
			for _, viewer := range viewers {
				viewer.ViewBlockUpdate(pos, secondLayer, 1)
			}
		}
	}

	if redstoneAfterRelevant || (needOldBlock && isRedstoneRelevant(oldBlock)) {
		w.redstone.forget(pos)
	}

	for _, viewer := range viewers {
		viewer.ViewBlockUpdate(pos, b, 0)
	}

	if !opts.DisableBlockUpdates {
		w.doBlockUpdatesAround(pos)
	}
	if !opts.DisableRedstoneUpdates {
		w.redstone.invalidateAroundBlockChange(pos, oldBlock, b, RedstoneUpdateCauseBlockUpdate, w.Range())
	}
}

func (tx *Tx) setBlockEntity(pos cube.Pos, b Block) {
	w := tx.World()
	if pos.OutOfBounds(w.Range()) {

		return
	}
	c := tx.chunk(chunkPosFromBlockPos(pos))

	rid := w.conf.Blocks.BlockRuntimeID(b)
	if !w.conf.Blocks.NBTBlock(rid) || c.Block(uint8(pos[0]), int16(pos[1]), uint8(pos[2]), 0) != rid {
		tx.setBlock(pos, b, nil)
		return
	}
	c.BlockEntities[pos] = b
	c.markModified()
}

func (tx *Tx) setBiome(pos cube.Pos, b Biome) {
	if pos.OutOfBounds(tx.Range()) {

		return
	}
	c := tx.chunk(chunkPosFromBlockPos(pos))
	c.markModified()
	c.SetBiome(uint8(pos[0]), int16(pos[1]), uint8(pos[2]), uint32(b.EncodeBiome()))
}

func (tx *Tx) buildStructure(pos cube.Pos, s Structure) {
	w := tx.World()
	dim := s.Dimensions()
	width, height, length := dim[0], dim[1], dim[2]
	maxX, maxY, maxZ := pos[0]+width, pos[1]+height, pos[2]+length
	f := func(x, y, z int) Block {
		return tx.block(cube.Pos{pos[0] + x, pos[1] + y, pos[2] + z})
	}

	for chunkX := pos[0] >> 4; chunkX <= maxX>>4; chunkX++ {
		for chunkZ := pos[2] >> 4; chunkZ <= maxZ>>4; chunkZ++ {
			chunkPos := ChunkPos{int32(chunkX), int32(chunkZ)}
			c := tx.chunk(chunkPos)

			baseX, baseZ := chunkX<<4, chunkZ<<4
			for i, sub := range c.Sub() {
				baseY := (i + (w.Range()[0] >> 4)) << 4
				if baseY>>4 < pos[1]>>4 {
					continue
				} else if baseY >= maxY {
					break
				}

				for localY := 0; localY < 16; localY++ {
					yOffset := baseY + localY
					if yOffset > w.Range()[1] || yOffset >= maxY {

						break
					} else if yOffset < w.Range()[0] || yOffset < pos[1] {

						continue
					}
					for localX := 0; localX < 16; localX++ {
						xOffset := baseX + localX
						if xOffset < pos[0] || xOffset >= maxX {
							continue
						}
						for localZ := 0; localZ < 16; localZ++ {
							zOffset := baseZ + localZ
							if zOffset < pos[2] || zOffset >= maxZ {
								continue
							}
							b, liq := s.At(xOffset-pos[0], yOffset-pos[1], zOffset-pos[2], f)
							if b != nil {
								rid := w.conf.Blocks.BlockRuntimeID(b)
								sub.SetBlock(uint8(xOffset), uint8(yOffset), uint8(zOffset), 0, rid)

								nbtPos := cube.Pos{xOffset, yOffset, zOffset}
								if w.conf.Blocks.NBTBlock(rid) {
									c.BlockEntities[nbtPos] = b
								} else {
									delete(c.BlockEntities, nbtPos)
								}
							}
							if liq != nil {
								sub.SetBlock(uint8(xOffset), uint8(yOffset), uint8(zOffset), 1, w.conf.Blocks.BlockRuntimeID(liq))
							} else if len(sub.Layers()) > 1 {
								sub.SetBlock(uint8(xOffset), uint8(yOffset), uint8(zOffset), 1, w.conf.Blocks.AirRuntimeID())
							}
						}
					}
				}
			}
			c.SetBlock(0, 0, 0, 0, c.Block(0, 0, 0, 0))
			c.markModified()

			for _, viewer := range c.viewers {
				viewer.ViewChunk(chunkPos, w.Dimension(), c)
			}
		}
	}
}

func (tx *Tx) liquid(pos cube.Pos) (Liquid, bool) {
	w := tx.World()
	if pos.OutOfBounds(w.Range()) {

		return nil, false
	}
	c := tx.chunk(chunkPosFromBlockPos(pos))
	x, y, z := uint8(pos[0]), int16(pos[1]), uint8(pos[2])

	id := c.Block(x, y, z, 0)
	b, ok := w.conf.Blocks.BlockByRuntimeID(id)
	if !ok {
		w.conf.Log.Error("Liquid: no block with runtime ID", "ID", id)
		return nil, false
	}
	if liq, ok := b.(Liquid); ok {
		return liq, true
	}
	id = c.Block(x, y, z, 1)

	b, ok = w.conf.Blocks.BlockByRuntimeID(id)
	if !ok {
		w.conf.Log.Error("Liquid: no block with runtime ID", "ID", id)
		return nil, false
	}
	liq, ok := b.(Liquid)
	return liq, ok
}

func (tx *Tx) setLiquid(pos cube.Pos, b Liquid) {
	w := tx.World()
	if pos.OutOfBounds(w.Range()) {

		return
	}
	chunkPos := chunkPosFromBlockPos(pos)
	c := tx.chunk(chunkPos)
	if b == nil {
		w.removeLiquids(c, pos)
		w.doBlockUpdatesAround(pos)
		w.redstone.invalidateAround(pos, pos, RedstoneUpdateCauseBlockUpdate, w.Range())
		return
	}
	x, y, z := uint8(pos[0]), int16(pos[1]), uint8(pos[2])
	if !replaceable(w, c, pos, b) {
		if displacer, ok := w.blockInChunk(c, pos).(LiquidDisplacer); !ok || !displacer.CanDisplace(b) {
			return
		}
	}
	rid := w.conf.Blocks.BlockRuntimeID(b)
	if w.removeLiquids(c, pos) {
		c.SetBlock(x, y, z, 0, rid)
		for _, v := range c.viewers {
			v.ViewBlockUpdate(pos, b, 0)
		}
	} else {
		c.SetBlock(x, y, z, 1, rid)
		for _, v := range c.viewers {
			v.ViewBlockUpdate(pos, b, 1)
		}
	}
	c.markModified()

	w.doBlockUpdatesAround(pos)
	w.redstone.invalidateAround(pos, pos, RedstoneUpdateCauseBlockUpdate, w.Range())
}

func (w *World) removeLiquids(c *Column, pos cube.Pos) bool {
	x, y, z := uint8(pos[0]), int16(pos[1]), uint8(pos[2])
	air := w.conf.Blocks.Air()

	noneLeft := false
	if noLeft, changed := w.removeLiquidOnLayer(c.Chunk, x, y, z, 0); noLeft {
		if changed {
			for _, v := range c.viewers {
				v.ViewBlockUpdate(pos, air, 0)
			}
		}
		noneLeft = true
	}
	if _, changed := w.removeLiquidOnLayer(c.Chunk, x, y, z, 1); changed {
		for _, v := range c.viewers {
			v.ViewBlockUpdate(pos, air, 1)
		}
	}
	return noneLeft
}

func (w *World) removeLiquidOnLayer(c *chunk.Chunk, x uint8, y int16, z, layer uint8) (bool, bool) {
	id := c.Block(x, y, z, layer)
	airRID := w.conf.Blocks.AirRuntimeID()

	b, ok := w.conf.Blocks.BlockByRuntimeID(id)
	if !ok {
		w.conf.Log.Error("removeLiquidOnLayer: no block with runtime ID", "ID", id)
		return false, false
	}
	if _, ok := b.(Liquid); ok {
		c.SetBlock(x, y, z, layer, airRID)
		return true, true
	}
	return id == airRID, false
}

func (tx *Tx) additionalLiquid(pos cube.Pos) (Liquid, bool) {
	w := tx.World()
	if pos.OutOfBounds(w.Range()) {

		return nil, false
	}
	c := tx.chunk(chunkPosFromBlockPos(pos))
	id := c.Block(uint8(pos[0]), int16(pos[1]), uint8(pos[2]), 1)

	b, ok := w.conf.Blocks.BlockByRuntimeID(id)
	if !ok {
		w.conf.Log.Error("additionalLiquid: no block with runtime ID", "ID", id)
		return nil, false
	}
	liq, ok := b.(Liquid)
	return liq, ok
}

func (tx *Tx) light(pos cube.Pos) uint8 {
	w := tx.World()
	if pos[1] < w.ra[0] {

		return 0
	}
	if pos[1] > w.ra[1] {

		return 15
	}
	c, ok := w.loadedChunk(chunkPosFromBlockPos(pos))
	if !ok {
		return 0
	}
	return c.Light(uint8(pos[0]), int16(pos[1]), uint8(pos[2]))
}

func (tx *Tx) skyLight(pos cube.Pos) uint8 {
	w := tx.World()
	if pos[1] < w.ra[0] {

		return 0
	}
	if pos[1] > w.ra[1] {

		return 15
	}
	return tx.chunk(chunkPosFromBlockPos(pos)).SkyLight(uint8(pos[0]), int16(pos[1]), uint8(pos[2]))
}

func (w *World) Time() int {
	if w == nil {
		return 0
	}
	w.set.Lock()
	defer w.set.Unlock()
	return int(w.set.Time)
}

func (w *World) SetTime(new int) {
	if w == nil {
		return
	}
	w.set.Lock()
	w.set.Time = int64(new)
	w.set.Unlock()

	viewers, _ := w.allViewers()
	for _, viewer := range viewers {
		viewer.ViewTime(new)
	}
}

func (w *World) StopTime() {
	w.enableTimeCycle(false)
}

func (w *World) StartTime() {
	w.enableTimeCycle(true)
}

func (w *World) TimeCycle() bool {
	if w == nil {
		return false
	}
	w.set.Lock()
	defer w.set.Unlock()
	return w.set.TimeCycle
}

func (w *World) enableTimeCycle(v bool) {
	if w == nil {
		return
	}
	w.set.Lock()
	defer w.set.Unlock()
	w.set.TimeCycle = v
	viewers, _ := w.allViewers()
	for _, viewer := range viewers {
		viewer.ViewTimeCycle(v)
	}
}

func (tx *Tx) temperature(pos cube.Pos) float64 {
	const (
		tempDrop = 1.0 / 600
		seaLevel = 64
	)
	diff := max(pos[1]-seaLevel, 0)
	return tx.biome(pos).Temperature() - float64(diff)*tempDrop
}

func (w *World) addParticle(pos mgl64.Vec3, p Particle) {
	p.Spawn(w, pos)
	for _, viewer := range w.viewersOf(pos) {
		viewer.ViewParticle(pos, p)
	}
}

func (w *World) playSound(tx *Tx, pos mgl64.Vec3, s Sound) {
	ctx := tx.Event()
	if w.Handler().HandleSound(ctx, s, pos); ctx.Cancelled() {
		return
	}
	s.Play(w, pos)
	for _, viewer := range w.viewersOf(pos) {
		viewer.ViewSound(pos, s)
	}
}

func (w *World) addEntity(tx *Tx, handle *EntityHandle) Entity {
	return w.addEntityAt(tx, handle, handle.data.Pos)
}

func (w *World) addEntityAt(tx *Tx, handle *EntityHandle, pos mgl64.Vec3) Entity {
	handle.setAndUnlockWorldAt(w, pos)
	chunkPos := chunkPosFromVec3(handle.data.Pos)
	w.entities[handle] = chunkPos

	c := tx.chunk(chunkPos)
	c.Entities, c.modified = append(c.Entities, handle), true

	e := handle.mustEntity(tx)
	for _, v := range c.viewers {

		showEntity(e, v)
	}
	w.Handler().HandleEntitySpawn(tx, e)
	handle.markWorldReady(w)
	return e
}

func (w *World) removeEntity(e Entity, tx *Tx) *EntityHandle {
	handle := e.H()
	pos, found := w.entities[handle]
	if !found {

		return nil
	}
	w.Handler().HandleEntityDespawn(tx, e)

	c := tx.chunk(pos)
	c.Entities, c.modified = sliceutil.DeleteVal(c.Entities, handle), true

	w.removeEntityFromViewLayers(e)
	for _, v := range c.viewers {
		v.HideEntity(e)
	}
	delete(w.entities, handle)
	handle.unsetAndLockWorld()
	return handle
}

func (w *World) removeEntityFromViewLayers(e Entity) {
	if _, ok := e.(viewLayerViewer); ok {
		return
	}
	viewers, _ := w.allViewers()
	for _, viewer := range viewers {
		v, ok := viewer.(viewLayerViewer)
		if !ok || v.ViewLayer() == nil {
			continue
		}
		v.ViewLayer().remove(e)
	}
}

func (w *World) entitiesWithin(tx *Tx, box cube.BBox) iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		minPos, maxPos := chunkPosFromVec3(box.Min()), chunkPosFromVec3(box.Max())

		for x := minPos[0]; x <= maxPos[0]; x++ {
			for z := minPos[1]; z <= maxPos[1]; z++ {
				c, ok := w.chunks[ChunkPos{x, z}]
				if !ok {

					continue
				}
				for _, handle := range slices.Clone(c.Entities) {
					if !box.Vec3Within(handle.data.Pos) {
						continue
					}
					ent, ok := handle.Entity(tx)
					if ok && !yield(ent) {
						return
					}
				}
			}
		}
	}
}

func (w *World) allEntities(tx *Tx) iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		for e := range w.entities {
			if ent := e.mustEntity(tx); !yield(ent) {
				return
			}
		}
	}
}

func (w *World) allPlayers(tx *Tx) iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		for e := range w.entities {
			if e.t.EncodeEntity() == "minecraft:player" {
				if ent := e.mustEntity(tx); !yield(ent) {
					return
				}
			}
		}
	}
}

func (w *World) Spawn() cube.Pos {
	if w == nil {
		return cube.Pos{}
	}

	if w.Dimension() == End {
		return cube.Pos{100, 50}
	} else if w.Dimension() == Nether {
		return cube.Pos{}
	}

	w.set.Lock()
	defer w.set.Unlock()
	return w.set.Spawn
}

func (w *World) SetSpawn(pos cube.Pos) {
	if w == nil {
		return
	}

	if w.Dimension() == Nether || w.Dimension() == End {
		return
	}

	w.set.Lock()
	w.set.Spawn = pos
	w.set.Unlock()

	viewers, _ := w.allViewers()
	for _, viewer := range viewers {
		viewer.ViewWorldSpawn(pos)
	}
}

func (w *World) PlayerSpawn(id uuid.UUID) cube.Pos {
	if w == nil {
		return cube.Pos{}
	}
	pos, exist, err := w.conf.Provider.LoadPlayerSpawnPosition(id)
	if err != nil {
		w.conf.Log.Error("load player spawn: "+err.Error(), "ID", id)
		return w.Spawn()
	}
	if !exist {
		return w.Spawn()
	}
	return pos
}

func (w *World) SetPlayerSpawn(id uuid.UUID, pos cube.Pos) {
	if w == nil {
		return
	}
	if err := w.conf.Provider.SavePlayerSpawnPosition(id, pos); err != nil {
		w.conf.Log.Error("save player spawn: "+err.Error(), "ID", id)
	}
}

func (w *World) SetRequiredSleepDuration(duration time.Duration) {
	if w == nil {
		return
	}
	w.set.Lock()
	defer w.set.Unlock()
	w.set.RequiredSleepTicks = duration.Milliseconds() / 50
}

func (w *World) DefaultGameMode() GameMode {
	if w == nil {
		return GameModeSurvival
	}
	w.set.Lock()
	defer w.set.Unlock()
	return w.set.DefaultGameMode
}

func (w *World) SetTickRange(v int) {
	if w == nil {
		return
	}
	w.set.Lock()
	defer w.set.Unlock()
	w.set.TickRange = int32(v)
}

func (w *World) tickRange() int {
	w.set.Lock()
	defer w.set.Unlock()
	return int(w.set.TickRange)
}

func (w *World) FallDamage() bool {
	if w == nil {
		return false
	}
	w.set.Lock()
	defer w.set.Unlock()
	return w.set.FallDamage
}

func (w *World) SetFallDamage(v bool) {
	if w == nil {
		return
	}
	w.set.Lock()
	defer w.set.Unlock()
	w.set.FallDamage = v
}

func (w *World) SetDefaultGameMode(mode GameMode) {
	if w == nil {
		return
	}
	w.set.Lock()
	defer w.set.Unlock()
	w.set.DefaultGameMode = mode
}

func (w *World) Difficulty() Difficulty {
	if w == nil {
		return DifficultyNormal
	}
	w.set.Lock()
	defer w.set.Unlock()
	return w.set.Difficulty
}

func (w *World) SetDifficulty(d Difficulty) {
	if w == nil {
		return
	}
	w.set.Lock()
	defer w.set.Unlock()
	w.set.Difficulty = d
}

func (w *World) scheduleBlockUpdate(pos cube.Pos, b Block, delay time.Duration) {
	if pos.OutOfBounds(w.Range()) {
		return
	}
	w.scheduledUpdates.schedule(w.conf.Blocks, pos, b, delay)
}

func (w *World) doBlockUpdatesAround(pos cube.Pos) {
	if w == nil || pos.OutOfBounds(w.Range()) {
		return
	}
	changed := pos

	w.updateNeighbour(pos, changed)
	pos.Neighbours(func(pos cube.Pos) {
		w.updateNeighbour(pos, changed)
	}, w.Range())
}

type neighbourUpdate struct {
	pos, neighbour cube.Pos
}

func (w *World) updateNeighbour(pos, changedNeighbour cube.Pos) {
	w.neighbourUpdates = append(w.neighbourUpdates, neighbourUpdate{pos: pos, neighbour: changedNeighbour})
}

func (w *World) Handle(h Handler) {
	if w == nil {
		return
	}
	if h == nil {
		h = NopHandler{}
	}
	w.handler.Store(&h)
}

func (w *World) viewersOf(pos mgl64.Vec3) []Viewer {
	c, ok := w.chunks[chunkPosFromVec3(pos)]
	if !ok {
		return nil
	}
	return c.viewers
}

func (w *World) PortalDestination(dim Dimension) *World {
	if w.conf.PortalDestination == nil {
		return w
	}
	if res := w.conf.PortalDestination(dim); res != nil {
		return res
	}
	return w
}

func (w *World) Save() {
	<-w.exec(w.save(w.saveChunk))
}

func (w *World) save(f func(*Tx, ChunkPos, *Column)) execFunc {
	return func(tx *Tx) {
		if w.conf.ReadOnly {
			return
		}
		w.conf.Log.Debug("Saving chunks in memory to disk...")
		for pos, c := range w.chunks {
			f(tx, pos, c)
		}
		w.conf.Log.Debug("Updating level.dat values...")
		w.conf.Provider.SaveSettings(w.set)
	}
}

func (w *World) saveChunk(_ *Tx, pos ChunkPos, c *Column) {
	if !w.conf.ReadOnly && c.modified {
		c.Compact()
		if err := w.conf.Provider.StoreColumn(pos, w.conf.Dim, w.columnTo(c, pos)); err != nil {
			w.conf.Log.Error("save chunk: "+err.Error(), "X", pos[0], "Z", pos[1])
		}
	}
}

func (w *World) closeChunk(tx *Tx, pos ChunkPos, c *Column) {
	w.saveChunk(tx, pos, c)
	w.scheduledUpdates.removeChunk(pos)
	w.redstone.removeChunk(pos)

	for _, e := range slices.Clone(c.Entities) {
		_ = e.mustEntity(tx).Close()
	}
	clear(c.Entities)
	delete(w.chunks, pos)
}

func (w *World) Close() error {
	w.o.Do(w.close)
	return nil
}

func (w *World) close() {
	w.scheduleMu.Lock()
	w.closed.Store(true)
	close(w.closeStarted)
	w.scheduleMu.Unlock()

	w.scheduling.Wait()
	w.scheduleMu.Lock()
	w.closeAcceptingEntityTasks.Store(true)
	w.scheduleMu.Unlock()
	<-w.exec(func(tx *Tx) {

		w.Handler().HandleClose(tx)
		tx.runDeferred()
		w.Handle(NopHandler{})

		w.save(w.closeChunk)(tx)
	})
	w.scheduleMu.Lock()
	w.closeAcceptingEntityTasks.Store(false)
	w.scheduleMu.Unlock()
	w.scheduling.Wait()

	close(w.closing)
	w.running.Wait()
	w.chunkWorkers.wg.Wait()

	close(w.queueClosing)
	w.queueing.Wait()

	if w.set.ref.Add(-1); !w.advance {
		return
	}
	w.conf.Log.Debug("Closing provider...")
	if err := w.conf.Provider.Close(); err != nil {
		w.conf.Log.Error("close world provider: " + err.Error())
	}
}

func (w *World) allViewers() ([]Viewer, []*Loader) {
	w.viewerMu.Lock()
	defer w.viewerMu.Unlock()

	viewers, loaders := make([]Viewer, 0, len(w.viewers)), make([]*Loader, 0, len(w.viewers))
	for k, v := range w.viewers {
		viewers = append(viewers, v)
		loaders = append(loaders, k)
	}
	return viewers, loaders
}

func (w *World) addWorldViewer(l *Loader) {
	w.viewerMu.Lock()
	w.viewers[l] = l.viewer
	w.viewerMu.Unlock()

	l.viewer.ViewTime(w.Time())
	l.viewer.ViewTimeCycle(w.TimeCycle())
	w.set.Lock()
	raining, thundering := w.set.Raining, w.set.Raining && w.set.Thundering
	w.set.Unlock()
	l.viewer.ViewWeather(raining, thundering)
	l.viewer.ViewWorldSpawn(w.Spawn())
}

func (w *World) addViewer(tx *Tx, c *Column, loader *Loader) {
	c.viewers = append(c.viewers, loader.viewer)
	c.loaders = append(c.loaders, loader)

	for _, entity := range c.Entities {
		showEntity(entity.mustEntity(tx), loader.viewer)
	}
}

func (w *World) removeViewer(tx *Tx, pos ChunkPos, loader *Loader) {
	if w == nil {
		return
	}
	c, ok := w.chunks[pos]
	if !ok {
		return
	}
	if i := slices.Index(c.loaders, loader); i != -1 {
		c.viewers = slices.Delete(c.viewers, i, i+1)
		c.loaders = slices.Delete(c.loaders, i, i+1)
	}

	for _, entity := range c.Entities {
		loader.viewer.HideEntity(entity.mustEntity(tx))
	}
}

func (w *World) Handler() Handler {
	if w == nil {
		return NopHandler{}
	}
	return *w.handler.Load()
}

func showEntity(e Entity, viewer Viewer) {
	viewer.ViewEntity(e)
	viewer.ViewEntityItems(e)
	viewer.ViewEntityArmour(e)
}

func (w *World) loadedChunk(pos ChunkPos) (*Column, bool) {
	c, ok := w.chunks[pos]
	return c, ok
}

func (tx *Tx) chunk(pos ChunkPos) *Column {
	w := tx.World()
	c, ok := w.chunks[pos]
	if ok {
		return c
	}
	c, ok = w.chunkFromAsyncPool(tx, pos)
	if ok {
		return c
	}
	col, err := w.loadChunk(pos)
	if err != nil {
		w.conf.Log.Error("load chunk: "+err.Error(), "X", pos[0], "Z", pos[1])
	}
	if col == nil {
		return w.emptyColumn()
	}
	return w.addChunk(pos, col)
}

func (w *World) loadChunk(pos ChunkPos) (*chunk.Column, error) {
	column, err := w.conf.Provider.LoadColumn(pos, w.conf.Dim)
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return nil, err
		}
		ch := chunk.New(w.conf.Blocks, w.Range())
		w.conf.Generator.GenerateChunk(pos, ch)
		column = &chunk.Column{Chunk: ch}
	}
	chunk.LightArea([]*chunk.Chunk{column.Chunk}, int(pos[0]), int(pos[1])).Fill()
	return column, nil
}

func (w *World) emptyColumn() *Column {
	return w.columnFrom(&chunk.Column{Chunk: chunk.New(w.conf.Blocks, w.Range())}, ChunkPos{})
}

func (w *World) loadChunkAsync(tx *Tx, pos ChunkPos, callback chunkCallback) bool {
	if c, ok := w.chunks[pos]; ok {
		callback(tx, c)
		return true
	}
	if w.conf.Synchronous {

		callback(tx, tx.chunk(pos))
		return true
	}
	if req, ok := w.chunkRequests[pos]; ok {
		req.callbacks = append(req.callbacks, callback)
		return true
	}
	req := &chunkRequest{pos: pos, done: make(chan struct{}), callbacks: []chunkCallback{callback}}
	if !w.chunkWorkers.schedule(req) {
		return false
	}
	w.chunkRequests[pos] = req
	return true
}

func (w *World) addChunk(pos ChunkPos, c *chunk.Column) *Column {
	column := w.columnFrom(c, pos)
	w.chunks[pos] = column
	for _, e := range column.Entities {
		w.entities[e] = pos
		e.setAndUnlockWorld(w)
		e.markWorldReady(w)
	}
	w.calculateLight(pos)
	return column
}

func (w *World) chunkFromAsyncPool(tx *Tx, pos ChunkPos) (*Column, bool) {
	req, ok := w.chunkRequests[pos]
	if ok {
		c := req.doImmediate(tx)
		return c, c != nil
	}
	return nil, false
}

func (w *World) calculateLight(centre ChunkPos) {
	for x := int32(-1); x <= 1; x++ {
		for z := int32(-1); z <= 1; z++ {

			pos := ChunkPos{centre[0] + x, centre[1] + z}
			if _, ok := w.chunks[pos]; ok {

				w.spreadLight(pos)
			}
		}
	}
}

func (w *World) spreadLight(pos ChunkPos) {
	c := make([]*chunk.Chunk, 0, 9)
	for z := int32(-1); z <= 1; z++ {
		for x := int32(-1); x <= 1; x++ {
			neighbour, ok := w.chunks[ChunkPos{pos[0] + x, pos[1] + z}]
			if !ok {

				return
			}
			c = append(c, neighbour.Chunk)
		}
	}

	chunk.LightArea(c, int(pos[0])-1, int(pos[1])-1).Spread()
}

func (w *World) autoSave() {
	save := &time.Ticker{C: make(<-chan time.Time)}
	if w.conf.SaveInterval > 0 {
		save = time.NewTicker(w.conf.SaveInterval)
		defer save.Stop()
	}
	closeUnused := time.NewTicker(w.conf.ChunkUnloadInterval)
	defer closeUnused.Stop()

	for {
		select {
		case <-closeUnused.C:
			<-w.exec(w.closeUnusedChunks)
		case <-save.C:
			w.Save()
		case <-w.closing:
			w.running.Done()
			return
		}
	}
}

func (w *World) closeUnusedChunks(tx *Tx) {
	for pos, c := range w.chunks {
		if len(c.viewers) == 0 {
			w.closeChunk(tx, pos, c)
		}
	}
}

type Column struct {
	modified bool

	net *NetworkChunk

	*chunk.Chunk
	Entities      []*EntityHandle
	BlockEntities map[cube.Pos]Block

	viewers []Viewer
	loaders []*Loader
}

func (w *World) columnTo(col *Column, pos ChunkPos) *chunk.Column {
	scheduled := w.scheduledUpdates.fromChunk(pos)
	c := &chunk.Column{
		Chunk:           col.Chunk,
		Entities:        make([]chunk.Entity, 0, len(col.Entities)),
		BlockEntities:   make([]chunk.BlockEntity, 0, len(col.BlockEntities)),
		ScheduledBlocks: make([]chunk.ScheduledBlockUpdate, 0, len(scheduled)),
		Tick:            w.scheduledUpdates.currentTick,
	}
	for _, e := range col.Entities {
		data := e.encodeNBT()
		maps.Copy(data, e.t.EncodeNBT(&e.data))
		data["identifier"] = e.t.EncodeEntity()
		c.Entities = append(c.Entities, chunk.Entity{ID: int64(binary.LittleEndian.Uint64(e.id[8:])), Data: data})
	}
	for pos, be := range col.BlockEntities {
		c.BlockEntities = append(c.BlockEntities, chunk.BlockEntity{Pos: pos, Data: be.(NBTer).EncodeNBT()})
	}
	for _, t := range scheduled {
		c.ScheduledBlocks = append(c.ScheduledBlocks, chunk.ScheduledBlockUpdate{Pos: t.pos, Block: w.conf.Blocks.BlockRuntimeID(t.b), Tick: t.t})
	}
	return c
}

func (w *World) columnFrom(c *chunk.Column, _ ChunkPos) *Column {
	col := &Column{
		Chunk:         c.Chunk,
		Entities:      make([]*EntityHandle, 0, len(c.Entities)),
		BlockEntities: make(map[cube.Pos]Block, len(c.BlockEntities)),
	}
	for _, e := range c.Entities {
		eid, ok := e.Data["identifier"].(string)
		if !ok {
			w.conf.Log.Error("read column: entity without identifier field", "ID", e.ID)
			continue
		}
		t, ok := w.conf.Entities.Lookup(eid)
		if !ok {

			w.conf.Log.Debug("read column: unknown entity type", "ID", e.ID, "type", eid)
			continue
		}
		handle := entityFromData(t, e.ID, e.Data)
		if handle == nil {
			w.conf.Log.Warn("read column: entity could not be decoded and is dropped", "ID", e.ID, "type", eid)
			continue
		}
		col.Entities = append(col.Entities, handle)
	}
	for _, be := range c.BlockEntities {
		rid := c.Chunk.Block(uint8(be.Pos[0]), int16(be.Pos[1]), uint8(be.Pos[2]), 0)
		b, ok := w.conf.Blocks.BlockByRuntimeID(rid)
		if !ok {
			w.conf.Log.Error("read column: no block with runtime ID", "ID", rid)
			continue
		}
		nb, ok := b.(NBTer)
		if !ok {
			w.conf.Log.Error("read column: block with nbt does not implement NBTer", "block", fmt.Sprintf("%#v", b))
			continue
		}
		col.BlockEntities[be.Pos] = nb.DecodeNBT(be.Data).(Block)
	}
	scheduled, savedTick := make([]scheduledTick, 0, len(c.ScheduledBlocks)), c.Tick
	for _, t := range c.ScheduledBlocks {
		bl := w.conf.Blocks.BlockByRuntimeIDOrAir(t.Block)
		scheduled = append(scheduled, scheduledTick{
			pos:   t.Pos,
			b:     bl,
			bhash: w.conf.Blocks.BlockHash(bl),
			t:     w.scheduledUpdates.currentTick + (t.Tick - savedTick),
		})
	}
	w.scheduledUpdates.add(scheduled)
	return col
}
