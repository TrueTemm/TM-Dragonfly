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

type UpdateBlock struct {
	Position protocol.BlockPos

	NewBlockRuntimeID uint32

	Flags uint32

	Layer uint32
}

func (*UpdateBlock) ID() uint32 {
	return IDUpdateBlock
}

func (pk *UpdateBlock) Marshal(io protocol.IO) {
	UBlockPos786(io, &pk.Position)
	io.Varuint32(&pk.NewBlockRuntimeID)
	io.Varuint32(&pk.Flags)
	io.Varuint32(&pk.Layer)
}

type UpdateBlockSynced struct {
	Position protocol.BlockPos

	NewBlockRuntimeID uint32

	Flags uint32

	Layer uint32

	EntityUniqueID uint64

	TransitionType uint64
}

func (*UpdateBlockSynced) ID() uint32 {
	return IDUpdateBlockSynced
}

func (pk *UpdateBlockSynced) Marshal(io protocol.IO) {
	UBlockPos786(io, &pk.Position)
	io.Varuint32(&pk.NewBlockRuntimeID)
	io.Varuint32(&pk.Flags)
	io.Varuint32(&pk.Layer)
	io.Varuint64(&pk.EntityUniqueID)
	io.Varuint64(&pk.TransitionType)
}
