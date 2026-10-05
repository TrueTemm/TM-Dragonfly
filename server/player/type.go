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

package player

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

var Type ptype

type ptype struct{}

func (t ptype) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	pd := data.Data.(*playerData)
	p := &Player{
		tx:         tx,
		handle:     handle,
		data:       data,
		playerData: pd,
	}

	pd.offHand.SlotValidatorFunc(func(s item.Stack, _ int) bool {
		if s.Empty() {
			return true
		}
		it, allowedInOffhand := s.Item().(item.OffHand)
		return allowedInOffhand && it.OffHand()
	})

	if pd.s != nil {
		pd.s.HandleInventories(tx, p, pd.inv, pd.offHand, pd.enderChest, pd.ui, pd.armour, pd.heldSlot)
	} else {
		pd.inv.SlotFunc(func(slot int, before, after item.Stack) {
			if slot == int(*p.heldSlot) {
				p.broadcastItems(slot, before, after)
			}
		})
		pd.offHand.SlotFunc(p.broadcastItems)
		pd.armour.Inventory().SlotFunc(p.broadcastArmour)
	}
	return p
}

func (ptype) EncodeEntity() string   { return "minecraft:player" }
func (ptype) NetworkOffset() float64 { return 1.621 }
func (ptype) BBox(e world.Entity) cube.BBox {
	p := e.(*Player)
	s := p.Scale()
	_, sleeping := p.Sleeping()
	switch {
	case sleeping:
		return cube.Box(-0.1*s, 0, -0.1*s, 0.1*s, 0.2*s, 0.1*s)
	case p.Gliding(), p.Swimming(), p.Crawling():
		return cube.Box(-0.3*s, 0, -0.3*s, 0.3*s, 0.6*s, 0.3*s)
	case p.Sneaking():
		return cube.Box(-0.3*s, 0, -0.3*s, 0.3*s, 1.49*s, 0.3*s)
	default:
		return cube.Box(-0.3*s, 0, -0.3*s, 0.3*s, 1.8*s, 0.3*s)
	}
}
func (t ptype) DecodeNBT(map[string]any, *world.EntityData) {}
func (t ptype) EncodeNBT(*world.EntityData) map[string]any  { return nil }
