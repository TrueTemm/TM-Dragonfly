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
)

type LingeringPotion struct {
	Type potion.Potion
}

func (l LingeringPotion) MaxCount() int {
	return 1
}

func (l LingeringPotion) Use(tx *world.Tx, user User, ctx *UseContext) bool {
	create := tx.World().EntityRegistry().Config().LingeringPotion
	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: throwableOffset(user.Rotation(), vanillaThrowableArc).Vec3().Mul(0.5)}
	tx.AddEntity(create(opts, l.Type, user))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}

func (l LingeringPotion) EncodeItem() (name string, meta int16) {
	return "minecraft:lingering_potion", int16(l.Type.Uint8())
}
