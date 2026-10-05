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

type BlastFurnace struct {
	solid
	bassDrum
	*smelter

	Facing cube.Direction

	Lit bool
}

func (BlastFurnace) ContainerSize() int { return 3 }

func NewBlastFurnace(face cube.Direction) BlastFurnace {
	return BlastFurnace{
		Facing:  face,
		smelter: newSmelter(),
	}
}

func (b BlastFurnace) Tick(_ int64, pos cube.Pos, tx *world.Tx) {
	if b.Lit && rand.Float64() <= 0.016 {
		tx.PlaySound(pos.Vec3Centre(), sound.BlastFurnaceCrackle{})
	}
	if lit := b.tickSmelting(time.Second*5, time.Millisecond*200, b.Lit, func(i item.SmeltInfo) bool {
		return i.Ores
	}); b.Lit != lit {
		b.Lit = lit
		tx.SetBlock(pos, b, nil)
	}
}

func (b BlastFurnace) LightEmissionLevel() uint8 {
	if b.Lit {
		return 13
	}
	return 0
}

func (b BlastFurnace) EncodeItem() (name string, meta int16) {
	return "minecraft:blast_furnace", 0
}

func (b BlastFurnace) EncodeBlock() (name string, properties map[string]interface{}) {
	if b.Lit {
		return "minecraft:lit_blast_furnace", map[string]interface{}{"minecraft:cardinal_direction": b.Facing.String()}
	}
	return "minecraft:blast_furnace", map[string]interface{}{"minecraft:cardinal_direction": b.Facing.String()}
}

func (b BlastFurnace) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, b)
	if !used {
		return false
	}

	place(tx, pos, NewBlastFurnace(user.Rotation().Direction().Opposite()), user, ctx)
	return placed(ctx)
}

func (b BlastFurnace) BreakInfo() BreakInfo {
	xp := 0
	if b.smelter != nil {
		xp = b.Experience()
	}
	return newBreakInfo(3.5, alwaysHarvestable, pickaxeEffective, oneOf(BlastFurnace{})).withXPDropRange(xp, xp).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		for _, i := range b.Inventory(tx, pos).Clear() {
			dropItem(tx, i, pos.Vec3())
		}
	})
}

func (b BlastFurnace) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (b BlastFurnace) EncodeNBT() map[string]interface{} {
	if b.smelter == nil {

		b = NewBlastFurnace(b.Facing)
	}
	remaining, maximum, cook := b.Durations()
	return map[string]interface{}{
		"BurnTime":     int16(remaining.Milliseconds() / 50),
		"CookTime":     int16(cook.Milliseconds() / 50),
		"BurnDuration": int16(maximum.Milliseconds() / 50),
		"StoredXPInt":  int16(b.Experience()),
		"Items":        nbtconv.InvToNBT(b.inventory),
		"id":           "BlastFurnace",
	}
}

func (b BlastFurnace) DecodeNBT(data map[string]interface{}) interface{} {
	remaining := nbtconv.TickDuration[int16](data, "BurnTime")
	maximum := nbtconv.TickDuration[int16](data, "BurnDuration")
	cook := nbtconv.TickDuration[int16](data, "CookTime")

	xp := int(nbtconv.Int16(data, "StoredXPInt"))
	lit := b.Lit

	b = NewBlastFurnace(b.Facing)
	b.Lit = lit
	b.setExperience(xp)
	b.setDurations(remaining, maximum, cook)
	nbtconv.InvFromNBT(b.inventory, nbtconv.Slice(data, "Items"))
	return b
}

func allBlastFurnaces() (furnaces []world.Block) {
	for _, face := range cube.Directions() {
		furnaces = append(furnaces, BlastFurnace{Facing: face})
		furnaces = append(furnaces, BlastFurnace{Facing: face, Lit: true})
	}
	return
}
