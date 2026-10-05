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

type Furnace struct {
	solid
	bassDrum
	*smelter

	Facing cube.Direction

	Lit bool
}

func NewFurnace(face cube.Direction) Furnace {
	return Furnace{
		Facing:  face,
		smelter: newSmelter(),
	}
}

func (Furnace) ContainerSize() int { return 3 }

func (f Furnace) Tick(_ int64, pos cube.Pos, tx *world.Tx) {
	if f.Lit && rand.Float64() <= 0.016 {
		tx.PlaySound(pos.Vec3Centre(), sound.FurnaceCrackle{})
	}
	if lit := f.tickSmelting(time.Second*10, time.Millisecond*100, f.Lit, func(item.SmeltInfo) bool {
		return true
	}); f.Lit != lit {
		f.Lit = lit
		tx.SetBlock(pos, f, nil)
	}
}

func (f Furnace) LightEmissionLevel() uint8 {
	if f.Lit {
		return 13
	}
	return 0
}

func (f Furnace) EncodeItem() (name string, meta int16) {
	return "minecraft:furnace", 0
}

func (f Furnace) EncodeBlock() (name string, properties map[string]interface{}) {
	if f.Lit {
		return "minecraft:lit_furnace", map[string]interface{}{"minecraft:cardinal_direction": f.Facing.String()}
	}
	return "minecraft:furnace", map[string]interface{}{"minecraft:cardinal_direction": f.Facing.String()}
}

func (f Furnace) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, f)
	if !used {
		return false
	}

	place(tx, pos, NewFurnace(user.Rotation().Direction().Opposite()), user, ctx)
	return placed(ctx)
}

func (f Furnace) BreakInfo() BreakInfo {
	xp := 0
	if f.smelter != nil {
		xp = f.Experience()
	}
	return newBreakInfo(3.5, alwaysHarvestable, pickaxeEffective, oneOf(Furnace{})).withXPDropRange(xp, xp).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		for _, i := range f.Inventory(tx, pos).Clear() {
			dropItem(tx, i, pos.Vec3())
		}
	})
}

func (f Furnace) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (f Furnace) EncodeNBT() map[string]interface{} {
	if f.smelter == nil {

		f = NewFurnace(f.Facing)
	}
	remaining, maximum, cook := f.Durations()
	return map[string]interface{}{
		"BurnTime":     int16(remaining.Milliseconds() / 50),
		"CookTime":     int16(cook.Milliseconds() / 50),
		"BurnDuration": int16(maximum.Milliseconds() / 50),
		"StoredXPInt":  int16(f.Experience()),
		"Items":        nbtconv.InvToNBT(f.inventory),
		"id":           "Furnace",
	}
}

func (f Furnace) DecodeNBT(data map[string]interface{}) interface{} {
	remaining := nbtconv.TickDuration[int16](data, "BurnTime")
	maximum := nbtconv.TickDuration[int16](data, "BurnDuration")
	cook := nbtconv.TickDuration[int16](data, "CookTime")

	xp := int(nbtconv.Int16(data, "StoredXPInt"))
	lit := f.Lit

	f = NewFurnace(f.Facing)
	f.Lit = lit
	f.setExperience(xp)
	f.setDurations(remaining, maximum, cook)
	nbtconv.InvFromNBT(f.inventory, nbtconv.Slice(data, "Items"))
	return f
}

func allFurnaces() (furnaces []world.Block) {
	for _, face := range cube.Directions() {
		furnaces = append(furnaces, Furnace{Facing: face})
		furnaces = append(furnaces, Furnace{Facing: face, Lit: true})
	}
	return
}
