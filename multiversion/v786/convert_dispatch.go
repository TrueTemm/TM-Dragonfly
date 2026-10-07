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
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func convertToLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *packet.ClientCacheStatus:

		_ = pk
		return []packet.Packet{&packet.ClientCacheStatus{Enabled: false}}, true
	case *ActorEvent:
		return []packet.Packet{toLatestActorEvent(pk)}, true
	case *AddActor:
		return []packet.Packet{toLatestAddActor(pk)}, true
	case *AddItemActor:
		return []packet.Packet{toLatestAddItemActor(pk)}, true
	case *AddPlayer:
		return []packet.Packet{toLatestAddPlayer(pk)}, true
	case *AddVolumeEntity:
		return []packet.Packet{toLatestAddVolumeEntity(pk)}, true
	case *RemoveVolumeEntity:
		return []packet.Packet{toLatestRemoveVolumeEntity(pk)}, true
	case *SetActorData:
		return []packet.Packet{toLatestSetActorData(pk)}, true
	case *HurtArmour:
		return []packet.Packet{toLatestHurtArmour(pk)}, true
	case *MobEffect:
		return []packet.Packet{toLatestMobEffect(pk)}, true
	case *ChangeMobProperty:
		return []packet.Packet{toLatestChangeMobProperty(pk)}, true
	case *PlayerUpdateEntityOverrides:
		return []packet.Packet{toLatestPlayerUpdateEntityOverrides(pk)}, true
	case *PlayerArmourDamage:
		return []packet.Packet{toLatestPlayerArmourDamage(pk)}, true
	case *MoveActorDelta:
		return []packet.Packet{toLatestMoveActorDelta(pk)}, true
	case *MovePlayer:
		return []packet.Packet{toLatestMovePlayer(pk)}, true
	case *PlayerAuthInput:
		return []packet.Packet{toLatestPlayerAuthInput(pk)}, true
	case *ClientMovementPredictionSync:
		return []packet.Packet{toLatestClientMovementPredictionSync(pk)}, true
	case *ClientCacheBlobStatus:
		return []packet.Packet{toLatestClientCacheBlobStatus(pk)}, true
	case *BiomeDefinitionList:
		return []packet.Packet{toLatestBiomeDefinitionList(pk)}, true
	case *LevelChunk:
		return []packet.Packet{toLatestLevelChunk(pk)}, true
	case *LevelSoundEvent:
		return []packet.Packet{ToLatestLevelSoundEvent786(proto, pk)}, true
	case *PlaySound:
		return []packet.Packet{toLatestPlaySound(pk)}, true
	case *JigsawStructureData:
		return []packet.Packet{toLatestJigsawStructureData(pk)}, true
	case *OnScreenTextureAnimation:
		return []packet.Packet{toLatestOnScreenTextureAnimation(pk)}, true
	case *ClientBoundMapItemData:
		return []packet.Packet{toLatestClientBoundMapItemData(pk)}, true
	case *Text:
		return []packet.Packet{toLatestText(pk)}, true
	case *BossEvent:
		return []packet.Packet{toLatestBossEvent(pk)}, true
	case *CameraAimAssist:
		return []packet.Packet{toLatestCameraAimAssist(pk)}, true
	case *CameraAimAssistPresets:
		return []packet.Packet{toLatestCameraAimAssistPresets(pk)}, true
	case *CameraInstruction:
		return []packet.Packet{toLatestCameraInstruction(pk)}, true
	case *CommandOutput:
		return []packet.Packet{toLatestCommandOutput(pk)}, true
	case *CommandRequest:
		return []packet.Packet{toLatestCommandRequest(pk)}, true
	case *AvailableCommands:
		return []packet.Packet{toLatestAvailableCommands(pk)}, true
	case *CreativeContent:
		return []packet.Packet{toLatestCreativeContent(pk)}, true
	case *InventoryTransaction:
		return []packet.Packet{toLatestInventoryTransaction(pk)}, true
	case *ItemStackRequest:
		return []packet.Packet{ToLatestItemStackRequest(pk)}, true
	case *PlayerSkin:
		return []packet.Packet{ToLatestPlayerSkin(pk)}, true
	case *SetScore:
		return []packet.Packet{toLatestSetScore(pk)}, true
	case *PlayerList:
		return []packet.Packet{toLatestPlayerList(pk)}, true
	case *SimpleEvent:
		return []packet.Packet{toLatestSimpleEvent(pk)}, true
	case *ShowStoreOffer:
		return []packet.Packet{toLatestShowStoreOffer(pk)}, true
	case *UpdateClientInputLocks:
		return []packet.Packet{toLatestUpdateClientInputLocks(pk)}, true
	case *UpdateClientOptions:
		return []packet.Packet{toLatestUpdateClientOptions(pk)}, true
	case *CommandBlockUpdate:
		return []packet.Packet{toLatestCommandBlockUpdate(pk)}, true
	case *StructureBlockUpdate:
		return []packet.Packet{toLatestStructureBlockUpdate(pk)}, true
	case *BookEdit:
		return []packet.Packet{toLatestBookEdit(pk)}, true
	case *LessonProgress:
		return []packet.Packet{toLatestLessonProgress(pk)}, true
	case *AnvilDamage:
		return []packet.Packet{toLatestAnvilDamage(pk)}, true
	case *ChangeDimension:
		return []packet.Packet{toLatestChangeDimension(pk)}, true
	case *SetSpawnPosition:
		return []packet.Packet{toLatestSetSpawnPosition(pk)}, true
	case *CorrectPlayerMovePrediction:
		return []packet.Packet{toLatestCorrectPlayerMovePrediction(pk)}, true
	case *Disconnect:
		return []packet.Packet{toLatestDisconnect(pk)}, true
	case *Event:
		return []packet.Packet{toLatestEvent(pk)}, true
	case *Interact:
		return []packet.Packet{toLatestInteract(pk)}, true
	case *PlayerAction:
		return []packet.Packet{toLatestPlayerAction(pk)}, true
	case *BlockActorData:
		return []packet.Packet{toLatestBlockActorData(pk)}, true
	case *ContainerOpen:
		return []packet.Packet{toLatestContainerOpen(pk)}, true
	case *ContainerClose:
		return []packet.Packet{toLatestContainerClose(pk)}, true
	case *UpdatePlayerGameType:
		return []packet.Packet{toLatestUpdatePlayerGameType(pk)}, true
	case *LecternUpdate:
		return []packet.Packet{toLatestLecternUpdate(pk)}, true
	case *OpenSign:
		return []packet.Packet{toLatestOpenSign(pk)}, true
	case *StructureTemplateDataRequest:
		return []packet.Packet{toLatestStructureTemplateDataRequest(pk)}, true
	case *InventorySlot:
		return []packet.Packet{toLatestInventorySlot(pk)}, true
	case *RequestChunkRadius:
		return []packet.Packet{toLatestRequestChunkRadius(pk)}, true
	case *RequestPermissions:
		return []packet.Packet{toLatestRequestPermissions(pk)}, true
	case *ResourcePackChunkRequest:
		return []packet.Packet{toLatestResourcePackChunkRequest(pk)}, true
	case *ResourcePackClientResponse:
		return []packet.Packet{toLatestResourcePackClientResponse(pk)}, true
	case *ResourcePackStack:
		return []packet.Packet{toLatestResourcePackStack(pk)}, true
	case *ResourcePacksInfo:
		return []packet.Packet{toLatestResourcePacksInfo(pk)}, true
	case *ServerBoundDiagnostics:
		return []packet.Packet{toLatestServerBoundDiagnostics(pk)}, true
	case *SetPlayerInventoryOptions:
		return []packet.Packet{toLatestSetPlayerInventoryOptions(pk)}, true
	case *Transfer:
		return []packet.Packet{toLatestTransfer(pk)}, true
	case *SubChunkRequest:
		return []packet.Packet{toLatestSubChunkRequest(pk)}, true
	case *UpdateSubChunkBlocks:
		return []packet.Packet{toLatestUpdateSubChunkBlocks(proto, pk)}, true
	case *UpdateBlock:
		return []packet.Packet{ToLatestUpdateBlock(pk, proto)}, true
	case *UpdateBlockSynced:
		return []packet.Packet{ToLatestUpdateBlockSynced(pk, proto)}, true
	case *Animate:
		return []packet.Packet{toLatestAnimate(pk)}, true
	case *CraftingData:
		return []packet.Packet{toLatestCraftingData(pk)}, true
	case *StartGame:
		return []packet.Packet{toLatestStartGame(pk)}, true
	case *SubChunk:
		return []packet.Packet{toLatestSubChunk(pk)}, true
	}
	return nil, false
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	return convertFromLatest(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	return convertToLatest(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *packet.ItemRegistry:
		_ = pk
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items786()}}, true
	case *packet.VoxelShapes:

		return nil, true
	case *packet.ActorEvent:
		return []packet.Packet{fromLatestActorEvent(pk)}, true
	case *packet.AddActor:
		return []packet.Packet{fromLatestAddActor(pk)}, true
	case *packet.AddItemActor:
		return []packet.Packet{fromLatestAddItemActor(pk)}, true
	case *packet.AddPlayer:
		return []packet.Packet{fromLatestAddPlayer(pk)}, true
	case *packet.AddVolumeEntity:
		return []packet.Packet{fromLatestAddVolumeEntity(pk)}, true
	case *packet.RemoveVolumeEntity:
		return []packet.Packet{fromLatestRemoveVolumeEntity(pk)}, true
	case *packet.SetActorData:
		return []packet.Packet{fromLatestSetActorData(pk)}, true
	case *packet.HurtArmour:
		return []packet.Packet{fromLatestHurtArmour(pk)}, true
	case *packet.MobEffect:
		return []packet.Packet{fromLatestMobEffect(pk)}, true
	case *packet.ChangeMobProperty:
		return []packet.Packet{fromLatestChangeMobProperty(pk)}, true
	case *packet.PlayerUpdateEntityOverrides:
		return []packet.Packet{fromLatestPlayerUpdateEntityOverrides(pk)}, true
	case *packet.PlayerArmourDamage:
		return []packet.Packet{fromLatestPlayerArmourDamage(pk)}, true
	case *packet.MoveActorDelta:
		return []packet.Packet{fromLatestMoveActorDelta(pk)}, true
	case *packet.MovePlayer:
		return []packet.Packet{fromLatestMovePlayer(pk)}, true
	case *packet.PlayerAuthInput:
		return []packet.Packet{fromLatestPlayerAuthInput(pk)}, true
	case *packet.ClientMovementPredictionSync:
		return []packet.Packet{fromLatestClientMovementPredictionSync(pk)}, true
	case *packet.BiomeDefinitionList:
		return []packet.Packet{fromLatestBiomeDefinitionList(pk)}, true
	case *packet.LevelChunk:
		return []packet.Packet{fromLatestLevelChunk(proto, pk)}, true
	case *packet.SubChunk:
		return []packet.Packet{fromLatestSubChunk(proto, pk)}, true
	case *packet.LevelSoundEvent:
		return []packet.Packet{FromLatestLevelSoundEvent786(proto, pk)}, true
	case *packet.PlaySound:
		return []packet.Packet{fromLatestPlaySound(pk)}, true
	case *packet.JigsawStructureData:
		return []packet.Packet{fromLatestJigsawStructureData(pk)}, true
	case *packet.OnScreenTextureAnimation:
		return []packet.Packet{fromLatestOnScreenTextureAnimation(pk)}, true
	case *packet.ClientBoundMapItemData:
		return []packet.Packet{fromLatestClientBoundMapItemData(pk)}, true
	case *packet.Text:
		return []packet.Packet{fromLatestText(pk)}, true
	case *packet.ModalFormRequest:
		return []packet.Packet{fromLatestModalFormRequest(proto, pk)}, true
	case *packet.BossEvent:
		return []packet.Packet{fromLatestBossEvent(pk)}, true
	case *packet.CameraAimAssist:
		return []packet.Packet{fromLatestCameraAimAssist(pk)}, true
	case *packet.CameraAimAssistPresets:
		return []packet.Packet{fromLatestCameraAimAssistPresets(pk)}, true
	case *packet.CameraInstruction:
		return []packet.Packet{fromLatestCameraInstruction(pk)}, true
	case *packet.CommandOutput:
		return []packet.Packet{fromLatestCommandOutput(pk)}, true
	case *packet.CommandRequest:
		return []packet.Packet{fromLatestCommandRequest(pk)}, true
	case *packet.AvailableCommands:
		return []packet.Packet{fromLatestAvailableCommands(pk)}, true
	case *packet.CreativeContent:
		return []packet.Packet{fromLatestCreativeContent(pk)}, true
	case *packet.InventoryTransaction:
		return []packet.Packet{fromLatestInventoryTransaction(pk)}, true
	case *packet.SetScore:
		return fromLatestSetScore(pk), true
	case *packet.PlayerList:
		return []packet.Packet{fromLatestPlayerList(pk)}, true
	case *packet.SimpleEvent:
		return []packet.Packet{fromLatestSimpleEvent(pk)}, true
	case *packet.ShowStoreOffer:
		return []packet.Packet{fromLatestShowStoreOffer(pk)}, true
	case *packet.UpdateClientInputLocks:
		return []packet.Packet{fromLatestUpdateClientInputLocks(pk)}, true
	case *packet.UpdateClientOptions:
		return []packet.Packet{fromLatestUpdateClientOptions(pk)}, true
	case *packet.CommandBlockUpdate:
		return []packet.Packet{fromLatestCommandBlockUpdate(pk)}, true
	case *packet.StructureBlockUpdate:
		return []packet.Packet{fromLatestStructureBlockUpdate(pk)}, true
	case *packet.BookEdit:
		return []packet.Packet{fromLatestBookEdit(pk)}, true
	case *packet.LessonProgress:
		return []packet.Packet{fromLatestLessonProgress(pk)}, true
	case *packet.AnvilDamage:
		return []packet.Packet{fromLatestAnvilDamage(pk)}, true
	case *packet.ChangeDimension:
		return []packet.Packet{fromLatestChangeDimension(pk)}, true
	case *packet.SetSpawnPosition:
		return []packet.Packet{fromLatestSetSpawnPosition(pk)}, true
	case *packet.CorrectPlayerMovePrediction:
		return []packet.Packet{fromLatestCorrectPlayerMovePrediction(pk)}, true
	case *packet.Disconnect:
		return []packet.Packet{fromLatestDisconnect(pk)}, true
	case *packet.Event:
		return []packet.Packet{fromLatestEvent(pk)}, true
	case *packet.Interact:
		return []packet.Packet{fromLatestInteract(pk)}, true
	case *packet.PlayerAction:
		return []packet.Packet{fromLatestPlayerAction(pk)}, true
	case *packet.BlockActorData:
		return []packet.Packet{fromLatestBlockActorData(pk)}, true
	case *packet.BlockEvent:
		return []packet.Packet{fromLatestBlockEvent(pk)}, true
	case *packet.ContainerOpen:
		return []packet.Packet{fromLatestContainerOpen(pk)}, true
	case *packet.ContainerClose:
		return []packet.Packet{fromLatestContainerClose(pk)}, true
	case *packet.UpdatePlayerGameType:
		return []packet.Packet{fromLatestUpdatePlayerGameType(pk)}, true
	case *packet.LecternUpdate:
		return []packet.Packet{fromLatestLecternUpdate(pk)}, true
	case *packet.OpenSign:
		return []packet.Packet{fromLatestOpenSign(pk)}, true
	case *packet.StructureTemplateDataRequest:
		return []packet.Packet{fromLatestStructureTemplateDataRequest(pk)}, true
	case *packet.InventorySlot:
		return []packet.Packet{fromLatestInventorySlot(pk)}, true
	case *packet.RequestChunkRadius:
		return []packet.Packet{fromLatestRequestChunkRadius(pk)}, true
	case *packet.RequestPermissions:
		return []packet.Packet{fromLatestRequestPermissions(pk)}, true
	case *packet.ResourcePackChunkRequest:
		return []packet.Packet{fromLatestResourcePackChunkRequest(pk)}, true
	case *packet.ResourcePackClientResponse:
		return []packet.Packet{fromLatestResourcePackClientResponse(pk)}, true
	case *packet.ResourcePackStack:
		return []packet.Packet{fromLatestResourcePackStack(pk)}, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{fromLatestResourcePacksInfo(pk)}, true
	case *packet.ServerBoundDiagnostics:
		return []packet.Packet{fromLatestServerBoundDiagnostics(pk)}, true
	case *packet.SetPlayerInventoryOptions:
		return []packet.Packet{fromLatestSetPlayerInventoryOptions(pk)}, true
	case *packet.Transfer:
		return []packet.Packet{fromLatestTransfer(pk)}, true
	case *packet.SubChunkRequest:
		return []packet.Packet{fromLatestSubChunkRequest(pk)}, true
	case *packet.UpdateSubChunkBlocks:
		return []packet.Packet{fromLatestUpdateSubChunkBlocks(proto, pk)}, true
	case *packet.LevelEvent:
		return []packet.Packet{FromLatestLevelEvent(pk, proto)}, true
	case *packet.ItemStackResponse:
		return []packet.Packet{FromLatestItemStackResponse(pk)}, true
	case *packet.PlayerSkin:
		return []packet.Packet{FromLatestPlayerSkin(pk)}, true
	case *packet.UpdateBlock:
		return []packet.Packet{FromLatestUpdateBlock(pk, proto)}, true
	case *packet.UpdateBlockSynced:
		return []packet.Packet{FromLatestUpdateBlockSynced(pk, proto)}, true
	case *packet.Animate:
		return []packet.Packet{fromLatestAnimate(pk)}, true
	case *packet.CraftingData:
		return []packet.Packet{fromLatestCraftingData(pk)}, true
	case *packet.StartGame:
		return []packet.Packet{fromLatestStartGame(pk)}, true
	}
	return nil, false
}
