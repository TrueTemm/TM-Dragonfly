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
	"bytes"
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
)

type lightArea struct {
	br           BlockRegistry
	baseX, baseZ int
	c            []*Chunk
	w            int
	r            cube.Range
}

type lightQueue struct {
	nodes []lightNode
	head  int
	tail  int
	size  int
}

const initialLightQueueCapacity = 1024

func newLightQueue(capacity int) *lightQueue {
	if capacity < 1 {
		capacity = 1
	}
	return &lightQueue{nodes: make([]lightNode, capacity)}
}

func (q *lightQueue) push(n lightNode) {
	if q.size == len(q.nodes) {
		q.grow()
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.size++
}

func (q *lightQueue) pop() (lightNode, bool) {
	if q.size == 0 {
		return lightNode{}, false
	}
	n := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.size--
	return n, true
}

func (q *lightQueue) empty() bool {
	return q.size == 0
}

func (q *lightQueue) grow() {
	nodes := make([]lightNode, len(q.nodes)<<1)
	if q.head < q.tail {
		copy(nodes, q.nodes[q.head:q.tail])
	} else {
		n := copy(nodes, q.nodes[q.head:])
		copy(nodes[n:], q.nodes[:q.tail])
	}
	q.head = 0
	q.tail = q.size
	q.nodes = nodes
}

func LightArea(c []*Chunk, baseX, baseZ int) *lightArea {
	w := int(math.Sqrt(float64(len(c))))
	if len(c) != w*w {
		panic("area must have a square chunk area")
	}
	return &lightArea{
		br:    c[0].br,
		c:     c,
		w:     w,
		baseX: baseX << 4,
		baseZ: baseZ << 4,
		r:     c[0].r,
	}
}

func (a *lightArea) Fill() {
	a.initialiseLightSlices()
	queue := newLightQueue(initialLightQueueCapacity)
	a.insertBlockLightNodes(queue)
	a.insertSkyLightNodes(queue)

	for !queue.empty() {
		a.propagate(queue)
	}
}

func (a *lightArea) Spread() {
	queue := newLightQueue(initialLightQueueCapacity)
	a.insertLightSpreadingNodes(queue, BlockLight)
	a.insertLightSpreadingNodes(queue, SkyLight)

	for !queue.empty() {
		a.propagate(queue)
	}
}

func (a *lightArea) light(pos cube.Pos, l light) uint8 {
	return l.light(a.sub(pos), uint8(pos[0]&0xf), uint8(pos[1]&0xf), uint8(pos[2]&0xf))
}

func (a *lightArea) setLight(pos cube.Pos, l light, v uint8) {
	l.setLight(a.sub(pos), uint8(pos[0]&0xf), uint8(pos[1]&0xf), uint8(pos[2]&0xf), v)
}

func (a *lightArea) iterSubChunks(filter func(sub *SubChunk) bool, f func(pos cube.Pos)) {
	for cx := 0; cx < a.w; cx++ {
		for cz := 0; cz < a.w; cz++ {
			baseX, baseZ, c := a.baseX+(cx<<4), a.baseZ+(cz<<4), a.c[a.chunkIndex(cx, cz)]

			for index, sub := range c.sub {
				if !filter(sub) {
					continue
				}
				baseY := int(c.SubY(int16(index)))
				a.iterSubChunk(func(x, y, z int) {
					f(cube.Pos{x + baseX, y + baseY, z + baseZ})
				})
			}
		}
	}
}

func (a *lightArea) iterEdges(filter func(a, b *SubChunk) bool, f func(a, b cube.Pos)) {
	minY, maxY := a.r[0]>>4, a.r[1]>>4

	for cu := 1; cu < a.w; cu++ {
		u := cu << 4
		for cv := 0; cv < a.w; cv++ {
			v := cv << 4
			for cy := minY; cy < maxY; cy++ {
				baseY := cy << 4

				xa, za := cube.Pos{a.baseX + u, baseY, a.baseZ + v}, cube.Pos{a.baseX + v, baseY, a.baseZ + u}
				xb, zb := xa.Side(cube.FaceWest), za.Side(cube.FaceNorth)

				addX, addZ := filter(a.sub(xa), a.sub(xb)), filter(a.sub(za), a.sub(zb))
				if !addX && !addZ {
					continue
				}

				for addV := 0; addV < 16; addV++ {
					for y := 0; y < 16; y++ {

						if addX {
							f(xa.Add(cube.Pos{0, y, addV}), xb.Add(cube.Pos{0, y, addV}))
						}
						if addZ {
							f(za.Add(cube.Pos{addV, y}), zb.Add(cube.Pos{addV, y}))
						}
					}
				}
			}
		}
	}
}

func (a *lightArea) iterHeightmap(f func(x, z int, height, highestNeighbour, highestY, lowestY int)) {
	m, highestY := a.c[0].HeightMap(), a.c[0].Range().Min()
	lowestY := highestY
	for index := range a.c[0].sub {
		if a.c[0].sub[index].Empty() {
			continue
		}
		highestY = int(a.c[0].SubY(int16(index))) + 15
	}
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			f(int(x)+a.baseX, int(z)+a.baseZ, int(m.At(x, z)), int(m.HighestNeighbour(x, z)), highestY, lowestY)
		}
	}
}

func (a *lightArea) iterSubChunk(f func(x, y, z int)) {
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			for z := 0; z < 16; z++ {
				f(x, y, z)
			}
		}
	}
}

func (a *lightArea) highest(pos cube.Pos, lightBlocking func(rid uint32) uint8) uint8 {
	x, y, z, sub := uint8(pos[0]&0xf), uint8(pos[1]&0xf), uint8(pos[2]&0xf), a.sub(pos)
	storages, l := sub.storages, len(sub.storages)

	switch l {
	case 0:
		return 0
	case 1:
		return lightBlocking(storages[0].At(x, y, z))
	default:
		level := lightBlocking(storages[0].At(x, y, z))
		if v := lightBlocking(storages[1].At(x, y, z)); v > level {
			return v
		}
		return level
	}
}

var (
	fullLight    = bytes.Repeat([]byte{0xff}, 2048)
	fullLightPtr = &fullLight[0]
	noLight      = make([]uint8, 2048)
	noLightPtr   = &noLight[0]
)

func (a *lightArea) initialiseLightSlices() {
	for _, c := range a.c {
		index := len(c.sub) - 1
		for index >= 0 {
			sub := c.sub[index]
			if !sub.Empty() {

				break
			}
			sub.skyLight = fullLight
			sub.blockLight = noLight
			index--
		}
		for index >= 0 {

			c.sub[index].skyLight = noLight
			c.sub[index].blockLight = noLight
			index--
		}
	}
}

func (a *lightArea) sub(pos cube.Pos) *SubChunk {
	return a.chunk(pos).SubChunk(int16(pos[1]))
}

func (a *lightArea) chunk(pos cube.Pos) *Chunk {
	x, z := pos[0]-a.baseX, pos[2]-a.baseZ
	return a.c[a.chunkIndex(x>>4, z>>4)]
}

func (a *lightArea) chunkIndex(x, z int) int {
	return x + (z * a.w)
}
