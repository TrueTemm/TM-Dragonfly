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
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"time"
)

type Smoker struct {
	solid
	bassDrum
	*smelter

	Facing cube.Direction

	Lit bool
}

func (Smoker) ContainerSize() int { return 3 }

func NewSmoker(face cube.Direction) Smoker {
	return Smoker{
		Facing:  face,
		smelter: newSmelter(),
	}
}

func (s Smoker) Tick(_ int64, pos cube.Pos, tx *world.Tx) {
	if s.Lit && rand.Float64() <= 0.016 {
		tx.PlaySound(pos.Vec3Centre(), sound.SmokerCrackle{})
	}
	if lit := s.tickSmelting(time.Second*5, time.Millisecond*200, s.Lit, func(i item.SmeltInfo) bool {
		return i.Food
	}); s.Lit != lit {
		s.Lit = lit
		tx.SetBlock(pos, s, nil)
	}
}

func (s Smoker) LightEmissionLevel() uint8 {
	if s.Lit {
		return 13
	}
	return 0
}

func (s Smoker) EncodeItem() (name string, meta int16) {
	return "minecraft:smoker", 0
}

func (s Smoker) EncodeBlock() (name string, properties map[string]interface{}) {
	if s.Lit {
		return "minecraft:lit_smoker", map[string]interface{}{"minecraft:cardinal_direction": s.Facing.String()}
	}
	return "minecraft:smoker", map[string]interface{}{"minecraft:cardinal_direction": s.Facing.String()}
}

func (s Smoker) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, s)
	if !used {
		return false
	}

	place(tx, pos, NewSmoker(user.Rotation().Direction().Opposite()), user, ctx)
	return placed(ctx)
}

func (s Smoker) BreakInfo() BreakInfo {
	xp := 0
	if s.smelter != nil {
		xp = s.Experience()
	}
	return newBreakInfo(3.5, alwaysHarvestable, pickaxeEffective, oneOf(Smoker{})).withXPDropRange(xp, xp).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		for _, i := range s.Inventory(tx, pos).Clear() {
			dropItem(tx, i, pos.Vec3())
		}
	})
}

func (s Smoker) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (s Smoker) EncodeNBT() map[string]interface{} {
	if s.smelter == nil {

		s = NewSmoker(s.Facing)
	}
	remaining, maximum, cook := s.Durations()
	return map[string]interface{}{
		"BurnTime":     int16(remaining.Milliseconds() / 50),
		"CookTime":     int16(cook.Milliseconds() / 50),
		"BurnDuration": int16(maximum.Milliseconds() / 50),
		"StoredXPInt":  int16(s.Experience()),
		"Items":        nbtconv.InvToNBT(s.inventory),
		"id":           "Smoker",
	}
}

func (s Smoker) DecodeNBT(data map[string]interface{}) interface{} {
	remaining := nbtconv.TickDuration[int16](data, "BurnTime")
	maximum := nbtconv.TickDuration[int16](data, "BurnDuration")
	cook := nbtconv.TickDuration[int16](data, "CookTime")

	xp := int(nbtconv.Int16(data, "StoredXPInt"))
	lit := s.Lit

	s = NewSmoker(s.Facing)
	s.Lit = lit
	s.setExperience(xp)
	s.setDurations(remaining, maximum, cook)
	nbtconv.InvFromNBT(s.inventory, nbtconv.Slice(data, "Items"))
	return s
}

func allSmokers() (smokers []world.Block) {
	for _, face := range cube.Directions() {
		smokers = append(smokers, Smoker{Facing: face})
		smokers = append(smokers, Smoker{Facing: face, Lit: true})
	}
	return
}
