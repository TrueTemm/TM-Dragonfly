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
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"strconv"
	"time"
)

type Campfire struct {
	transparent
	bass
	sourceWaterDisplacer

	Items [4]CampfireItem

	Facing cube.Direction

	Extinguished bool

	Type FireType
}

type CampfireItem struct {
	Item item.Stack

	Time time.Duration
}

func (Campfire) Model() world.BlockModel {
	return model.Campfire{}
}

func (Campfire) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c Campfire) BreakInfo() BreakInfo {
	return newBreakInfo(2, alwaysHarvestable, axeEffective, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(Campfire{Type: c.Type}, 1)}
		}
		switch c.Type {
		case NormalFire():
			return []item.Stack{item.NewStack(item.Charcoal{}, 2)}
		case SoulFire():
			return []item.Stack{item.NewStack(SoulSoil{}, 1)}
		}
		panic("should never happen")
	}).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		for _, v := range c.Items {
			if !v.Item.Empty() {
				dropItem(tx, v.Item, pos.Vec3Centre())
			}
		}
	})
}

func (c Campfire) LightEmissionLevel() uint8 {
	if c.Extinguished {
		return 0
	}
	return c.Type.LightLevel()
}

func (c Campfire) Ignite(pos cube.Pos, tx *world.Tx, _ world.Entity) bool {
	tx.PlaySound(pos.Vec3(), sound.Ignite{})
	if !c.Extinguished {
		return false
	}
	if _, ok := tx.Liquid(pos); ok {
		return false
	}

	c.Extinguished = false
	tx.SetBlock(pos, c, nil)
	return true
}

func (c Campfire) Splash(tx *world.Tx, pos cube.Pos) {
	if c.Extinguished {
		return
	}

	c.extinguish(pos, tx)
}

func (c Campfire) extinguish(pos cube.Pos, tx *world.Tx) {
	tx.PlaySound(pos.Vec3Centre(), sound.FireExtinguish{})
	c.Extinguished = true

	for i := range c.Items {
		c.Items[i].Time = time.Second * 30
	}

	tx.SetBlock(pos, c, nil)
}

func (c Campfire) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, ctx *item.UseContext) bool {
	held, _ := u.HeldItems()
	if held.Empty() {
		return false
	}

	if _, ok := held.Item().(item.Shovel); ok && !c.Extinguished {
		c.extinguish(pos, tx)
		ctx.DamageItem(1)
		return true
	}

	rawFood, ok := held.Item().(item.Smeltable)
	if !ok || !rawFood.SmeltInfo().Food {
		return false
	}

	if _, ok = tx.Liquid(pos); ok {
		return false
	}

	for i, it := range c.Items {
		if it.Item.Empty() {
			c.Items[i] = CampfireItem{
				Item: held.Grow(-held.Count() + 1),
				Time: time.Second * 30,
			}

			ctx.SubtractFromCount(1)

			tx.PlaySound(pos.Vec3Centre(), sound.ItemAdd{})
			tx.SetBlock(pos, c, nil)
			return true
		}
	}
	return false
}

func (c Campfire) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, c)
	if !used {
		return
	}
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Campfire); ok {
		return false
	}
	c.Facing = user.Rotation().Direction().Opposite()
	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (c Campfire) Tick(_ int64, pos cube.Pos, tx *world.Tx) {
	if c.Extinguished {

		return
	}
	if rand.Float64() <= 0.016 {
		tx.PlaySound(pos.Vec3Centre(), sound.CampfireCrackle{})
	}

	updated := false
	for i, it := range c.Items {
		if it.Item.Empty() {
			continue
		}

		updated = true
		if it.Time > 0 {
			c.Items[i].Time = it.Time - time.Millisecond*50
			continue
		}

		if food, ok := it.Item.Item().(item.Smeltable); ok {
			dropItem(tx, food.SmeltInfo().Product, pos.Vec3Middle())
		}
		c.Items[i].Item = item.Stack{}
	}
	if updated {
		tx.SetBlock(pos, c, nil)
	}
}

func (c Campfire) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Liquid(pos); ok {
		var updated bool
		for i, it := range c.Items {
			if !it.Item.Empty() {
				dropItem(tx, it.Item, pos.Vec3Centre())
				c.Items[i].Item, updated = item.Stack{}, true
			}
		}
		if !c.Extinguished {
			c.extinguish(pos, tx)
		} else if updated {
			tx.SetBlock(pos, c, nil)
		}
		return
	}
	if liquid, ok := tx.Liquid(pos.Side(cube.FaceUp)); ok && liquid.LiquidType() == "water" && !c.Extinguished {
		c.extinguish(pos, tx)
	}
}

func (c Campfire) EntityInside(pos cube.Pos, tx *world.Tx, e world.Entity) {
	if flammable, ok := e.(flammableEntity); ok {
		if flammable.OnFireDuration() > 0 && c.Extinguished {
			c.Extinguished = false
			tx.PlaySound(pos.Vec3(), sound.Ignite{})
			tx.SetBlock(pos, c, nil)
		}
		if !c.Extinguished {
			if l, ok := e.(livingEntity); ok {
				l.Hurt(c.Type.Damage(), FireDamageSource{})
			}
		}
	}
}

func (c Campfire) EncodeNBT() map[string]any {
	m := map[string]any{"id": "Campfire"}
	for i, v := range c.Items {
		id := strconv.Itoa(i + 1)
		if !v.Item.Empty() {
			m["Item"+id] = item.WriteNBT(v.Item, true)
			m["ItemTime"+id] = int32(v.Time.Milliseconds() / 50)
		}
	}
	return m
}

func (c Campfire) DecodeNBT(data map[string]any) any {
	for i := 0; i < 4; i++ {
		id := strconv.Itoa(i + 1)
		c.Items[i] = CampfireItem{
			Item: item.MapNBT(data, "Item"+id),
			Time: time.Duration(nbtconv.Int32(data, "ItemTime"+id)) * time.Millisecond * 50,
		}
	}
	return c
}

func (c Campfire) EncodeItem() (name string, meta int16) {
	switch c.Type {
	case NormalFire():
		return "minecraft:campfire", 0
	case SoulFire():
		return "minecraft:soul_campfire", 0
	}
	panic("invalid fire type")
}

func (c Campfire) EncodeBlock() (name string, properties map[string]any) {
	switch c.Type {
	case NormalFire():
		name = "minecraft:campfire"
	case SoulFire():
		name = "minecraft:soul_campfire"
	}
	return name, map[string]any{
		"minecraft:cardinal_direction": c.Facing.String(),
		"extinguished":                 c.Extinguished,
	}
}

func allCampfires() (campfires []world.Block) {
	for _, d := range cube.Directions() {
		for _, f := range FireTypes() {
			campfires = append(campfires, Campfire{Facing: d, Type: f, Extinguished: true})
			campfires = append(campfires, Campfire{Facing: d, Type: f})
		}
	}
	return campfires
}
