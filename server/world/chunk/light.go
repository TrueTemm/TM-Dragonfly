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

import "github.com/df-mc/dragonfly/server/block/cube"

func (a *lightArea) insertBlockLightNodes(queue *lightQueue) {
	a.iterSubChunks(a.anyLightBlocks, func(pos cube.Pos) {
		if level := a.highest(pos, a.br.LightBlock); level > 0 {
			queue.push(node(pos, level, BlockLight))
		}
	})
}

func (a *lightArea) anyLightBlocks(sub *SubChunk) bool {
	for _, layer := range sub.storages {
		for _, id := range layer.palette.values {
			if a.br.LightBlock(id) != 0 {
				return true
			}
		}
	}
	return false
}

func (a *lightArea) insertSkyLightNodes(queue *lightQueue) {
	a.iterHeightmap(func(x, z int, height, highestNeighbour, highestY, lowestY int) {
		pos := cube.Pos{x, height, z}
		if height <= a.r.Max() {

			a.setLight(pos, SkyLight, 15)

			if pos[1] > lowestY {
				if level := a.highest(pos.Sub(cube.Pos{0, 1}), a.br.FilteringBlock); level != 15 && level != 0 {

					queue.push(node(pos, 15, SkyLight))
				}
			}
		}
		for y := pos[1]; y < highestY; y++ {

			if pos[1]++; pos[1] < highestNeighbour {
				queue.push(node(pos, 15, SkyLight))
				continue
			}

			a.setLight(pos, SkyLight, 15)
		}
	})
}

func (a *lightArea) insertLightSpreadingNodes(queue *lightQueue, lt light) {
	a.iterEdges(a.nodesNeeded(lt), func(pa, pb cube.Pos) {
		la, lb := a.light(pa, lt), a.light(pb, lt)
		if la == lb || la-1 == lb || lb-1 == la {

			return
		}
		if filter := a.highest(pb, a.br.FilteringBlock) + 1; la > filter && la-filter > lb {
			queue.push(node(pb, la-filter, lt))
		} else if filter = a.highest(pa, a.br.FilteringBlock) + 1; lb > filter && lb-filter > la {
			queue.push(node(pa, lb-filter, lt))
		}
	})
}

func (a *lightArea) nodesNeeded(lt light) func(sa, sb *SubChunk) bool {
	if lt == SkyLight {
		return func(sa, sb *SubChunk) bool {
			return &sa.skyLight[0] != &sb.skyLight[0]
		}
	}
	return func(sa, sb *SubChunk) bool {

		return &sa.blockLight[0] != &sb.blockLight[0]
	}
}

func (a *lightArea) propagate(queue *lightQueue) {
	n, ok := queue.pop()
	if !ok {
		return
	}
	if a.light(n.pos, n.lt) >= n.level {
		return
	}
	a.setLight(n.pos, n.lt, n.level)

	x, y, z := n.pos[0], n.pos[1], n.pos[2]
	a.propagateNeighbour(queue, n.lt, n.level, x+1, y, z)
	a.propagateNeighbour(queue, n.lt, n.level, x-1, y, z)
	a.propagateNeighbour(queue, n.lt, n.level, x, y+1, z)
	a.propagateNeighbour(queue, n.lt, n.level, x, y-1, z)
	a.propagateNeighbour(queue, n.lt, n.level, x, y, z+1)
	a.propagateNeighbour(queue, n.lt, n.level, x, y, z-1)
}

func (a *lightArea) propagateNeighbour(queue *lightQueue, lt light, level uint8, x, y, z int) {
	if y < a.r.Min() || y > a.r.Max() || x < a.baseX || z < a.baseZ || x >= a.baseX+a.w*16 || z >= a.baseZ+a.w*16 {
		return
	}
	pos := cube.Pos{x, y, z}
	filter := a.highest(pos, a.br.FilteringBlock) + 1
	if level > filter && a.light(pos, lt) < level-filter {
		queue.push(node(pos, level-filter, lt))
	}
}

type lightNode struct {
	pos   cube.Pos
	lt    light
	level uint8
}

func node(pos cube.Pos, level uint8, lt light) lightNode {
	return lightNode{pos: pos, level: level, lt: lt}
}
