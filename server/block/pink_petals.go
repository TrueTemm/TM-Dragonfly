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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type PinkPetals struct {
	empty
	transparent

	AdditionalCount int

	Facing cube.Direction
}

func (p PinkPetals) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	if p.AdditionalCount < 3 {
		p.AdditionalCount++
		tx.SetBlock(pos, p, nil)
		return item.BoneMealResultSmall
	}
	dropItem(tx, item.NewStack(p, 1), pos.Vec3Centre())
	return item.BoneMealResultSmall
}

func (p PinkPetals) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if existing, ok := tx.Block(pos).(PinkPetals); ok {
		if existing.AdditionalCount >= 3 {
			return false
		}

		existing.AdditionalCount++
		place(tx, pos, existing, user, ctx)
		return placed(ctx)
	}

	pos, _, used := firstReplaceable(tx, pos, face, p)
	if !used {
		return false
	}
	if !supportsVegetation(p, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	p.Facing = user.Rotation().Direction().Opposite()
	place(tx, pos, p, user, ctx)
	return placed(ctx)
}

func (p PinkPetals) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !supportsVegetation(p, tx.Block(pos.Side(cube.FaceDown))) {
		breakBlock(p, pos, tx)
	}
}

func (PinkPetals) HasLiquidDrops() bool {
	return true
}

func (p PinkPetals) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, simpleDrops(item.NewStack(PinkPetals{}, p.AdditionalCount+1)))
}

func (p PinkPetals) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 100, true)
}

func (PinkPetals) CompostChance() float64 {
	return 0.3
}

func (PinkPetals) EncodeItem() (name string, meta int16) {
	return "minecraft:pink_petals", 0
}

func (p PinkPetals) EncodeBlock() (string, map[string]any) {
	return "minecraft:pink_petals", map[string]any{"growth": int32(p.AdditionalCount), "minecraft:cardinal_direction": p.Facing.String()}
}

func allPinkPetals() (b []world.Block) {
	for i := 0; i <= 7; i++ {
		for _, d := range cube.Directions() {
			b = append(b, PinkPetals{AdditionalCount: i, Facing: d})
		}
	}
	return
}
