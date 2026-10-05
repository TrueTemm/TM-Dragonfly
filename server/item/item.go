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

package item

import (
	"encoding/binary"
	"image/color"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type MaxCounter interface {
	MaxCount() int
}

type UsableOnBlock interface {
	UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user User, ctx *UseContext) bool
}

type UsableOnEntity interface {
	UseOnEntity(e world.Entity, tx *world.Tx, user User, ctx *UseContext) bool
}

type Usable interface {
	Use(tx *world.Tx, user User, ctx *UseContext) bool
}

type Throwable interface {
	SwingAnimation() bool
}

type OffHand interface {
	OffHand() bool
}

type Consumable interface {
	AlwaysConsumable() bool

	ConsumeDuration() time.Duration

	Consume(tx *world.Tx, c Consumer) Stack
}

type Consumer interface {
	User

	Saturate(food int, saturation float64)

	AddEffect(e effect.Effect)

	RemoveEffect(e effect.Type)

	Effects() []effect.Effect

	Absorption() float64

	SetAbsorption(health float64)
}

const DefaultConsumeDuration = (time.Second * 161) / 100

type Drinkable interface {
	Drinkable() bool
}

type Glinted interface {
	Glinted() bool
}

type HandEquipped interface {
	HandEquipped() bool
}

type Weapon interface {
	AttackDamage() float64
}

type Cooldown interface {
	Cooldown() time.Duration
}

type nameable interface {
	WithName(a ...any) world.Item
}

type Releaser interface {
	User

	GameMode() world.GameMode

	PlaySound(sound world.Sound)
}

type Releasable interface {
	Release(releaser Releaser, tx *world.Tx, ctx *UseContext, duration time.Duration)

	Requirements() []Stack
}

type Chargeable interface {
	Charge(releaser Releaser, tx *world.Tx, ctx *UseContext, duration time.Duration) bool

	ContinueCharge(releaser Releaser, tx *world.Tx, ctx *UseContext, duration time.Duration)

	ReleaseCharge(releaser Releaser, tx *world.Tx, ctx *UseContext) bool

	CanCharge(releaser Releaser, tx *world.Tx, ctx *UseContext) bool
}

type User interface {
	Carrier
	SetHeldItems(mainHand, offHand Stack)

	UsingItem() bool
	ReleaseItem()
	UseItem()
}

type Carrier interface {
	world.Entity

	HeldItems() (mainHand, offHand Stack)
}

type BeaconPayment interface {
	PayableForBeacon() bool
}

type Compostable interface {
	CompostChance() float64
}

type nopReleasable struct{}

func (nopReleasable) Release(Releaser, *world.Tx, *UseContext, time.Duration) {}
func (nopReleasable) Requirements() []Stack {
	return []Stack{}
}

type defaultFood struct{}

func (defaultFood) AlwaysConsumable() bool {
	return false
}

func (d defaultFood) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}

func eyePosition(e world.Entity) mgl64.Vec3 {
	pos := e.Position()
	if eyed, ok := e.(interface{ EyeHeight() float64 }); ok {
		pos = pos.Add(mgl64.Vec3{0, eyed.EyeHeight()})
	}
	return pos
}

func torsoPosition(e world.Entity) mgl64.Vec3 {
	pos := e.Position()
	if torso, ok := e.(interface{ TorsoHeight() float64 }); ok {
		pos = pos.Add(mgl64.Vec3{0, torso.TorsoHeight()})
	}
	return pos
}

func int32FromRGBA(x color.RGBA) int32 {
	if x.R == 0 && x.G == 0 && x.B == 0 {

		return int32(-0x1000000)
	}
	return int32(binary.BigEndian.Uint32([]byte{x.A, x.R, x.G, x.B}))
}

func rgbaFromInt32(x int32) color.RGBA {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(x))

	return color.RGBA{A: b[0], R: b[1], G: b[2], B: b[3]}
}

func boolByte(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
