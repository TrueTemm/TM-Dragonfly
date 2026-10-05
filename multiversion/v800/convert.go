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

package v800

import (
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func toLatestResourcePackClientResponse800(pk *v786.ResourcePackClientResponse) *packet.ResourcePackClientResponse {
	return &packet.ResourcePackClientResponse{Response: uint32(pk.Response) - 1, PacksToDownload: pk.PacksToDownload}
}
func fromLatestResourcePackClientResponse800(pk *packet.ResourcePackClientResponse) *v786.ResourcePackClientResponse {
	return &v786.ResourcePackClientResponse{Response: byte(pk.Response) + 1, PacksToDownload: pk.PacksToDownload}
}

func toLatestClientMovementPredictionSync(pk *ClientMovementPredictionSync) *packet.ClientMovementPredictionSync {
	return &packet.ClientMovementPredictionSync{
		ActorFlags: pk.ActorFlags, BoundingBoxScale: pk.BoundingBoxScale, BoundingBoxWidth: pk.BoundingBoxWidth,
		BoundingBoxHeight: pk.BoundingBoxHeight, MovementSpeed: pk.MovementSpeed,
		UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed, LavaMovementSpeed: pk.LavaMovementSpeed,
		JumpStrength: pk.JumpStrength, Health: pk.Health, Hunger: pk.Hunger, EntityUniqueID: pk.EntityUniqueID,
		Flying: pk.Flying,
	}
}
func fromLatestClientMovementPredictionSync(pk *packet.ClientMovementPredictionSync) *ClientMovementPredictionSync {
	return &ClientMovementPredictionSync{
		ActorFlags: pk.ActorFlags, BoundingBoxScale: pk.BoundingBoxScale, BoundingBoxWidth: pk.BoundingBoxWidth,
		BoundingBoxHeight: pk.BoundingBoxHeight, MovementSpeed: pk.MovementSpeed,
		UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed, LavaMovementSpeed: pk.LavaMovementSpeed,
		JumpStrength: pk.JumpStrength, Health: pk.Health, Hunger: pk.Hunger, EntityUniqueID: pk.EntityUniqueID,
		Flying: pk.Flying,
	}
}

func toLatestPlayerLocation(pk *PlayerLocation) *packet.PlayerLocation {
	return &packet.PlayerLocation{Type: pk.Type, EntityUniqueID: pk.EntityUniqueID, Position: pk.Position}
}
func fromLatestPlayerLocation(pk *packet.PlayerLocation) *PlayerLocation {
	return &PlayerLocation{Type: pk.Type, EntityUniqueID: pk.EntityUniqueID, Position: pk.Position}
}

func toLatestStartGame800(pk *v786.StartGame) *packet.StartGame {
	return v786.ToLatestStartGame786(pk)
}
func fromLatestStartGame800(pk *packet.StartGame) *v786.StartGame {
	sg := v786.FromLatestStartGame786(pk)
	sg.BaseGameVersion = "1.21.80"
	sg.GameVersion = "1.21.80"
	return sg
}

func convertToLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *packet.ClientCacheStatus:

		_ = pk
		return []packet.Packet{&packet.ClientCacheStatus{Enabled: false}}, true
	case *ClientMovementPredictionSync:
		return []packet.Packet{toLatestClientMovementPredictionSync(pk)}, true
	case *PlayerLocation:
		return []packet.Packet{toLatestPlayerLocation(pk)}, true
	case *v786.ResourcePackClientResponse:
		return []packet.Packet{toLatestResourcePackClientResponse800(pk)}, true
	case *v786.ResourcePacksInfo:
		return []packet.Packet{v786.ToLatestResourcePacksInfo786(pk)}, true
	case *v786.ResourcePackStack:
		return []packet.Packet{v786.ToLatestResourcePackStack786(pk)}, true
	case *v786.StartGame:
		return []packet.Packet{toLatestStartGame800(pk)}, true
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

		return []packet.Packet{FromLatestPlayerList800(pk)}, true
	case *packet.BiomeDefinitionList:

		return []packet.Packet{fromLatestBiomeDefinitionList800(pk)}, true

	case *packet.ItemRegistry:
		_ = pk
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items800()}}, true
	case *packet.VoxelShapes:

		return nil, true
	case *packet.ClientMovementPredictionSync:
		return []packet.Packet{fromLatestClientMovementPredictionSync(pk)}, true
	case *packet.PlayerLocation:
		return []packet.Packet{fromLatestPlayerLocation(pk)}, true
	case *packet.ResourcePackClientResponse:
		return []packet.Packet{fromLatestResourcePackClientResponse800(pk)}, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{v786.FromLatestResourcePacksInfo786(pk)}, true
	case *packet.ResourcePackStack:
		return []packet.Packet{v786.FromLatestResourcePackStack786(pk)}, true
	case *packet.StartGame:
		return []packet.Packet{fromLatestStartGame800(pk)}, true
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
