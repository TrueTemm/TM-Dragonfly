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

const PlayerAuthInputBitsetSize = 65

const (
	InputFlagAscend = iota
	InputFlagDescend
	InputFlagNorthJump
	InputFlagJumpDown
	InputFlagSprintDown
	InputFlagChangeHeight
	InputFlagJumping
	InputFlagAutoJumpingInWater
	InputFlagSneaking
	InputFlagSneakDown
	InputFlagUp
	InputFlagDown
	InputFlagLeft
	InputFlagRight
	InputFlagUpLeft
	InputFlagUpRight
	InputFlagWantUp
	InputFlagWantDown
	InputFlagWantDownSlow
	InputFlagWantUpSlow
	InputFlagSprinting
	InputFlagAscendBlock
	InputFlagDescendBlock
	InputFlagSneakToggleDown
	InputFlagPersistSneak
	InputFlagStartSprinting
	InputFlagStopSprinting
	InputFlagStartSneaking
	InputFlagStopSneaking
	InputFlagStartSwimming
	InputFlagStopSwimming
	InputFlagStartJumping
	InputFlagStartGliding
	InputFlagStopGliding
	InputFlagPerformItemInteraction
	InputFlagPerformBlockActions
	InputFlagPerformItemStackRequest
	InputFlagHandledTeleport
	InputFlagEmoting
	InputFlagMissedSwing
	InputFlagStartCrawling
	InputFlagStopCrawling
	InputFlagStartFlying
	InputFlagStopFlying
	InputFlagClientAckServerData
	InputFlagClientPredictedVehicle
	InputFlagPaddlingLeft
	InputFlagPaddlingRight
	InputFlagBlockBreakingDelayEnabled
	InputFlagHorizontalCollision
	InputFlagVerticalCollision
	InputFlagDownLeft
	InputFlagDownRight
	InputFlagStartUsingItem
	InputFlagCameraRelativeMovementEnabled
	InputFlagRotControlledByMoveDirection
	InputFlagStartSpinAttack
	InputFlagStopSpinAttack
	InputFlagIsHotbarTouchOnly
	InputFlagJumpReleasedRaw
	InputFlagJumpPressedRaw
	InputFlagJumpCurrentRaw
	InputFlagSneakReleasedRaw
	InputFlagSneakPressedRaw
	InputFlagSneakCurrentRaw
)

const (
	InputModeMouse = iota + 1
	InputModeTouch
	InputModeGamePad
	InputModeMotionController
)

const (
	PlayModeNormal = iota
	PlayModeTeaser
	PlayModeScreen
	PlayModeViewer
	PlayModeReality
	PlayModePlacement
	PlayModeLivingRoom
	PlayModeExitLevel
	PlayModeExitLevelLivingRoom
	PlayModeNumModes
)

const (
	InteractionModelTouch = iota
	InteractionModelCrosshair
	InteractionModelClassic
)

type PlayerAuthInput struct {
	Pitch, Yaw float32

	Position mgl32.Vec3

	MoveVector mgl32.Vec2

	HeadYaw float32

	InputData protocol.Bitset

	InputMode uint32

	PlayMode uint32

	InteractionModel uint32

	InteractPitch, InteractYaw float32

	Tick uint64

	Delta mgl32.Vec3

	ItemInteractionData protocol.UseItemTransactionData

	ItemStackRequest protocol.ItemStackRequest

	BlockActions []protocol.PlayerBlockAction

	VehicleRotation mgl32.Vec2

	ClientPredictedVehicle int64

	AnalogueMoveVector mgl32.Vec2

	CameraOrientation mgl32.Vec3

	RawMoveVector mgl32.Vec2
}

func (pk *PlayerAuthInput) ID() uint32 {
	return IDPlayerAuthInput
}

func (pk *PlayerAuthInput) Marshal(io protocol.IO) {
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Vec3(&pk.Position)
	io.Vec2(&pk.MoveVector)
	io.Float32(&pk.HeadYaw)
	io.Bitset(&pk.InputData, PlayerAuthInputBitsetSize)
	io.Varuint32(&pk.InputMode)
	io.Varuint32(&pk.PlayMode)
	if p := ProtoOf(io); p != 0 && p < 671 {
		v := int32(pk.InteractionModel) // signed below 671
		io.Varint32(&v)
		pk.InteractionModel = uint32(v)
	} else {
		io.Varuint32(&pk.InteractionModel)
	}
	io.Float32(&pk.InteractPitch)
	io.Float32(&pk.InteractYaw)
	io.Varuint64(&pk.Tick)
	io.Vec3(&pk.Delta)

	if pk.InputData.Load(InputFlagPerformItemInteraction) {
		io.PlayerInventoryAction(&pk.ItemInteractionData)
	}

	if pk.InputData.Load(InputFlagPerformItemStackRequest) {

		marshalItemStackRequest786(io, &pk.ItemStackRequest)
	}

	if pk.InputData.Load(InputFlagPerformBlockActions) {
		playerBlockActions786(io, &pk.BlockActions)
	}

	if pk.InputData.Load(InputFlagClientPredictedVehicle) {
		io.Vec2(&pk.VehicleRotation)
		io.Varint64(&pk.ClientPredictedVehicle)
	}

	io.Vec2(&pk.AnalogueMoveVector)
	io.Vec3(&pk.CameraOrientation)
	io.Vec2(&pk.RawMoveVector)
}

var playerBlockActionsHasPosition = map[int32]bool{
	protocol.PlayerActionStartBreak:           true,
	protocol.PlayerActionAbortBreak:           true,
	protocol.PlayerActionCrackBreak:           true,
	protocol.PlayerActionPredictDestroyBlock:  true,
	protocol.PlayerActionContinueDestroyBlock: true,
}

func playerBlockActions786(io protocol.IO, x *[]protocol.PlayerBlockAction) {
	count := int32(len(*x))
	io.Varint32(&count)
	if count < 0 {
		io.InvalidValue(count, "block actions", "negative length")
		return
	}
	if uint32(count) != uint32(len(*x)) {
		*x = make([]protocol.PlayerBlockAction, count)
	}
	for i := range *x {
		a := &(*x)[i]
		io.Varint32(&a.Action)
		if playerBlockActionsHasPosition[a.Action] {
			io.BlockPos(&a.BlockPos)
			io.Varint32(&a.Face)
		} else {
			a.BlockPos = protocol.BlockPos{}
			a.Face = 0
		}
	}
}
