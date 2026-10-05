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

package v818

import (
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v800"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func toLatestResourcePackClientResponse818(pk *v786.ResourcePackClientResponse) *packet.ResourcePackClientResponse {
	return &packet.ResourcePackClientResponse{Response: uint32(pk.Response) - 1, PacksToDownload: pk.PacksToDownload}
}
func fromLatestResourcePackClientResponse818(pk *packet.ResourcePackClientResponse) *v786.ResourcePackClientResponse {
	return &v786.ResourcePackClientResponse{Response: byte(pk.Response) + 1, PacksToDownload: pk.PacksToDownload}
}

func toLatestResourcePacksInfo(pk *ResourcePacksInfo) *packet.ResourcePacksInfo {
	return &packet.ResourcePacksInfo{
		TexturePackRequired: pk.TexturePackRequired, HasAddons: pk.HasAddons, HasScripts: pk.HasScripts,
		ForceDisableVibrantVisuals: pk.ForceDisableVibrantVisuals, WorldTemplateUUID: pk.WorldTemplateUUID,
		WorldTemplateVersion: pk.WorldTemplateVersion, TexturePacks: pk.TexturePacks,
	}
}
func fromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	return &ResourcePacksInfo{
		TexturePackRequired: pk.TexturePackRequired, HasAddons: pk.HasAddons, HasScripts: pk.HasScripts,
		ForceDisableVibrantVisuals: pk.ForceDisableVibrantVisuals, WorldTemplateUUID: pk.WorldTemplateUUID,
		WorldTemplateVersion: pk.WorldTemplateVersion, TexturePacks: pk.TexturePacks,
	}
}

func FromLatestResourcePacksInfo818(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	return fromLatestResourcePacksInfo(pk)
}
func ToLatestResourcePacksInfo818(pk *ResourcePacksInfo) *packet.ResourcePacksInfo {
	return toLatestResourcePacksInfo(pk)
}

func toLatestStartGame(pk *StartGame) *packet.StartGame {
	return &packet.StartGame{
		AchievementsDisabled:           pk.AchievementsDisabled,
		BaseGameVersion:                pk.BaseGameVersion,
		Blocks:                         pk.Blocks,
		BonusChestEnabled:              pk.BonusChestEnabled,
		ChatRestrictionLevel:           pk.ChatRestrictionLevel,
		ClientSideGeneration:           pk.ClientSideGeneration,
		CommandsEnabled:                pk.CommandsEnabled,
		ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
		CreatedInEditor:                pk.CreatedInEditor,
		CustomSkinsDisabled:            pk.CustomSkinsDisabled,
		DayCycleLockTime:               pk.DayCycleLockTime,
		Difficulty:                     pk.Difficulty,
		Dimension:                      pk.Dimension,
		DisablePlayerInteractions:      pk.DisablePlayerInteractions,
		EditorWorldType:                pk.EditorWorldType,
		EducationEditionOffer:          uint32(pk.EducationEditionOffer),
		EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
		EducationProductID:             pk.EducationProductID,
		EducationSharedResourceURI:     pk.EducationSharedResourceURI,
		EmoteChatMuted:                 pk.EmoteChatMuted,
		EnchantmentSeed:                pk.EnchantmentSeed,
		EntityRuntimeID:                pk.EntityRuntimeID,
		EntityUniqueID:                 pk.EntityUniqueID,
		Experiments:                    pk.Experiments,
		ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
		ExportedFromEditor:             pk.ExportedFromEditor,
		ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
		FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
		FromWorldTemplate:              pk.FromWorldTemplate,
		GameRules:                      pk.GameRules,
		GameVersion:                    pk.GameVersion,
		Generator:                      pk.Generator,
		Hardcore:                       pk.Hardcore,
		HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
		HasLockedTexturePack:           pk.HasLockedTexturePack,
		LANBroadcastEnabled:            pk.LANBroadcastEnabled,
		LevelID:                        pk.LevelID,
		LightningLevel:                 pk.LightningLevel,
		LimitedWorldDepth:              pk.LimitedWorldDepth,
		LimitedWorldWidth:              pk.LimitedWorldWidth,
		MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
		MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
		MultiPlayerGame:                pk.MultiPlayerGame,
		NewNether:                      pk.NewNether,
		OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
		PersonaDisabled:                pk.PersonaDisabled,
		Pitch:                          pk.Pitch,
		PlatformBroadcastMode:          pk.PlatformBroadcastMode,
		PlayerGameMode:                 pk.PlayerGameMode,
		PlayerMovementSettings:         pk.PlayerMovementSettings,
		PlayerPermissions:              byte(pk.PlayerPermissions),
		PlayerPosition:                 pk.PlayerPosition,
		PropertyData:                   pk.PropertyData,
		RainLevel:                      pk.RainLevel,
		ScenarioID:                     pk.ScenarioID,
		ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
		ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
		ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
		ServerChunkTickRadius:          pk.ServerChunkTickRadius,
		ServerID:                       pk.ServerID,
		SpawnBiomeType:                 pk.SpawnBiomeType,
		StartWithMapEnabled:            pk.StartWithMapEnabled,
		TemplateContentIdentity:        pk.TemplateContentIdentity,
		TexturePackRequired:            pk.TexturePackRequired,
		Time:                           pk.Time,
		Trial:                          pk.Trial,
		UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
		UserDefinedBiomeName:           pk.UserDefinedBiomeName,
		WorldGameMode:                  pk.WorldGameMode,
		WorldID:                        pk.WorldID,
		WorldName:                      pk.WorldName,
		WorldSeed:                      pk.WorldSeed,
		WorldSpawn:                     pk.WorldSpawn,
		WorldTemplateID:                pk.WorldTemplateID,
		WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
		XBLBroadcastMode:               pk.XBLBroadcastMode,
		Yaw:                            pk.Yaw,
		OwnerID:                        pk.OwnerID,
	}
}

func fromLatestStartGame(pk *packet.StartGame) *StartGame {
	return &StartGame{
		AchievementsDisabled:           pk.AchievementsDisabled,
		BaseGameVersion:                "1.21.90",
		Blocks:                         pk.Blocks,
		BonusChestEnabled:              pk.BonusChestEnabled,
		ChatRestrictionLevel:           pk.ChatRestrictionLevel,
		ClientSideGeneration:           pk.ClientSideGeneration,
		CommandsEnabled:                pk.CommandsEnabled,
		ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
		CreatedInEditor:                pk.CreatedInEditor,
		CustomSkinsDisabled:            pk.CustomSkinsDisabled,
		DayCycleLockTime:               pk.DayCycleLockTime,
		Difficulty:                     pk.Difficulty,
		Dimension:                      pk.Dimension,
		DisablePlayerInteractions:      pk.DisablePlayerInteractions,
		EditorWorldType:                pk.EditorWorldType,
		EducationEditionOffer:          int32(pk.EducationEditionOffer),
		EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
		EducationProductID:             pk.EducationProductID,
		EducationSharedResourceURI:     pk.EducationSharedResourceURI,
		EmoteChatMuted:                 pk.EmoteChatMuted,
		EnchantmentSeed:                pk.EnchantmentSeed,
		EntityRuntimeID:                pk.EntityRuntimeID,
		EntityUniqueID:                 pk.EntityUniqueID,
		Experiments:                    pk.Experiments,
		ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
		ExportedFromEditor:             pk.ExportedFromEditor,
		ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
		FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
		FromWorldTemplate:              pk.FromWorldTemplate,
		GameRules:                      pk.GameRules,
		GameVersion:                    "1.21.90",
		Generator:                      pk.Generator,
		Hardcore:                       pk.Hardcore,
		HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
		HasLockedTexturePack:           pk.HasLockedTexturePack,
		LANBroadcastEnabled:            pk.LANBroadcastEnabled,
		LevelID:                        pk.LevelID,
		LightningLevel:                 pk.LightningLevel,
		LimitedWorldDepth:              pk.LimitedWorldDepth,
		LimitedWorldWidth:              pk.LimitedWorldWidth,
		MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
		MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
		MultiPlayerGame:                pk.MultiPlayerGame,
		NewNether:                      pk.NewNether,
		OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
		PersonaDisabled:                pk.PersonaDisabled,
		Pitch:                          pk.Pitch,
		PlatformBroadcastMode:          pk.PlatformBroadcastMode,
		PlayerGameMode:                 pk.PlayerGameMode,
		PlayerMovementSettings:         pk.PlayerMovementSettings,
		PlayerPermissions:              int32(pk.PlayerPermissions),
		PlayerPosition:                 pk.PlayerPosition,
		PropertyData:                   pk.PropertyData,
		RainLevel:                      pk.RainLevel,
		ScenarioID:                     pk.ScenarioID,
		ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
		ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
		ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
		ServerChunkTickRadius:          pk.ServerChunkTickRadius,
		ServerID:                       pk.ServerID,
		SpawnBiomeType:                 pk.SpawnBiomeType,
		StartWithMapEnabled:            pk.StartWithMapEnabled,
		TemplateContentIdentity:        pk.TemplateContentIdentity,
		TexturePackRequired:            pk.TexturePackRequired,
		Time:                           pk.Time,
		Trial:                          pk.Trial,
		UseBlockNetworkIDHashes:        true,
		UserDefinedBiomeName:           pk.UserDefinedBiomeName,
		WorldGameMode:                  pk.WorldGameMode,
		WorldID:                        pk.WorldID,
		WorldName:                      pk.WorldName,
		WorldSeed:                      pk.WorldSeed,
		WorldSpawn:                     pk.WorldSpawn,
		WorldTemplateID:                pk.WorldTemplateID,
		WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
		XBLBroadcastMode:               pk.XBLBroadcastMode,
		Yaw:                            pk.Yaw,
		OwnerID:                        pk.OwnerID,
	}
}

func convertToLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *packet.ClientCacheStatus:

		_ = pk
		return []packet.Packet{&packet.ClientCacheStatus{Enabled: false}}, true
	case *v800.ClientMovementPredictionSync:
		return []packet.Packet{v800.ToLatestClientMovementPredictionSync(pk)}, true
	case *v800.PlayerLocation:
		return []packet.Packet{v800.ToLatestPlayerLocation(pk)}, true
	case *ResourcePacksInfo:
		return []packet.Packet{toLatestResourcePacksInfo(pk)}, true
	case *StartGame:
		return []packet.Packet{toLatestStartGame(pk)}, true
	case *v786.ResourcePackClientResponse:
		return []packet.Packet{toLatestResourcePackClientResponse818(pk)}, true
	case *v786.ResourcePackStack:
		return []packet.Packet{v786.ToLatestResourcePackStack786(pk)}, true
	case *v786.ClientCacheBlobStatus:
		return []packet.Packet{toLatestClientCacheBlobStatus(pk)}, true
	case *v786.CommandRequest:
		return []packet.Packet{v786.ToLatestCommandRequest786(pk)}, true
	case *v786.CommandOutput:
		return []packet.Packet{v786.ToLatestCommandOutput786(pk)}, true
	case *v786.AvailableCommands:
		return []packet.Packet{v786.ToLatestAvailableCommands786(pk)}, true
	case *v786.CreativeContent:
		return []packet.Packet{v786.ToLatestCreativeContent786(pk)}, true
	case *v786.LevelSoundEvent:
		return []packet.Packet{v786.ToLatestLevelSoundEvent786(proto, pk)}, true
	case *v786.ItemStackRequest:
		return []packet.Packet{v786.ToLatestItemStackRequest(pk)}, true
	case *v786.PlayerSkin:
		return []packet.Packet{v786.ToLatestPlayerSkin(pk)}, true
	case *v786.UpdateBlock:
		return []packet.Packet{v786.ToLatestUpdateBlock(pk, proto)}, true
	case *v786.UpdateBlockSynced:
		return []packet.Packet{v786.ToLatestUpdateBlockSynced(pk, proto)}, true
	case *v786.UpdateSubChunkBlocks:
		return []packet.Packet{v786.ToLatestUpdateSubChunkBlocks(pk, proto)}, true
	}
	return nil, false
}

func toLatestClientCacheBlobStatus(pk *v786.ClientCacheBlobStatus) *packet.ClientCacheBlobStatus {
	return &packet.ClientCacheBlobStatus{MissHashes: pk.MissHashes, HitHashes: pk.HitHashes}
}

func convertFromLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {

	case *packet.CameraAimAssistPresets:
		return []packet.Packet{pk}, true
	case *packet.PlayerList:

		return []packet.Packet{v800.FromLatestPlayerList800(pk)}, true
	case *packet.BiomeDefinitionList:

		return []packet.Packet{v800.FromLatestBiomeDefinitionList(pk)}, true
	case *packet.ClientMovementPredictionSync:
		return []packet.Packet{v800.FromLatestClientMovementPredictionSync(pk)}, true
	case *packet.PlayerLocation:
		return []packet.Packet{v800.FromLatestPlayerLocation(pk)}, true

	case *packet.ItemRegistry:
		_ = pk

		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items818()}}, true
	case *packet.VoxelShapes:

		return nil, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{fromLatestResourcePacksInfo(pk)}, true
	case *packet.StartGame:
		return []packet.Packet{fromLatestStartGame(pk)}, true
	case *packet.ResourcePackClientResponse:
		return []packet.Packet{fromLatestResourcePackClientResponse818(pk)}, true
	case *packet.ResourcePackStack:
		return []packet.Packet{v786.FromLatestResourcePackStack786(pk)}, true
	case *packet.CommandRequest:
		return []packet.Packet{v786.FromLatestCommandRequest786(pk)}, true
	case *packet.CommandOutput:
		return []packet.Packet{v786.FromLatestCommandOutput786(pk)}, true
	case *packet.AvailableCommands:
		return []packet.Packet{v786.FromLatestAvailableCommands786(pk)}, true
	case *packet.CreativeContent:
		return []packet.Packet{v786.FromLatestCreativeContent786(pk)}, true
	case *packet.LevelSoundEvent:
		return []packet.Packet{v786.FromLatestLevelSoundEvent786(proto, pk)}, true
	case *packet.LevelEvent:
		return []packet.Packet{v786.FromLatestLevelEvent(pk, proto)}, true
	case *packet.ItemStackResponse:
		return []packet.Packet{v786.FromLatestItemStackResponse(pk)}, true
	case *packet.PlayerSkin:
		return []packet.Packet{v786.FromLatestPlayerSkin(pk)}, true
	case *packet.UpdateBlock:
		return []packet.Packet{v786.FromLatestUpdateBlock(pk, proto)}, true
	case *packet.UpdateBlockSynced:
		return []packet.Packet{v786.FromLatestUpdateBlockSynced(pk, proto)}, true
	case *packet.UpdateSubChunkBlocks:
		return []packet.Packet{v786.FromLatestUpdateSubChunkBlocks(pk, proto)}, true
	}
	return nil, false
}
