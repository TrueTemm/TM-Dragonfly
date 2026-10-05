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

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type NetherWart struct {
	transparent
	empty

	Age int
}

func (n NetherWart) HasLiquidDrops() bool {
	return true
}

func (n NetherWart) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if n.Age < 3 && r.Float64() < 0.1 {
		n.Age++
		tx.SetBlock(pos, n, nil)
	}
}

func (n NetherWart) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, n)
	if !used {
		return false
	}
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(SoulSand); !ok {
		return false
	}

	place(tx, pos, n, user, ctx)
	return placed(ctx)
}

func (n NetherWart) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(SoulSand); !ok {
		breakBlock(n, pos, tx)
	}
}

func (n NetherWart) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if n.Age < 3 {
			return []item.Stack{item.NewStack(n, 1)}
		}
		return []item.Stack{item.NewStack(n, fortuneDiscreteCount(2, 4, 7, enchantments))}
	})
}

func (NetherWart) CompostChance() float64 {
	return 0.65
}

func (NetherWart) EncodeItem() (name string, meta int16) {
	return "minecraft:nether_wart", 0
}

func (n NetherWart) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:nether_wart", map[string]any{"age": int32(n.Age)}
}

func allNetherWart() (wart []world.Block) {
	for i := 0; i < 4; i++ {
		wart = append(wart, NetherWart{Age: i})
	}
	return
}
