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
	PredictionTypePlayer = iota
	PredictionTypeVehicle
)

type CorrectPlayerMovePrediction struct {
	PredictionType byte

	Position mgl32.Vec3

	Delta mgl32.Vec3

	Rotation mgl32.Vec2

	VehicleAngularVelocity protocol.Optional[float32]

	OnGround bool

	Tick uint64
}

func (*CorrectPlayerMovePrediction) ID() uint32 {
	return IDCorrectPlayerMovePrediction
}

func (pk *CorrectPlayerMovePrediction) Marshal(io protocol.IO) {
	legacy := false
	if p := ProtoOf(io); p != 0 && p < 671 {
		legacy = true // 662 had none of the vehicle fields
	}
	if !legacy {
		io.Uint8(&pk.PredictionType)
	}
	io.Vec3(&pk.Position)
	io.Vec3(&pk.Delta)
	if !legacy && pk.PredictionType == PredictionTypeVehicle {
		io.Vec2(&pk.Rotation)
		protocol.OptionalFunc(io, &pk.VehicleAngularVelocity, io.Float32)
	}
	io.Bool(&pk.OnGround)
	io.Varuint64(&pk.Tick)
}
