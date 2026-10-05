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

type PlayerAuthInput struct {
	packet.PlayerAuthInput
}

func (pk *PlayerAuthInput) Marshal(io protocol.IO) {
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Vec3(&pk.Position)
	io.Vec2(&pk.MoveVector)
	io.Float32(&pk.HeadYaw)
	inputFlagList(io, &pk.InputData, packet.InputFlagCount)
	io.Varuint32(&pk.InputMode)
	io.Varuint32(&pk.PlayMode)
	io.Varint32(&pk.InteractionModel)
	io.Float32(&pk.InteractPitch)
	io.Float32(&pk.InteractYaw)
	io.Varuint64(&pk.Tick)
	io.Vec3(&pk.Delta)
	DoubleOptionalFunc(io, &pk.ItemInteractionData, io.PlayerInventoryAction)
	DoubleOptionalFunc(io, &pk.ItemStackRequest, func(x *protocol.ItemStackRequest) {
		x.Marshal(io)
	})
	DoubleOptionalFunc(io, &pk.BlockActions, func(x *[]protocol.PlayerBlockAction) {
		protocol.Slice(io, x)
	})
	DoubleOptionalFunc(io, &pk.VehicleRotation, io.Vec2)
	DoubleOptionalFunc(io, &pk.ClientPredictedVehicle, io.ActorUniqueID)
	io.Vec2(&pk.AnalogueMoveVector)
	io.Vec3(&pk.CameraOrientation)
	io.Vec2(&pk.RawMoveVector)
}

func toLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	out := pk.PlayerAuthInput
	return &out
}

func fromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	return &PlayerAuthInput{PlayerAuthInput: *pk}
}
