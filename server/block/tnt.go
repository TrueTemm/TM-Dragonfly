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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

type TNT struct {
	solid
}

var _ world.RedstonePowerAction = TNT{}

func (TNT) RedstoneNonConductive() {}

func (t TNT) RedstonePowerAction(pos cube.Pos, tx *world.Tx, oldPower, newPower int) {
	if oldPower > 0 || newPower == 0 {
		return
	}
	t.Ignite(pos, tx, nil)
}

func (t TNT) ProjectileHit(pos cube.Pos, tx *world.Tx, e world.Entity, _ cube.Face) {
	if f, ok := e.(flammableEntity); ok && f.OnFireDuration() > 0 {
		t.Ignite(pos, tx, nil)
	}
}

func (t TNT) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, ctx *item.UseContext) bool {
	held, _ := u.HeldItems()
	if _, ok := held.Enchantment(enchantment.FireAspect); ok {
		t.Ignite(pos, tx, nil)
		ctx.DamageItem(1)
		return true
	}
	return false
}

func (t TNT) Ignite(pos cube.Pos, tx *world.Tx, _ world.Entity) bool {
	spawnTnt(pos, tx, time.Second*4)
	return true
}

func (t TNT) Explode(_ world.ExplosionSource, pos cube.Pos, tx *world.Tx) {
	spawnTnt(pos, tx, time.Second/2+time.Duration(rand.IntN(int(time.Second+time.Second/2))))
}

func (t TNT) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(t))
}

func (t TNT) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(15, 100, true)
}

func (t TNT) EncodeItem() (name string, meta int16) {
	return "minecraft:tnt", 0
}

func (t TNT) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:tnt", map[string]interface{}{"explode_bit": false}
}

func spawnTnt(pos cube.Pos, tx *world.Tx, fuse time.Duration) {
	tx.PlaySound(pos.Vec3Centre(), sound.TNT{})
	tx.SetBlock(pos, nil, nil)
	opts := world.EntitySpawnOpts{Position: pos.Vec3Centre()}
	tx.AddEntity(tx.World().EntityRegistry().Config().TNT(opts, fuse))
}
