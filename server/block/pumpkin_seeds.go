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

type PumpkinSeeds struct {
	crop

	Direction cube.Face
}

func (PumpkinSeeds) SameCrop(c Crop) bool {
	_, ok := c.(PumpkinSeeds)
	return ok
}

func (p PumpkinSeeds) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		breakBlock(p, pos, tx)
	} else if p.Direction != cube.FaceDown {
		if pumpkin, ok := tx.Block(pos.Side(p.Direction)).(Pumpkin); !ok || pumpkin.Carved {
			p.Direction = cube.FaceDown
			tx.SetBlock(pos, p, nil)
		}
	}
}

func (p PumpkinSeeds) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if r.Float64() <= p.CalculateGrowthChance(pos, tx) && tx.Light(pos) >= 8 {
		if p.Growth < 7 {
			p.Growth++
			tx.SetBlock(pos, p, nil)
		} else {
			directions := []cube.Direction{cube.North, cube.South, cube.West, cube.East}
			for _, i := range directions {
				if _, ok := tx.Block(pos.Side(i.Face())).(Pumpkin); ok {
					return
				}
			}
			direction := directions[r.IntN(len(directions))].Face()
			stemPos := pos.Side(direction)
			if _, ok := tx.Block(stemPos).(Air); ok {
				switch tx.Block(stemPos.Side(cube.FaceDown)).(type) {
				case Farmland, Dirt, Grass:
					p.Direction = direction
					tx.SetBlock(pos, p, nil)
					tx.SetBlock(stemPos, Pumpkin{}, nil)
				}
			}
		}
	}
}

func (p PumpkinSeeds) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if p.Growth == 7 {
		return item.BoneMealResultNone
	}
	p.Growth = min(p.Growth+rand.IntN(4)+2, 7)
	tx.SetBlock(pos, p, nil)
	return item.BoneMealResultSmall
}

func (p PumpkinSeeds) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
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

func (p PumpkinSeeds) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(p))
}

func (PumpkinSeeds) CompostChance() float64 {
	return 0.3
}

func (p PumpkinSeeds) EncodeItem() (name string, meta int16) {
	return "minecraft:pumpkin_seeds", 0
}

func (p PumpkinSeeds) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:pumpkin_stem", map[string]any{"facing_direction": int32(p.Direction), "growth": int32(p.Growth)}
}

func allPumpkinStems() (stems []world.Block) {
	for i := 0; i <= 7; i++ {
		for j := cube.Face(0); j <= 5; j++ {
			stems = append(stems, PumpkinSeeds{Direction: j, crop: crop{Growth: i}})
		}
	}
	return
}
