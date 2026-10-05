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
	"sync"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

type FishingRod struct{}

var FishingHookSpeed = 1.5

var hooks sync.Map

func (FishingRod) Use(tx *world.Tx, user User, ctx *UseContext) bool {
	key := user.H()
	if out, ok := hooks.LoadAndDelete(key); ok {
		if hook, ok := out.(*world.EntityHandle).Entity(tx); ok {
			_ = hook.Close()
			ctx.DamageItem(1)
			return true
		}
	}
	create := tx.World().EntityRegistry().Config().FishingHook
	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: user.Rotation().Vec3().Mul(FishingHookSpeed)}
	hook := create(opts, user)
	tx.AddEntity(hook)
	hooks.Store(key, hook)
	tx.PlaySound(user.Position(), sound.ItemThrow{})
	ctx.DamageItem(1)
	return true
}

func (FishingRod) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    384,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 2,
		BreakDurability:  2,
	}
}

func (FishingRod) MaxCount() int {
	return 1
}

func (FishingRod) EncodeItem() (name string, meta int16) {
	return "minecraft:fishing_rod", 0
}
