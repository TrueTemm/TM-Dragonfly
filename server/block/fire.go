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
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/portal"
)

type Fire struct {
	replaceable
	transparent
	empty

	Type FireType

	Age int
}

func flammableBlock(block world.Block) bool {
	flammable, ok := block.(Flammable)
	return ok && flammable.FlammabilityInfo().Encouragement > 0
}

func neighboursFlammable(pos cube.Pos, tx *world.Tx) bool {
	for _, i := range cube.Faces() {
		if flammableBlock(tx.Block(pos.Side(i))) {
			return true
		}
	}
	return false
}

func infinitelyBurning(pos cube.Pos, tx *world.Tx) bool {
	switch block := tx.Block(pos.Side(cube.FaceDown)).(type) {

	case Netherrack:
		return true
	case Bedrock:
		return block.InfiniteBurning
	}
	return false
}

func (f Fire) burn(from, to cube.Pos, tx *world.Tx, r *rand.Rand, chanceBound int) {
	if flammable, ok := tx.Block(to).(Flammable); ok && r.IntN(chanceBound) < flammable.FlammabilityInfo().Flammability {
		if t, ok := flammable.(TNT); ok {
			t.Ignite(to, tx, nil)
			return
		}
		if r.IntN(f.Age+10) < 5 && !rainingAround(to, tx) {
			f.spread(from, to, tx, r)
			return
		}
		tx.SetBlock(to, nil, nil)
	}
}

func rainingAround(pos cube.Pos, tx *world.Tx) bool {
	raining := tx.RainingAt(pos)
	for _, face := range cube.HorizontalFaces() {
		if raining {
			break
		}
		raining = tx.RainingAt(pos.Side(face))
	}
	return raining
}

func (f Fire) tick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if f.Type == SoulFire() {
		return
	}
	infinitelyBurns := infinitelyBurning(pos, tx)
	if !infinitelyBurns && (20+f.Age*3) > r.IntN(100) && rainingAround(pos, tx) {

		tx.SetBlock(pos, nil, nil)
		return
	}

	if f.Age < 15 && r.IntN(3) == 0 {
		f.Age++
		tx.SetBlock(pos, f, nil)
	}

	tx.ScheduleBlockUpdate(pos, f, time.Duration(30+r.IntN(10))*time.Second/20)

	if !infinitelyBurns {
		_, waterBelow := tx.Block(pos.Side(cube.FaceDown)).(Water)
		if waterBelow {
			tx.SetBlock(pos, nil, nil)
			return
		}
		if !neighboursFlammable(pos, tx) {
			if !tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos, cube.FaceUp, tx) || f.Age > 3 {
				tx.SetBlock(pos, nil, nil)
			}
			return
		}
		if !flammableBlock(tx.Block(pos.Side(cube.FaceDown))) && f.Age == 15 && r.IntN(4) == 0 {
			tx.SetBlock(pos, nil, nil)
			return
		}
	}

	humid := tx.Biome(pos).Rainfall() > 0.85

	s := 0
	if humid {
		s = 50
	}
	for _, face := range cube.Faces() {
		if face == cube.FaceUp || face == cube.FaceDown {
			f.burn(pos, pos.Side(face), tx, r, 300-s)
		} else {
			f.burn(pos, pos.Side(face), tx, r, 250-s)
		}
	}

	for y := -1; y <= 4; y++ {
		randomBound := 100
		if y > 1 {
			randomBound += (y - 1) * 100
		}

		for x := -1; x <= 1; x++ {
			for z := -1; z <= 1; z++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				blockPos := pos.Add(cube.Pos{x, y, z})
				block := tx.Block(blockPos)
				if _, ok := block.(Air); !ok {
					continue
				}

				encouragement := 0
				blockPos.Neighbours(func(neighbour cube.Pos) {
					if flammable, ok := tx.Block(neighbour).(Flammable); ok {
						encouragement = max(encouragement, flammable.FlammabilityInfo().Encouragement)
					}
				}, tx.Range())
				if encouragement <= 0 {
					continue
				}

				maxChance := (encouragement + 40 + tx.World().Difficulty().FireSpreadIncrease()) / (f.Age + 30)
				if humid {
					maxChance /= 2
				}

				if maxChance > 0 && r.IntN(randomBound) <= maxChance && !rainingAround(blockPos, tx) {
					f.spread(pos, blockPos, tx, r)
				}
			}
		}
	}
}

func (f Fire) spread(from, to cube.Pos, tx *world.Tx, r *rand.Rand) {
	if _, air := tx.Block(to).(Air); !air {
		ctx := tx.Event()
		if tx.World().Handler().HandleBlockBurn(ctx, to); ctx.Cancelled() {
			return
		}
	}
	ctx := tx.Event()
	if tx.World().Handler().HandleFireSpread(ctx, from, to); ctx.Cancelled() {
		return
	}
	spread := Fire{Type: f.Type, Age: min(15, f.Age+r.IntN(5)/4)}
	if spread.Type == NormalFire() && portal.ActivateNetherPortal(tx, to) {
		return
	}
	tx.SetBlock(to, spread, nil)
	tx.ScheduleBlockUpdate(to, spread, time.Duration(30+r.IntN(10))*time.Second/20)
}

func (f Fire) EntityInside(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if flammable, ok := e.(flammableEntity); ok {
		if l, ok := e.(livingEntity); ok {
			l.Hurt(f.Type.Damage(), FireDamageSource{})
		}
		if flammable.OnFireDuration() < time.Second*8 {
			flammable.SetOnFire(8 * time.Second)
		}
	}
}

func (f Fire) ScheduledTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	f.tick(pos, tx, r)
}

func (f Fire) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	f.tick(pos, tx, r)
}

func (f Fire) NeighbourUpdateTick(pos, changedNeighbour cube.Pos, tx *world.Tx) {
	below := tx.Block(pos.Side(cube.FaceDown))
	if diffuser, ok := below.(LightDiffuser); (ok && diffuser.LightDiffusionLevel() != 15) && (!neighboursFlammable(pos, tx) || f.Type == SoulFire()) {
		tx.SetBlock(pos, nil, nil)
		return
	}
	switch below.(type) {
	case SoulSand, SoulSoil:
		f.Type = SoulFire()
		tx.SetBlock(pos, f, nil)
	case Water:
		if changedNeighbour == pos {
			tx.SetBlock(pos, nil, nil)
		}
	default:
		if f.Type == SoulFire() {
			tx.SetBlock(pos, nil, nil)
			return
		}
	}
}

func (f Fire) HasLiquidDrops() bool {
	return false
}

func (f Fire) PortalInterior(target world.Dimension) bool {
	return target == world.Nether && f.Type == NormalFire()
}

func (f Fire) LightEmissionLevel() uint8 {
	return f.Type.LightLevel()
}

func (f Fire) EncodeBlock() (name string, properties map[string]any) {
	switch f.Type {
	case NormalFire():
		return "minecraft:fire", map[string]any{"age": int32(f.Age)}
	case SoulFire():
		return "minecraft:soul_fire", map[string]any{"age": int32(f.Age)}
	}
	panic("unknown fire type")
}

func (f Fire) Start(tx *world.Tx, pos cube.Pos) {
	b := tx.Block(pos)
	_, air := b.(Air)
	_, shortGrass := b.(ShortGrass)
	_, fern := b.(Fern)
	if air || shortGrass || fern {
		below := tx.Block(pos.Side(cube.FaceDown))
		if below.Model().FaceSolid(pos, cube.FaceUp, tx) || neighboursFlammable(pos, tx) {
			if portal.ActivateNetherPortal(tx, pos) {
				return
			}
			f := Fire{}
			tx.SetBlock(pos, f, nil)
			tx.ScheduleBlockUpdate(pos, f, time.Duration(30+rand.IntN(10))*time.Second/20)
		}
	}
}

func allFire() (b []world.Block) {
	for i := 0; i < 16; i++ {
		b = append(b, Fire{Age: i, Type: NormalFire()})
		b = append(b, Fire{Age: i, Type: SoulFire()})
	}
	return
}

type FireDamageSource struct{}

func (FireDamageSource) ReducedByResistance() bool { return true }
func (FireDamageSource) ReducedByArmour() bool     { return true }
func (FireDamageSource) Fire() bool                { return true }
func (FireDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool {
	return e == enchantment.FireProtection
}
func (FireDamageSource) IgnoreTotem() bool { return false }
