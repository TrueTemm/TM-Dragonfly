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

func toLatestMoveActorDelta(pk *MoveActorDelta) *packet.MoveActorDelta {
	out := &packet.MoveActorDelta{EntityRuntimeID: pk.EntityRuntimeID}
	if pk.Flags&MoveActorDeltaFlagHasX != 0 {
		out.PositionX = protocol.Option(pk.Position[0])
	}
	if pk.Flags&MoveActorDeltaFlagHasY != 0 {
		out.PositionY = protocol.Option(pk.Position[1])
	}
	if pk.Flags&MoveActorDeltaFlagHasZ != 0 {
		out.PositionZ = protocol.Option(pk.Position[2])
	}
	if pk.Flags&MoveActorDeltaFlagHasRotX != 0 {

		_ = pk.Rotation[0]
	}
	if pk.Flags&MoveActorDeltaFlagHasRotY != 0 {
		out.RotationY = protocol.Option(pk.Rotation[1])
	}
	if pk.Flags&MoveActorDeltaFlagHasRotZ != 0 {
		out.RotationYHead = protocol.Option(pk.Rotation[2])
	}
	out.OnGround = pk.Flags&MoveActorDeltaFlagOnGround != 0
	out.ForceMove = pk.Flags&MoveActorDeltaFlagForceMove != 0
	return out
}

func fromLatestMoveActorDelta(pk *packet.MoveActorDelta) *MoveActorDelta {
	out := &MoveActorDelta{EntityRuntimeID: pk.EntityRuntimeID}
	if v, ok := pk.PositionX.Value(); ok {
		out.Flags |= MoveActorDeltaFlagHasX
		out.Position[0] = v
	}
	if v, ok := pk.PositionY.Value(); ok {
		out.Flags |= MoveActorDeltaFlagHasY
		out.Position[1] = v
	}
	if v, ok := pk.PositionZ.Value(); ok {
		out.Flags |= MoveActorDeltaFlagHasZ
		out.Position[2] = v
	}
	if v, ok := pk.RotationY.Value(); ok {
		out.Flags |= MoveActorDeltaFlagHasRotY
		out.Rotation[1] = v
	}
	if v, ok := pk.RotationYHead.Value(); ok {
		out.Flags |= MoveActorDeltaFlagHasRotZ
		out.Rotation[2] = v
	}
	if pk.OnGround {
		out.Flags |= MoveActorDeltaFlagOnGround
	}
	if pk.ForceMove {
		out.Flags |= MoveActorDeltaFlagForceMove
	}

	return out
}

func toLatestMovePlayer(pk *MovePlayer) *packet.MovePlayer {
	out := &packet.MovePlayer{
		EntityRuntimeID: pk.EntityRuntimeID, Position: pk.Position, Pitch: pk.Pitch, Yaw: pk.Yaw,
		HeadYaw: pk.HeadYaw, Mode: pk.Mode, OnGround: pk.OnGround,
		RiddenEntityRuntimeID: pk.RiddenEntityRuntimeID, Tick: pk.Tick,
	}
	if pk.Mode == MoveModeTeleport {
		out.TeleportData = protocol.Option(protocol.TeleportData{
			TeleportCause: pk.TeleportCause, TeleportSourceEntityType: pk.TeleportSourceEntityType,
		})
	}
	return out
}
func fromLatestMovePlayer(pk *packet.MovePlayer) *MovePlayer {
	out := &MovePlayer{
		EntityRuntimeID: pk.EntityRuntimeID, Position: pk.Position, Pitch: pk.Pitch, Yaw: pk.Yaw,
		HeadYaw: pk.HeadYaw, Mode: pk.Mode, OnGround: pk.OnGround,
		RiddenEntityRuntimeID: pk.RiddenEntityRuntimeID, Tick: pk.Tick,
	}
	if td, ok := pk.TeleportData.Value(); ok {
		out.TeleportCause, out.TeleportSourceEntityType = td.TeleportCause, td.TeleportSourceEntityType
	}
	return out
}

func toLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	out := &packet.PlayerAuthInput{
		Pitch: pk.Pitch, Yaw: pk.Yaw, Position: pk.Position, MoveVector: pk.MoveVector, HeadYaw: pk.HeadYaw,
		InputData: inputFlagsFromBitset786(pk.InputData), InputMode: pk.InputMode, PlayMode: pk.PlayMode,
		InteractionModel: int32(pk.InteractionModel), InteractPitch: pk.InteractPitch, InteractYaw: pk.InteractYaw,
		Tick: pk.Tick, Delta: pk.Delta, AnalogueMoveVector: pk.AnalogueMoveVector,
		CameraOrientation: pk.CameraOrientation, RawMoveVector: pk.RawMoveVector,
	}
	if pk.InputData.Load(InputFlagPerformItemInteraction) {
		out.ItemInteractionData = protocol.Option(pk.ItemInteractionData)
	}
	if pk.InputData.Load(InputFlagPerformItemStackRequest) {
		out.ItemStackRequest = protocol.Option(pk.ItemStackRequest)
	}
	if pk.InputData.Load(InputFlagPerformBlockActions) {
		out.BlockActions = protocol.Option(pk.BlockActions)
	}
	if pk.InputData.Load(InputFlagClientPredictedVehicle) {
		out.VehicleRotation = protocol.Option(pk.VehicleRotation)
		out.ClientPredictedVehicle = protocol.Option(pk.ClientPredictedVehicle)
	}
	return out
}

func fromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	out := &PlayerAuthInput{
		Pitch: pk.Pitch, Yaw: pk.Yaw, Position: pk.Position, MoveVector: pk.MoveVector, HeadYaw: pk.HeadYaw,
		InputData: bitsetFromInputFlags786(pk.InputData), InputMode: pk.InputMode, PlayMode: pk.PlayMode,
		InteractionModel: uint32(pk.InteractionModel), InteractPitch: pk.InteractPitch, InteractYaw: pk.InteractYaw,
		Tick: pk.Tick, Delta: pk.Delta, AnalogueMoveVector: pk.AnalogueMoveVector,
		CameraOrientation: pk.CameraOrientation, RawMoveVector: pk.RawMoveVector,
	}
	if v, ok := pk.ItemInteractionData.Value(); ok {
		out.ItemInteractionData = v
	}
	if v, ok := pk.ItemStackRequest.Value(); ok {
		out.ItemStackRequest = v
	}
	if v, ok := pk.BlockActions.Value(); ok {
		out.BlockActions = v
	}
	if v, ok := pk.VehicleRotation.Value(); ok {
		out.VehicleRotation = v
	}
	if v, ok := pk.ClientPredictedVehicle.Value(); ok {
		out.ClientPredictedVehicle = v
	}
	return out
}

func inputFlagsFromBitset786(old protocol.Bitset) protocol.InputFlags {
	flags := protocol.NewInputFlags(packet.InputFlagCount)
	n := old.Len()
	if n > packet.InputFlagCount {
		n = packet.InputFlagCount
	}
	for i := 0; i < n; i++ {
		if old.Load(i) {
			flags.Set(i)
		}
	}
	return flags
}

func bitsetFromInputFlags786(new protocol.InputFlags) protocol.Bitset {
	old := protocol.NewBitset(PlayerAuthInputBitsetSize)
	if !new.Present() {
		return old
	}
	n := new.Len()
	if n > PlayerAuthInputBitsetSize {
		n = PlayerAuthInputBitsetSize
	}
	for i := 0; i < n; i++ {
		if new.Load(i) {
			old.Set(i)
		}
	}
	return old
}

func toLatestClientMovementPredictionSync(pk *ClientMovementPredictionSync) *packet.ClientMovementPredictionSync {
	return &packet.ClientMovementPredictionSync{
		ActorFlags: pk.ActorFlags, BoundingBoxScale: pk.BoundingBoxScale, BoundingBoxWidth: pk.BoundingBoxWidth,
		BoundingBoxHeight: pk.BoundingBoxHeight, MovementSpeed: pk.MovementSpeed,
		UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed, LavaMovementSpeed: pk.LavaMovementSpeed,
		JumpStrength: pk.JumpStrength, Health: pk.Health, Hunger: pk.Hunger, EntityUniqueID: pk.EntityUniqueID,
	}
}
func fromLatestClientMovementPredictionSync(pk *packet.ClientMovementPredictionSync) *ClientMovementPredictionSync {
	return &ClientMovementPredictionSync{
		ActorFlags: pk.ActorFlags, BoundingBoxScale: pk.BoundingBoxScale, BoundingBoxWidth: pk.BoundingBoxWidth,
		BoundingBoxHeight: pk.BoundingBoxHeight, MovementSpeed: pk.MovementSpeed,
		UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed, LavaMovementSpeed: pk.LavaMovementSpeed,
		JumpStrength: pk.JumpStrength, Health: pk.Health, Hunger: pk.Hunger, EntityUniqueID: pk.EntityUniqueID,
	}
}
