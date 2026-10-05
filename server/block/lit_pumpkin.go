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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type LitPumpkin struct {
	solid

	Facing cube.Direction
}

func (l LitPumpkin) LightEmissionLevel() uint8 {
	return 15
}

func (l LitPumpkin) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, l)
	if !used {
		return
	}
	l.Facing = user.Rotation().Direction().Opposite()

	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

func (l LitPumpkin) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, axeEffective, oneOf(l))
}

func (l LitPumpkin) EncodeItem() (name string, meta int16) {
	return "minecraft:lit_pumpkin", 0
}

func (l LitPumpkin) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:lit_pumpkin", map[string]any{"minecraft:cardinal_direction": l.Facing.String()}
}

func allLitPumpkins() (pumpkins []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		pumpkins = append(pumpkins, LitPumpkin{Facing: i})
	}
	return
}
