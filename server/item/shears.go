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

type Shears struct{}

func (s Shears) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	if face == cube.FaceUp || face == cube.FaceDown {

		return false
	}
	if c, ok := tx.Block(pos).(carvable); ok {
		if res, ok := c.Carve(face); ok {

			tx.SetBlock(pos, res, nil)

			ctx.DamageItem(1)
			return true
		}
	}
	return false
}

type carvable interface {
	Carve(f cube.Face) (world.Block, bool)
}

func (s Shears) ToolType() ToolType {
	return TypeShears
}

func (s Shears) HarvestLevel() int {
	return 1
}

func (s Shears) BaseMiningEfficiency(world.Block) float64 {
	return 1.5
}

func (s Shears) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    238,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 0,
		BreakDurability:  1,
	}
}

func (s Shears) MaxCount() int {
	return 1
}

func (s Shears) EncodeItem() (name string, meta int16) {
	return "minecraft:shears", 0
}
