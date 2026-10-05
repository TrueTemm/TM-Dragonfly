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
	"log/slog"
	"math/rand/v2"
	"time"
)

type blockRegistrySetter interface {
	SetBlockRegistry(BlockRegistry)
}

type Config struct {
	Log *slog.Logger

	Dim Dimension

	PortalDestination func(dim Dimension) *World

	Provider Provider

	Generator Generator

	ReadOnly bool

	SaveInterval time.Duration

	ChunkUnloadInterval time.Duration

	ChunkLoadWorkers int

	RandomTickSpeed int

	RandSource rand.Source

	Entities EntityRegistry

	Blocks BlockRegistry

	Synchronous bool
}

func (conf Config) New() *World {
	if conf.Log == nil {
		conf.Log = slog.Default()
	}
	if conf.Dim == nil {
		conf.Dim = Overworld
	}
	if conf.SaveInterval == 0 {
		conf.SaveInterval = time.Minute * 10
	}
	if conf.ChunkUnloadInterval <= 0 {
		conf.ChunkUnloadInterval = time.Minute * 2
	}
	if conf.ChunkLoadWorkers <= 0 {
		conf.ChunkLoadWorkers = defaultChunkLoadWorkers
	}
	if conf.Generator == nil {
		conf.Generator = NopGenerator{}
	}
	if conf.Provider == nil {

		s := defaultSettings()
		s.Spawn = conf.Generator.DefaultSpawn(conf.Dim)
		conf.Provider = NopProvider{Set: s}
	}
	if conf.RandomTickSpeed == 0 {
		conf.RandomTickSpeed = 3
	}
	if conf.Blocks == nil {
		conf.Blocks = DefaultBlockRegistry
	}

	conf.Blocks.Finalize()
	DefaultBlockRegistry.Finalize()
	if provider, ok := conf.Provider.(blockRegistrySetter); ok {
		provider.SetBlockRegistry(conf.Blocks)
	}

	if conf.RandSource == nil {
		t := uint64(time.Now().UnixNano())
		conf.RandSource = rand.NewPCG(t, t)
	}
	s := conf.Provider.Settings()

	conf.Provider = &lockedProvider{p: conf.Provider}
	if conf.ChunkLoadWorkers == 1 {
		conf.Generator = &lockedGenerator{g: conf.Generator}
	}
	w := &World{
		scheduledUpdates: newScheduledTickQueue(s.CurrentTick),
		redstone:         newRedstoneEngine(s.CurrentTick),
		entities:         make(map[*EntityHandle]ChunkPos),
		viewers:          make(map[*Loader]Viewer),
		chunks:           make(map[ChunkPos]*Column),
		chunkRequests:    make(map[ChunkPos]*chunkRequest),
		queueClosing:     make(chan struct{}),
		closeStarted:     make(chan struct{}),
		closing:          make(chan struct{}),
		queue:            make(chan transaction, 128),
		r:                rand.New(conf.RandSource),
		advance:          s.ref.Add(1) == 1,
		conf:             conf,
		ra:               conf.Dim.Range(),
		set:              s,
	}
	w.chunkWorkers = newChunkWorkerPool(w)
	w.weather = weather{w: w}
	var h Handler = NopHandler{}
	w.handler.Store(&h)

	t := ticker{interval: time.Second / 20}
	if !conf.Synchronous {
		w.queueing.Add(1)
		w.running.Add(2)

		go t.tickLoop(w)
		go w.autoSave()
		go w.handleTransactions()
		w.chunkWorkers.wg.Add(conf.ChunkLoadWorkers)
		for range conf.ChunkLoadWorkers {
			go w.chunkWorkers.handle()
		}
	}

	<-w.exec(t.tick)
	return w
}
