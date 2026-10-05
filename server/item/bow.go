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
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"math"
	"time"
)

var ArrowSpeed = 2.8

type Bow struct{}

func (Bow) MaxCount() int {
	return 1
}

func (Bow) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 385,
		BrokenItem:    simpleItem(Stack{}),
	}
}

func (Bow) FuelInfo() FuelInfo {
	return newFuelInfo(time.Second * 10)
}

func (Bow) Release(releaser Releaser, tx *world.Tx, ctx *UseContext, duration time.Duration) {
	creative := releaser.GameMode().CreativeInventory()
	ticks := duration.Milliseconds() / 50
	if ticks < 3 {

		return
	}

	d := float64(ticks) / 20
	force := math.Min((d*d+d*2)/3, 1)
	if force < 0.1 {

		return
	}

	arrow, ok := ctx.FirstFunc(func(stack Stack) bool {
		_, ok := stack.Item().(Arrow)
		return ok
	})
	if !ok && !creative {

		return
	}

	var tip potion.Potion
	if !arrow.Empty() {

		tip = arrow.Item().(Arrow).Tip
	}

	held, _ := releaser.HeldItems()
	damage, punchLevel, burnDuration, consume := 2.0, 0, time.Duration(0), !creative
	for _, enchant := range held.Enchantments() {
		if f, ok := enchant.Type().(interface{ BurnDuration() time.Duration }); ok {
			burnDuration = f.BurnDuration()
		}
		if _, ok := enchant.Type().(interface{ KnockBackMultiplier() float64 }); ok {
			punchLevel = enchant.Level()
		}
		if p, ok := enchant.Type().(interface{ PowerDamage(int) float64 }); ok {
			damage += p.PowerDamage(enchant.Level())
		}
		if i, ok := enchant.Type().(interface{ ConsumesArrows() bool }); ok && !i.ConsumesArrows() {
			consume = false
		}
	}

	create := tx.World().EntityRegistry().Config().Arrow
	opts := world.EntitySpawnOpts{
		Position: eyePosition(releaser),
		Velocity: releaser.Rotation().Vec3().Mul(force * ArrowSpeed),
		Rotation: releaser.Rotation().Neg(),
	}
	projectile := tx.AddEntity(create(opts, world.ArrowSpawnConfig{
		Damage:              damage,
		Owner:               releaser,
		Critical:            force >= 1,
		ObtainArrowOnPickup: !creative && consume,
		PunchLevel:          punchLevel,
		Tip:                 tip,
	}))
	if f, ok := projectile.(interface{ SetOnFire(duration time.Duration) }); ok {
		f.SetOnFire(burnDuration)
	}

	ctx.DamageItem(1)
	if consume {
		ctx.Consume(arrow.Grow(-arrow.Count() + 1))
	}

	tx.PlaySound(releaser.Position(), sound.BowShoot{})
}

func (Bow) EnchantmentValue() int {
	return 1
}

func (Bow) Requirements() []Stack {
	return []Stack{NewStack(Arrow{}, 1)}
}

func (Bow) EncodeItem() (name string, meta int16) {
	return "minecraft:bow", 0
}
