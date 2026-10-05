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

type AddActor struct {
	EntityUniqueID int64

	EntityRuntimeID uint64

	EntityType string

	Position mgl32.Vec3

	Velocity mgl32.Vec3

	Pitch float32

	Yaw float32

	HeadYaw float32

	BodyYaw float32

	Attributes []protocol.AttributeValue

	EntityMetadata protocol.EntityMetadata

	EntityProperties protocol.EntityProperties

	EntityLinks []protocol.EntityLink
}

func (*AddActor) ID() uint32 {
	return IDAddActor
}

func (pk *AddActor) Marshal(io protocol.IO) {
	io.Varint64(&pk.EntityUniqueID)
	io.Varuint64(&pk.EntityRuntimeID)
	io.String(&pk.EntityType)
	io.Vec3(&pk.Position)
	io.Vec3(&pk.Velocity)
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Float32(&pk.HeadYaw)
	io.Float32(&pk.BodyYaw)
	protocol.Slice(io, &pk.Attributes)
	io.EntityMetadata(&pk.EntityMetadata)
	protocol.Single(io, &pk.EntityProperties)
	EntityLinks786(io, &pk.EntityLinks)
}
