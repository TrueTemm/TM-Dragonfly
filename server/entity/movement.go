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

package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type MovementComputer struct {
	Gravity, Drag     float64
	DragBeforeGravity bool

	onGround bool
}

type Movement struct {
	v                    []world.Viewer
	e                    world.Entity
	pos, vel, dpos, dvel mgl64.Vec3
	rot                  cube.Rotation
	onGround             bool
}

func (m *Movement) Send() {
	posChanged := !m.dpos.ApproxEqualThreshold(zeroVec3, epsilon)
	velChanged := !m.dvel.ApproxEqualThreshold(zeroVec3, epsilon)

	for _, v := range m.v {
		if posChanged {
			v.ViewEntityMovement(m.e, m.pos, m.rot, m.onGround)
		}
		if velChanged {
			v.ViewEntityVelocity(m.e, m.vel)
		}
	}
}

func (m *Movement) Position() mgl64.Vec3 {
	return m.pos
}

func (m *Movement) Velocity() mgl64.Vec3 {
	return m.vel
}

func (m *Movement) Rotation() cube.Rotation {
	return m.rot
}

func (c *MovementComputer) TickMovement(e world.Entity, pos, vel mgl64.Vec3, rot cube.Rotation, tx *world.Tx) *Movement {
	viewers := tx.Viewers(pos)

	velBefore := vel
	vel = c.applyHorizontalForces(tx, pos, c.applyVerticalForces(vel))
	dPos, vel := c.CheckCollision(tx, e, pos, vel)

	return &Movement{v: viewers, e: e,
		pos: pos.Add(dPos), vel: vel, dpos: dPos, dvel: vel.Sub(velBefore),
		rot: rot, onGround: c.onGround,
	}
}

func (c *MovementComputer) TickMovementVanilla(e world.Entity, pos, vel mgl64.Vec3, rot cube.Rotation, tx *world.Tx) *Movement {
	viewers := tx.Viewers(pos)

	velBefore := vel
	dPos, vel := c.CheckCollision(tx, e, pos, vel)
	vel = c.applyVerticalForces(vel)
	friction := playerAirFriction
	if c.onGround {
		friction *= blockFriction(tx, pos.Add(dPos))
	}
	vel[0] *= friction
	vel[2] *= friction

	return &Movement{v: viewers, e: e,
		pos: pos.Add(dPos), vel: vel, dpos: dPos, dvel: vel.Sub(velBefore),
		rot: rot, onGround: c.onGround,
	}
}

func (c *MovementComputer) OnGround() bool {
	return c.onGround
}

func (c *MovementComputer) Airborne() {
	c.onGround = false
}

var zeroVec3 mgl64.Vec3

const epsilon = 0.001

func (c *MovementComputer) applyVerticalForces(vel mgl64.Vec3) mgl64.Vec3 {
	if c.DragBeforeGravity {
		vel[1] *= 1 - c.Drag
	}
	vel[1] -= c.Gravity
	if !c.DragBeforeGravity {
		vel[1] *= 1 - c.Drag
	}
	return vel
}

func (c *MovementComputer) applyHorizontalForces(tx *world.Tx, pos, vel mgl64.Vec3) mgl64.Vec3 {
	friction := 1 - c.Drag
	if c.onGround {
		friction *= blockFriction(tx, pos)
	}
	vel[0] *= friction
	vel[2] *= friction
	return vel
}

const playerAirFriction = 0.91

func blockFriction(tx *world.Tx, pos mgl64.Vec3) float64 {
	if f, ok := tx.Block(cube.PosFromVec3(pos).Side(cube.FaceDown)).(interface {
		Friction() float64
	}); ok {
		return f.Friction()
	}
	return 0.6
}

func (c *MovementComputer) CheckCollision(tx *world.Tx, e world.Entity, pos, vel mgl64.Vec3) (mgl64.Vec3, mgl64.Vec3) {

	deltaX, deltaY, deltaZ := vel[0], vel[1], vel[2]

	entityBBox := e.H().Type().BBox(e).Translate(pos)
	blocks := blockBBoxsAround(tx, entityBBox.Extend(vel))

	if !mgl64.FloatEqualThreshold(deltaY, 0, epsilon) {

		for _, blockBBox := range blocks {
			deltaY = entityBBox.YOffset(blockBBox, deltaY)
		}
		entityBBox = entityBBox.Translate(mgl64.Vec3{0, deltaY})
	}
	if !mgl64.FloatEqualThreshold(deltaX, 0, epsilon) {

		for _, blockBBox := range blocks {
			deltaX = entityBBox.XOffset(blockBBox, deltaX)
		}
		entityBBox = entityBBox.Translate(mgl64.Vec3{deltaX})
	}
	if !mgl64.FloatEqualThreshold(deltaZ, 0, epsilon) {

		for _, blockBBox := range blocks {
			deltaZ = entityBBox.ZOffset(blockBBox, deltaZ)
		}
	}
	if !mgl64.FloatEqual(vel[1], 0) {

		c.onGround = false
	}
	if !mgl64.FloatEqual(deltaX, vel[0]) {
		vel[0] = 0
	}
	if !mgl64.FloatEqual(deltaY, vel[1]) {

		if vel[1] < 0 {

			c.onGround = true
		}
		vel[1] = 0
	}
	if !mgl64.FloatEqual(deltaZ, vel[2]) {
		vel[2] = 0
	}
	return mgl64.Vec3{deltaX, deltaY, deltaZ}, vel
}

func blockBBoxsAround(tx *world.Tx, box cube.BBox) []cube.BBox {
	grown := box.Grow(0.25)
	min, max := grown.Min(), grown.Max()
	minX, minY, minZ := int(math.Floor(min[0])), int(math.Floor(min[1])), int(math.Floor(min[2]))

	maxX, maxY, maxZ := int(math.Ceil(max[0])), int(math.Ceil(max[1])), int(math.Ceil(max[2]))

	predicted := (maxX-minX)*(maxY-minY)*(maxZ-minZ) + 2
	var blockBBoxs []cube.BBox
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			for z := minZ; z < maxZ; z++ {
				pos := cube.Pos{x, y, z}
				boxes := tx.Block(pos).Model().BBox(pos, tx)
				if len(boxes) != 0 && blockBBoxs == nil {
					blockBBoxs = make([]cube.BBox, 0, predicted)
				}
				for _, box := range boxes {
					blockBBoxs = append(blockBBoxs, box.Translate(mgl64.Vec3{float64(x), float64(y), float64(z)}))
				}
			}
		}
	}
	return blockBBoxs
}
