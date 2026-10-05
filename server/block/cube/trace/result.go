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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"iter"
	"math"
)

type Result interface {
	BBox() cube.BBox

	Position() mgl64.Vec3

	Face() cube.Face
}

type EntityFilter func(iter.Seq[world.Entity]) iter.Seq[world.Entity]

func Perform(start, end mgl64.Vec3, tx *world.Tx, box cube.BBox, filter EntityFilter) (hit Result, ok bool) {

	TraverseBlocks(start, end, func(pos cube.Pos) (cont bool) {
		b := tx.Block(pos)

		if result, ok := BlockIntercept(pos, tx, b, start, end); ok {
			hit = result
			end = hit.Position()
			return false
		}
		return true
	})

	dist := math.MaxFloat64
	bb := box.Translate(start).Extend(end.Sub(start))
	entities := tx.EntitiesWithin(bb.Grow(8.0))
	if filter != nil {
		entities = filter(entities)
	}
	for entity := range entities {
		if !entity.H().Type().BBox(entity).Translate(entity.Position()).IntersectsWith(bb) {
			continue
		}

		result, ok := EntityIntercept(entity, start, end)
		if !ok {
			continue
		}

		if distance := start.Sub(result.Position()).LenSqr(); distance < dist {
			dist = distance
			hit = result
		}
	}

	return hit, hit != nil
}
