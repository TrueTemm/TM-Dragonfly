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

type MelonSeeds struct {
	crop

	Direction cube.Face
}

func (MelonSeeds) SameCrop(c Crop) bool {
	_, ok := c.(MelonSeeds)
	return ok
}

func (m MelonSeeds) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		breakBlock(m, pos, tx)
	} else if m.Direction != cube.FaceDown {
		if _, ok := tx.Block(pos.Side(m.Direction)).(Melon); !ok {
			m.Direction = cube.FaceDown
			tx.SetBlock(pos, m, nil)
		}
	}
}

func (m MelonSeeds) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if r.Float64() <= m.CalculateGrowthChance(pos, tx) && tx.Light(pos) >= 8 {
		if m.Growth < 7 {
			m.Growth++
			tx.SetBlock(pos, m, nil)
		} else {
			directions := cube.Directions()
			for _, i := range directions {
				if _, ok := tx.Block(pos.Side(i.Face())).(Melon); ok {
					return
				}
			}
			direction := directions[r.IntN(len(directions))].Face()
			stemPos := pos.Side(direction)
			if _, ok := tx.Block(stemPos).(Air); ok {
				switch tx.Block(stemPos.Side(cube.FaceDown)).(type) {
				case Farmland, Dirt, Grass:
					m.Direction = direction
					tx.SetBlock(pos, m, nil)
					tx.SetBlock(stemPos, Melon{}, nil)
				}
			}
		}
	}
}

func (m MelonSeeds) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if m.Growth == 7 {
		return item.BoneMealResultNone
	}
	m.Growth = min(m.Growth+rand.IntN(4)+2, 7)
	tx.SetBlock(pos, m, nil)
	return item.BoneMealResultSmall
}

func (m MelonSeeds) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, m)
	if !used {
		return false
	}

	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		return false
	}

	place(tx, pos, m, user, ctx)
	return placed(ctx)
}

func (m MelonSeeds) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(m))
}

func (MelonSeeds) CompostChance() float64 {
	return 0.3
}

func (m MelonSeeds) EncodeItem() (name string, meta int16) {
	return "minecraft:melon_seeds", 0
}

func (m MelonSeeds) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:melon_stem", map[string]any{"facing_direction": int32(m.Direction), "growth": int32(m.Growth)}
}

func allMelonStems() (stems []world.Block) {
	for i := 0; i <= 7; i++ {
		for j := cube.Face(0); j <= 5; j++ {
			stems = append(stems, MelonSeeds{crop: crop{Growth: i}, Direction: j})
		}
	}
	return
}
