package v2193

import (
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddActor struct{ packet.AddActor }

func (pk *AddActor) Marshal(io protocol.IO) {
	io.ActorUniqueID(&pk.EntityUniqueID)
	io.ActorRuntimeID(&pk.EntityRuntimeID)
	io.String(&pk.EntityType)
	io.Vec3(&pk.Position)
	io.Vec3(&pk.Velocity)
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Float32(&pk.HeadYaw)
	io.Float32(&pk.BodyYaw)
	protocol.Slice(io, &pk.Attributes)
	io.EntityMetadata(&pk.EntityMetadata)
	protocol.Single(io, &pk.EntityProperties)
	protocol.Slice(io, &pk.EntityLinks)
}

type AddPlayer struct{ packet.AddPlayer }

func (pk *AddPlayer) Marshal(io protocol.IO) {
	io.UUID(&pk.UUID)
	io.String(&pk.Username)
	io.ActorRuntimeID(&pk.EntityRuntimeID)
	io.String(&pk.PlatformChatID)
	io.Vec3(&pk.Position)
	io.Vec3(&pk.Velocity)
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Float32(&pk.HeadYaw)
	io.ItemInstance(&pk.HeldItem)
	io.Varint32(&pk.GameType)
	io.EntityMetadata(&pk.EntityMetadata)
	protocol.Single(io, &pk.EntityProperties)
	protocol.Single(io, &pk.AbilityData)
	protocol.Slice(io, &pk.EntityLinks)
	io.String(&pk.DeviceID)
	io.Int32(&pk.BuildPlatform)
}

type Animate struct{ packet.Animate }

func (pk *Animate) Marshal(io protocol.IO) {
	var swingSource protocol.Optional[string]
	if pk.SwingSource != 0 {
		swingSource = protocol.Option(swingSourceToString(pk.SwingSource))
	}
	io.Uint8(&pk.ActionType)
	io.ActorRuntimeID(&pk.EntityRuntimeID)
	io.Float32(&pk.Data)
	protocol.OptionalFunc(io, &swingSource, io.String)
	if val, ok := swingSource.Value(); ok {
		swingSourceFromString(io, &pk.SwingSource, val)
	}
}

type LevelChunk struct{ packet.LevelChunk }

func (pk *LevelChunk) Marshal(io protocol.IO) {
	io.ChunkPos(&pk.Position)
	io.Varint32(&pk.Dimension)
	io.Varuint32(&pk.SubChunkCount)
	if pk.SubChunkCount > 64 {
		io.InvalidValue(pk.SubChunkCount, "level chunk sub-chunk count", "must not exceed 64")
	}
	protocol.OptionalFunc(io, &pk.SubChunkLimit, io.Varint32)
	io.Bool(&pk.CacheEnabled)
	protocol.FuncSlice(io, &pk.BlobHashes, io.Uint64)
	io.ByteSlice(&pk.RawPayload)
}

type StartGame struct{ packet.StartGame }

func (pk *StartGame) Marshal(io protocol.IO) {
	io.ActorUniqueID(&pk.EntityUniqueID)
	io.ActorRuntimeID(&pk.EntityRuntimeID)
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
	io.BlockPos(&pk.WorldSpawn)
	io.Bool(&pk.AchievementsDisabled)
	io.Varint32(&pk.EditorWorldType)
	io.Bool(&pk.CreatedInEditor)
	io.Bool(&pk.ExportedFromEditor)
	io.Varint32(&pk.DayCycleLockTime)
	io.Varuint32(&pk.EducationEditionOffer)
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
	io.Uint8(&pk.PlayerPermissions)
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
	io.Varint32(&pk.ServerEditorConnectionPolicy)
	io.Bool(&pk.AllowAnonymousBlockDropsInEditorWorlds)
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

func swingSourceFromString(io protocol.IO, x *uint8, s string) {
	switch s {
	case "none":
		*x = packet.AnimateSwingSourceNone
	case "build":
		*x = packet.AnimateSwingSourceBuild
	case "mine":
		*x = packet.AnimateSwingSourceMine
	case "interact":
		*x = packet.AnimateSwingSourceInteract
	case "attack":
		*x = packet.AnimateSwingSourceAttack
	case "useitem":
		*x = packet.AnimateSwingSourceUseItem
	case "throwitem":
		*x = packet.AnimateSwingSourceThrowItem
	case "dropitem":
		*x = packet.AnimateSwingSourceDropItem
	case "event":
		*x = packet.AnimateSwingSourceEvent
	default:
		io.InvalidValue(s, "swingSource", "unknown source")
	}
}

func swingSourceToString(x uint8) string {
	switch x {
	case packet.AnimateSwingSourceNone:
		return "none"
	case packet.AnimateSwingSourceBuild:
		return "build"
	case packet.AnimateSwingSourceMine:
		return "mine"
	case packet.AnimateSwingSourceInteract:
		return "interact"
	case packet.AnimateSwingSourceAttack:
		return "attack"
	case packet.AnimateSwingSourceUseItem:
		return "useitem"
	case packet.AnimateSwingSourceThrowItem:
		return "throwitem"
	case packet.AnimateSwingSourceDropItem:
		return "dropitem"
	case packet.AnimateSwingSourceEvent:
		return "event"
	default:
		return "unknown"
	}
}
