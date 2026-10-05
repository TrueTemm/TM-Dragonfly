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
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v944"
)

func applyDeltas975(p packet.Pool) {
	p[packet.IDInventorySlot] = func() packet.Packet { return &InventorySlot{} }
	p[packet.IDMobEquipment] = func() packet.Packet { return &MobEquipment{} }
	p[packet.IDLevelSoundEvent] = func() packet.Packet { return &LevelSoundEvent{} }
	p[packet.IDPlaySound] = func() packet.Packet { return &PlaySound{} }
	p[packet.IDServerBoundDiagnostics] = func() packet.Packet { return &ServerBoundDiagnostics{} }

	p[packet.IDActorEvent] = func() packet.Packet { return &packet.ActorEvent{} }
	p[packet.IDClientMovementPredictionSync] = func() packet.Packet { return &packet.ClientMovementPredictionSync{} }
	p[packet.IDJigsawStructureData] = func() packet.Packet { return &packet.JigsawStructureData{} }
	p[packet.IDUpdateClientOptions] = func() packet.Packet { return &packet.UpdateClientOptions{} }
	p[packet.IDPartyChanged] = func() packet.Packet { return &packet.PartyChanged{} }

	p[packet.IDServerStoreInfo] = func() packet.Packet { return &packet.ServerStoreInfo{} }
	p[packet.IDServerPresenceInfo] = func() packet.Packet { return &packet.ServerPresenceInfo{} }
}

func NewClientPool() packet.Pool {
	p := v944.NewClientPool()
	applyDeltas975(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v944.NewServerPool()
	applyDeltas975(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v944.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v944.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.InventorySlot:
		return []packet.Packet{&InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, Container: pk.Container,
			StorageItem: pk.StorageItem, NewItem: pk.NewItem}}, true
	case *packet.MobEquipment:
		return []packet.Packet{&MobEquipment{EntityRuntimeID: pk.EntityRuntimeID, NewItem: pk.NewItem,
			InventorySlot: pk.InventorySlot, HotBarSlot: pk.HotBarSlot, WindowID: pk.WindowID}}, true
	case *packet.LevelSoundEvent:
		return []packet.Packet{&LevelSoundEvent{LevelSoundEvent: *v786.FromLatestLevelSoundEvent786(proto, pk)}}, true
	case *packet.PlaySound:
		return []packet.Packet{&PlaySound{PlaySound: v786.PlaySound{SoundName: pk.SoundName, Position: pk.Position,
			Volume: pk.Volume, Pitch: pk.Pitch}}}, true
	case *packet.ActorEvent, *packet.ClientMovementPredictionSync, *packet.JigsawStructureData,
		*packet.UpdateClientOptions, *packet.PartyChanged:
		return []packet.Packet{pk}, true
	case *packet.PrimitiveShapes, *packet.ServerStoreInfo, *packet.ServerPresenceInfo:

		return nil, true
	case *packet.StartGame:

		sg := v944.FromLatestStartGame944(pk)
		sg.BaseGameVersion, sg.GameVersion = "1.26.20", "1.26.20"
		return []packet.Packet{sg}, true
	case *packet.ItemRegistry:
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items975()}}, true
	case *packet.ServerBoundDiagnostics:
		return []packet.Packet{fromLatestServerBoundDiagnostics(pk)}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *InventorySlot:
		return []packet.Packet{&packet.InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, Container: pk.Container,
			StorageItem: pk.StorageItem, NewItem: pk.NewItem}}, true
	case *MobEquipment:
		return []packet.Packet{&packet.MobEquipment{EntityRuntimeID: pk.EntityRuntimeID, NewItem: pk.NewItem,
			InventorySlot: pk.InventorySlot, HotBarSlot: pk.HotBarSlot, WindowID: pk.WindowID}}, true
	case *LevelSoundEvent:
		out := v786.ToLatestLevelSoundEvent786(proto, &pk.LevelSoundEvent)
		out.FireAtPosition = pk.FireAtPosition
		return []packet.Packet{out}, true
	case *PlaySound:
		return []packet.Packet{&packet.PlaySound{SoundName: pk.SoundName, Position: pk.Position, Volume: pk.Volume,
			Pitch: pk.Pitch}}, true
	case *packet.ActorEvent, *packet.ClientMovementPredictionSync, *packet.JigsawStructureData,
		*packet.UpdateClientOptions, *packet.PartyChanged, *packet.ServerStoreInfo, *packet.ServerPresenceInfo:
		return []packet.Packet{pk}, true
	case *ServerBoundDiagnostics:
		return []packet.Packet{toLatestServerBoundDiagnostics(pk)}, true
	}
	return nil, false
}

func fromLatestServerBoundDiagnostics(pk *packet.ServerBoundDiagnostics) *ServerBoundDiagnostics {
	return &ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
		AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
		AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
		AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
		AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
		MemoryCategoryValues: pk.MemoryCategoryValues, EntityDiagnostics: pk.EntityDiagnostics, SystemDiagnostics: pk.SystemDiagnostics}
}

func toLatestServerBoundDiagnostics(pk *ServerBoundDiagnostics) *packet.ServerBoundDiagnostics {
	return &packet.ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
		AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
		AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
		AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
		AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
		MemoryCategoryValues: pk.MemoryCategoryValues, EntityDiagnostics: pk.EntityDiagnostics, SystemDiagnostics: pk.SystemDiagnostics}
}
