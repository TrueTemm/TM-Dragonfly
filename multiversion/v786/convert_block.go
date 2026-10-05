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
	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func FromLatestUpdateBlock(pk *packet.UpdateBlock, protocolID uint32) *UpdateBlock {
	translated := blockpalette.HashFor(protocolID, pk.NewBlockRuntimeID)
	return &UpdateBlock{
		Position:          pk.Position,
		NewBlockRuntimeID: translated,
		Flags:             pk.Flags,
		Layer:             pk.Layer,
	}
}

func ToLatestUpdateBlock(pk *UpdateBlock, protocolID uint32) *packet.UpdateBlock {
	return &packet.UpdateBlock{
		Position:          pk.Position,
		NewBlockRuntimeID: blockpalette.RuntimeIDFor(protocolID, pk.NewBlockRuntimeID),
		Flags:             pk.Flags,
		Layer:             pk.Layer,
	}
}

func FromLatestUpdateBlockSynced(pk *packet.UpdateBlockSynced, protocolID uint32) *UpdateBlockSynced {
	return &UpdateBlockSynced{
		Position:          pk.Position,
		NewBlockRuntimeID: blockpalette.HashFor(protocolID, pk.NewBlockRuntimeID),
		Flags:             pk.Flags,
		Layer:             pk.Layer,
		EntityUniqueID:    pk.EntityUniqueID,
		TransitionType:    pk.TransitionType,
	}
}

func ToLatestUpdateBlockSynced(pk *UpdateBlockSynced, protocolID uint32) *packet.UpdateBlockSynced {
	return &packet.UpdateBlockSynced{
		Position:          pk.Position,
		NewBlockRuntimeID: blockpalette.RuntimeIDFor(protocolID, pk.NewBlockRuntimeID),
		Flags:             pk.Flags,
		Layer:             pk.Layer,
		EntityUniqueID:    pk.EntityUniqueID,
		TransitionType:    pk.TransitionType,
	}
}

func FromLatestLevelEvent(pk *packet.LevelEvent, protocolID uint32) *packet.LevelEvent {
	switch pk.EventType {
	case packet.LevelEventParticlesDestroyBlock:
		pk.EventData = int32(blockpalette.HashFor(protocolID, uint32(pk.EventData)))
	case packet.LevelEventParticlesCrackBlock:

		pk.EventData = int32(blockpalette.HashFor(protocolID, uint32(pk.EventData)&0x00FFFFFF))
	}
	return pk
}

func TranslateBlockChangeEntries(blocks []protocol.BlockChangeEntry, protocolID uint32, fromLatest bool) []protocol.BlockChangeEntry {
	if len(blocks) == 0 {
		return blocks
	}
	out := make([]protocol.BlockChangeEntry, len(blocks))
	for i, e := range blocks {
		if fromLatest {
			e.BlockRuntimeID = blockpalette.HashFor(protocolID, e.BlockRuntimeID)
		} else {
			e.BlockRuntimeID = blockpalette.RuntimeIDFor(protocolID, e.BlockRuntimeID)
		}
		out[i] = e
	}
	return out
}
