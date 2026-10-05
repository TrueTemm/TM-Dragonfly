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
	"math"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

type Loader struct {
	r      int
	w      *World
	viewer Viewer

	mu        sync.RWMutex
	pos       ChunkPos
	loadQueue []ChunkPos
	loaded    map[ChunkPos]*Column
	pending   map[ChunkPos]struct{}

	closed bool
}

func NewLoader(chunkRadius int, world *World, v Viewer) *Loader {
	l := &Loader{r: chunkRadius, loaded: make(map[ChunkPos]*Column), pending: make(map[ChunkPos]struct{}), viewer: v}
	l.world(world)
	return l
}

func (l *Loader) World() *World {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.w
}

func (l *Loader) ChangeWorld(tx *Tx, new *World) {
	l.mu.Lock()
	defer l.mu.Unlock()

	loaded := maps.Clone(l.loaded)
	l.w.exec(func(tx *Tx) {
		for pos := range loaded {
			tx.World().removeViewer(tx, pos, l)
		}
	})
	clear(l.loaded)
	clear(l.pending)
	l.w.viewerMu.Lock()
	delete(l.w.viewers, l)
	l.w.viewerMu.Unlock()

	l.world(new)
}

func (l *Loader) ChangeRadius(tx *Tx, new int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.r = new
	l.evictUnused(tx)
	l.populateLoadQueue()
}

func (l *Loader) Move(tx *Tx, pos mgl64.Vec3) {
	l.mu.Lock()
	defer l.mu.Unlock()

	chunkPos := chunkPosFromVec3(pos)
	if chunkPos == l.pos {
		return
	}
	l.pos = chunkPos
	l.evictUnused(tx)
	l.populateLoadQueue()
}

func (l *Loader) Load(tx *Tx, n int) {
	for i := 0; i < n; i++ {
		l.mu.Lock()
		if l.closed || l.w == nil {
			l.mu.Unlock()
			return
		}
		if len(l.loadQueue) == 0 {
			l.mu.Unlock()
			break
		}
		pos := l.loadQueue[0]
		w := tx.World()
		l.pending[pos] = struct{}{}

		l.loadQueue = l.loadQueue[1:]
		l.mu.Unlock()

		if !w.loadChunkAsync(tx, pos, func(tx2 *Tx, col *Column) {
			l.viewChunk(tx2, pos, col)
		}) {
			l.mu.Lock()
			delete(l.pending, pos)
			l.queueLoad(pos)
			l.mu.Unlock()
		}
	}
}

func (l *Loader) viewChunk(tx *Tx, pos ChunkPos, c *Column) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed || l.viewer == nil || l.w == nil || l.w != tx.World() {
		return
	}
	delete(l.pending, pos)
	if c == nil {
		l.queueLoad(pos)
		return
	}
	if _, ok := l.loaded[pos]; ok {
		return
	}
	if !l.withinLoadRadius(pos) {
		return
	}
	l.viewer.ViewChunk(pos, l.w.Dimension(), c)
	l.w.addViewer(tx, c, l)

	l.loaded[pos] = c
}

func (l *Loader) Chunk(pos ChunkPos) (*Column, bool) {
	l.mu.RLock()
	c, ok := l.loaded[pos]
	l.mu.RUnlock()
	return c, ok
}

func (l *Loader) Close(tx *Tx) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for pos := range l.loaded {
		tx.World().removeViewer(tx, pos, l)
	}
	l.loaded = map[ChunkPos]*Column{}
	clear(l.pending)

	l.w.viewerMu.Lock()
	delete(l.w.viewers, l)
	l.w.viewerMu.Unlock()

	l.closed = true
	l.viewer = nil
}

func (l *Loader) world(new *World) {
	l.w = new
	l.w.addWorldViewer(l)
	l.populateLoadQueue()
}

func (l *Loader) evictUnused(tx *Tx) {
	for pos := range l.loaded {
		if !l.withinLoadRadius(pos) {
			delete(l.loaded, pos)
			l.w.removeViewer(tx, pos, l)
		}
	}
}

func (l *Loader) withinLoadRadius(pos ChunkPos) bool {
	return chunkDistance(pos, l.pos) <= int32(l.r)
}

func chunkDistance(a, b ChunkPos) int32 {
	diffX, diffZ := float64(a[0])-float64(b[0]), float64(a[1])-float64(b[1])
	return int32(math.Round(math.Sqrt(diffX*diffX + diffZ*diffZ)))
}

func (l *Loader) queueLoad(pos ChunkPos) {
	if l.closed || l.w == nil || !l.withinLoadRadius(pos) {
		return
	}
	if _, ok := l.loaded[pos]; ok {
		return
	}
	if _, ok := l.pending[pos]; ok {
		return
	}
	for _, queued := range l.loadQueue {
		if queued == pos {
			return
		}
	}
	l.loadQueue = append(l.loadQueue, pos)
}

func (l *Loader) populateLoadQueue() {

	queue := map[int32][]ChunkPos{}

	r := int32(l.r)
	for x := -r; x <= r; x++ {
		for z := -r; z <= r; z++ {
			pos := ChunkPos{x + l.pos[0], z + l.pos[1]}
			dist := chunkDistance(pos, l.pos)
			if dist > r {

				continue
			}
			if _, ok := l.loaded[pos]; ok {

				continue
			}
			if _, ok := l.pending[pos]; ok {

				continue
			}
			queue[dist] = append(queue[dist], pos)
		}
	}

	l.loadQueue = l.loadQueue[:0]
	for i := int32(0); i <= r; i++ {
		l.loadQueue = append(l.loadQueue, queue[i]...)
	}
}
