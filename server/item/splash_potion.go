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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"math"
)

var SplashPotionSpeed = 0.5

var SplashPotionArc = 0.0

type SplashPotion struct {
	Type potion.Potion
}

func (s SplashPotion) MaxCount() int {
	return 1
}

func (s SplashPotion) Use(tx *world.Tx, user User, ctx *UseContext) bool {
	create := tx.World().EntityRegistry().Config().SplashPotion

	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: throwableOffset(user.Rotation(), SplashPotionArc).Vec3().Mul(SplashPotionSpeed)}
	tx.AddEntity(create(opts, s.Type, user))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}

const vanillaThrowableArc = 26.5

func throwableOffset(r cube.Rotation, arc float64) cube.Rotation {
	r[1] = max(min(r[1], 89.9), -89.9)

	r[1] -= math.Sqrt(max(89.9*89.9-r[1]*r[1], 0)) * (arc / 89.9)
	r[1] = max(min(r[1], 89.9), -89.9)

	return r
}

func (s SplashPotion) EncodeItem() (name string, meta int16) {
	return "minecraft:splash_potion", int16(s.Type.Uint8())
}
