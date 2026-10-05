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

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Generator interface {
	GenerateChunk(pos ChunkPos, chunk *chunk.Chunk)

	DefaultSpawn(dim Dimension) cube.Pos
}

type NopGenerator struct{}

func (NopGenerator) GenerateChunk(ChunkPos, *chunk.Chunk) {}

func (NopGenerator) DefaultSpawn(Dimension) cube.Pos { return cube.Pos{} }

type lockedGenerator struct {
	mu sync.Mutex
	g  Generator
}

func (l *lockedGenerator) GenerateChunk(pos ChunkPos, c *chunk.Chunk) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.g.GenerateChunk(pos, c)
}

func (l *lockedGenerator) DefaultSpawn(dim Dimension) cube.Pos {
	return l.g.DefaultSpawn(dim)
}
