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

func applyDeltas2169(p packet.Pool) {
	p[packet.IDBossEvent] = func() packet.Packet { return &BossEvent{} }
	p[packet.IDMoveActorDelta] = func() packet.Packet { return &MoveActorDelta{} }
	p[packet.IDPlaySound] = func() packet.Packet { return &PlaySound{} }
	p[packet.IDSubChunk] = func() packet.Packet { return &SubChunk{} }
	p[packet.IDItemStackResponse] = func() packet.Packet { return &ItemStackResponse{} }
	p[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	p[packet.IDInventoryTransaction] = func() packet.Packet { return &InventoryTransaction{} }
	p[packet.IDServerBoundDiagnostics] = func() packet.Packet { return &ServerBoundDiagnostics{} }
	delete(p, packet.IDSetPlayerFurnaceOptions)
	delete(p, packet.IDRecordStarted)
}

func NewClientPool() packet.Pool {
	p := packet.NewClientPool()
	applyDeltas2169(p)
	return p
}

func NewServerPool() packet.Pool {
	p := packet.NewServerPool()
	applyDeltas2169(p)
	return p
}
