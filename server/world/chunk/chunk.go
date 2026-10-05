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

package chunk

import (
	"slices"

	"github.com/df-mc/dragonfly/server/block/cube"
)

type Chunk struct {
	r cube.Range

	br BlockRegistry

	air uint32

	recalculateHeightMap bool

	heightMap HeightMap

	sub []*SubChunk

	biomes []*PalettedStorage
}

func New(br BlockRegistry, r cube.Range) *Chunk {
	n := (r.Height() >> 4) + 1
	sub, biomes := make([]*SubChunk, n), make([]*PalettedStorage, n)
	air := br.AirRuntimeID()
	for i := 0; i < n; i++ {
		sub[i] = NewSubChunk(air)
		biomes[i] = emptyStorage(0)
	}
	return &Chunk{
		r:                    r,
		br:                   br,
		air:                  air,
		sub:                  sub,
		biomes:               biomes,
		recalculateHeightMap: true,
		heightMap:            make(HeightMap, 256),
	}
}

func (chunk *Chunk) Clone() *Chunk {
	clone := &Chunk{
		r:                    chunk.r,
		br:                   chunk.br,
		air:                  chunk.air,
		recalculateHeightMap: chunk.recalculateHeightMap,
		heightMap:            slices.Clone(chunk.heightMap),
		sub:                  make([]*SubChunk, len(chunk.sub)),
		biomes:               make([]*PalettedStorage, len(chunk.biomes)),
	}
	for i, sub := range chunk.sub {
		clone.sub[i] = sub.Clone()
	}
	for i, biomes := range chunk.biomes {
		clone.biomes[i] = biomes.Clone()
	}
	return clone
}

func (chunk *Chunk) Equals(c *Chunk) bool {
	if !chunk.recalculateHeightMap && !c.recalculateHeightMap && !slices.Equal(c.heightMap, chunk.heightMap) {
		return false
	}

	if c.r != chunk.r || c.air != chunk.air || len(c.sub) != len(chunk.sub) {
		return false
	}

	for i, s := range c.sub {
		if !s.Equals(chunk.sub[i]) {
			return false
		}
	}

	return true
}

func (chunk *Chunk) Range() cube.Range {
	return chunk.r
}

func (chunk *Chunk) Sub() []*SubChunk {
	return chunk.sub
}

func (chunk *Chunk) Block(x uint8, y int16, z uint8, layer uint8) uint32 {
	sub := chunk.SubChunk(y)
	if sub.Empty() || uint8(len(sub.storages)) <= layer {
		return chunk.air
	}
	return sub.storages[layer].At(x, uint8(y), z)
}

func (chunk *Chunk) SetBlock(x uint8, y int16, z uint8, layer uint8, block uint32) {
	sub := chunk.sub[chunk.SubIndex(y)]
	if uint8(len(sub.storages)) <= layer && block == chunk.air {

		return
	}
	sub.Layer(layer).Set(x, uint8(y), z, block)
	chunk.recalculateHeightMap = true
}

func (chunk *Chunk) Biome(x uint8, y int16, z uint8) uint32 {
	return chunk.biomes[chunk.SubIndex(y)].At(x, uint8(y), z)
}

func (chunk *Chunk) SetBiome(x uint8, y int16, z uint8, biome uint32) {
	chunk.biomes[chunk.SubIndex(y)].Set(x, uint8(y), z, biome)
}

func (chunk *Chunk) Light(x uint8, y int16, z uint8) uint8 {
	ux, uy, uz, sub := x&0xf, uint8(y&0xf), z&0xf, chunk.SubChunk(y)
	sky := sub.SkyLight(ux, uy, uz)
	if sky == 15 {

		return sky
	}
	if block := sub.BlockLight(ux, uy, uz); block > sky {
		return block
	}
	return sky
}

func (chunk *Chunk) SkyLight(x uint8, y int16, z uint8) uint8 {
	return chunk.SubChunk(y).SkyLight(x&15, uint8(y&15), z&15)
}

func (chunk *Chunk) HighestLightBlocker(x, z uint8) int16 {
	return chunk.highestLightBlocker(x, z, false)
}

func (chunk *Chunk) highestLightBlocker(x, z uint8, addOne bool) int16 {
	var plus int16
	if addOne {
		plus++
	}
	for index := int16(len(chunk.sub) - 1); index >= 0; index-- {
		if sub := chunk.sub[index]; !sub.Empty() {
			for y := 15; y >= 0; y-- {
				if chunk.br.FilteringBlock(sub.storages[0].At(x, uint8(y), z)) == 15 {
					return int16(y) | chunk.SubY(index) + plus
				}
			}
		}
	}
	return int16(chunk.r[0])
}

func (chunk *Chunk) HighestBlock(x, z uint8) int16 {
	for index := int16(len(chunk.sub) - 1); index >= 0; index-- {
		if sub := chunk.sub[index]; !sub.Empty() {
			for y := 15; y >= 0; y-- {
				if rid := sub.storages[0].At(x, uint8(y), z); rid != chunk.air {
					return int16(y) | chunk.SubY(index)
				}
			}
		}
	}
	return int16(chunk.r[0])
}

func (chunk *Chunk) HeightMap() HeightMap {
	if chunk.recalculateHeightMap {
		for x := uint8(0); x < 16; x++ {
			for z := uint8(0); z < 16; z++ {
				chunk.heightMap.Set(x, z, chunk.highestLightBlocker(x, z, true))
			}
		}
		chunk.recalculateHeightMap = false
	}
	return chunk.heightMap
}

func (chunk *Chunk) Compact() {
	for i := range chunk.sub {
		chunk.sub[i].compact()
	}
}

func (chunk *Chunk) SubChunk(y int16) *SubChunk {
	return chunk.sub[chunk.SubIndex(y)]
}

func (chunk *Chunk) SubIndex(y int16) int16 {
	return (y - int16(chunk.r[0])) >> 4
}

func (chunk *Chunk) SubY(index int16) int16 {
	return (index << 4) + int16(chunk.r[0])
}

func (chunk *Chunk) HighestFilledSubChunk() uint16 {
	for i, sub := range slices.Backward(chunk.sub) {
		if !sub.Empty() {
			return uint16(i + 1)
		}
	}
	return 0
}
