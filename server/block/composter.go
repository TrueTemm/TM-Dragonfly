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
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"math/rand/v2"
	"time"
)

type Composter struct {
	bass
	transparent
	sourceWaterDisplacer

	Level int
}

func (c Composter) InsertItem(h Hopper, pos cube.Pos, tx *world.Tx) bool {
	if c.Level >= 7 || h.Facing != cube.FaceDown {
		return false
	}

	for sourceSlot, sourceStack := range h.inventory.Slots() {
		if sourceStack.Empty() {
			continue
		}

		if c.fill(sourceStack, pos, tx) {
			_ = h.inventory.SetItem(sourceSlot, sourceStack.Grow(-1))
			return true
		}
	}

	return false
}

func (c Composter) ExtractItem(h Hopper, pos cube.Pos, tx *world.Tx) bool {
	if c.Level == 8 {
		_, err := h.inventory.AddItem(item.NewStack(item.BoneMeal{}, 1))
		if err != nil {

			return false
		}

		c.Level = 0
		tx.SetBlock(pos, c, nil)
		tx.PlaySound(pos.Vec3(), sound.ComposterEmpty{})
		return true
	}

	return false
}

func (c Composter) Model() world.BlockModel {
	return model.Composter{Level: c.Level}
}

func (c Composter) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (c Composter) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(5, 20, true)
}

func (c Composter) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c Composter) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, axeEffective, oneOf(Composter{})).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		if c.Level == 8 {
			dropItem(tx, item.NewStack(item.BoneMeal{}, 1), pos.Side(cube.FaceUp).Vec3Middle())
		}
	})
}

func (c Composter) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, ctx *item.UseContext) bool {
	if c.Level >= 7 {
		if c.Level == 8 {
			c.Level = 0
			tx.SetBlock(pos, c, nil)
			dropItem(tx, item.NewStack(item.BoneMeal{}, 1), pos.Side(cube.FaceUp).Vec3Middle())
			tx.PlaySound(pos.Vec3(), sound.ComposterEmpty{})

			return true
		}
		return false
	}
	it, _ := u.HeldItems()
	if c.fill(it, pos, tx) {
		ctx.SubtractFromCount(1)
		return true
	}
	return false
}

func (c Composter) fill(it item.Stack, pos cube.Pos, tx *world.Tx) bool {
	compostable, ok := it.Item().(item.Compostable)
	if !ok {
		return false
	}
	tx.AddParticle(pos.Vec3(), particle.BoneMeal{})
	if rand.Float64() > compostable.CompostChance() {
		tx.PlaySound(pos.Vec3(), sound.ComposterFill{})
		return true
	}
	c.Level++
	tx.SetBlock(pos, c, nil)
	tx.PlaySound(pos.Vec3(), sound.ComposterFillLayer{})
	if c.Level == 7 {
		tx.ScheduleBlockUpdate(pos, c, time.Second)
	}

	return true
}

func (c Composter) ScheduledTick(pos cube.Pos, tx *world.Tx, _ *rand.Rand) {
	if c.Level == 7 {
		c.Level = 8
		tx.SetBlock(pos, c, nil)
		tx.PlaySound(pos.Vec3(), sound.ComposterReady{})
	}
}

func (c Composter) EncodeItem() (name string, meta int16) {
	return "minecraft:composter", 0
}

func (c Composter) EncodeBlock() (string, map[string]any) {
	return "minecraft:composter", map[string]any{"composter_fill_level": int32(c.Level)}
}

func allComposters() (all []world.Block) {
	for i := 0; i < 9; i++ {
		all = append(all, Composter{Level: i})
	}
	return
}
