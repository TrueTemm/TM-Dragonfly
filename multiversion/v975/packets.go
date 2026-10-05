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

package v975

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

type InventorySlot struct {
	WindowID    uint32
	Slot        uint32
	Container   protocol.Optional[protocol.FullContainerName]
	StorageItem protocol.Optional[protocol.ItemInstance]
	NewItem     protocol.ItemInstance
}

func (*InventorySlot) ID() uint32 { return packet.IDInventorySlot }

func (pk *InventorySlot) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	io.Varuint32(&pk.Slot)
	protocol.OptionalMarshaler(io, &pk.Container)
	protocol.OptionalFunc(io, &pk.StorageItem, func(i *protocol.ItemInstance) { v786.ItemInstanceNew(io, i) })
	v786.ItemInstanceNew(io, &pk.NewItem)
}

type MobEquipment struct {
	EntityRuntimeID uint64
	NewItem         protocol.ItemInstance
	InventorySlot   byte
	HotBarSlot      byte
	WindowID        byte
}

func (*MobEquipment) ID() uint32 { return packet.IDMobEquipment }

func (pk *MobEquipment) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	v786.ItemInstanceNew(io, &pk.NewItem)
	io.Uint8(&pk.InventorySlot)
	io.Uint8(&pk.HotBarSlot)
	io.Uint8(&pk.WindowID)
}

type LevelSoundEvent struct {
	v786.LevelSoundEvent
	FireAtPosition protocol.Optional[mgl32.Vec3]
}

func (*LevelSoundEvent) ID() uint32 { return packet.IDLevelSoundEvent }

func (pk *LevelSoundEvent) Marshal(io protocol.IO) {
	pk.LevelSoundEvent.Marshal(io)
	protocol.OptionalFunc(io, &pk.FireAtPosition, io.Vec3)
}

type PlaySound struct {
	v786.PlaySound
	Handle protocol.Optional[uint64]
}

func (*PlaySound) ID() uint32 { return packet.IDPlaySound }

func (pk *PlaySound) Marshal(io protocol.IO) {
	pk.PlaySound.Marshal(io)
	protocol.OptionalFunc(io, &pk.Handle, io.Uint64)
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
}
