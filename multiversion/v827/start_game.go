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

package v827

import (
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	SpawnBiomeTypeDefault = iota
	SpawnBiomeTypeUserDefined
)

const (
	ChatRestrictionLevelNone     = 0
	ChatRestrictionLevelDropped  = 1
	ChatRestrictionLevelDisabled = 2
)

const (
	EditorWorldTypeNotEditor = iota
	EditorWorldTypeProject
	EditorWorldTypeTestLevel
)

type StartGame struct {
	EntityUniqueID int64

	EntityRuntimeID uint64

	PlayerGameMode int32

	PlayerPosition mgl32.Vec3

	Pitch float32

	Yaw float32

	WorldSeed int64

	SpawnBiomeType int16

	UserDefinedBiomeName string

	Dimension int32

	Generator int32

	WorldGameMode int32

	Hardcore bool

	Difficulty int32

	WorldSpawn protocol.BlockPos

	AchievementsDisabled bool

	EditorWorldType int32

	CreatedInEditor bool

	ExportedFromEditor bool

	DayCycleLockTime int32

	EducationEditionOffer int32

	EducationFeaturesEnabled bool

	EducationProductID string

	RainLevel float32

	LightningLevel float32

	ConfirmedPlatformLockedContent bool

	MultiPlayerGame bool

	LANBroadcastEnabled bool

	XBLBroadcastMode int32

	PlatformBroadcastMode int32

	CommandsEnabled bool

	TexturePackRequired bool

	GameRules []protocol.GameRule

	Experiments []protocol.ExperimentData

	ExperimentsPreviouslyToggled bool

	BonusChestEnabled bool

	StartWithMapEnabled bool

	PlayerPermissions int32

	ServerChunkTickRadius int32

	HasLockedBehaviourPack bool

	HasLockedTexturePack bool

	FromLockedWorldTemplate bool

	MSAGamerTagsOnly bool

	FromWorldTemplate bool

	WorldTemplateSettingsLocked bool

	OnlySpawnV1Villagers bool

	PersonaDisabled bool

	CustomSkinsDisabled bool

	EmoteChatMuted bool

	BaseGameVersion string

	LimitedWorldWidth, LimitedWorldDepth int32

	NewNether bool

	EducationSharedResourceURI protocol.EducationSharedResourceURI

	ForceExperimentalGameplay protocol.Optional[bool]

	LevelID string

	WorldName string

	TemplateContentIdentity string

	Trial bool

	PlayerMovementSettings protocol.PlayerMovementSettings

	Time int64

	EnchantmentSeed int32

	Blocks []protocol.BlockEntry

	MultiPlayerCorrelationID string

	ServerAuthoritativeInventory bool

	GameVersion string

	PropertyData map[string]any

	ServerBlockStateChecksum uint64

	ClientSideGeneration bool

	WorldTemplateID uuid.UUID

	ChatRestrictionLevel uint8

	DisablePlayerInteractions bool

	ServerID string

	WorldID string

	ScenarioID string

	OwnerID string

	TickDeathSystemsEnabled bool

	UseBlockNetworkIDHashes bool

	ServerAuthoritativeSound bool
}

func (*StartGame) ID() uint32 {
	return IDStartGame
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
	protocol.FuncSlice(io, &pk.GameRules, io.GameRule)
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
	protocol.OptionalFunc(io, &pk.ForceExperimentalGameplay, io.Bool)
	io.Uint8(&pk.ChatRestrictionLevel)
	io.Bool(&pk.DisablePlayerInteractions)
	io.String(&pk.ServerID)
	io.String(&pk.WorldID)
	io.String(&pk.ScenarioID)
	io.String(&pk.OwnerID)
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

	io.Bool(&pk.TickDeathSystemsEnabled)
	io.Bool(&pk.ServerAuthoritativeSound)
}
