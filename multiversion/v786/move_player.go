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
	MoveModeNormal = iota
	MoveModeReset
	MoveModeTeleport
	MoveModeRotation
)

const (
	TeleportCauseUnknown = iota
	TeleportCauseProjectile
	TeleportCauseChorusFruit
	TeleportCauseCommand
	TeleportCauseBehaviour
)

type MovePlayer struct {
	EntityRuntimeID uint64

	Position mgl32.Vec3

	Pitch float32

	Yaw float32

	HeadYaw float32

	Mode byte

	OnGround bool

	RiddenEntityRuntimeID uint64

	TeleportCause int32

	TeleportSourceEntityType int32

	Tick uint64
}

func (*MovePlayer) ID() uint32 {
	return IDMovePlayer
}

func (pk *MovePlayer) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.Vec3(&pk.Position)
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Float32(&pk.HeadYaw)
	io.Uint8(&pk.Mode)
	io.Bool(&pk.OnGround)
	io.Varuint64(&pk.RiddenEntityRuntimeID)
	if pk.Mode == MoveModeTeleport {
		io.Int32(&pk.TeleportCause)
		io.Int32(&pk.TeleportSourceEntityType)
	}
	io.Varuint64(&pk.Tick)
}
