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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type InkSac struct {
	Glowing bool
}

func (i InkSac) UseOnBlock(pos cube.Pos, _ cube.Face, _ mgl64.Vec3, tx *world.Tx, user User, ctx *UseContext) bool {
	if in, ok := tx.Block(pos).(inkable); ok {
		if res, ok := in.Ink(pos, user.Position(), i.Glowing); ok {
			tx.SetBlock(pos, res, nil)
			ctx.SubtractFromCount(1)
			return true
		}
	}
	return false
}

type inkable interface {
	Ink(pos cube.Pos, userPos mgl64.Vec3, glowing bool) (world.Block, bool)
}

func (i InkSac) EncodeItem() (name string, meta int16) {
	if i.Glowing {
		return "minecraft:glow_ink_sac", 0
	}
	return "minecraft:ink_sac", 0
}
