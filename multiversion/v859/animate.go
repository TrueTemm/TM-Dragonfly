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

package v859

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	animateActionRowRight = 128
	animateActionRowLeft  = 129
)

type Animate struct {
	ActionType      int32
	EntityRuntimeID uint64
	Data            float32
	RowingTime      float32
}

func (*Animate) ID() uint32 {
	return IDAnimate
}

func (pk *Animate) Marshal(io protocol.IO) {
	io.Varint32(&pk.ActionType)
	io.Varuint64(&pk.EntityRuntimeID)
	io.Float32(&pk.Data)
	if pk.ActionType == animateActionRowLeft || pk.ActionType == animateActionRowRight {
		io.Float32(&pk.RowingTime)
	}
}

func toLatestAnimate859(pk *Animate) *packet.Animate {
	return &packet.Animate{ActionType: uint8(pk.ActionType), EntityRuntimeID: pk.EntityRuntimeID, Data: pk.Data}
}

func fromLatestAnimate859(pk *packet.Animate) *Animate {
	return &Animate{ActionType: int32(pk.ActionType), EntityRuntimeID: pk.EntityRuntimeID, Data: pk.Data}
}
