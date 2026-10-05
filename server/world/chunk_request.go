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
	"sync"

	"github.com/df-mc/dragonfly/server/world/chunk"
)

type chunkRequest struct {
	pos       ChunkPos
	callbacks []chunkCallback
	signalled bool

	done   chan struct{}
	col    *chunk.Column
	err    error
	result *Column
}

const defaultChunkLoadWorkers = 1

type chunkCallback = func(tx *Tx, col *Column)

type chunkWorkerPool struct {
	w     *World
	queue chan *chunkRequest
	wg    sync.WaitGroup

	mu     sync.Mutex
	closed bool
}

func newChunkWorkerPool(w *World) *chunkWorkerPool {
	return &chunkWorkerPool{w: w, queue: make(chan *chunkRequest, 4096)}
}

func (r *chunkRequest) doImmediate(tx *Tx) *Column {
	<-r.done
	r.signal(tx)
	return r.result
}

func (r *chunkRequest) load(w *World) {
	r.col, r.err = w.loadChunk(r.pos)
	close(r.done)
	w.Do(r.signal)
}

func (r *chunkRequest) abort() {
	close(r.done)
}

func (p *chunkWorkerPool) schedule(r *chunkRequest) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed || p.w.closed.Load() {
		p.closed = true
		return false
	}
	select {
	case p.queue <- r:
		return true
	default:
		return false
	}
}

func (p *chunkWorkerPool) handle() {
	defer p.wg.Done()
	for {
		if p.w.closed.Load() {
			p.drainAndAbort()
			return
		}
		select {
		case r := <-p.queue:
			r.load(p.w)
		case <-p.w.closeStarted:
			p.drainAndAbort()
			return
		}
	}
}

func (p *chunkWorkerPool) drainAndAbort() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	for {
		select {
		case r := <-p.queue:
			r.abort()
		default:
			return
		}
	}
}

func (r *chunkRequest) signal(tx *Tx) {
	if r.signalled {
		return
	}
	r.signalled = true

	w := tx.World()
	pos := r.pos

	delete(w.chunkRequests, pos)
	if w.closed.Load() {
		return
	}
	if r.err != nil {
		w.conf.Log.Error("load chunk: "+r.err.Error(), "X", pos[0], "Z", pos[1])
		for _, recv := range r.callbacks {
			recv(tx, nil)
		}
		return
	}
	r.result = w.addChunk(pos, r.col)
	if w.closed.Load() {
		return
	}
	for _, recv := range r.callbacks {
		recv(tx, r.result)
	}
}
