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
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/customblock"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Activatable interface {
	Activate(pos cube.Pos, clickedFace cube.Face, tx *world.Tx, u item.User, ctx *item.UseContext) bool
}

type Pickable interface {
	Pick() item.Stack
}

type Punchable interface {
	Punch(pos cube.Pos, clickedFace cube.Face, tx *world.Tx, u item.User)
}

type LightEmitter interface {
	LightEmissionLevel() uint8
}

type LightDiffuser interface {
	LightDiffusionLevel() uint8
}

type NonSuffocating interface {
	PreventsSuffocation() bool
}

type RedstoneWireStepDowner interface {
	CanRedstoneWireStepDown(pos, from cube.Pos, tx *world.Tx) bool
}

type Replaceable interface {
	ReplaceableBy(b world.Block) bool
}

type EntityLander interface {
	EntityLand(pos cube.Pos, tx *world.Tx, e world.Entity, distance *float64)
}

type EntityInsider interface {
	EntityInside(pos cube.Pos, tx *world.Tx, e world.Entity)
}

type EntityStepper interface {
	EntityStepOn(pos cube.Pos, tx *world.Tx, e world.Entity)
}

type ProjectileHitter interface {
	ProjectileHit(pos cube.Pos, tx *world.Tx, e world.Entity, face cube.Face)
}

type Frictional interface {
	Friction() float64
}

type Permutable interface {
	States() map[string][]any

	Permutations() []customblock.Permutation
}

var unknownFace = cube.Face(len(cube.Faces()))

var unknownDirection = cube.Direction(len(cube.Directions()))

func calculateFace(user item.User, placePos cube.Pos) cube.Face {
	userPos := user.Position()
	pos := cube.PosFromVec3(userPos)
	if abs(pos[0]-placePos[0]) < 2 && abs(pos[2]-placePos[2]) < 2 {
		y := userPos[1]
		if eyed, ok := user.(interface{ EyeHeight() float64 }); ok {
			y += eyed.EyeHeight()
		}

		if y-float64(placePos[1]) > 2.0 {
			return cube.FaceUp
		} else if float64(placePos[1])-y > 0.0 {
			return cube.FaceDown
		}
	}
	return user.Rotation().Direction().Opposite().Face()
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func replaceableWith(tx *world.Tx, pos cube.Pos, with world.Block) bool {
	if pos.OutOfBounds(tx.Range()) {
		return false
	}
	b := tx.Block(pos)
	if replaceable, ok := b.(Replaceable); ok {
		if !replaceable.ReplaceableBy(with) || b == with {
			return false
		}
		if liquid, ok := tx.Liquid(pos); ok {
			replaceable, ok := liquid.(Replaceable)
			return ok && replaceable.ReplaceableBy(with)
		}
		return true
	}
	return false
}

func firstReplaceable(tx *world.Tx, pos cube.Pos, face cube.Face, with world.Block) (cube.Pos, cube.Face, bool) {
	if replaceableWith(tx, pos, with) {

		return pos, cube.FaceUp, true
	}
	side := pos.Side(face)
	if replaceableWith(tx, side, with) {
		return side, face, true
	}
	return pos, face, false
}

func place(tx *world.Tx, pos cube.Pos, b world.Block, user item.User, ctx *item.UseContext) {
	if placer, ok := user.(Placer); ok {
		placer.PlaceBlock(pos, b, ctx)
		return
	}
	tx.SetBlock(pos, b, nil)
	tx.PlaySound(pos.Vec3(), sound.BlockPlace{Block: b})
}

func horizontalDirection(d cube.Direction) cube.Direction {
	switch d {
	case cube.South:
		return cube.North
	case cube.West:
		return cube.South
	case cube.North:
		return cube.West
	case cube.East:
		return cube.East
	}
	panic("invalid direction")
}

func placed(ctx *item.UseContext) bool {
	return ctx.CountSub > 0
}

func boolByte(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

type replaceable struct{}

func (replaceable) ReplaceableBy(world.Block) bool {
	return true
}

type transparent struct{}

func (transparent) LightDiffusionLevel() uint8 {
	return 0
}

type gravityAffected struct{}

func (g gravityAffected) Solidifies(cube.Pos, *world.Tx) bool {
	return false
}

func (g gravityAffected) fall(b world.Block, pos cube.Pos, tx *world.Tx) {
	if replaceableWith(tx, pos.Side(cube.FaceDown), b) {
		tx.SetBlock(pos, nil, nil)
		opts := world.EntitySpawnOpts{Position: pos.Vec3Centre()}
		tx.AddEntity(tx.World().EntityRegistry().Config().FallingBlock(opts, b))
	}
}

type Flammable interface {
	FlammabilityInfo() FlammabilityInfo
}

type FlammabilityInfo struct {
	Encouragement int

	Flammability int

	LavaFlammable bool
}

func newFlammabilityInfo(encouragement, flammability int, lavaFlammable bool) FlammabilityInfo {
	return FlammabilityInfo{
		Encouragement: encouragement,
		Flammability:  flammability,
		LavaFlammable: lavaFlammable,
	}
}

type livingEntity interface {
	Hurt(damage float64, src world.DamageSource) (n float64, vulnerable bool)
}

type flammableEntity interface {
	OnFireDuration() time.Duration

	SetOnFire(duration time.Duration)

	Extinguish()
}

func dropItem(tx *world.Tx, it item.Stack, pos mgl64.Vec3) {
	create := tx.World().EntityRegistry().Config().Item
	opts := world.EntitySpawnOpts{Position: pos, Velocity: mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1}}
	tx.AddEntity(create(opts, it))
}

type bass struct{}

func (bass) Instrument() sound.Instrument {
	return sound.Bass()
}

type snare struct{}

func (snare) Instrument() sound.Instrument {
	return sound.Snare()
}

type clicksAndSticks struct{}

func (clicksAndSticks) Instrument() sound.Instrument {
	return sound.ClicksAndSticks()
}

type bassDrum struct{}

func (bassDrum) Instrument() sound.Instrument {
	return sound.BassDrum()
}

type flute struct{}

func (flute) Instrument() sound.Instrument {
	return sound.Flute()
}

func newSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
	}
}

func newFoodSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
		Food:       true,
	}
}

func newOreSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
		Ores:       true,
	}
}

func newFuelInfo(duration time.Duration) item.FuelInfo {
	return item.FuelInfo{Duration: duration}
}
