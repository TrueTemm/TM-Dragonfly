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
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v898"
	"github.com/df-mc/dragonfly/multiversion/v975"
)

func applyDeltas1001(p packet.Pool) {
	p[v898.IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[packet.IDInventoryContent] = func() packet.Packet { return &InventoryContent{} }
	p[packet.IDMobArmourEquipment] = func() packet.Packet { return &MobArmourEquipment{} }
	p[packet.IDServerBoundDiagnostics] = func() packet.Packet { return &ServerBoundDiagnostics{} }
	p[packet.IDSubChunkRequest] = func() packet.Packet { return &SubChunkRequest{} }
	p[packet.IDInventoryTransaction] = func() packet.Packet { return &InventoryTransaction{} }
	p[packet.IDServerBoundDataDrivenScreenClosed] = func() packet.Packet { return &ServerBoundDataDrivenScreenClosed{} }

	p[packet.IDLevelSoundEvent] = func() packet.Packet { return &packet.LevelSoundEvent{} }
	p[packet.IDBossEvent] = func() packet.Packet { return &packet.BossEvent{} }
	p[packet.IDClientCacheBlobStatus] = func() packet.Packet { return &packet.ClientCacheBlobStatus{} }
	p[packet.IDResourcePackChunkRequest] = func() packet.Packet { return &packet.ResourcePackChunkRequest{} }
	p[packet.IDGraphicsOverrideParameter] = func() packet.Packet { return &packet.GraphicsOverrideParameter{} }

	p[packet.IDSendPartyDestinationCookie] = func() packet.Packet { return &packet.SendPartyDestinationCookie{} }
	p[packet.IDPartyDestinationCookieResponse] = func() packet.Packet { return &packet.PartyDestinationCookieResponse{} }
}

func NewClientPool() packet.Pool {
	p := v975.NewClientPool()
	applyDeltas1001(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v975.NewServerPool()
	applyDeltas1001(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v975.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v975.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		return []packet.Packet{FromLatestStartGame1001(pk)}, true
	case *packet.InventoryContent:
		return []packet.Packet{&InventoryContent{WindowID: pk.WindowID, Content: pk.Content, Container: pk.Container,
			StorageItem: pk.StorageItem}}, true
	case *packet.MobArmourEquipment:
		return []packet.Packet{&MobArmourEquipment{EntityRuntimeID: pk.EntityRuntimeID, Helmet: pk.Helmet,
			Chestplate: pk.Chestplate, Leggings: pk.Leggings, Boots: pk.Boots, Body: pk.Body}}, true
	case *packet.LevelSoundEvent, *packet.BossEvent, *packet.ClientCacheBlobStatus, *packet.ResourcePackChunkRequest,
		*packet.GraphicsOverrideParameter, *packet.SendPartyDestinationCookie, *packet.PartyDestinationCookieResponse:
		return []packet.Packet{pk}, true
	case *packet.ClientboundUpdateSoundData:
		return nil, true
	case *packet.ItemRegistry:
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items1001()}}, true
	case *packet.ServerBoundDiagnostics:
		return []packet.Packet{&ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
			AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
			AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
			AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
			AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
			MemoryCategoryValues: pk.MemoryCategoryValues, EntityDiagnostics: pk.EntityDiagnostics,
			SystemDiagnostics: pk.SystemDiagnostics, WhiskerScopes: pk.WhiskerScopes}}, true
	case *packet.SubChunkRequest:
		return []packet.Packet{&SubChunkRequest{Dimension: pk.Dimension, Offsets: pk.Offsets, Position: pk.Position}}, true
	case *packet.ServerBoundDataDrivenScreenClosed:
		return []packet.Packet{&ServerBoundDataDrivenScreenClosed{FormID: protocol.Option(pk.FormID), CloseReason: pk.CloseReason}}, true
	case *packet.InventoryTransaction:

		out := &InventoryTransaction{LegacyRequestID: pk.LegacyRequestID, LegacySetItemSlots: pk.LegacySetItemSlots,
			TransactionData: &normalTx{}}
		out.Actions = make([]InventoryAction, len(pk.Actions))
		for i, a := range pk.Actions {
			na := InventoryAction{SourceType: a.SourceType, InventorySlot: a.InventorySlot, OldItem: a.OldItem, NewItem: a.NewItem}
			if v, ok := a.WindowID.Value(); ok {
				na.WindowID = v
			}
			if v, ok := a.SourceFlags.Value(); ok {
				na.SourceFlags = v
			}
			out.Actions[i] = na
		}
		return []packet.Packet{out}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *StartGame:
		return []packet.Packet{ToLatestStartGame1001(pk)}, true
	case *InventoryContent:
		return []packet.Packet{&packet.InventoryContent{WindowID: pk.WindowID, Content: pk.Content, Container: pk.Container,
			StorageItem: pk.StorageItem}}, true
	case *MobArmourEquipment:
		return []packet.Packet{&packet.MobArmourEquipment{EntityRuntimeID: pk.EntityRuntimeID, Helmet: pk.Helmet,
			Chestplate: pk.Chestplate, Leggings: pk.Leggings, Boots: pk.Boots, Body: pk.Body}}, true
	case *packet.LevelSoundEvent, *packet.BossEvent, *packet.ClientCacheBlobStatus, *packet.ResourcePackChunkRequest,
		*packet.GraphicsOverrideParameter, *packet.SendPartyDestinationCookie, *packet.PartyDestinationCookieResponse:
		return []packet.Packet{pk}, true
	case *ServerBoundDiagnostics:
		return []packet.Packet{&packet.ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
			AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
			AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
			AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
			AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
			MemoryCategoryValues: pk.MemoryCategoryValues, EntityDiagnostics: pk.EntityDiagnostics,
			SystemDiagnostics: pk.SystemDiagnostics, WhiskerScopes: pk.WhiskerScopes}}, true
	case *SubChunkRequest:
		return []packet.Packet{&packet.SubChunkRequest{Dimension: pk.Dimension, Offsets: pk.Offsets, Position: pk.Position}}, true
	case *ServerBoundDataDrivenScreenClosed:
		out := &packet.ServerBoundDataDrivenScreenClosed{CloseReason: pk.CloseReason}
		if v, ok := pk.FormID.Value(); ok {
			out.FormID = v
		}
		return []packet.Packet{out}, true
	case *InventoryTransaction:
		return []packet.Packet{toLatestInventoryTransaction(pk)}, true
	}
	return nil, false
}

func toLatestInventoryTransaction(pk *InventoryTransaction) *packet.InventoryTransaction {
	out := &packet.InventoryTransaction{LegacyRequestID: pk.LegacyRequestID, LegacySetItemSlots: pk.LegacySetItemSlots}
	out.Actions = make([]protocol.InventoryAction, len(pk.Actions))
	for i, a := range pk.Actions {
		la := protocol.InventoryAction{SourceType: a.SourceType, InventorySlot: a.InventorySlot, OldItem: a.OldItem, NewItem: a.NewItem}
		switch a.SourceType {
		case protocol.InventoryActionSourceContainer, protocol.InventoryActionSourceTODO:
			la.WindowID = protocol.Option(a.WindowID)
		case protocol.InventoryActionSourceWorld:
			la.SourceFlags = protocol.Option(a.SourceFlags)
		}
		out.Actions[i] = la
	}
	switch d := pk.TransactionData.(type) {
	case *useItemTx:
		out.TransactionData = &protocol.UseItemTransactionData{
			ActionType: d.ActionType, TriggerType: d.TriggerType, BlockPosition: d.BlockPosition,
			BlockFace: d.BlockFace, HotBarSlot: d.HotBarSlot, HeldItem: d.HeldItem, Position: d.Position,
			ClickedPosition: d.ClickedPosition, BlockRuntimeID: d.BlockRuntimeID,
			ClientPrediction: d.ClientPrediction, ClientCooldownState: d.ClientCooldownState,
		}
	case *useItemOnEntityTx:
		out.TransactionData = &protocol.UseItemOnEntityTransactionData{
			TargetEntityRuntimeID: d.TargetEntityRuntimeID, ActionType: d.ActionType, HotBarSlot: d.HotBarSlot,
			HeldItem: d.HeldItem, Position: d.Position, ClickedPosition: d.ClickedPosition,
		}
	case *releaseItemTx:
		out.TransactionData = &protocol.ReleaseItemTransactionData{
			ActionType: d.ActionType, HotBarSlot: d.HotBarSlot, HeldItem: d.HeldItem, HeadPosition: d.HeadPosition,
		}
	case *mismatchTx:
		out.TransactionData = &protocol.MismatchTransactionData{}
	default:
		out.TransactionData = &protocol.NormalTransactionData{}
	}
	return out
}
