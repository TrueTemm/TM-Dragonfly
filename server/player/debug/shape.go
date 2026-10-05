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

package debug

import (
	"image/color"
	"sync/atomic"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

var nextShapeID atomic.Int32

type Shape interface {
	ShapeID() int
}

type shape struct {
	id atomic.Int32
}

func (s *shape) ShapeID() int {
	if id := s.id.Load(); id != 0 {
		return int(id)
	}
	s.id.CompareAndSwap(0, nextShapeID.Add(1))
	return int(s.id.Load())
}

type Arrow struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	EndPosition mgl64.Vec3

	HeadLength float64

	HeadRadius float64

	HeadSegments int

	Entity *world.EntityHandle
}

type Box struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	Bounds mgl64.Vec3

	Entity *world.EntityHandle
}

type Circle struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	Segments int

	Entity *world.EntityHandle
}

type Line struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	EndPosition mgl64.Vec3

	Entity *world.EntityHandle
}

type Sphere struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	Segments int

	Entity *world.EntityHandle
}

type Text struct {
	shape

	Colour color.RGBA

	BackgroundColour color.RGBA

	HideBackground bool

	Position mgl64.Vec3

	Rotation mgl64.Vec3

	Scale float64

	Text string

	LockRotation bool

	DisableDepthTest bool

	HideBackface bool

	HideBackfaceText bool

	Entity *world.EntityHandle
}

type Cylinder struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	BaseRadius mgl64.Vec2

	TopRadius mgl64.Vec2

	Height float64

	Segments int

	Entity *world.EntityHandle
}

type Pyramid struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	Width float64

	Depth float64

	Height float64

	Entity *world.EntityHandle
}

type Ellipsoid struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	Radii mgl64.Vec3

	SegmentsPerAxis int

	Entity *world.EntityHandle
}

type Cone struct {
	shape

	Colour color.RGBA

	Position mgl64.Vec3

	Scale float64

	Radii mgl64.Vec2

	Height float64

	Segments int

	Entity *world.EntityHandle
}
