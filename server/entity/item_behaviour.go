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
	"math"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type ItemBehaviourConfig struct {
	Item item.Stack

	Gravity float64

	Drag float64

	ExistenceDuration time.Duration

	PickupDelay time.Duration
}

func (conf ItemBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf ItemBehaviourConfig) New() *ItemBehaviour {
	i := conf.Item
	if i.Count() > i.MaxCount() {
		i = i.Grow(i.MaxCount() - i.Count())
	}
	i = item.ReadNBT(item.WriteNBT(i, true), nil)

	if conf.PickupDelay == 0 {
		conf.PickupDelay = time.Second / 2
	}
	if conf.ExistenceDuration == 0 {
		conf.ExistenceDuration = time.Minute * 5
	}

	b := &ItemBehaviour{conf: conf, i: i, pickupDelay: conf.PickupDelay}
	b.passive = PassiveBehaviourConfig{
		Gravity:           conf.Gravity,
		Drag:              conf.Drag,
		ExistenceDuration: conf.ExistenceDuration,
		Tick:              b.tick,
	}.New()
	return b
}

type ItemBehaviour struct {
	conf    ItemBehaviourConfig
	passive *PassiveBehaviour
	i       item.Stack

	pickupDelay time.Duration
}

func (i *ItemBehaviour) PortalTravelComputer() *PortalTravelComputer {
	return i.passive.PortalTravelComputer()
}

func (i *ItemBehaviour) Item() item.Stack {
	return i.i
}

func (i *ItemBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	pos := cube.PosFromVec3(e.Position())
	blockPos := pos.Side(cube.FaceDown)

	bl, ok := tx.Block(blockPos).(block.Hopper)
	if ok && !bl.Powered && bl.CollectCooldown <= 0 {
		addedCount, err := bl.Inventory(tx, blockPos).AddItem(i.i)
		if err != nil {
			if addedCount == 0 {
				return i.passive.Tick(e, tx)
			}

			opts := world.EntitySpawnOpts{Position: pos.Vec3Centre()}
			tx.AddEntity(NewItem(opts, i.Item().Grow(-addedCount)))
		}

		_ = e.Close()
		bl.CollectCooldown = 8
		tx.SetBlock(blockPos, bl, nil)
		return nil
	}
	return i.passive.Tick(e, tx)
}

func (i *ItemBehaviour) Explode(e *Ent, _ world.ExplosionSource, impact float64) {
	if impact > 0 {
		if expl, ok := i.Item().Item().(interface{ BlastProof() bool }); ok && expl.BlastProof() {
			return
		}
		_ = e.Close()
	}
}

func (i *ItemBehaviour) tick(e *Ent, tx *world.Tx) {
	if i.pickupDelay == 0 {
		i.checkNearby(e, tx)
	} else if i.pickupDelay < math.MaxInt16*(time.Second/20) {
		i.pickupDelay -= time.Second / 20
	}
}

func (i *ItemBehaviour) checkNearby(e *Ent, tx *world.Tx) {
	pos := e.Position()
	bbox := e.H().Type().BBox(e)
	grown := bbox.GrowVec3(mgl64.Vec3{1, 0.5, 1}).Translate(pos)

	for other := range tx.EntitiesWithin(bbox.Translate(pos).Grow(2)) {
		if e.H() == other.H() || !other.H().Type().BBox(other).Translate(other.Position()).IntersectsWith(grown) {
			continue
		}
		if collector, ok := other.(Collector); ok {

			i.collect(e, collector, tx)
			return
		} else if other.H().Type() == ItemType {

			if i.merge(e, other.(*Ent), tx) {
				return
			}
		}
	}
}

func (i *ItemBehaviour) merge(e *Ent, other *Ent, tx *world.Tx) bool {
	pos := e.Position()
	otherBehaviour := other.Behaviour().(*ItemBehaviour)
	if otherBehaviour.i.Count() == otherBehaviour.i.MaxCount() || i.i.Count() == i.i.MaxCount() || !i.i.Comparable(otherBehaviour.i) {

		return false
	}
	a, b := otherBehaviour.i.AddStack(i.i)

	tx.AddEntity(NewItem(world.EntitySpawnOpts{Position: other.Position(), Velocity: other.Velocity()}, a))
	if !b.Empty() {
		tx.AddEntity(NewItem(world.EntitySpawnOpts{Position: pos, Velocity: e.Velocity()}, b))
	}
	_ = e.Close()
	_ = other.Close()
	return true
}

func (i *ItemBehaviour) collect(e *Ent, collector Collector, tx *world.Tx) {
	pos := e.Position()
	n, _ := collector.Collect(i.i)
	if n == 0 {
		return
	}
	for _, viewer := range tx.Viewers(pos) {
		viewer.ViewEntityAction(e, PickedUpAction{Collector: collector})
	}

	if n == i.i.Count() {

		_ = e.Close()
		return
	}

	tx.AddEntity(NewItem(world.EntitySpawnOpts{Position: pos}, i.i.Grow(-n)))
	_ = e.Close()
}

type Collector interface {
	world.Entity

	Collect(stack item.Stack) (n int, ok bool)
}
