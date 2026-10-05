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
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type MoveActorDelta struct {
	packet.MoveActorDelta
}

func (pk *MoveActorDelta) Marshal(io protocol.IO) {
	io.ActorRuntimeID(&pk.EntityRuntimeID)
	protocol.OptionalFunc(io, &pk.PositionX, io.Float32)
	protocol.OptionalFunc(io, &pk.PositionY, io.Float32)
	protocol.OptionalFunc(io, &pk.PositionZ, io.Float32)
	protocol.OptionalFunc(io, &pk.RotationX, io.ByteFloat)
	protocol.OptionalFunc(io, &pk.RotationY, io.ByteFloat)
	protocol.OptionalFunc(io, &pk.RotationYHead, io.ByteFloat)
	io.Bool(&pk.OnGround)
	io.Bool(&pk.ForceMove)
	io.Bool(&pk.ForceMoveLocalEntity)
	io.Bool(&pk.ForceCompletion)
}

func fromLatestMoveActorDelta(pk *packet.MoveActorDelta) *MoveActorDelta {
	return &MoveActorDelta{MoveActorDelta: *pk}
}

func toLatestMoveActorDelta(pk *MoveActorDelta) *packet.MoveActorDelta {
	out := pk.MoveActorDelta
	return &out
}
