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
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

type Water struct {
	empty

	Still bool

	Depth int

	Falling bool
}

func (w Water) ReplaceableBy(b world.Block) bool {
	if _, ok := b.(LiquidRemovable); ok {
		_, displacer := b.(world.LiquidDisplacer)
		_, liquid := b.(world.Liquid)
		return displacer || liquid
	}
	return true
}

func (w Water) EntityInside(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if fallEntity, ok := e.(fallDistanceEntity); ok {
		fallEntity.ResetFallDistance()
	}
	if flammable, ok := e.(flammableEntity); ok {
		flammable.Extinguish()
	}
}

func (w Water) FillBottle() (world.Block, item.Stack, bool) {
	if w.Depth == 8 {
		return w, item.NewStack(item.Potion{Type: potion.Water()}, 1), true
	}
	return nil, item.Stack{}, false
}

func (w Water) LiquidDepth() int {
	return w.Depth
}

func (Water) SpreadDecay() int {
	return 1
}

func (w Water) WithDepth(depth int, falling bool) world.Liquid {
	w.Depth = depth
	w.Falling = falling
	w.Still = false
	return w
}

func (w Water) LiquidFalling() bool {
	return w.Falling
}

func (Water) BlastResistance() float64 {
	return 100
}

func (Water) HasLiquidDrops() bool {
	return false
}

func (Water) LiquidRemoveBlock(pos cube.Pos, tx *world.Tx, removed world.Block) {
	r, ok := removed.(LiquidRemovable)
	if !ok || !r.HasLiquidDrops() {
		return
	}
	b, ok := removed.(Breakable)
	if !ok {
		panic("liquid drops should always implement breakable")
	}
	for _, d := range b.BreakInfo().Drops(item.ToolNone{}, nil) {
		dropItem(tx, d, pos.Vec3Centre())
	}
}

func (Water) LightDiffusionLevel() uint8 {
	return 2
}

func (w Water) ScheduledTick(pos cube.Pos, tx *world.Tx, _ *rand.Rand) {
	if w.Depth == 7 {

		count := 0
		pos.Neighbours(func(neighbour cube.Pos) {
			if neighbour[1] == pos[1] {
				if liquid, ok := tx.Liquid(neighbour); ok {
					if water, ok := liquid.(Water); ok && water.Depth == 8 && !water.Falling {
						count++
					}
				}
			}
		}, tx.Range())
		if count >= 2 {
			if !canFlowInto(w, tx, pos.Side(cube.FaceDown), true) {

				res := Water{Depth: 8, Still: true}
				ctx := tx.Event()
				if tx.World().Handler().HandleLiquidFlow(ctx, pos, pos, res, w); ctx.Cancelled() {
					return
				}
				tx.SetLiquid(pos, res)
			}
		}
	}
	tickLiquid(w, pos, tx)
}

func (w Water) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if tx.World().Dimension().WaterEvaporates() {

		tx.SetLiquid(pos, nil)
		return
	}
	tx.ScheduleBlockUpdate(pos, w, time.Second/4)
}

func (Water) LiquidType() string {
	return "water"
}

func (w Water) Harden(pos cube.Pos, tx *world.Tx, flownIntoBy *cube.Pos) bool {
	if flownIntoBy == nil {
		return false
	}
	if lava, ok := tx.Block(pos.Side(cube.FaceUp)).(Lava); ok {
		ctx := tx.Event()
		if tx.World().Handler().HandleLiquidHarden(ctx, pos, w, lava, Stone{}); ctx.Cancelled() {
			return false
		}
		tx.SetBlock(pos, Stone{}, nil)
		tx.PlaySound(pos.Vec3Centre(), sound.Fizz{})
		return true
	} else if lava, ok := tx.Block(*flownIntoBy).(Lava); ok {
		ctx := tx.Event()
		if tx.World().Handler().HandleLiquidHarden(ctx, pos, w, lava, Cobblestone{}); ctx.Cancelled() {
			return false
		}
		tx.SetBlock(*flownIntoBy, Cobblestone{}, nil)
		tx.PlaySound(pos.Vec3Centre(), sound.Fizz{})
		return true
	}
	return false
}

func (w Water) EncodeBlock() (name string, properties map[string]any) {
	if w.Depth < 1 || w.Depth > 8 {
		panic("invalid water depth, must be between 1 and 8")
	}
	v := 8 - w.Depth
	if w.Falling {
		v += 8
	}
	if w.Still {
		return "minecraft:water", map[string]any{"liquid_depth": int32(v)}
	}
	return "minecraft:flowing_water", map[string]any{"liquid_depth": int32(v)}
}

func allWater() (b []world.Block) {
	f := func(still, falling bool) {
		b = append(b, Water{Still: still, Falling: falling, Depth: 8})
		b = append(b, Water{Still: still, Falling: falling, Depth: 7})
		b = append(b, Water{Still: still, Falling: falling, Depth: 6})
		b = append(b, Water{Still: still, Falling: falling, Depth: 5})
		b = append(b, Water{Still: still, Falling: falling, Depth: 4})
		b = append(b, Water{Still: still, Falling: falling, Depth: 3})
		b = append(b, Water{Still: still, Falling: falling, Depth: 2})
		b = append(b, Water{Still: still, Falling: falling, Depth: 1})
	}
	f(true, true)
	f(true, false)
	f(false, false)
	f(false, true)
	return
}
