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

type Kelp struct {
	empty
	transparent
	sourceWaterDisplacer

	Age int
}

func (k Kelp) SmeltInfo() item.SmeltInfo {
	return newFoodSmeltInfo(item.NewStack(item.DriedKelp{}, 1), 0.1)
}

func (k Kelp) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	for y := pos.Y(); y <= tx.Range()[1]; y++ {
		currentPos := cube.Pos{pos.X(), y, pos.Z()}
		block := tx.Block(currentPos)
		if kelp, ok := block.(Kelp); ok {
			if kelp.Age == 25 {
				break
			}
			continue
		}
		if water, ok := block.(Water); ok && water.Depth == 8 {
			tx.SetBlock(currentPos, Kelp{Age: k.Age + 1}, nil)
			return item.BoneMealResultSmall
		}
		break
	}
	return item.BoneMealResultNone
}

func (k Kelp) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(k))
}

func (Kelp) CompostChance() float64 {
	return 0.3
}

func (Kelp) EncodeItem() (name string, meta int16) {
	return "minecraft:kelp", 0
}

func (k Kelp) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:kelp", map[string]any{"kelp_age": int32(k.Age)}
}

func (Kelp) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (k Kelp) withRandomAge() Kelp {
	k.Age = rand.IntN(25)
	return k
}

func (k Kelp) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, k)
	if !used {
		return
	}

	below := pos.Side(cube.FaceDown)
	belowBlock := tx.Block(below)
	if _, kelp := belowBlock.(Kelp); !kelp {
		if !belowBlock.Model().FaceSolid(below, cube.FaceUp, tx) {
			return false
		}
	}

	liquid, ok := tx.Liquid(pos)
	if !ok {
		return false
	} else if _, ok := liquid.(Water); !ok || liquid.LiquidDepth() < 8 {
		return false
	}

	place(tx, pos, k.withRandomAge(), user, ctx)
	return placed(ctx)
}

func (k Kelp) NeighbourUpdateTick(pos, changedNeighbour cube.Pos, tx *world.Tx) {
	if _, ok := tx.Liquid(pos); !ok {
		breakBlock(k, pos, tx)
		return
	}
	if changedNeighbour[1]-1 == pos.Y() {

		tx.SetBlock(pos, k.withRandomAge(), nil)
	}

	below := pos.Side(cube.FaceDown)
	belowBlock := tx.Block(below)
	if _, kelp := belowBlock.(Kelp); !kelp {
		if !belowBlock.Model().FaceSolid(below, cube.FaceUp, tx) {
			breakBlock(k, pos, tx)
		}
	}
}

func (k Kelp) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {

	if r.IntN(100) < 15 && k.Age < 25 {
		abovePos := pos.Side(cube.FaceUp)

		liquid, ok := tx.Liquid(abovePos)

		if !ok {
			return
		} else if _, ok := liquid.(Water); ok {
			switch tx.Block(abovePos).(type) {
			case Air, Water:
				tx.SetBlock(abovePos, Kelp{Age: k.Age + 1}, nil)
				if liquid.LiquidDepth() < 8 {

					tx.SetLiquid(abovePos, Water{Still: true, Depth: 8, Falling: false})
				}
			}
		}
	}
}

func allKelp() (b []world.Block) {
	for i := 0; i < 26; i++ {
		b = append(b, Kelp{Age: i})
	}
	return
}
