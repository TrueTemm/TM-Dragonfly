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

package block

import (
	"math"
	"sync"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type LiquidRemovable interface {
	HasLiquidDrops() bool
}

type sourceWaterDisplacer struct{}

func (s sourceWaterDisplacer) CanDisplace(b world.Liquid) bool {
	w, ok := b.(Water)
	return ok && !w.Falling && w.Depth == 8
}

type flowingWaterDisplacer struct{}

func (s flowingWaterDisplacer) CanDisplace(b world.Liquid) bool {
	_, ok := b.(Water)
	return ok
}

func tickLiquid(b world.Liquid, pos cube.Pos, tx *world.Tx) {
	if !source(b) && !sourceAround(b, pos, tx) {
		var res world.Liquid
		if b.LiquidDepth()-4 > 0 {
			res = b.WithDepth(b.LiquidDepth()-2*b.SpreadDecay(), false)
		}
		ctx := tx.Event()
		if tx.World().Handler().HandleLiquidDecay(ctx, pos, b, res); ctx.Cancelled() {
			return
		}
		tx.SetLiquid(pos, res)
		return
	}
	displacer, _ := tx.Block(pos).(world.LiquidDisplacer)

	canFlowBelow := canFlowInto(b, tx, pos.Side(cube.FaceDown), false)
	if b.LiquidFalling() && !canFlowBelow {
		b = b.WithDepth(8, true)
	} else if canFlowBelow {
		below := pos.Side(cube.FaceDown)
		if displacer == nil || !displacer.SideClosed(pos, below, tx) {
			flowInto(b.WithDepth(8, true), pos, below, tx, true)
		}
	}

	depth, decay := b.LiquidDepth(), b.SpreadDecay()
	if depth <= decay {

		return
	}
	if source(b) || !canFlowBelow {
		paths := calculateLiquidPaths(b, pos, tx, displacer)
		if len(paths) == 0 {
			spreadOutwards(b, pos, tx, displacer)
			return
		}

		smallestLen := len(paths[0])
		for _, path := range paths {
			if len(path) <= smallestLen {
				flowInto(b, pos, path[0], tx, false)
			}
		}
	}
}

func source(b world.Liquid) bool {
	return b.LiquidDepth() == 8 && !b.LiquidFalling()
}

func spreadOutwards(b world.Liquid, pos cube.Pos, tx *world.Tx, displacer world.LiquidDisplacer) {
	pos.Neighbours(func(neighbour cube.Pos) {
		if neighbour[1] == pos[1] {
			if displacer == nil || !displacer.SideClosed(pos, neighbour, tx) {
				flowInto(b, pos, neighbour, tx, false)
			}
		}
	}, tx.Range())
}

func sourceAround(b world.Liquid, pos cube.Pos, tx *world.Tx) (sourcePresent bool) {
	pos.Neighbours(func(neighbour cube.Pos) {
		if neighbour[1] == pos[1]-1 {

			return
		}
		side, ok := tx.Liquid(neighbour)
		if !ok || side.LiquidType() != b.LiquidType() {
			return
		}
		if displacer, ok := tx.Block(neighbour).(world.LiquidDisplacer); ok && displacer.SideClosed(neighbour, pos, tx) {

			return
		}
		if neighbour[1] == pos[1]+1 || source(side) || side.LiquidDepth() > b.LiquidDepth() {
			sourcePresent = true
		}
	}, tx.Range())
	return
}

func flowInto(b world.Liquid, src, pos cube.Pos, tx *world.Tx, falling bool) bool {
	newDepth := b.LiquidDepth() - b.SpreadDecay()
	if falling {
		newDepth = b.LiquidDepth()
	}
	if newDepth <= 0 && !falling {
		return false
	}
	existing := tx.Block(pos)
	if existingLiquid, alsoLiquid := existing.(world.Liquid); alsoLiquid && existingLiquid.LiquidType() == b.LiquidType() {
		if existingLiquid.LiquidDepth() >= newDepth || existingLiquid.LiquidFalling() {

			return true
		}
		ctx := tx.Event()
		if tx.World().Handler().HandleLiquidFlow(ctx, src, pos, b.WithDepth(newDepth, falling), existing); ctx.Cancelled() {
			return false
		}
		tx.SetLiquid(pos, b.WithDepth(newDepth, falling))
		return true
	} else if alsoLiquid {
		existingLiquid.Harden(pos, tx, &src)
		return false
	}
	displacer, isDisplacer := existing.(world.LiquidDisplacer)
	if isDisplacer {
		if _, ok := tx.Liquid(pos); ok {

			return false
		}
	}
	_, isRemovable := existing.(LiquidRemovable)
	if !isRemovable && (!isDisplacer || !displacer.CanDisplace(b.WithDepth(newDepth, falling))) {

		return false
	}
	ctx := tx.Event()
	if tx.World().Handler().HandleLiquidFlow(ctx, src, pos, b.WithDepth(newDepth, falling), existing); ctx.Cancelled() {
		return false
	}

	if isRemovable {
		if _, air := existing.(Air); !air {
			tx.SetBlock(pos, nil, nil)
			b.LiquidRemoveBlock(pos, tx, existing)
		}
	}
	tx.SetLiquid(pos, b.WithDepth(newDepth, falling))
	return true
}

type liquidPath []cube.Pos

func calculateLiquidPaths(b world.Liquid, pos cube.Pos, tx *world.Tx, displacer world.LiquidDisplacer) []liquidPath {
	queue := liquidQueuePool.Get().(*liquidQueue)
	defer func() {
		queue.Reset()
		liquidQueuePool.Put(queue)
	}()
	queue.PushBack(liquidNode{x: pos[0], z: pos[2], depth: int8(b.LiquidDepth())})
	decay := int8(b.SpreadDecay())

	paths := make([]liquidPath, 0, 3)
	first := true

	for queue.Len() != 0 {
		node := queue.Front()
		neighA, neighB, neighC, neighD := node.neighbours(decay * 2)
		if !first || (displacer == nil || !displacer.SideClosed(pos, cube.Pos{neighA.x, pos[1], neighA.z}, tx)) {
			if spreadNeighbour(b, pos, tx, neighA, queue) {
				queue.shortestPath = neighA.Len()
				paths = append(paths, neighA.Path(pos))
			}
		}
		if !first || (displacer == nil || !displacer.SideClosed(pos, cube.Pos{neighB.x, pos[1], neighB.z}, tx)) {
			if spreadNeighbour(b, pos, tx, neighB, queue) {
				queue.shortestPath = neighB.Len()
				paths = append(paths, neighB.Path(pos))
			}
		}
		if !first || (displacer == nil || !displacer.SideClosed(pos, cube.Pos{neighC.x, pos[1], neighC.z}, tx)) {
			if spreadNeighbour(b, pos, tx, neighC, queue) {
				queue.shortestPath = neighC.Len()
				paths = append(paths, neighC.Path(pos))
			}
		}
		if !first || (displacer == nil || !displacer.SideClosed(pos, cube.Pos{neighD.x, pos[1], neighD.z}, tx)) {
			if spreadNeighbour(b, pos, tx, neighD, queue) {
				queue.shortestPath = neighD.Len()
				paths = append(paths, neighD.Path(pos))
			}
		}
		first = false
	}
	return paths
}

func spreadNeighbour(b world.Liquid, src cube.Pos, tx *world.Tx, node liquidNode, queue *liquidQueue) bool {
	if node.depth+3 <= 0 {

		return false
	}
	if node.Len() > queue.shortestPath {

		return false
	}
	pos := cube.Pos{node.x, src[1], node.z}
	if !canFlowInto(b, tx, pos, true) {

		return false
	}
	pos[1]--
	if canFlowInto(b, tx, pos, false) {
		return true
	}
	queue.PushBack(node)
	return false
}

func canFlowInto(b world.Liquid, tx *world.Tx, pos cube.Pos, sideways bool) bool {
	bl := tx.Block(pos)
	if _, air := bl.(Air); air {

		return true
	}
	if _, ok := bl.(LiquidRemovable); ok {
		if liq, ok := bl.(world.Liquid); ok && sideways {
			if (liq.LiquidDepth() == 8 && !liq.LiquidFalling()) || liq.LiquidType() != b.LiquidType() {

				return false
			}
		}
		return true
	}
	if dis, ok := bl.(world.LiquidDisplacer); ok {
		res := b.WithDepth(b.LiquidDepth()-b.SpreadDecay(), !sideways)
		if dis.CanDisplace(res) {
			return true
		}
	}
	return false
}

type liquidNode struct {
	x, z     int
	depth    int8
	previous *liquidNode
}

func (node liquidNode) neighbours(decay int8) (a, b, c, d liquidNode) {
	return liquidNode{x: node.x - 1, z: node.z, depth: node.depth - decay, previous: &node},
		liquidNode{x: node.x + 1, z: node.z, depth: node.depth - decay, previous: &node},
		liquidNode{x: node.x, z: node.z - 1, depth: node.depth - decay, previous: &node},
		liquidNode{x: node.x, z: node.z + 1, depth: node.depth - decay, previous: &node}
}

func (node liquidNode) Len() int {
	i := 1
	for {
		if node.previous == nil {
			return i - 1
		}

		node = *node.previous
		i++
	}
}

func (node liquidNode) Path(src cube.Pos) liquidPath {
	l := node.Len()
	path := make(liquidPath, l)
	i := l - 1
	for {
		if node.previous == nil {
			return path
		}
		path[i] = cube.Pos{node.x, src[1], node.z}

		node = *node.previous
		i--
	}
}

var liquidQueuePool = sync.Pool{
	New: func() any {
		return &liquidQueue{
			nodes:        make([]liquidNode, 0, 64),
			shortestPath: math.MaxInt8,
		}
	},
}

type liquidQueue struct {
	nodes        []liquidNode
	i            int
	shortestPath int
}

func (q *liquidQueue) PushBack(node liquidNode) {
	q.nodes = append(q.nodes, node)
}

func (q *liquidQueue) Front() liquidNode {
	v := q.nodes[q.i]
	q.i++
	return v
}

func (q *liquidQueue) Len() int {
	return len(q.nodes) - q.i
}

func (q *liquidQueue) Reset() {
	q.nodes = q.nodes[:0]
	q.i = 0
	q.shortestPath = math.MaxInt8
}
