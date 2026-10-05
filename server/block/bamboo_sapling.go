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
)

type BambooSapling struct {
	empty
	transparent
	bass

	Ready bool
}

var (
	_ item.BoneMealAffected = BambooSapling{}
	_ Flammable             = BambooSapling{}
)

func (b BambooSapling) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if b.grow(pos, tx) {
		return item.BoneMealResultSmall
	}
	return item.BoneMealResultNone
}

func (b BambooSapling) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 60, true)
}

func (b BambooSapling) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	down := tx.Block(pos.Side(cube.FaceDown))
	if supportsVegetation(b, down) {
		return
	}
	breakBlock(b, pos, tx)
}

func (b BambooSapling) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if tx.Light(pos) >= 9 && r.IntN(3) == 0 {
		b.grow(pos, tx)
	}
}

func (b BambooSapling) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, axeEffective, oneOf(Bamboo{})).withBlastResistance(1)
}

func (b BambooSapling) HasLiquidDrops() bool {
	return true
}

func (b BambooSapling) EncodeBlock() (string, map[string]any) {
	return "minecraft:bamboo_sapling", map[string]any{"age_bit": boolByte(b.Ready)}
}

func (b BambooSapling) grow(pos cube.Pos, tx *world.Tx) bool {
	if !replaceableWith(tx, pos.Side(cube.FaceUp), b) {
		return false
	}

	tx.SetBlock(pos, Bamboo{}, nil)
	tx.SetBlock(pos.Side(cube.FaceUp), Bamboo{LeafSize: BambooSizeSmallLeaves()}, nil)
	return true
}

func allBambooSaplings() (saplings []world.Block) {
	saplings = append(saplings, BambooSapling{Ready: false})
	saplings = append(saplings, BambooSapling{Ready: true})
	return
}
