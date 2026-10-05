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

package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func NewArrow(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	return NewTippedArrowWithDamage(opts, 2.0, owner, potion.Potion{})
}

func NewArrowWithDamage(opts world.EntitySpawnOpts, damage float64, owner world.Entity) *world.EntityHandle {
	return NewTippedArrowWithDamage(opts, damage, owner, potion.Potion{})
}

func NewTippedArrow(opts world.EntitySpawnOpts, owner world.Entity, tip potion.Potion) *world.EntityHandle {
	return NewTippedArrowWithDamage(opts, 2.0, owner, tip)
}

func NewTippedArrowWithDamage(opts world.EntitySpawnOpts, damage float64, owner world.Entity, tip potion.Potion) *world.EntityHandle {
	conf := arrowConf
	conf.Damage = damage
	conf.Potion = tip
	conf.Owner = owner.H()
	return opts.New(ArrowType, conf)
}

var arrowConf = ProjectileBehaviourConfig{
	Gravity:               0.05,
	Drag:                  0.01,
	Damage:                2.0,
	Sound:                 sound.ArrowHit{},
	SurviveBlockCollision: true,
}

func boolByte(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

var ArrowType arrowType

type arrowType struct{}

func (t arrowType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (arrowType) EncodeEntity() string { return "minecraft:arrow" }
func (arrowType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (arrowType) DecodeNBT(m map[string]any, data *world.EntityData) {
	conf := arrowConf
	conf.Damage = float64(nbtconv.Float32(m, "Damage"))
	conf.Potion = potion.From(nbtconv.Int32(m, "auxValue") - 1)
	conf.DisablePickup = !nbtconv.Bool(m, "player")
	if !nbtconv.Bool(m, "isCreative") {
		conf.PickupItem = item.NewStack(item.Arrow{Tip: conf.Potion}, 1)
	}
	conf.KnockBackForceAddend = enchantment.Punch.KnockBackMultiplier() * float64(nbtconv.Uint8(m, "enchantPunch"))
	conf.CollisionPosition = nbtconv.Pos(m, "StuckToBlockPos")

	data.Data = conf.New()
}

func (arrowType) EncodeNBT(data *world.EntityData) map[string]any {
	b := data.Data.(*ProjectileBehaviour)
	m := map[string]any{
		"Damage":       float32(b.conf.Damage),
		"enchantPunch": byte(b.conf.KnockBackForceAddend / enchantment.Punch.KnockBackMultiplier()),
		"auxValue":     int32(b.conf.Potion.Uint8() + 1),
		"player":       boolByte(!b.conf.DisablePickup),
		"isCreative":   boolByte(b.conf.PickupItem.Empty()),
	}

	if b.collided {
		m["StuckToBlockPos"] = nbtconv.PosToInt32Slice(b.collisionPos)
	}
	return m
}
