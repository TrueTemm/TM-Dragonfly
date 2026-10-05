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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	MoveActorDeltaFlagHasX = 1 << iota
	MoveActorDeltaFlagHasY
	MoveActorDeltaFlagHasZ
	MoveActorDeltaFlagHasRotX
	MoveActorDeltaFlagHasRotY
	MoveActorDeltaFlagHasRotZ
	MoveActorDeltaFlagOnGround
	MoveActorDeltaFlagTeleport
	MoveActorDeltaFlagForceMove
)

type MoveActorDelta struct {
	Flags uint16

	EntityRuntimeID uint64

	Position mgl32.Vec3

	Rotation mgl32.Vec3
}

func (*MoveActorDelta) ID() uint32 {
	return IDMoveActorDelta
}

func (pk *MoveActorDelta) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.Uint16(&pk.Flags)
	if pk.Flags&MoveActorDeltaFlagHasX != 0 {
		io.Float32(&pk.Position[0])
	} else {
		pk.Position[0] = 0
	}
	if pk.Flags&MoveActorDeltaFlagHasY != 0 {
		io.Float32(&pk.Position[1])
	} else {
		pk.Position[1] = 0
	}
	if pk.Flags&MoveActorDeltaFlagHasZ != 0 {
		io.Float32(&pk.Position[2])
	} else {
		pk.Position[2] = 0
	}
	if pk.Flags&MoveActorDeltaFlagHasRotX != 0 {
		io.ByteFloat(&pk.Rotation[0])
	} else {
		pk.Rotation[0] = 0
	}
	if pk.Flags&MoveActorDeltaFlagHasRotY != 0 {
		io.ByteFloat(&pk.Rotation[1])
	} else {
		pk.Rotation[1] = 0
	}
	if pk.Flags&MoveActorDeltaFlagHasRotZ != 0 {
		io.ByteFloat(&pk.Rotation[2])
	} else {
		pk.Rotation[2] = 0
	}
}
