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

package trace

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type BlockResult struct {
	bb   cube.BBox
	pos  mgl64.Vec3
	face cube.Face

	blockPos cube.Pos
}

func (r BlockResult) BBox() cube.BBox {
	return r.bb
}

func (r BlockResult) Position() mgl64.Vec3 {
	return r.pos
}

func (r BlockResult) Face() cube.Face {
	return r.face
}

func (r BlockResult) BlockPosition() cube.Pos {
	return r.blockPos
}

func BlockIntercept(pos cube.Pos, src world.BlockSource, b world.Block, start, end mgl64.Vec3) (result BlockResult, ok bool) {
	bbs := b.Model().BBox(pos, src)
	if len(bbs) == 0 {
		return
	}

	var (
		hit  Result
		dist = math.MaxFloat64
	)

	for _, bb := range bbs {
		next, ok := BBoxIntercept(bb.Translate(pos.Vec3()), start, end)
		if !ok {
			continue
		}

		nextDist := next.Position().Sub(start).LenSqr()
		if nextDist < dist {
			hit = next
			dist = nextDist
		}
	}

	if hit == nil {
		return result, false
	}

	return BlockResult{bb: hit.BBox(), pos: hit.Position(), face: hit.Face(), blockPos: pos}, true
}

func BlockIntersects(pos cube.Pos, src world.BlockSource, b world.Block, start, end mgl64.Vec3) bool {
	m := b.Model()
	switch m.(type) {
	case model.Empty:
		return false
	case model.Solid:
		return BBoxIntersects(cube.Box(0, 0, 0, 1, 1, 1).Translate(pos.Vec3()), start, end)
	}

	for _, bb := range m.BBox(pos, src) {
		if BBoxIntersects(bb.Translate(pos.Vec3()), start, end) {
			return true
		}
	}
	return false
}
