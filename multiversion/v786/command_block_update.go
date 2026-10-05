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
	CommandBlockImpulse = iota
	CommandBlockRepeating
	CommandBlockChain
)

type CommandBlockUpdate struct {
	Block bool

	Position protocol.BlockPos

	Mode uint32

	NeedsRedstone bool

	Conditional bool

	MinecartEntityRuntimeID uint64

	Command string

	LastOutput string

	Name string

	FilteredName string

	ShouldTrackOutput bool

	TickDelay int32

	ExecuteOnFirstTick bool
}

func (*CommandBlockUpdate) ID() uint32 {
	return IDCommandBlockUpdate
}

func (pk *CommandBlockUpdate) Marshal(io protocol.IO) {
	io.Bool(&pk.Block)
	if pk.Block {
		UBlockPos786(io, &pk.Position)
		io.Varuint32(&pk.Mode)
		io.Bool(&pk.NeedsRedstone)
		io.Bool(&pk.Conditional)
	} else {
		io.Varuint64(&pk.MinecartEntityRuntimeID)
	}
	io.String(&pk.Command)
	io.String(&pk.LastOutput)
	io.String(&pk.Name)
	io.String(&pk.FilteredName)
	io.Bool(&pk.ShouldTrackOutput)
	io.Int32(&pk.TickDelay)
	io.Bool(&pk.ExecuteOnFirstTick)
}
