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
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

type Crossbow struct {
	Item Stack
}

func (c Crossbow) Charge(releaser Releaser, _ *world.Tx, ctx *UseContext, duration time.Duration) bool {
	if !c.Item.Empty() {
		return false
	}

	creative := releaser.GameMode().CreativeInventory()
	held, left := releaser.HeldItems()

	if chargeDuration, _ := c.chargeDuration(held); duration < chargeDuration {
		return false
	}
	projectileItem, ok := c.findProjectile(releaser, ctx)
	if !ok {
		return false
	}
	c.Item = projectileItem.Grow(-projectileItem.Count() + 1)
	if !creative {
		ctx.Consume(c.Item)
	}

	releaser.SetHeldItems(held.WithItem(c), left)
	return true
}

func (c Crossbow) ContinueCharge(releaser Releaser, tx *world.Tx, ctx *UseContext, duration time.Duration) {
	if !c.Item.Empty() {
		return
	}

	held, _ := releaser.HeldItems()
	if _, ok := c.findProjectile(releaser, ctx); !ok {
		return
	}

	chargeDuration, qcLevel := c.chargeDuration(held)
	if duration.Seconds() <= 0.1 {
		tx.PlaySound(releaser.Position(), sound.CrossbowLoad{Stage: sound.CrossbowLoadingStart, QuickCharge: qcLevel > 0})
	}

	multiplier := 25.0 / float64(25-(5*qcLevel))

	adjustedTicks := int(float64(duration.Milliseconds()) / (50 / multiplier))

	if adjustedTicks%16 == 0 {
		tx.PlaySound(releaser.Position(), sound.CrossbowLoad{Stage: sound.CrossbowLoadingMiddle, QuickCharge: qcLevel > 0})
	}

	if progress := float64(duration) / float64(chargeDuration); progress >= 1 {
		tx.PlaySound(releaser.Position(), sound.CrossbowLoad{Stage: sound.CrossbowLoadingEnd, QuickCharge: qcLevel > 0})
	}
}

func (c Crossbow) chargeDuration(s Stack) (dur time.Duration, quickChargeLvl int) {
	dur, lvl := time.Duration(1.25*float64(time.Second)), 0
	for _, enchant := range s.Enchantments() {
		if q, ok := enchant.Type().(interface{ ChargeDuration(int) time.Duration }); ok {
			dur = min(dur, q.ChargeDuration(enchant.Level()))
			lvl = enchant.Level()
		}
	}
	return dur, lvl
}

func (c Crossbow) findProjectile(r Releaser, ctx *UseContext) (Stack, bool) {
	_, left := r.HeldItems()
	_, isFirework := left.Item().(Firework)
	_, isArrow := left.Item().(Arrow)
	if isFirework || isArrow {
		return left, true
	}
	if res, ok := ctx.FirstFunc(func(stack Stack) bool {
		_, ok := stack.Item().(Arrow)
		return ok
	}); ok {
		return res, true
	}
	if r.GameMode().CreativeInventory() {

		return NewStack(Arrow{}, 1), true
	}
	return Stack{}, false
}

func (c Crossbow) ReleaseCharge(releaser Releaser, tx *world.Tx, ctx *UseContext) bool {
	if c.Item.Empty() {
		return false
	}

	held, _ := releaser.HeldItems()
	creative := releaser.GameMode().CreativeInventory()

	pierceLevel, multishot := 0, false
	for _, enchant := range held.Enchantments() {
		if _, ok := enchant.Type().(interface{ MultipleProjectiles() bool }); ok {
			multishot = true
		}
		if _, ok := enchant.Type().(interface{ Pierces() bool }); ok {
			pierceLevel = enchant.Level()
		}
	}

	arrowConf := world.ArrowSpawnConfig{
		Damage:              9,
		Owner:               releaser,
		Critical:            true,
		ObtainArrowOnPickup: !creative,
		PiercingLevel:       pierceLevel,
	}
	c.shoot(releaser, tx, 0, arrowConf)
	if multishot {
		arrowConf.ObtainArrowOnPickup = false
		c.shoot(releaser, tx, -10, arrowConf)
		c.shoot(releaser, tx, 10, arrowConf)
	}
	c.applyDamage(ctx)

	c.Item = Stack{}
	held, left := releaser.HeldItems()
	crossbow := held.WithItem(c)
	releaser.SetHeldItems(crossbow, left)
	tx.PlaySound(releaser.Position(), sound.CrossbowShoot{})
	return true
}

func (c Crossbow) CanCharge(releaser Releaser, _ *world.Tx, ctx *UseContext) bool {
	_, found := c.findProjectile(releaser, ctx)
	return found && c.Item.Empty()
}

func (c Crossbow) shoot(releaser Releaser, tx *world.Tx, offsetAngle float64, arrowConf world.ArrowSpawnConfig) {
	rot := releaser.Rotation()
	dirVec := cube.Rotation{rot[0] + offsetAngle, rot[1]}.Vec3()

	if firework, ok := c.Item.Item().(Firework); ok {
		createFirework := tx.World().EntityRegistry().Config().Firework
		projectile := createFirework(world.EntitySpawnOpts{
			Position: torsoPosition(releaser),
			Velocity: dirVec.Mul(0.8),
			Rotation: rot.Neg(),
		}, firework, releaser, 1.0, 0, false)
		tx.AddEntity(projectile)
	} else {
		createArrow := tx.World().EntityRegistry().Config().Arrow
		arrowConf.Tip = c.Item.Item().(Arrow).Tip
		arrow := createArrow(world.EntitySpawnOpts{
			Position: torsoPosition(releaser),
			Velocity: dirVec.Mul(5.15),
			Rotation: rot.Neg(),
		}, arrowConf)
		tx.AddEntity(arrow)
	}
}

func (c Crossbow) applyDamage(ctx *UseContext) {
	if _, ok := c.Item.Item().(Firework); ok {
		ctx.DamageItem(3)
	} else {
		ctx.DamageItem(1)
	}
}

func (Crossbow) MaxCount() int {
	return 1
}

func (Crossbow) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 464,
		BrokenItem:    simpleItem(Stack{}),
	}
}

func (Crossbow) FuelInfo() FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (Crossbow) EnchantmentValue() int {
	return 1
}

func (Crossbow) EncodeItem() (name string, meta int16) {
	return "minecraft:crossbow", 0
}

func (c Crossbow) DecodeNBT(data map[string]any) any {
	c.Item = MapNBT(data, "chargedItem")
	return c
}

func (c Crossbow) EncodeNBT() map[string]any {
	if !c.Item.Empty() {
		return map[string]any{"chargedItem": WriteNBT(c.Item, true)}
	}
	return nil
}
