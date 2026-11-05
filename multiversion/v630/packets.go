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

package v630

import (
	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type CraftingEvent struct {
	WindowID     byte
	CraftingType int32
	RecipeUUID   uuid.UUID
	Input        []protocol.ItemInstance
	Output       []protocol.ItemInstance
}

func (*CraftingEvent) ID() uint32 { return 53 }

func (pk *CraftingEvent) Marshal(io protocol.IO) {
	io.Uint8(&pk.WindowID)
	io.Varint32(&pk.CraftingType)
	io.UUID(&pk.RecipeUUID)
	protocol.FuncSlice(io, &pk.Input, io.ItemInstance)
	protocol.FuncSlice(io, &pk.Output, io.ItemInstance)
}

type LevelChunk struct{ v786.LevelChunk }

func (pk *LevelChunk) Marshal(io protocol.IO) {
	io.ChunkPos(&pk.Position)
	io.Varuint32(&pk.SubChunkCount)
	if pk.SubChunkCount == v786.SubChunkRequestModeLimited786 {
		io.Uint16(&pk.HighestSubChunk)
	}
	io.Bool(&pk.CacheEnabled)
	if pk.CacheEnabled {
		protocol.FuncSlice(io, &pk.BlobHashes, io.Uint64)
	}
	io.ByteSlice(&pk.RawPayload)
}

type PlayerAuthInput struct{ v729.PlayerAuthInput }

func (pk *PlayerAuthInput) Marshal(io protocol.IO) {
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Vec3(&pk.Position)
	io.Vec2(&pk.MoveVector)
	io.Float32(&pk.HeadYaw)
	io.Varuint64(&pk.InputData)
	io.Varuint32(&pk.InputMode)
	io.Varuint32(&pk.PlayMode)
	io.Varuint32(&pk.InteractionModel)
	if pk.PlayMode == packet.PlayModeReality {
		io.Vec3(&pk.GazeDirection)
	}
	io.Varuint64(&pk.Tick)
	io.Vec3(&pk.Delta)
	if pk.InputData&(1<<34) != 0 {
		io.PlayerInventoryAction(&pk.ItemInteractionData)
	}
	if pk.InputData&(1<<36) != 0 {
		v786.MarshalItemStackRequest786(io, &pk.ItemStackRequest)
	}
	if pk.InputData&(1<<35) != 0 {
		v786.PlayerBlockActions786(io, &pk.BlockActions)
	}
	io.Vec2(&pk.AnalogueMoveVector)
}

type PlayerList struct{ v786.PlayerList }

func (pk *PlayerList) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	switch pk.ActionType {
	case v786.PlayerListActionAdd:
		protocol.FuncIOSlice(io, &pk.Entries, playerListEntry)
	case v786.PlayerListActionRemove:
		protocol.FuncIOSlice(io, &pk.Entries, v786.PlayerListRemoveEntry786)
	default:
		io.UnknownEnumOption(pk.ActionType, "player list action type")
	}
	if pk.ActionType == v786.PlayerListActionAdd {
		for i := range pk.Entries {
			io.Bool(&pk.Entries[i].Skin.Trusted)
		}
	}
}

func playerListEntry(io protocol.IO, x *v786.PlayerListEntry786) {
	io.UUID(&x.UUID)
	io.Varint64(&x.EntityUniqueID)
	io.String(&x.Username)
	io.String(&x.XUID)
	io.String(&x.PlatformChatID)
	io.Int32(&x.BuildPlatform)
	protocol.Single(io, &x.Skin)
	io.Bool(&x.Teacher)
	io.Bool(&x.Host)
}
