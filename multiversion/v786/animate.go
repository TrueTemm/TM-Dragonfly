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
)

const (
	AnimateActionSwingArm = iota + 1
	_
	AnimateActionStopSleep
	AnimateActionCriticalHit
	AnimateActionMagicCriticalHit
)

const (
	AnimateActionRowRight = iota + 128
	AnimateActionRowLeft
)

type Animate struct {
	ActionType int32

	EntityRuntimeID uint64

	BoatRowingTime float32
}

func (*Animate) ID() uint32 {
	return IDAnimate
}

func (pk *Animate) Marshal(io protocol.IO) {
	io.Varint32(&pk.ActionType)
	io.Varuint64(&pk.EntityRuntimeID)
	if pk.ActionType&0x80 != 0 {
		io.Float32(&pk.BoatRowingTime)
	}
}
