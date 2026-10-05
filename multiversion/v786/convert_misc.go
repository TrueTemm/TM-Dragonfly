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
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func toLatestSetSpawnPosition(pk *SetSpawnPosition) *packet.SetSpawnPosition {
	return &packet.SetSpawnPosition{SpawnType: pk.SpawnType, Position: pk.Position, Dimension: pk.Dimension, SpawnPosition: pk.SpawnPosition}
}
func fromLatestSetSpawnPosition(pk *packet.SetSpawnPosition) *SetSpawnPosition {
	return &SetSpawnPosition{SpawnType: pk.SpawnType, Position: pk.Position, Dimension: pk.Dimension, SpawnPosition: pk.SpawnPosition}
}

func toLatestChangeDimension(pk *ChangeDimension) *packet.ChangeDimension {
	return &packet.ChangeDimension{Dimension: pk.Dimension, Position: pk.Position, Respawn: pk.Respawn, LoadingScreenID: pk.LoadingScreenID}
}
func fromLatestChangeDimension(pk *packet.ChangeDimension) *ChangeDimension {
	return &ChangeDimension{Dimension: pk.Dimension, Position: pk.Position, Respawn: pk.Respawn, LoadingScreenID: pk.LoadingScreenID}
}

func toLatestCorrectPlayerMovePrediction(pk *CorrectPlayerMovePrediction) *packet.CorrectPlayerMovePrediction {
	return &packet.CorrectPlayerMovePrediction{
		PredictionType: pk.PredictionType, Position: pk.Position, Delta: pk.Delta, Rotation: pk.Rotation,
		VehicleAngularVelocity: pk.VehicleAngularVelocity, OnGround: pk.OnGround, Tick: pk.Tick,
	}
}
func fromLatestCorrectPlayerMovePrediction(pk *packet.CorrectPlayerMovePrediction) *CorrectPlayerMovePrediction {
	return &CorrectPlayerMovePrediction{
		PredictionType: pk.PredictionType, Position: pk.Position, Delta: pk.Delta, Rotation: pk.Rotation,
		VehicleAngularVelocity: pk.VehicleAngularVelocity, OnGround: pk.OnGround, Tick: pk.Tick,
	}
}

func toLatestDisconnect(pk *Disconnect) *packet.Disconnect {
	return &packet.Disconnect{Reason: pk.Reason, HideDisconnectionScreen: pk.HideDisconnectionScreen, Message: pk.Message, FilteredMessage: pk.FilteredMessage}
}
func fromLatestDisconnect(pk *packet.Disconnect) *Disconnect {
	return &Disconnect{Reason: pk.Reason, HideDisconnectionScreen: pk.HideDisconnectionScreen, Message: pk.Message, FilteredMessage: pk.FilteredMessage}
}

func toLatestEvent(pk *Event) *packet.Event {
	return &packet.Event{EntityRuntimeID: pk.EntityRuntimeID, UsePlayerID: pk.UsePlayerID != 0, Event: pk.Event}
}
func fromLatestEvent(pk *packet.Event) *Event {
	var b byte
	if pk.UsePlayerID {
		b = 1
	}
	return &Event{EntityRuntimeID: pk.EntityRuntimeID, UsePlayerID: b, Event: pk.Event}
}

func toLatestInteract(pk *Interact) *packet.Interact {
	out := &packet.Interact{ActionType: pk.ActionType, TargetEntityRuntimeID: pk.TargetEntityRuntimeID}
	if pk.ActionType == InteractActionMouseOverEntity || pk.ActionType == InteractActionLeaveVehicle {
		out.Position = protocol.Option(pk.Position)
	}
	return out
}
func fromLatestInteract(pk *packet.Interact) *Interact {
	out := &Interact{ActionType: pk.ActionType, TargetEntityRuntimeID: pk.TargetEntityRuntimeID}
	if v, ok := pk.Position.Value(); ok {
		out.Position = v
	}
	return out
}

func toLatestInventorySlot(pk *InventorySlot) *packet.InventorySlot {
	return &packet.InventorySlot{
		WindowID: pk.WindowID, Slot: pk.Slot, Container: protocol.Option(pk.Container),
		StorageItem: protocol.Option(pk.StorageItem), NewItem: pk.NewItem,
	}
}
func fromLatestInventorySlot(pk *packet.InventorySlot) *InventorySlot {
	out := &InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, NewItem: pk.NewItem}
	if v, ok := pk.Container.Value(); ok {
		out.Container = v
	}
	if v, ok := pk.StorageItem.Value(); ok {
		out.StorageItem = v
	}
	return out
}

func toLatestClientCacheBlobStatus(pk *ClientCacheBlobStatus) *packet.ClientCacheBlobStatus {
	return &packet.ClientCacheBlobStatus{MissHashes: pk.MissHashes, HitHashes: pk.HitHashes}
}

func toLatestRequestChunkRadius(pk *RequestChunkRadius) *packet.RequestChunkRadius {
	return &packet.RequestChunkRadius{ChunkRadius: pk.ChunkRadius, MaxChunkRadius: uint8(pk.MaxChunkRadius)}
}
func fromLatestRequestChunkRadius(pk *packet.RequestChunkRadius) *RequestChunkRadius {
	return &RequestChunkRadius{ChunkRadius: pk.ChunkRadius, MaxChunkRadius: int32(pk.MaxChunkRadius)}
}

func toLatestRequestPermissions(pk *RequestPermissions) *packet.RequestPermissions {
	return &packet.RequestPermissions{EntityUniqueID: pk.EntityUniqueID, PermissionLevel: int32(pk.PermissionLevel), RequestedPermissions: pk.RequestedPermissions}
}
func fromLatestRequestPermissions(pk *packet.RequestPermissions) *RequestPermissions {
	return &RequestPermissions{EntityUniqueID: pk.EntityUniqueID, PermissionLevel: uint8(pk.PermissionLevel), RequestedPermissions: pk.RequestedPermissions}
}

func toLatestResourcePackChunkRequest(pk *ResourcePackChunkRequest) *packet.ResourcePackChunkRequest {
	return &packet.ResourcePackChunkRequest{UUID: pk.UUID, ChunkIndex: int32(pk.ChunkIndex)}
}
func fromLatestResourcePackChunkRequest(pk *packet.ResourcePackChunkRequest) *ResourcePackChunkRequest {
	return &ResourcePackChunkRequest{UUID: pk.UUID, ChunkIndex: uint32(pk.ChunkIndex)}
}

func toLatestResourcePackClientResponse(pk *ResourcePackClientResponse) *packet.ResourcePackClientResponse {
	return &packet.ResourcePackClientResponse{Response: uint32(pk.Response) - 1, PacksToDownload: pk.PacksToDownload}
}
func fromLatestResourcePackClientResponse(pk *packet.ResourcePackClientResponse) *ResourcePackClientResponse {
	return &ResourcePackClientResponse{Response: byte(pk.Response) + 1, PacksToDownload: pk.PacksToDownload}
}

func toLatestResourcePackStack(pk *ResourcePackStack) *packet.ResourcePackStack {
	return &packet.ResourcePackStack{
		TexturePackRequired: pk.TexturePackRequired,
		TexturePacks:        append(append([]protocol.StackResourcePack{}, pk.TexturePacks...), pk.BehaviourPacks...),
		BaseGameVersion:     pk.BaseGameVersion, Experiments: pk.Experiments,
		ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled, IncludeEditorPacks: pk.IncludeEditorPacks,
	}
}
func fromLatestResourcePackStack(pk *packet.ResourcePackStack) *ResourcePackStack {
	return &ResourcePackStack{
		TexturePackRequired: pk.TexturePackRequired, TexturePacks: pk.TexturePacks,
		BaseGameVersion: pk.BaseGameVersion, Experiments: pk.Experiments,
		ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled, IncludeEditorPacks: pk.IncludeEditorPacks,
	}
}

func toLatestResourcePacksInfo(pk *ResourcePacksInfo) *packet.ResourcePacksInfo {
	return &packet.ResourcePacksInfo{
		TexturePackRequired: pk.TexturePackRequired, HasAddons: pk.HasAddons, HasScripts: pk.HasScripts,
		WorldTemplateUUID: pk.WorldTemplateUUID, WorldTemplateVersion: pk.WorldTemplateVersion, TexturePacks: pk.TexturePacks,
	}
}
func fromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	return &ResourcePacksInfo{
		TexturePackRequired: pk.TexturePackRequired, HasAddons: pk.HasAddons, HasScripts: pk.HasScripts,
		WorldTemplateUUID: pk.WorldTemplateUUID, WorldTemplateVersion: pk.WorldTemplateVersion, TexturePacks: pk.TexturePacks,
	}
}

func FromLatestResourcePacksInfo786(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	return fromLatestResourcePacksInfo(pk)
}
func ToLatestResourcePacksInfo786(pk *ResourcePacksInfo) *packet.ResourcePacksInfo {
	return toLatestResourcePacksInfo(pk)
}
func FromLatestResourcePackStack786(pk *packet.ResourcePackStack) *ResourcePackStack {
	return fromLatestResourcePackStack(pk)
}
func ToLatestResourcePackStack786(pk *ResourcePackStack) *packet.ResourcePackStack {
	return toLatestResourcePackStack(pk)
}

func toLatestServerBoundDiagnostics(pk *ServerBoundDiagnostics) *packet.ServerBoundDiagnostics {
	return &packet.ServerBoundDiagnostics{
		AverageFramesPerSecond: pk.AverageFramesPerSecond, AverageServerSimTickTime: pk.AverageServerSimTickTime,
		AverageClientSimTickTime: pk.AverageClientSimTickTime, AverageBeginFrameTime: pk.AverageBeginFrameTime,
		AverageInputTime: pk.AverageInputTime, AverageRenderTime: pk.AverageRenderTime,
		AverageEndFrameTime: pk.AverageEndFrameTime, AverageRemainderTimePercent: pk.AverageRemainderTimePercent,
		AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
	}
}
func fromLatestServerBoundDiagnostics(pk *packet.ServerBoundDiagnostics) *ServerBoundDiagnostics {
	return &ServerBoundDiagnostics{
		AverageFramesPerSecond: pk.AverageFramesPerSecond, AverageServerSimTickTime: pk.AverageServerSimTickTime,
		AverageClientSimTickTime: pk.AverageClientSimTickTime, AverageBeginFrameTime: pk.AverageBeginFrameTime,
		AverageInputTime: pk.AverageInputTime, AverageRenderTime: pk.AverageRenderTime,
		AverageEndFrameTime: pk.AverageEndFrameTime, AverageRemainderTimePercent: pk.AverageRemainderTimePercent,
		AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
	}
}

func toLatestSetPlayerInventoryOptions(pk *SetPlayerInventoryOptions) *packet.SetPlayerInventoryOptions {
	return &packet.SetPlayerInventoryOptions{
		LeftInventoryTab: int32(pk.LeftInventoryTab), RightInventoryTab: int32(pk.RightInventoryTab),
		Filtering: pk.Filtering, InventoryLayout: int32(pk.InventoryLayout), CraftingLayout: int32(pk.CraftingLayout),
	}
}
func fromLatestSetPlayerInventoryOptions(pk *packet.SetPlayerInventoryOptions) *SetPlayerInventoryOptions {
	return &SetPlayerInventoryOptions{
		LeftInventoryTab: byte(pk.LeftInventoryTab), RightInventoryTab: byte(pk.RightInventoryTab),
		Filtering: pk.Filtering, InventoryLayout: byte(pk.InventoryLayout), CraftingLayout: byte(pk.CraftingLayout),
	}
}

func toLatestTransfer(pk *Transfer) *packet.Transfer {
	return &packet.Transfer{Address: pk.Address, Port: pk.Port, ReloadWorld: pk.ReloadWorld}
}
func fromLatestTransfer(pk *packet.Transfer) *Transfer {
	return &Transfer{Address: pk.Address, Port: pk.Port, ReloadWorld: pk.ReloadWorld}
}

func toLatestSubChunkRequest(pk *SubChunkRequest) *packet.SubChunkRequest {
	return &packet.SubChunkRequest{Dimension: pk.Dimension, Offsets: pk.Offsets, Position: pk.Position}
}
func fromLatestSubChunkRequest(pk *packet.SubChunkRequest) *SubChunkRequest {
	return &SubChunkRequest{Dimension: pk.Dimension, Offsets: pk.Offsets, Position: pk.Position}
}

func toLatestUpdateSubChunkBlocks(proto uint32, pk *UpdateSubChunkBlocks) *packet.UpdateSubChunkBlocks {
	return ToLatestUpdateSubChunkBlocks(pk, proto)
}
func fromLatestUpdateSubChunkBlocks(proto uint32, pk *packet.UpdateSubChunkBlocks) *UpdateSubChunkBlocks {
	return FromLatestUpdateSubChunkBlocks(pk, proto)
}

func ToLatestUpdateSubChunkBlocks(pk *UpdateSubChunkBlocks, protocolID uint32) *packet.UpdateSubChunkBlocks {
	pos := protocol.BlockPos{pk.Position[0] * 16, pk.Position[1] * 16, pk.Position[2] * 16}
	return &packet.UpdateSubChunkBlocks{
		Position: pos,
		Blocks:   TranslateBlockChangeEntries(pk.Blocks, protocolID, false),
		Extra:    TranslateBlockChangeEntries(pk.Extra, protocolID, false),
	}
}
func FromLatestUpdateSubChunkBlocks(pk *packet.UpdateSubChunkBlocks, protocolID uint32) *UpdateSubChunkBlocks {
	pos := protocol.SubChunkPos{pk.Position[0] / 16, pk.Position[1] / 16, pk.Position[2] / 16}
	return &UpdateSubChunkBlocks{
		Position: pos,
		Blocks:   TranslateBlockChangeEntries(pk.Blocks, protocolID, true),
		Extra:    TranslateBlockChangeEntries(pk.Extra, protocolID, true),
	}
}

func toLatestAnimate(pk *Animate) *packet.Animate {
	return &packet.Animate{ActionType: uint8(pk.ActionType), EntityRuntimeID: pk.EntityRuntimeID, Data: pk.BoatRowingTime}
}
func fromLatestAnimate(pk *packet.Animate) *Animate {
	out := &Animate{ActionType: int32(pk.ActionType), EntityRuntimeID: pk.EntityRuntimeID}
	if out.ActionType&0x80 != 0 {
		out.BoatRowingTime = pk.Data
	}
	return out
}

func toLatestCraftingData(pk *CraftingData) *packet.CraftingData {
	return &packet.CraftingData{
		PotionRecipes: pk.PotionRecipes, PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
		MaterialReducers: pk.MaterialReducers, ClearRecipes: pk.ClearRecipes,
	}
}
func fromLatestCraftingData(pk *packet.CraftingData) *CraftingData {
	return &CraftingData{
		PotionRecipes: pk.PotionRecipes, PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
		MaterialReducers: pk.MaterialReducers, ClearRecipes: pk.ClearRecipes,
	}
}

func toLatestInventoryTransaction(pk *InventoryTransaction) *packet.InventoryTransaction {
	out := &packet.InventoryTransaction{
		LegacyRequestID:    pk.LegacyRequestID,
		LegacySetItemSlots: pk.LegacySetItemSlots,
	}
	out.Actions = make([]protocol.InventoryAction, len(pk.Actions))
	for i, a := range pk.Actions {
		la := protocol.InventoryAction{SourceType: a.SourceType, InventorySlot: a.InventorySlot, OldItem: a.OldItem, NewItem: a.NewItem}
		switch a.SourceType {
		case protocol.InventoryActionSourceContainer, protocol.InventoryActionSourceTODO:
			la.WindowID = protocol.Option(int8(a.WindowID))
		case protocol.InventoryActionSourceWorld:
			la.SourceFlags = protocol.Option(a.SourceFlags)
		}
		out.Actions[i] = la
	}
	switch d := pk.TransactionData.(type) {
	case *useItemTx786:
		out.TransactionData = &protocol.UseItemTransactionData{
			ActionType: d.ActionType, TriggerType: d.TriggerType, BlockPosition: d.BlockPosition,
			BlockFace: d.BlockFace, HotBarSlot: d.HotBarSlot, HeldItem: d.HeldItem, Position: d.Position,
			ClickedPosition: d.ClickedPosition, BlockRuntimeID: d.BlockRuntimeID,
			ClientPrediction:    uint8(d.ClientPrediction),
			ClientCooldownState: d.ClientCooldownState,
		}
	case *useItemOnEntityTx786:
		out.TransactionData = &protocol.UseItemOnEntityTransactionData{
			TargetEntityRuntimeID: d.TargetEntityRuntimeID, ActionType: int32(d.ActionType), HotBarSlot: d.HotBarSlot,
			HeldItem: d.HeldItem, Position: d.Position, ClickedPosition: d.ClickedPosition,
		}
	case *releaseItemTx786:
		out.TransactionData = &protocol.ReleaseItemTransactionData{
			ActionType: int32(d.ActionType), HotBarSlot: d.HotBarSlot, HeldItem: d.HeldItem, HeadPosition: d.HeadPosition,
		}
	case *mismatchTx786:
		out.TransactionData = &protocol.MismatchTransactionData{}
	default:
		out.TransactionData = &protocol.NormalTransactionData{}
	}
	return out
}
func fromLatestInventoryTransaction(pk *packet.InventoryTransaction) *InventoryTransaction {
	out := &InventoryTransaction{LegacyRequestID: pk.LegacyRequestID, LegacySetItemSlots: pk.LegacySetItemSlots}
	out.Actions = make([]InventoryAction786, len(pk.Actions))
	for i, a := range pk.Actions {
		na := InventoryAction786{SourceType: a.SourceType, InventorySlot: a.InventorySlot, OldItem: a.OldItem, NewItem: a.NewItem}
		if v, ok := a.WindowID.Value(); ok {
			na.WindowID = int32(v)
		}
		if v, ok := a.SourceFlags.Value(); ok {
			na.SourceFlags = v
		}
		out.Actions[i] = na
	}

	out.TransactionData = &normalTx786{}
	return out
}
