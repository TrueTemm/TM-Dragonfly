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
	CommandOutputTypeNone = iota
	CommandOutputTypeLastOutput
	CommandOutputTypeSilent
	CommandOutputTypeAllOutput
	CommandOutputTypeDataSet
)

type CommandOutput struct {
	CommandOrigin protocol.CommandOrigin

	OutputType byte

	SuccessCount uint32

	OutputMessages []protocol.CommandOutputMessage

	DataSet string
}

func (*CommandOutput) ID() uint32 {
	return IDCommandOutput
}

func (pk *CommandOutput) Marshal(io protocol.IO) {

	CommandOriginData786(io, &pk.CommandOrigin)
	io.Uint8(&pk.OutputType)
	io.Varuint32(&pk.SuccessCount)
	protocol.FuncIOSlice(io, &pk.OutputMessages, CommandOutputMessage786)
	if pk.OutputType == CommandOutputTypeDataSet {
		io.String(&pk.DataSet)
	}
}
