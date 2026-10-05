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
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
)

type Cactus struct {
	transparent

	Age int
}

func (c Cactus) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, c)
	if !used || !c.canGrowHere(pos, tx, true) {
		return false
	}

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (c Cactus) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !c.canGrowHere(pos, tx, true) {
		breakBlock(c, pos, tx)
	}
}

func (c Cactus) RandomTick(pos cube.Pos, tx *world.Tx, _ *rand.Rand) {
	if c.Age < 15 {
		c.Age++
	} else if c.Age == 15 {
		c.Age = 0
		if c.canGrowHere(pos.Side(cube.FaceDown), tx, false) {
			for y := 1; y < 3; y++ {
				if _, ok := tx.Block(pos.Add(cube.Pos{0, y})).(Air); ok {
					tx.SetBlock(pos.Add(cube.Pos{0, y}), Cactus{Age: 0}, nil)
					break
				} else if _, ok := tx.Block(pos.Add(cube.Pos{0, y})).(Cactus); !ok {
					break
				}
			}
		}
	}
	tx.SetBlock(pos, c, nil)
}

func (c Cactus) canGrowHere(pos cube.Pos, tx *world.Tx, recursive bool) bool {
	for _, face := range cube.HorizontalFaces() {
		if _, ok := tx.Block(pos.Side(face)).(Air); !ok {
			return false
		}
	}
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Cactus); ok && recursive {
		return c.canGrowHere(pos.Side(cube.FaceDown), tx, recursive)
	}
	return supportsVegetation(c, tx.Block(pos.Sub(cube.Pos{0, 1})))
}

func (c Cactus) EntityInside(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if l, ok := e.(livingEntity); ok {
		l.Hurt(0.5, DamageSource{Block: c})
	}
}

func (c Cactus) BreakInfo() BreakInfo {
	return newBreakInfo(0.4, alwaysHarvestable, nothingEffective, oneOf(c))
}

func (Cactus) CompostChance() float64 {
	return 0.5
}

func (c Cactus) EncodeItem() (name string, meta int16) {
	return "minecraft:cactus", 0
}

func (c Cactus) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:cactus", map[string]any{"age": int32(c.Age)}
}

func (c Cactus) Model() world.BlockModel {
	return model.Cactus{}
}

func allCactus() (b []world.Block) {
	for i := 0; i < 16; i++ {
		b = append(b, Cactus{Age: i})
	}
	return
}

type DamageSource struct {
	Block world.Block
}

func (DamageSource) ReducedByResistance() bool { return true }
func (DamageSource) ReducedByArmour() bool     { return true }
func (DamageSource) Fire() bool                { return false }
func (DamageSource) IgnoreTotem() bool         { return false }
