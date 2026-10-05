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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	return convertFromLatest(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	return convertToLatest(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *packet.BossEvent:
		return []packet.Packet{fromLatestBossEvent(pk)}, true
	case *packet.MoveActorDelta:
		return []packet.Packet{fromLatestMoveActorDelta(pk)}, true
	case *packet.PlaySound:
		return []packet.Packet{fromLatestPlaySound(pk)}, true
	case *packet.SubChunk:
		return []packet.Packet{fromLatestSubChunk(pk)}, true
	case *packet.ItemStackResponse:
		return []packet.Packet{fromLatestItemStackResponse(pk)}, true
	case *packet.PlayerAuthInput:
		return []packet.Packet{fromLatestPlayerAuthInput(pk)}, true
	case *packet.InventoryTransaction:
		return []packet.Packet{fromLatestInventoryTransaction(pk)}, true
	case *packet.ServerBoundDiagnostics:
		return []packet.Packet{fromLatestServerBoundDiagnostics(pk)}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *BossEvent:
		return []packet.Packet{toLatestBossEvent(pk)}, true
	case *MoveActorDelta:
		return []packet.Packet{toLatestMoveActorDelta(pk)}, true
	case *PlaySound:
		return []packet.Packet{toLatestPlaySound(pk)}, true
	case *SubChunk:
		return []packet.Packet{toLatestSubChunk(pk)}, true
	case *ItemStackResponse:
		return []packet.Packet{toLatestItemStackResponse(pk)}, true
	case *PlayerAuthInput:
		return []packet.Packet{toLatestPlayerAuthInput(pk)}, true
	case *InventoryTransaction:
		return []packet.Packet{toLatestInventoryTransaction(pk)}, true
	case *ServerBoundDiagnostics:
		return []packet.Packet{toLatestServerBoundDiagnostics(pk)}, true
	}
	return nil, false
}
