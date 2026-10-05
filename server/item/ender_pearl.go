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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"time"
)

var EnderPearlSpeed = 2.2

type EnderPearl struct{}

func (e EnderPearl) Use(tx *world.Tx, user User, ctx *UseContext) bool {
	create := tx.World().EntityRegistry().Config().EnderPearl
	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: user.Rotation().Vec3().Mul(EnderPearlSpeed)}
	tx.AddEntity(create(opts, user))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}

func (EnderPearl) Cooldown() time.Duration {
	return EnderPearlCooldown
}

var EnderPearlCooldown = 16 * time.Second

func (EnderPearl) MaxCount() int {
	return 16
}

func (EnderPearl) EncodeItem() (name string, meta int16) {
	return "minecraft:ender_pearl", 0
}
