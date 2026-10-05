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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func (r Reader786) PlayerInventoryAction(x *protocol.UseItemTransactionData) {
	playerInventoryAction786(r, r.Proto, x)
}

func (w Writer786) PlayerInventoryAction(x *protocol.UseItemTransactionData) {
	playerInventoryAction786(w, w.Proto, x)
}

func playerInventoryAction786(r protocol.IO, proto uint32, x *protocol.UseItemTransactionData) {
	r.Varint32(&x.LegacyRequestID)
	if x.LegacyRequestID < -1 && (x.LegacyRequestID&1) == 0 {
		slots, _ := x.LegacySetItemSlots.Value()
		protocol.Slice(r, &slots)
		x.LegacySetItemSlots = protocol.Option(slots)
	}

	actions := make([]InventoryAction786, len(x.Actions))
	{
		for i, la := range x.Actions {
			a := InventoryAction786{SourceType: la.SourceType, InventorySlot: la.InventorySlot, OldItem: la.OldItem, NewItem: la.NewItem}
			if w, ok := la.WindowID.Value(); ok {
				a.WindowID = int32(w)
			}
			if f, ok := la.SourceFlags.Value(); ok {
				a.SourceFlags = f
			}
			actions[i] = a
		}
	}
	protocol.Slice(r, &actions)
	latest := make([]protocol.InventoryAction, len(actions))
	for i, a := range actions {
		la := protocol.InventoryAction{SourceType: a.SourceType, InventorySlot: a.InventorySlot, OldItem: a.OldItem, NewItem: a.NewItem}
		switch a.SourceType {
		case protocol.InventoryActionSourceContainer, protocol.InventoryActionSourceTODO:
			la.WindowID = protocol.Option(int8(a.WindowID))
		case protocol.InventoryActionSourceWorld:
			la.SourceFlags = protocol.Option(a.SourceFlags)
		}
		latest[i] = la
	}
	x.Actions = latest

	r.Varuint32(&x.ActionType)
	if proto == 0 || proto >= 712 {
		r.Varuint32(&x.TriggerType)
	}
	r.BlockPos(&x.BlockPosition)
	r.Varint32(&x.BlockFace)
	r.Varint32(&x.HotBarSlot)
	r.ItemInstance(&x.HeldItem)
	r.Vec3(&x.Position)
	r.Vec3(&x.ClickedPosition)
	r.Varuint32(&x.BlockRuntimeID)
	switch {
	case proto >= 1001:
		r.Uint8(&x.ClientPrediction)
	case proto == 0 || proto >= 712:
		prediction := uint32(x.ClientPrediction)
		r.Varuint32(&prediction)
		x.ClientPrediction = uint8(prediction)
	}
	if proto >= 944 {
		r.Uint8(&x.ClientCooldownState)
	}
}
