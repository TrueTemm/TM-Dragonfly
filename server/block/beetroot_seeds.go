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

type BeetrootSeeds struct {
	crop
}

func (BeetrootSeeds) SameCrop(c Crop) bool {
	_, ok := c.(BeetrootSeeds)
	return ok
}

func (b BeetrootSeeds) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if b.Growth == 7 {
		return item.BoneMealResultNone
	}
	if rand.Float64() < 0.75 {
		b.Growth++
		tx.SetBlock(pos, b, nil)
		return item.BoneMealResultSmall
	}
	return item.BoneMealResultNone
}

func (b BeetrootSeeds) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, b)
	if !used {
		return false
	}

	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		return false
	}

	place(tx, pos, b, user, ctx)
	return placed(ctx)
}

func (b BeetrootSeeds) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, cropSeedDrops(b, item.Beetroot{}, b.Growth))
}

func (BeetrootSeeds) CompostChance() float64 {
	return 0.3
}

func (b BeetrootSeeds) EncodeItem() (name string, meta int16) {
	return "minecraft:beetroot_seeds", 0
}

func (b BeetrootSeeds) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if tx.Light(pos) < 8 {
		breakBlock(b, pos, tx)
	} else if b.Growth < 7 && r.IntN(3) > 0 && r.Float64() <= b.CalculateGrowthChance(pos, tx) {
		b.Growth++
		tx.SetBlock(pos, b, nil)
	}
}

func (b BeetrootSeeds) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:beetroot", map[string]any{"growth": int32(b.Growth)}
}

func allBeetroot() (beetroot []world.Block) {
	for i := 0; i <= 7; i++ {
		beetroot = append(beetroot, BeetrootSeeds{crop: crop{Growth: i}})
	}
	return
}
