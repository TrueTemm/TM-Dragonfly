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

func toLatestPlayerMovementSettings786(x PlayerMovementSettings786) protocol.PlayerMovementSettings {

	return protocol.PlayerMovementSettings{
		RewindHistorySize: x.RewindHistorySize, ServerAuthoritativeBlockBreaking: x.ServerAuthoritativeBlockBreaking,
	}
}
func fromLatestPlayerMovementSettings786(x protocol.PlayerMovementSettings) PlayerMovementSettings786 {

	return PlayerMovementSettings786{
		MovementType: PlayerMovementModeServer786, RewindHistorySize: x.RewindHistorySize,
		ServerAuthoritativeBlockBreaking: x.ServerAuthoritativeBlockBreaking,
	}
}

func toLatestStartGame(pk *StartGame) *packet.StartGame {
	return &packet.StartGame{

		EducationEditionOffer:          uint32(pk.EducationEditionOffer),
		PlayerPermissions:              byte(pk.PlayerPermissions),
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
		PlayerMovementSettings:         toLatestPlayerMovementSettings786(pk.PlayerMovementSettings),
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
	}
}

func fromLatestStartGame(pk *packet.StartGame) *StartGame {
	return &StartGame{
		EducationEditionOffer: int32(pk.EducationEditionOffer),
		PlayerPermissions:     int32(pk.PlayerPermissions),
		AchievementsDisabled:  pk.AchievementsDisabled,

		BaseGameVersion:                "1.21.73",
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
		GameVersion:                    "1.21.73",
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
		PlayerMovementSettings:         fromLatestPlayerMovementSettings786(pk.PlayerMovementSettings),
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
	}
}

func ToLatestStartGame786(pk *StartGame) *packet.StartGame   { return toLatestStartGame(pk) }
func FromLatestStartGame786(pk *packet.StartGame) *StartGame { return fromLatestStartGame(pk) }
