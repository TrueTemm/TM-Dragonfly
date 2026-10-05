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

package v827

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type CameraInstruction struct {
	Set protocol.Optional[protocol.CameraInstructionSet]

	Clear protocol.Optional[bool]

	Fade protocol.Optional[protocol.CameraInstructionFade]

	Target protocol.Optional[protocol.CameraInstructionTarget]

	RemoveTarget protocol.Optional[bool]

	FieldOfView protocol.Optional[protocol.CameraInstructionFieldOfView]
}

func (*CameraInstruction) ID() uint32 {
	return IDCameraInstruction
}

func (pk *CameraInstruction) Marshal(io protocol.IO) {
	protocol.OptionalMarshaler(io, &pk.Set)
	protocol.OptionalFunc(io, &pk.Clear, io.Bool)
	protocol.OptionalMarshaler(io, &pk.Fade)
	protocol.OptionalMarshaler(io, &pk.Target)
	protocol.OptionalFunc(io, &pk.RemoveTarget, io.Bool)
	protocol.OptionalMarshaler(io, &pk.FieldOfView)
}
