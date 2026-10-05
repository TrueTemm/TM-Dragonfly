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

type PlayerAction struct {
	EntityRuntimeID uint64
	ActionType      int32
	BlockPosition   protocol.BlockPos
	ResultPosition  protocol.BlockPos
	BlockFace       int32
}

func (*PlayerAction) ID() uint32 {
	return IDPlayerAction
}

func (pk *PlayerAction) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.Varint32(&pk.ActionType)
	UBlockPos786(io, &pk.BlockPosition)
	UBlockPos786(io, &pk.ResultPosition)
	io.Varint32(&pk.BlockFace)
}

func toLatestPlayerAction(pk *PlayerAction) *packet.PlayerAction {
	return &packet.PlayerAction{
		EntityRuntimeID: pk.EntityRuntimeID, ActionType: pk.ActionType,
		BlockPosition: pk.BlockPosition, ResultPosition: pk.ResultPosition, BlockFace: pk.BlockFace,
	}
}

func fromLatestPlayerAction(pk *packet.PlayerAction) *PlayerAction {
	return &PlayerAction{
		EntityRuntimeID: pk.EntityRuntimeID, ActionType: pk.ActionType,
		BlockPosition: pk.BlockPosition, ResultPosition: pk.ResultPosition, BlockFace: pk.BlockFace,
	}
}
