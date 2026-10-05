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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Reader2169 struct {
	*protocol.Reader
}

type Writer2169 struct {
	*protocol.Writer
}

func DoubleOptionalFunc[T any](r protocol.IO, x *protocol.Optional[T], f func(*T)) {
	outer := true
	r.Bool(&outer)
	if outer {
		protocol.OptionalFunc(r, x, f)
		return
	}
	*x = protocol.Optional[T]{}
}

type InventoryAction struct {
	protocol.InventoryAction
}

func (x *InventoryAction) Marshal(r protocol.IO) {
	r.Varuint32(&x.SourceType)
	DoubleOptionalFunc(r, &x.WindowID, r.Int8)
	DoubleOptionalFunc(r, &x.SourceFlags, r.Varuint32)
	r.Varuint32(&x.InventorySlot)
	r.ItemInstance(&x.OldItem)
	r.ItemInstance(&x.NewItem)
}

func inventoryActions(actions []protocol.InventoryAction) []InventoryAction {
	out := make([]InventoryAction, len(actions))
	for i, a := range actions {
		out[i] = InventoryAction{InventoryAction: a}
	}
	return out
}

func latestInventoryActions(actions []InventoryAction) []protocol.InventoryAction {
	out := make([]protocol.InventoryAction, len(actions))
	for i, a := range actions {
		out[i] = a.InventoryAction
	}
	return out
}

func (r Reader2169) PlayerInventoryAction(x *protocol.UseItemTransactionData) {
	playerInventoryAction(r, x)
}

func (w Writer2169) PlayerInventoryAction(x *protocol.UseItemTransactionData) {
	playerInventoryAction(w, x)
}

func playerInventoryAction(r protocol.IO, x *protocol.UseItemTransactionData) {
	r.Varint32(&x.LegacyRequestID)
	protocol.OptionalFunc(r, &x.LegacySetItemSlots, func(slots *[]protocol.LegacySetItemSlot) {
		protocol.Slice(r, slots)
	})

	var actions protocol.Optional[[]InventoryAction]
	if len(x.Actions) != 0 {
		actions = protocol.Option(inventoryActions(x.Actions))
	}
	DoubleOptionalFunc(r, &actions, func(actions *[]InventoryAction) {
		protocol.Slice(r, actions)
	})
	value, _ := actions.Value()
	x.Actions = latestInventoryActions(value)

	protocol.IntegerFunc(&x.ActionType, r.Varint32)
	protocol.IntegerFunc(&x.TriggerType, r.Uint8)
	r.BlockPos(&x.BlockPosition)
	protocol.IntegerFunc(&x.BlockFace, r.Uint8)
	r.Varint32(&x.HotBarSlot)

	r.ItemInstance(&x.HeldItem)
	r.Vec3(&x.Position)
	r.Vec3(&x.ClickedPosition)
	r.Varuint32(&x.BlockRuntimeID)
	r.Uint8(&x.ClientPrediction)
	r.Uint8(&x.ClientCooldownState)
}

func useItemTransactionData(r protocol.IO, data *protocol.UseItemTransactionData) {
	protocol.IntegerFunc(&data.ActionType, r.Varint32)
	protocol.IntegerFunc(&data.TriggerType, r.Uint8)
	r.BlockPos(&data.BlockPosition)
	protocol.IntegerFunc(&data.BlockFace, r.Uint8)
	r.Varint32(&data.HotBarSlot)
	r.ItemInstance(&data.HeldItem)
	r.Vec3(&data.Position)
	r.Vec3(&data.ClickedPosition)
	r.Varuint32(&data.BlockRuntimeID)
	r.Uint8(&data.ClientPrediction)
	r.Uint8(&data.ClientCooldownState)
}

func inputFlagList(r protocol.IO, x *protocol.InputFlags, size int) {
	present := x.Present()
	r.Bool(&present)
	if !present {
		*x = protocol.InputFlags{}
		return
	}
	protocol.InputFlagList(r, x, size)
}
