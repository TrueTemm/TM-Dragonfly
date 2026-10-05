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

package native

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

func convertFromLatest(pk packet.Packet) []packet.Packet {
	switch pk := pk.(type) {
	case *packet.StartGame:

		out := *pk
		out.UseBlockNetworkIDHashes = true
		return []packet.Packet{&out}
	case *packet.LevelChunk:
		out := *pk
		out.RawPayload = blockpalette.RemapChunkBlockPalette(pk.RawPayload, pk.SubChunkCount, proto)
		return []packet.Packet{&out}
	case *packet.SubChunk:
		return []packet.Packet{remapSubChunk(pk)}
	case *packet.UpdateBlock:
		out := *pk
		out.NewBlockRuntimeID = blockpalette.HashFor(proto, pk.NewBlockRuntimeID)
		return []packet.Packet{&out}
	case *packet.UpdateBlockSynced:
		out := *pk
		out.NewBlockRuntimeID = blockpalette.HashFor(proto, pk.NewBlockRuntimeID)
		return []packet.Packet{&out}
	case *packet.UpdateSubChunkBlocks:
		out := *pk
		out.Blocks = v786.TranslateBlockChangeEntries(pk.Blocks, proto, true)
		out.Extra = v786.TranslateBlockChangeEntries(pk.Extra, proto, true)
		return []packet.Packet{&out}
	case *packet.LevelEvent:
		return []packet.Packet{v786.FromLatestLevelEvent(pk, proto)}
	case *packet.LevelSoundEvent:
		if !v786.BlockSoundCarriesRuntimeID(pk.SoundType) {
			return []packet.Packet{pk}
		}
		out := *pk
		out.ExtraData = int32(blockpalette.HashFor(proto, uint32(pk.ExtraData)))
		return []packet.Packet{&out}
	}
	return []packet.Packet{pk}
}

func convertToLatest(pk packet.Packet) []packet.Packet {
	switch pk := pk.(type) {
	case *packet.ClientCacheStatus:

		_ = pk
		return []packet.Packet{&packet.ClientCacheStatus{Enabled: false}}
	case *packet.UpdateSubChunkBlocks:
		out := *pk
		out.Blocks = v786.TranslateBlockChangeEntries(pk.Blocks, proto, false)
		out.Extra = v786.TranslateBlockChangeEntries(pk.Extra, proto, false)
		return []packet.Packet{&out}
	case *packet.LevelSoundEvent:
		if !v786.BlockSoundCarriesRuntimeID(pk.SoundType) {
			return []packet.Packet{pk}
		}
		out := *pk
		out.ExtraData = int32(blockpalette.RuntimeIDFor(proto, uint32(pk.ExtraData)))
		return []packet.Packet{&out}
	}
	return []packet.Packet{pk}
}

func remapSubChunk(pk *packet.SubChunk) *packet.SubChunk {
	out := *pk
	out.SubChunkEntries = make([]protocol.SubChunkEntry, len(pk.SubChunkEntries))
	for i, e := range pk.SubChunkEntries {
		if payload, ok := e.RawPayload.Value(); ok {
			e.RawPayload = protocol.Option(blockpalette.RemapChunkBlockPalette(payload, 1, proto))
		}
		out.SubChunkEntries[i] = e
	}
	return &out
}
