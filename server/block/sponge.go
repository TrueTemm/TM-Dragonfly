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
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

type Sponge struct {
	solid

	Wet bool
}

func (s Sponge) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, nothingEffective, oneOf(s))
}

func (s Sponge) SmeltInfo() item.SmeltInfo {
	if s.Wet {
		return newSmeltInfo(item.NewStack(Sponge{}, 1), 0.15)
	}
	return item.SmeltInfo{}
}

func (s Sponge) EncodeItem() (name string, meta int16) {
	if s.Wet {
		return "minecraft:wet_sponge", 0
	}
	return "minecraft:sponge", 0
}

func (s Sponge) EncodeBlock() (string, map[string]any) {
	if s.Wet {
		return "minecraft:wet_sponge", nil
	}
	return "minecraft:sponge", nil
}

func (s Sponge) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	var particles = false
	pos, _, used = firstReplaceable(tx, pos, face, s)
	if !used {
		return
	}

	if tx.World().Dimension().WaterEvaporates() && s.Wet {
		s.Wet = false
		particles = true
	}

	place(tx, pos, s, user, ctx)
	if particles && placed(ctx) {
		tx.AddParticle(pos.Side(cube.FaceUp).Vec3(), particle.Evaporate{})
	}
	return placed(ctx)
}

func (s Sponge) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {

	if !s.Wet {
		if s.absorbWater(pos, tx) > 0 {

			s.setWet(pos, tx)
		}
	}
}

func (s Sponge) setWet(pos cube.Pos, tx *world.Tx) {
	s.Wet = true
	tx.SetBlock(pos, s, nil)
	tx.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: Water{Depth: 1}})
}

func (s Sponge) absorbWater(pos cube.Pos, tx *world.Tx) int {

	type distanceToSponge struct {
		block    cube.Pos
		distance int32
	}

	queue := make([]distanceToSponge, 0)
	queue = append(queue, distanceToSponge{pos, 0})

	replaced := 0
	for replaced < 65 {
		if len(queue) == 0 {
			break
		}

		next := queue[0]
		queue = queue[1:]

		next.block.Neighbours(func(neighbour cube.Pos) {
			liquid, found := tx.Liquid(neighbour)
			if found {
				if _, isWater := liquid.(Water); isWater {
					tx.SetLiquid(neighbour, nil)
					replaced++
					if next.distance < 7 {
						queue = append(queue, distanceToSponge{neighbour, next.distance + 1})
					}
				}
			}
		}, tx.Range())
	}

	return replaced
}
