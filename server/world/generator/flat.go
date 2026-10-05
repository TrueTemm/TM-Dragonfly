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

package generator

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Flat struct {
	biome uint32

	layers []uint32
}

func NewFlat(biome world.Biome, layers []world.Block) Flat {
	return NewFlatWithRegistry(biome, layers, world.DefaultBlockRegistry)
}

func NewFlatWithRegistry(biome world.Biome, layers []world.Block, br world.BlockRegistry) Flat {
	f := Flat{
		biome:  uint32(biome.EncodeBiome()),
		layers: make([]uint32, len(layers)),
	}
	for i, b := range layers {
		f.layers[i] = br.BlockRuntimeID(b)
	}
	return f
}

func (f Flat) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	min, max := int16(chunk.Range().Min()), int16(chunk.Range().Max())
	n := int16(len(f.layers))

	for x := range uint8(16) {
		for z := range uint8(16) {
			for y := int16(0); y <= max; y++ {
				if y < n {
					chunk.SetBlock(x, min+y, z, 0, f.layers[n-y-1])
				}
				chunk.SetBiome(x, min+y, z, f.biome)
			}
		}
	}
}

func (f Flat) DefaultSpawn(dim world.Dimension) cube.Pos {
	return cube.Pos{0, dim.Range().Min() + len(f.layers) + 1, 0}
}
