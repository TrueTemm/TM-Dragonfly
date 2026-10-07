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
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type BlockActorData struct {
	Position protocol.BlockPos
	NBTData  map[string]any
}

func (*BlockActorData) ID() uint32 { return IDBlockActorData }
func (pk *BlockActorData) Marshal(io protocol.IO) {
	UBlockPos786(io, &pk.Position)
	io.NBT(&pk.NBTData, nbt.NetworkLittleEndian)
}
func toLatestBlockActorData(pk *BlockActorData) *packet.BlockActorData {
	return &packet.BlockActorData{Position: pk.Position, NBTData: pk.NBTData}
}
func fromLatestBlockActorData(pk *packet.BlockActorData) *BlockActorData {
	return &BlockActorData{Position: pk.Position, NBTData: pk.NBTData}
}

type ContainerOpen struct {
	WindowID                byte
	ContainerType           byte
	ContainerPosition       protocol.BlockPos
	ContainerEntityUniqueID int64
}

func (*ContainerOpen) ID() uint32 { return IDContainerOpen }
func (pk *ContainerOpen) Marshal(io protocol.IO) {
	io.Uint8(&pk.WindowID)
	io.Uint8(&pk.ContainerType)
	UBlockPos786(io, &pk.ContainerPosition)
	io.Varint64(&pk.ContainerEntityUniqueID)
}
func toLatestContainerOpen(pk *ContainerOpen) *packet.ContainerOpen {
	return &packet.ContainerOpen{
		WindowID: pk.WindowID, ContainerType: pk.ContainerType,
		ContainerPosition: pk.ContainerPosition, ContainerEntityUniqueID: pk.ContainerEntityUniqueID,
	}
}
func fromLatestContainerOpen(pk *packet.ContainerOpen) *ContainerOpen {
	return &ContainerOpen{
		WindowID: pk.WindowID, ContainerType: pk.ContainerType,
		ContainerPosition: pk.ContainerPosition, ContainerEntityUniqueID: pk.ContainerEntityUniqueID,
	}
}

type ContainerClose struct {
	WindowID      byte
	ContainerType byte
	ServerSide    bool
}

func (*ContainerClose) ID() uint32 { return IDContainerClose }
func (pk *ContainerClose) Marshal(io protocol.IO) {
	io.Uint8(&pk.WindowID)
	if p := ProtoOf(io); p == 0 || p >= 685 { // added in 685
		io.Uint8(&pk.ContainerType)
	}
	io.Bool(&pk.ServerSide)
}
func toLatestContainerClose(pk *ContainerClose) *packet.ContainerClose {
	return &packet.ContainerClose{WindowID: pk.WindowID, ContainerType: pk.ContainerType, ServerSide: pk.ServerSide}
}
func fromLatestContainerClose(pk *packet.ContainerClose) *ContainerClose {
	return &ContainerClose{WindowID: pk.WindowID, ContainerType: pk.ContainerType, ServerSide: pk.ServerSide}
}

type LecternUpdate struct {
	Page      byte
	PageCount byte
	Position  protocol.BlockPos
}

func (*LecternUpdate) ID() uint32 { return IDLecternUpdate }
func (pk *LecternUpdate) Marshal(io protocol.IO) {
	io.Uint8(&pk.Page)
	io.Uint8(&pk.PageCount)
	UBlockPos786(io, &pk.Position)
}
func toLatestLecternUpdate(pk *LecternUpdate) *packet.LecternUpdate {
	return &packet.LecternUpdate{Page: pk.Page, PageCount: pk.PageCount, Position: pk.Position}
}
func fromLatestLecternUpdate(pk *packet.LecternUpdate) *LecternUpdate {
	return &LecternUpdate{Page: pk.Page, PageCount: pk.PageCount, Position: pk.Position}
}

type OpenSign struct {
	Position  protocol.BlockPos
	FrontSide bool
}

func (*OpenSign) ID() uint32 { return IDOpenSign }
func (pk *OpenSign) Marshal(io protocol.IO) {
	UBlockPos786(io, &pk.Position)
	io.Bool(&pk.FrontSide)
}
func toLatestOpenSign(pk *OpenSign) *packet.OpenSign {
	return &packet.OpenSign{Position: pk.Position, FrontSide: pk.FrontSide}
}
func fromLatestOpenSign(pk *packet.OpenSign) *OpenSign {
	return &OpenSign{Position: pk.Position, FrontSide: pk.FrontSide}
}

type StructureTemplateDataRequest struct {
	StructureName string
	Position      protocol.BlockPos
	Settings      protocol.StructureSettings
	RequestType   byte
}

func (*StructureTemplateDataRequest) ID() uint32 { return IDStructureTemplateDataRequest }
func (pk *StructureTemplateDataRequest) Marshal(io protocol.IO) {
	io.String(&pk.StructureName)
	UBlockPos786(io, &pk.Position)
	protocol.Single(io, &pk.Settings)
	io.Uint8(&pk.RequestType)
}
func toLatestStructureTemplateDataRequest(pk *StructureTemplateDataRequest) *packet.StructureTemplateDataRequest {
	return &packet.StructureTemplateDataRequest{
		StructureName: pk.StructureName, Position: pk.Position, Settings: pk.Settings, RequestType: pk.RequestType,
	}
}
func fromLatestStructureTemplateDataRequest(pk *packet.StructureTemplateDataRequest) *StructureTemplateDataRequest {
	return &StructureTemplateDataRequest{
		StructureName: pk.StructureName, Position: pk.Position, Settings: pk.Settings, RequestType: pk.RequestType,
	}
}

type BlockEvent struct {
	Position  protocol.BlockPos
	EventType int32
	EventData int32
}

func (*BlockEvent) ID() uint32 { return IDBlockEvent }
func (pk *BlockEvent) Marshal(io protocol.IO) {
	UBlockPos786(io, &pk.Position)
	io.Varint32(&pk.EventType)
	io.Varint32(&pk.EventData)
}
func fromLatestBlockEvent(pk *packet.BlockEvent) *BlockEvent {
	return &BlockEvent{Position: pk.Position, EventType: pk.EventType, EventData: pk.EventData}
}
