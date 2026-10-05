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

package v800

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	PlayerLocationTypeCoordinates = iota
	PlayerLocationTypeHide
)

type PlayerLocation struct {
	Type int32

	EntityUniqueID int64

	Position mgl32.Vec3
}

func (*PlayerLocation) ID() uint32 {
	return IDPlayerLocation
}

func (pk *PlayerLocation) Marshal(io protocol.IO) {
	io.Int32(&pk.Type)
	io.Varint64(&pk.EntityUniqueID)
	if pk.Type == PlayerLocationTypeCoordinates {
		io.Vec3(&pk.Position)
	}
}
