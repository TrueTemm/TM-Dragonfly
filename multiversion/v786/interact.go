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
	_ = iota + 1
	_
	InteractActionLeaveVehicle
	InteractActionMouseOverEntity
	InteractActionNPCOpen
	InteractActionOpenInventory
)

type Interact struct {
	ActionType byte

	TargetEntityRuntimeID uint64

	Position mgl32.Vec3
}

func (*Interact) ID() uint32 {
	return IDInteract
}

func (pk *Interact) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	io.Varuint64(&pk.TargetEntityRuntimeID)
	if pk.ActionType == InteractActionMouseOverEntity || pk.ActionType == InteractActionLeaveVehicle {
		io.Vec3(&pk.Position)
	}
}
