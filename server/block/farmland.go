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
	"github.com/df-mc/dragonfly/server/world"
)

type Farmland struct {
	tilledGrass

	Hydration int
}

func (f Farmland) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, NetherSprouts, PinkPetals, DeadBush:
		return true
	}
	return false
}

func (f Farmland) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if solid := tx.Block(pos.Side(cube.FaceUp)).Model().FaceSolid(pos.Side(cube.FaceUp), cube.FaceDown, tx); solid {
		tx.SetBlock(pos, Dirt{}, nil)
	}
}

func (f Farmland) RandomTick(pos cube.Pos, tx *world.Tx, _ *rand.Rand) {
	if !f.hydrated(pos, tx) {
		if f.Hydration > 0 {
			f.Hydration--
			tx.SetBlock(pos, f, nil)
		} else {
			blockAbove := tx.Block(pos.Side(cube.FaceUp))
			if _, cropAbove := blockAbove.(Crop); !cropAbove {
				tx.SetBlock(pos, Dirt{}, nil)
			}
		}
	} else {
		f.Hydration = 7
		tx.SetBlock(pos, f, nil)
	}
}

func (f Farmland) hydrated(pos cube.Pos, tx *world.Tx) bool {
	posX, posY, posZ := pos.X(), pos.Y(), pos.Z()
	for y := 0; y <= 1; y++ {
		for x := -4; x <= 4; x++ {
			for z := -4; z <= 4; z++ {
				if liquid, ok := tx.Liquid(cube.Pos{posX + x, posY + y, posZ + z}); ok {
					if _, ok := liquid.(Water); ok {
						return true
					}
				}
			}
		}
	}
	return false
}

func (f Farmland) EntityLand(pos cube.Pos, tx *world.Tx, e world.Entity, _ *float64) {
	if living, ok := e.(livingEntity); ok {
		if fall, ok := living.(fallDistanceEntity); ok && rand.Float64() < fall.FallDistance()-0.5 {
			ctx := tx.Event()
			if tx.World().Handler().HandleCropTrample(ctx, pos); !ctx.Cancelled() {
				tx.SetBlock(pos, Dirt{}, nil)
			}
		}
	}
}

type fallDistanceEntity interface {
	ResetFallDistance()

	FallDistance() float64
}

func (f Farmland) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, shovelEffective, oneOf(Dirt{}))
}

func (f Farmland) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:farmland", map[string]any{"moisturized_amount": int32(f.Hydration)}
}

func (f Farmland) EncodeItem() (name string, meta int16) {
	return "minecraft:farmland", 0
}

func allFarmland() (b []world.Block) {
	for i := 0; i <= 7; i++ {
		b = append(b, Farmland{Hydration: i})
	}
	return
}
