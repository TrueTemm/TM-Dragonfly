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

package v924

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v844"
	"github.com/df-mc/dragonfly/multiversion/v898"
)

type StartGame struct {
	EntityUniqueID                       int64
	EntityRuntimeID                      uint64
	PlayerGameMode                       int32
	PlayerPosition                       mgl32.Vec3
	Pitch                                float32
	Yaw                                  float32
	WorldSeed                            int64
	SpawnBiomeType                       int16
	UserDefinedBiomeName                 string
	Dimension                            int32
	Generator                            int32
	WorldGameMode                        int32
	Hardcore                             bool
	Difficulty                           int32
	WorldSpawn                           protocol.BlockPos
	AchievementsDisabled                 bool
	EditorWorldType                      int32
	CreatedInEditor                      bool
	ExportedFromEditor                   bool
	DayCycleLockTime                     int32
	EducationEditionOffer                int32
	EducationFeaturesEnabled             bool
	EducationProductID                   string
	RainLevel                            float32
	LightningLevel                       float32
	ConfirmedPlatformLockedContent       bool
	MultiPlayerGame                      bool
	LANBroadcastEnabled                  bool
	XBLBroadcastMode                     int32
	PlatformBroadcastMode                int32
	CommandsEnabled                      bool
	TexturePackRequired                  bool
	GameRules                            []protocol.GameRule
	Experiments                          []protocol.ExperimentData
	ExperimentsPreviouslyToggled         bool
	BonusChestEnabled                    bool
	StartWithMapEnabled                  bool
	PlayerPermissions                    int32
	ServerChunkTickRadius                int32
	HasLockedBehaviourPack               bool
	HasLockedTexturePack                 bool
	FromLockedWorldTemplate              bool
	MSAGamerTagsOnly                     bool
	FromWorldTemplate                    bool
	WorldTemplateSettingsLocked          bool
	OnlySpawnV1Villagers                 bool
	PersonaDisabled                      bool
	CustomSkinsDisabled                  bool
	EmoteChatMuted                       bool
	BaseGameVersion                      string
	LimitedWorldWidth, LimitedWorldDepth int32
	NewNether                            bool
	EducationSharedResourceURI           protocol.EducationSharedResourceURI
	ForceExperimentalGameplay            bool
	LevelID                              string
	WorldName                            string
	TemplateContentIdentity              string
	Trial                                bool
	PlayerMovementSettings               protocol.PlayerMovementSettings
	Time                                 int64
	EnchantmentSeed                      int32
	Blocks                               []protocol.BlockEntry
	MultiPlayerCorrelationID             string
	ServerAuthoritativeInventory         bool
	GameVersion                          string
	PropertyData                         map[string]any
	ServerBlockStateChecksum             uint64
	ClientSideGeneration                 bool
	WorldTemplateID                      uuid.UUID
	ChatRestrictionLevel                 uint8
	DisablePlayerInteractions            bool
	ServerID                             string
	WorldID                              string
	ScenarioID                           string
	OwnerID                              string
	UseBlockNetworkIDHashes              bool

	ServerJoinInformation    protocol.Optional[protocol.ServerJoinInformation]
	ServerAuthoritativeSound bool
}

func (*StartGame) ID() uint32 {
	return v898.IDStartGame
}

func (pk *StartGame) Marshal(io protocol.IO) {
	io.Varint64(&pk.EntityUniqueID)
	io.Varuint64(&pk.EntityRuntimeID)
	io.Varint32(&pk.PlayerGameMode)
	io.Vec3(&pk.PlayerPosition)
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Int64(&pk.WorldSeed)
	io.Int16(&pk.SpawnBiomeType)
	io.String(&pk.UserDefinedBiomeName)
	io.Varint32(&pk.Dimension)
	io.Varint32(&pk.Generator)
	io.Varint32(&pk.WorldGameMode)
	io.Bool(&pk.Hardcore)
	io.Varint32(&pk.Difficulty)
	v786.UBlockPos786(io, &pk.WorldSpawn)
	io.Bool(&pk.AchievementsDisabled)
	io.Varint32(&pk.EditorWorldType)
	io.Bool(&pk.CreatedInEditor)
	io.Bool(&pk.ExportedFromEditor)
	io.Varint32(&pk.DayCycleLockTime)
	io.Varint32(&pk.EducationEditionOffer)
	io.Bool(&pk.EducationFeaturesEnabled)
	io.String(&pk.EducationProductID)
	io.Float32(&pk.RainLevel)
	io.Float32(&pk.LightningLevel)
	io.Bool(&pk.ConfirmedPlatformLockedContent)
	io.Bool(&pk.MultiPlayerGame)
	io.Bool(&pk.LANBroadcastEnabled)
	io.Varint32(&pk.XBLBroadcastMode)
	io.Varint32(&pk.PlatformBroadcastMode)
	io.Bool(&pk.CommandsEnabled)
	io.Bool(&pk.TexturePackRequired)
	protocol.FuncSlice(io, &pk.GameRules, func(x *protocol.GameRule) { v844.GameRuleLegacy(io, x) })
	protocol.SliceUint32Length(io, &pk.Experiments)
	io.Bool(&pk.ExperimentsPreviouslyToggled)
	io.Bool(&pk.BonusChestEnabled)
	io.Bool(&pk.StartWithMapEnabled)
	io.Varint32(&pk.PlayerPermissions)
	io.Int32(&pk.ServerChunkTickRadius)
	io.Bool(&pk.HasLockedBehaviourPack)
	io.Bool(&pk.HasLockedTexturePack)
	io.Bool(&pk.FromLockedWorldTemplate)
	io.Bool(&pk.MSAGamerTagsOnly)
	io.Bool(&pk.FromWorldTemplate)
	io.Bool(&pk.WorldTemplateSettingsLocked)
	io.Bool(&pk.OnlySpawnV1Villagers)
	io.Bool(&pk.PersonaDisabled)
	io.Bool(&pk.CustomSkinsDisabled)
	io.Bool(&pk.EmoteChatMuted)
	io.String(&pk.BaseGameVersion)
	io.Int32(&pk.LimitedWorldWidth)
	io.Int32(&pk.LimitedWorldDepth)
	io.Bool(&pk.NewNether)
	protocol.Single(io, &pk.EducationSharedResourceURI)
	io.Bool(&pk.ForceExperimentalGameplay)
	io.Uint8(&pk.ChatRestrictionLevel)
	io.Bool(&pk.DisablePlayerInteractions)
	io.String(&pk.LevelID)
	io.String(&pk.WorldName)
	io.String(&pk.TemplateContentIdentity)
	io.Bool(&pk.Trial)
	protocol.PlayerMoveSettings(io, &pk.PlayerMovementSettings)
	io.Int64(&pk.Time)
	io.Varint32(&pk.EnchantmentSeed)
	protocol.Slice(io, &pk.Blocks)
	io.String(&pk.MultiPlayerCorrelationID)
	io.Bool(&pk.ServerAuthoritativeInventory)
	io.String(&pk.GameVersion)
	io.NBT(&pk.PropertyData, nbt.NetworkLittleEndian)
	io.Uint64(&pk.ServerBlockStateChecksum)
	io.UUID(&pk.WorldTemplateID)
	io.Bool(&pk.ClientSideGeneration)
	io.Bool(&pk.UseBlockNetworkIDHashes)
	io.Bool(&pk.ServerAuthoritativeSound)

	protocol.OptionalMarshaler(io, &pk.ServerJoinInformation)
	io.String(&pk.ServerID)
	io.String(&pk.ScenarioID)
	io.String(&pk.WorldID)
	io.String(&pk.OwnerID)
}

func ToLatestStartGame924(pk *StartGame) *packet.StartGame {
	return &packet.StartGame{
		EntityUniqueID: pk.EntityUniqueID, EntityRuntimeID: pk.EntityRuntimeID, PlayerGameMode: pk.PlayerGameMode,
		PlayerPosition: pk.PlayerPosition, Pitch: pk.Pitch, Yaw: pk.Yaw, WorldSeed: pk.WorldSeed,
		SpawnBiomeType: pk.SpawnBiomeType, UserDefinedBiomeName: pk.UserDefinedBiomeName, Dimension: pk.Dimension,
		Generator: pk.Generator, WorldGameMode: pk.WorldGameMode, Hardcore: pk.Hardcore, Difficulty: pk.Difficulty,
		WorldSpawn: pk.WorldSpawn, AchievementsDisabled: pk.AchievementsDisabled, EditorWorldType: pk.EditorWorldType,
		CreatedInEditor: pk.CreatedInEditor, ExportedFromEditor: pk.ExportedFromEditor, DayCycleLockTime: pk.DayCycleLockTime,
		EducationEditionOffer: uint32(pk.EducationEditionOffer), EducationFeaturesEnabled: pk.EducationFeaturesEnabled,
		EducationProductID: pk.EducationProductID, RainLevel: pk.RainLevel, LightningLevel: pk.LightningLevel,
		ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent, MultiPlayerGame: pk.MultiPlayerGame,
		LANBroadcastEnabled: pk.LANBroadcastEnabled, XBLBroadcastMode: pk.XBLBroadcastMode,
		PlatformBroadcastMode: pk.PlatformBroadcastMode, CommandsEnabled: pk.CommandsEnabled,
		TexturePackRequired: pk.TexturePackRequired, GameRules: pk.GameRules, Experiments: pk.Experiments,
		ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled, BonusChestEnabled: pk.BonusChestEnabled,
		StartWithMapEnabled: pk.StartWithMapEnabled, PlayerPermissions: byte(pk.PlayerPermissions),
		ServerChunkTickRadius: pk.ServerChunkTickRadius, HasLockedBehaviourPack: pk.HasLockedBehaviourPack,
		HasLockedTexturePack: pk.HasLockedTexturePack, FromLockedWorldTemplate: pk.FromLockedWorldTemplate,
		MSAGamerTagsOnly: pk.MSAGamerTagsOnly, FromWorldTemplate: pk.FromWorldTemplate,
		WorldTemplateSettingsLocked: pk.WorldTemplateSettingsLocked, OnlySpawnV1Villagers: pk.OnlySpawnV1Villagers,
		PersonaDisabled: pk.PersonaDisabled, CustomSkinsDisabled: pk.CustomSkinsDisabled, EmoteChatMuted: pk.EmoteChatMuted,
		BaseGameVersion: pk.BaseGameVersion, LimitedWorldWidth: pk.LimitedWorldWidth, LimitedWorldDepth: pk.LimitedWorldDepth,
		NewNether: pk.NewNether, EducationSharedResourceURI: pk.EducationSharedResourceURI,
		ForceExperimentalGameplay: protocol.Option(pk.ForceExperimentalGameplay), LevelID: pk.LevelID, WorldName: pk.WorldName,
		TemplateContentIdentity: pk.TemplateContentIdentity, Trial: pk.Trial,
		PlayerMovementSettings: pk.PlayerMovementSettings, Time: pk.Time, EnchantmentSeed: pk.EnchantmentSeed,
		Blocks: pk.Blocks, MultiPlayerCorrelationID: pk.MultiPlayerCorrelationID,
		ServerAuthoritativeInventory: pk.ServerAuthoritativeInventory, GameVersion: pk.GameVersion,
		PropertyData: pk.PropertyData, ServerBlockStateChecksum: pk.ServerBlockStateChecksum,
		ClientSideGeneration: pk.ClientSideGeneration, WorldTemplateID: pk.WorldTemplateID,
		ChatRestrictionLevel: pk.ChatRestrictionLevel, DisablePlayerInteractions: pk.DisablePlayerInteractions,
		ServerID: pk.ServerID, WorldID: pk.WorldID, ScenarioID: pk.ScenarioID, OwnerID: pk.OwnerID,
		UseBlockNetworkIDHashes:  pk.UseBlockNetworkIDHashes,
		ServerAuthoritativeSound: pk.ServerAuthoritativeSound,
	}
}

func FromLatestStartGame924(pk *packet.StartGame) *StartGame {
	forceExperimental, _ := pk.ForceExperimentalGameplay.Value()
	return &StartGame{
		EntityUniqueID: pk.EntityUniqueID, EntityRuntimeID: pk.EntityRuntimeID, PlayerGameMode: pk.PlayerGameMode,
		PlayerPosition: pk.PlayerPosition, Pitch: pk.Pitch, Yaw: pk.Yaw, WorldSeed: pk.WorldSeed,
		SpawnBiomeType: pk.SpawnBiomeType, UserDefinedBiomeName: pk.UserDefinedBiomeName, Dimension: pk.Dimension,
		Generator: pk.Generator, WorldGameMode: pk.WorldGameMode, Hardcore: pk.Hardcore, Difficulty: pk.Difficulty,
		WorldSpawn: pk.WorldSpawn, AchievementsDisabled: pk.AchievementsDisabled, EditorWorldType: pk.EditorWorldType,
		CreatedInEditor: pk.CreatedInEditor, ExportedFromEditor: pk.ExportedFromEditor, DayCycleLockTime: pk.DayCycleLockTime,
		EducationEditionOffer: int32(pk.EducationEditionOffer), EducationFeaturesEnabled: pk.EducationFeaturesEnabled,
		EducationProductID: pk.EducationProductID, RainLevel: pk.RainLevel, LightningLevel: pk.LightningLevel,
		ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent, MultiPlayerGame: pk.MultiPlayerGame,
		LANBroadcastEnabled: pk.LANBroadcastEnabled, XBLBroadcastMode: pk.XBLBroadcastMode,
		PlatformBroadcastMode: pk.PlatformBroadcastMode, CommandsEnabled: pk.CommandsEnabled,
		TexturePackRequired: pk.TexturePackRequired, GameRules: pk.GameRules, Experiments: pk.Experiments,
		ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled, BonusChestEnabled: pk.BonusChestEnabled,
		StartWithMapEnabled: pk.StartWithMapEnabled, PlayerPermissions: int32(pk.PlayerPermissions),
		ServerChunkTickRadius: pk.ServerChunkTickRadius, HasLockedBehaviourPack: pk.HasLockedBehaviourPack,
		HasLockedTexturePack: pk.HasLockedTexturePack, FromLockedWorldTemplate: pk.FromLockedWorldTemplate,
		MSAGamerTagsOnly: pk.MSAGamerTagsOnly, FromWorldTemplate: pk.FromWorldTemplate,
		WorldTemplateSettingsLocked: pk.WorldTemplateSettingsLocked, OnlySpawnV1Villagers: pk.OnlySpawnV1Villagers,
		PersonaDisabled: pk.PersonaDisabled, CustomSkinsDisabled: pk.CustomSkinsDisabled, EmoteChatMuted: pk.EmoteChatMuted,
		BaseGameVersion: "1.26.0", LimitedWorldWidth: pk.LimitedWorldWidth, LimitedWorldDepth: pk.LimitedWorldDepth,
		NewNether: pk.NewNether, EducationSharedResourceURI: pk.EducationSharedResourceURI,
		ForceExperimentalGameplay: forceExperimental, LevelID: pk.LevelID, WorldName: pk.WorldName,
		TemplateContentIdentity: pk.TemplateContentIdentity, Trial: pk.Trial,
		PlayerMovementSettings: pk.PlayerMovementSettings, Time: pk.Time, EnchantmentSeed: pk.EnchantmentSeed,
		Blocks: pk.Blocks, MultiPlayerCorrelationID: pk.MultiPlayerCorrelationID,
		ServerAuthoritativeInventory: pk.ServerAuthoritativeInventory, GameVersion: "1.26.0",
		PropertyData: pk.PropertyData, ServerBlockStateChecksum: pk.ServerBlockStateChecksum,
		ClientSideGeneration: pk.ClientSideGeneration, WorldTemplateID: pk.WorldTemplateID,
		ChatRestrictionLevel: pk.ChatRestrictionLevel, DisablePlayerInteractions: pk.DisablePlayerInteractions,
		ServerID: pk.ServerID, WorldID: pk.WorldID, ScenarioID: pk.ScenarioID, OwnerID: pk.OwnerID,

		UseBlockNetworkIDHashes:  true,
		ServerAuthoritativeSound: pk.ServerAuthoritativeSound,
	}
}
