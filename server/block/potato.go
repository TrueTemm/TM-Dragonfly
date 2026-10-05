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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Potato struct {
	crop
}

func (p Potato) SmeltInfo() item.SmeltInfo {
	return newFoodSmeltInfo(item.NewStack(item.BakedPotato{}, 1), 0.35)
}

func (Potato) SameCrop(c Crop) bool {
	_, ok := c.(Potato)
	return ok
}

func (p Potato) AlwaysConsumable() bool {
	return false
}

func (p Potato) ConsumeDuration() time.Duration {
	return item.DefaultConsumeDuration
}

func (p Potato) Consume(_ *world.Tx, c item.Consumer) item.Stack {
	c.Saturate(1, 0.6)
	return item.Stack{}
}

func (p Potato) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if p.Growth == 7 {
		return item.BoneMealResultNone
	}
	p.Growth = min(p.Growth+rand.IntN(4)+2, 7)
	tx.SetBlock(pos, p, nil)
	return item.BoneMealResultSmall
}

func (p Potato) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, p)
	if !used {
		return false
	}

	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		return false
	}

	place(tx, pos, p, user, ctx)
	return placed(ctx)
}

func (p Potato) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if p.Growth < 7 {
			return []item.Stack{item.NewStack(p, 1)}
		}
		fortune := fortuneLevel(enchantments)
		count := rand.IntN(fortune+1) + 1 + fortuneBinomial(3+fortune)
		if rand.Float64() < 0.02 {
			return []item.Stack{item.NewStack(p, count), item.NewStack(item.PoisonousPotato{}, 1)}
		}
		return []item.Stack{item.NewStack(p, count)}
	})
}

func (Potato) CompostChance() float64 {
	return 0.65
}

func (p Potato) EncodeItem() (name string, meta int16) {
	return "minecraft:potato", 0
}

func (p Potato) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if tx.Light(pos) < 8 {
		breakBlock(p, pos, tx)
	} else if p.Growth < 7 && r.Float64() <= p.CalculateGrowthChance(pos, tx) {
		p.Growth++
		tx.SetBlock(pos, p, nil)
	}
}

func (p Potato) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:potatoes", map[string]any{"growth": int32(p.Growth)}
}

func allPotato() (potato []world.Block) {
	for i := 0; i <= 7; i++ {
		potato = append(potato, Potato{crop{Growth: i}})
	}
	return
}
