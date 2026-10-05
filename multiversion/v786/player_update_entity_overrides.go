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
	PlayerUpdateEntityOverridesTypeClearAll = iota
	PlayerUpdateEntityOverridesTypeRemove
	PlayerUpdateEntityOverridesTypeInt
	PlayerUpdateEntityOverridesTypeFloat
)

type PlayerUpdateEntityOverrides struct {
	EntityRuntimeID uint64

	PropertyIndex uint32

	Type byte

	IntValue int32

	FloatValue float32
}

func (*PlayerUpdateEntityOverrides) ID() uint32 {
	return IDPlayerUpdateEntityOverrides
}

func (pk *PlayerUpdateEntityOverrides) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.Varuint32(&pk.PropertyIndex)
	io.Uint8(&pk.Type)
	if pk.Type == PlayerUpdateEntityOverridesTypeInt {
		io.Int32(&pk.IntValue)
	} else if pk.Type == PlayerUpdateEntityOverridesTypeFloat {
		io.Float32(&pk.FloatValue)
	}
}
