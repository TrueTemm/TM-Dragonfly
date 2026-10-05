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

package v1001

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

type InventoryContent struct {
	WindowID    uint32
	Content     []protocol.ItemInstance
	Container   protocol.FullContainerName
	StorageItem protocol.ItemInstance
}

func (*InventoryContent) ID() uint32 { return packet.IDInventoryContent }

func (pk *InventoryContent) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	protocol.FuncSlice(io, &pk.Content, func(i *protocol.ItemInstance) { v786.ItemInstanceNew(io, i) })
	protocol.Single(io, &pk.Container)
	v786.ItemInstanceNew(io, &pk.StorageItem)
}

type MobArmourEquipment struct {
	EntityRuntimeID uint64
	Helmet          protocol.ItemInstance
	Chestplate      protocol.ItemInstance
	Leggings        protocol.ItemInstance
	Boots           protocol.ItemInstance
	Body            protocol.ItemInstance
}

func (*MobArmourEquipment) ID() uint32 { return packet.IDMobArmourEquipment }

func (pk *MobArmourEquipment) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	v786.ItemInstanceNew(io, &pk.Helmet)
	v786.ItemInstanceNew(io, &pk.Chestplate)
	v786.ItemInstanceNew(io, &pk.Leggings)
	v786.ItemInstanceNew(io, &pk.Boots)
	v786.ItemInstanceNew(io, &pk.Body)
}

type ServerBoundDiagnostics struct {
	AverageFramesPerSecond        float32
	AverageServerSimTickTime      float32
	AverageClientSimTickTime      float32
	AverageBeginFrameTime         float32
	AverageInputTime              float32
	AverageRenderTime             float32
	AverageEndFrameTime           float32
	AverageRemainderTimePercent   float32
	AverageUnaccountedTimePercent float32
	MemoryCategoryValues          []protocol.MemoryCategoryCounter
	EntityDiagnostics             []protocol.EntityDiagnosticTimingInfo
	SystemDiagnostics             []protocol.SystemDiagnosticTimingInfo
	WhiskerScopes                 []protocol.WhiskerScopeDataSummary
}

func (*ServerBoundDiagnostics) ID() uint32 { return packet.IDServerBoundDiagnostics }

func (pk *ServerBoundDiagnostics) Marshal(io protocol.IO) {
	io.Float32(&pk.AverageFramesPerSecond)
	io.Float32(&pk.AverageServerSimTickTime)
	io.Float32(&pk.AverageClientSimTickTime)
	io.Float32(&pk.AverageBeginFrameTime)
	io.Float32(&pk.AverageInputTime)
	io.Float32(&pk.AverageRenderTime)
	io.Float32(&pk.AverageEndFrameTime)
	io.Float32(&pk.AverageRemainderTimePercent)
	io.Float32(&pk.AverageUnaccountedTimePercent)
	protocol.Slice(io, &pk.MemoryCategoryValues)
	protocol.Slice(io, &pk.EntityDiagnostics)
	protocol.Slice(io, &pk.SystemDiagnostics)
	protocol.Slice(io, &pk.WhiskerScopes)
}

type SubChunkRequest struct {
	Dimension int32
	Offsets   []protocol.SubChunkOffset
	Position  protocol.SubChunkPos
}

func (*SubChunkRequest) ID() uint32 { return packet.IDSubChunkRequest }

func (pk *SubChunkRequest) Marshal(io protocol.IO) {
	io.Varint32(&pk.Dimension)
	protocol.Slice(io, &pk.Offsets)
	io.Int32(&pk.Position[0])
	io.Int32(&pk.Position[1])
	io.Int32(&pk.Position[2])
}

type ServerBoundDataDrivenScreenClosed struct {
	FormID      protocol.Optional[uint32]
	CloseReason string
}

func (*ServerBoundDataDrivenScreenClosed) ID() uint32 {
	return packet.IDServerBoundDataDrivenScreenClosed
}

func (pk *ServerBoundDataDrivenScreenClosed) Marshal(io protocol.IO) {
	protocol.OptionalFunc(io, &pk.FormID, io.Uint32)
	io.String(&pk.CloseReason)
}

type InventoryTransaction struct {
	LegacyRequestID    int32
	LegacySetItemSlots []protocol.LegacySetItemSlot
	Actions            []InventoryAction
	TransactionData    txData
}

func (*InventoryTransaction) ID() uint32 { return packet.IDInventoryTransaction }

func (pk *InventoryTransaction) Marshal(io protocol.IO) {
	io.Varint32(&pk.LegacyRequestID)
	hasLegacy := pk.LegacyRequestID < -1 && (pk.LegacyRequestID&1) == 0
	io.Bool(&hasLegacy)
	if hasLegacy {
		protocol.Slice(io, &pk.LegacySetItemSlots)
	}
	hasType := true
	io.Bool(&hasType)
	if !hasType {
		io.InvalidValue(hasType, "InventoryTransaction transaction type", "expected presence bool to be true")
	}
	transactionDataType(io, &pk.TransactionData)
	hasActions := true
	io.Bool(&hasActions)
	if !hasActions {
		io.InvalidValue(hasActions, "InventoryTransaction actions", "expected presence bool to be true")
	}
	protocol.Slice(io, &pk.Actions)
	if pk.TransactionData != nil {
		pk.TransactionData.Marshal(io)
	}
}

type InventoryAction struct {
	SourceType    uint32
	WindowID      int8
	SourceFlags   uint32
	InventorySlot uint32
	OldItem       protocol.ItemInstance
	NewItem       protocol.ItemInstance
}

func (x *InventoryAction) Marshal(r protocol.IO) {
	r.Varuint32(&x.SourceType)
	present := true
	r.Bool(&present)
	hasContainerID := x.SourceType == protocol.InventoryActionSourceContainer || x.SourceType == protocol.InventoryActionSourceTODO
	r.Bool(&hasContainerID)
	if hasContainerID {
		r.Int8(&x.WindowID)
	}
	r.Bool(&present)
	hasFlags := x.SourceType == protocol.InventoryActionSourceWorld
	r.Bool(&hasFlags)
	if hasFlags {
		r.Varuint32(&x.SourceFlags)
	}
	r.Varuint32(&x.InventorySlot)
	v786.ItemInstanceNew(r, &x.OldItem)
	v786.ItemInstanceNew(r, &x.NewItem)
}

type txData interface {
	Marshal(r protocol.IO)
	id() uint32
}

func transactionDataType(r protocol.IO, x *txData) {
	var id uint32
	if *x != nil {
		id = (*x).id()
	}
	r.Varuint32(&id)
	if *x == nil {
		switch id {
		case protocol.InventoryTransactionTypeNormal:
			*x = &normalTx{}
		case protocol.InventoryTransactionTypeMismatch:
			*x = &mismatchTx{}
		case protocol.InventoryTransactionTypeUseItem:
			*x = &useItemTx{}
		case protocol.InventoryTransactionTypeUseItemOnEntity:
			*x = &useItemOnEntityTx{}
		case protocol.InventoryTransactionTypeReleaseItem:
			*x = &releaseItemTx{}
		default:
			r.UnknownEnumOption(id, "v1001 inventory transaction data type")
		}
	}
}

type normalTx struct{}

func (*normalTx) id() uint32          { return protocol.InventoryTransactionTypeNormal }
func (*normalTx) Marshal(protocol.IO) {}

type mismatchTx struct{}

func (*mismatchTx) id() uint32          { return protocol.InventoryTransactionTypeMismatch }
func (*mismatchTx) Marshal(protocol.IO) {}

type useItemTx struct {
	ActionType          uint32
	TriggerType         uint32
	BlockPosition       protocol.BlockPos
	BlockFace           int32
	HotBarSlot          int32
	HeldItem            protocol.ItemInstance
	Position            mgl32.Vec3
	ClickedPosition     mgl32.Vec3
	BlockRuntimeID      uint32
	ClientPrediction    uint8
	ClientCooldownState byte
}

func (*useItemTx) id() uint32 { return protocol.InventoryTransactionTypeUseItem }
func (d *useItemTx) Marshal(r protocol.IO) {
	protocol.IntegerFunc(&d.ActionType, r.Varint32)
	protocol.IntegerFunc(&d.TriggerType, r.Uint8)
	r.BlockPos(&d.BlockPosition)
	protocol.IntegerFunc(&d.BlockFace, r.Uint8)
	r.Varint32(&d.HotBarSlot)
	v786.ItemInstanceNew(r, &d.HeldItem)
	r.Vec3(&d.Position)
	r.Vec3(&d.ClickedPosition)
	r.Varuint32(&d.BlockRuntimeID)
	r.Uint8(&d.ClientPrediction)
	r.Uint8(&d.ClientCooldownState)
}

type useItemOnEntityTx struct {
	TargetEntityRuntimeID uint64
	ActionType            int32
	HotBarSlot            int32
	HeldItem              protocol.ItemInstance
	Position              mgl32.Vec3
	ClickedPosition       mgl32.Vec3
}

func (*useItemOnEntityTx) id() uint32 { return protocol.InventoryTransactionTypeUseItemOnEntity }
func (d *useItemOnEntityTx) Marshal(r protocol.IO) {
	r.Varuint64(&d.TargetEntityRuntimeID)
	r.Varint32(&d.ActionType)
	r.Varint32(&d.HotBarSlot)
	v786.ItemInstanceNew(r, &d.HeldItem)
	r.Vec3(&d.Position)
	r.Vec3(&d.ClickedPosition)
}

type releaseItemTx struct {
	ActionType   int32
	HotBarSlot   int32
	HeldItem     protocol.ItemInstance
	HeadPosition mgl32.Vec3
}

func (*releaseItemTx) id() uint32 { return protocol.InventoryTransactionTypeReleaseItem }
func (d *releaseItemTx) Marshal(r protocol.IO) {
	r.Varint32(&d.ActionType)
	r.Varint32(&d.HotBarSlot)
	v786.ItemInstanceNew(r, &d.HeldItem)
	r.Vec3(&d.HeadPosition)
}
