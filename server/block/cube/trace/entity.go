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
)

type EntityResult struct {
	bb   cube.BBox
	pos  mgl64.Vec3
	face cube.Face

	entity world.Entity
}

func (r EntityResult) BBox() cube.BBox {
	return r.bb
}

func (r EntityResult) Position() mgl64.Vec3 {
	return r.pos
}

func (r EntityResult) Face() cube.Face {
	return r.face
}

func (r EntityResult) Entity() world.Entity {
	return r.entity
}

func EntityIntercept(e world.Entity, start, end mgl64.Vec3) (result EntityResult, ok bool) {
	bb := e.H().Type().BBox(e).Translate(e.Position()).Grow(0.3)
	if v, ok := e.(interface{ Velocity() mgl64.Vec3 }); ok {

		bb = bb.Extend(v.Velocity())
	}

	r, ok := BBoxIntercept(bb, start, end)
	if !ok {
		return
	}

	return EntityResult{bb: bb, pos: r.Position(), face: r.Face(), entity: e}, true
}
