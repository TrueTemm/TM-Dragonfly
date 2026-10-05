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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type InventoryTransaction struct {
	LegacyRequestID    int32
	LegacySetItemSlots []protocol.LegacySetItemSlot
	Actions            []InventoryAction786
	TransactionData    txData786
}

func (*InventoryTransaction) ID() uint32 {
	return IDInventoryTransaction
}

func (pk *InventoryTransaction) Marshal(io protocol.IO) {
	io.Varint32(&pk.LegacyRequestID)
	if pk.LegacyRequestID != 0 {
		protocol.Slice(io, &pk.LegacySetItemSlots)
	}
	transactionDataType786(io, &pk.TransactionData)
	protocol.Slice(io, &pk.Actions)
	if pk.TransactionData != nil {
		pk.TransactionData.Marshal(io)
	}
}

type InventoryAction786 struct {
	SourceType    uint32
	WindowID      int32
	SourceFlags   uint32
	InventorySlot uint32
	OldItem       protocol.ItemInstance
	NewItem       protocol.ItemInstance
}

func (x *InventoryAction786) Marshal(r protocol.IO) {
	r.Varuint32(&x.SourceType)
	switch x.SourceType {
	case protocol.InventoryActionSourceContainer, protocol.InventoryActionSourceTODO:
		r.Varint32(&x.WindowID)
	case protocol.InventoryActionSourceWorld:
		r.Varuint32(&x.SourceFlags)
	}
	r.Varuint32(&x.InventorySlot)
	r.ItemInstance(&x.OldItem)
	r.ItemInstance(&x.NewItem)
}

type txData786 interface {
	Marshal(r protocol.IO)
	id() uint32
}

func transactionDataType786(r protocol.IO, x *txData786) {
	var id uint32
	if *x != nil {
		id = (*x).id()
	}
	r.Varuint32(&id)
	if *x == nil {
		switch id {
		case protocol.InventoryTransactionTypeNormal:
			*x = &normalTx786{}
		case protocol.InventoryTransactionTypeMismatch:
			*x = &mismatchTx786{}
		case protocol.InventoryTransactionTypeUseItem:
			*x = &useItemTx786{}
		case protocol.InventoryTransactionTypeUseItemOnEntity:
			*x = &useItemOnEntityTx786{}
		case protocol.InventoryTransactionTypeReleaseItem:
			*x = &releaseItemTx786{}
		default:
			r.UnknownEnumOption(id, "v786 inventory transaction data type")
		}
	}
}

type normalTx786 struct{}

func (*normalTx786) id() uint32          { return protocol.InventoryTransactionTypeNormal }
func (*normalTx786) Marshal(protocol.IO) {}

type mismatchTx786 struct{}

func (*mismatchTx786) id() uint32          { return protocol.InventoryTransactionTypeMismatch }
func (*mismatchTx786) Marshal(protocol.IO) {}

type useItemTx786 struct {
	ActionType          uint32
	TriggerType         uint32
	BlockPosition       protocol.BlockPos
	BlockFace           int32
	HotBarSlot          int32
	HeldItem            protocol.ItemInstance
	Position            mgl32.Vec3
	ClickedPosition     mgl32.Vec3
	BlockRuntimeID      uint32
	ClientPrediction    uint32
	ClientCooldownState byte
}

func (*useItemTx786) id() uint32 { return protocol.InventoryTransactionTypeUseItem }
func (d *useItemTx786) Marshal(r protocol.IO) {

	p := ProtoOf(r)
	r.Varuint32(&d.ActionType)
	if p == 0 || p >= 712 {
		r.Varuint32(&d.TriggerType)
	}
	UBlockPos786(r, &d.BlockPosition)
	r.Varint32(&d.BlockFace)
	r.Varint32(&d.HotBarSlot)
	r.ItemInstance(&d.HeldItem)
	r.Vec3(&d.Position)
	r.Vec3(&d.ClickedPosition)
	r.Varuint32(&d.BlockRuntimeID)
	if p == 0 || p >= 712 {
		r.Varuint32(&d.ClientPrediction)
	}
	if p >= 944 {
		r.Uint8(&d.ClientCooldownState)
	}
}

type useItemOnEntityTx786 struct {
	TargetEntityRuntimeID uint64
	ActionType            uint32
	HotBarSlot            int32
	HeldItem              protocol.ItemInstance
	Position              mgl32.Vec3
	ClickedPosition       mgl32.Vec3
}

func (*useItemOnEntityTx786) id() uint32 { return protocol.InventoryTransactionTypeUseItemOnEntity }
func (d *useItemOnEntityTx786) Marshal(r protocol.IO) {
	r.Varuint64(&d.TargetEntityRuntimeID)
	r.Varuint32(&d.ActionType)
	r.Varint32(&d.HotBarSlot)
	r.ItemInstance(&d.HeldItem)
	r.Vec3(&d.Position)
	r.Vec3(&d.ClickedPosition)
}

type releaseItemTx786 struct {
	ActionType   uint32
	HotBarSlot   int32
	HeldItem     protocol.ItemInstance
	HeadPosition mgl32.Vec3
}

func (*releaseItemTx786) id() uint32 { return protocol.InventoryTransactionTypeReleaseItem }
func (d *releaseItemTx786) Marshal(r protocol.IO) {
	r.Varuint32(&d.ActionType)
	r.Varint32(&d.HotBarSlot)
	r.ItemInstance(&d.HeldItem)
	r.Vec3(&d.HeadPosition)
}
