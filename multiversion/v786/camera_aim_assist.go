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
	CameraAimAssistActionSet = iota
	CameraAimAssistActionClear
)

type CameraAimAssist struct {
	Preset string

	Angle mgl32.Vec2

	Distance float32

	TargetMode byte

	Action byte
}

func (*CameraAimAssist) ID() uint32 {
	return IDCameraAimAssist
}

func (pk *CameraAimAssist) Marshal(io protocol.IO) {
	io.String(&pk.Preset)
	io.Vec2(&pk.Angle)
	io.Float32(&pk.Distance)
	io.Uint8(&pk.TargetMode)
	io.Uint8(&pk.Action)
}
