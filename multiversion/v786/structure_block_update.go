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
	StructureBlockData = iota
	StructureBlockSave
	StructureBlockLoad
	StructureBlockCorner
	StructureBlockInvalid
	StructureBlockExport
)

const (
	StructureRedstoneSaveModeMemory = iota
	StructureRedstoneSaveModeDisk
)

type StructureBlockUpdate struct {
	Position protocol.BlockPos

	StructureName string

	FilteredStructureName string

	DataField string

	IncludePlayers bool

	ShowBoundingBox bool

	StructureBlockType int32

	Settings protocol.StructureSettings

	RedstoneSaveMode int32

	ShouldTrigger bool

	Waterlogged bool
}

func (*StructureBlockUpdate) ID() uint32 {
	return IDStructureBlockUpdate
}

func (pk *StructureBlockUpdate) Marshal(io protocol.IO) {
	UBlockPos786(io, &pk.Position)
	io.String(&pk.StructureName)
	io.String(&pk.FilteredStructureName)
	io.String(&pk.DataField)
	io.Bool(&pk.IncludePlayers)
	io.Bool(&pk.ShowBoundingBox)
	io.Varint32(&pk.StructureBlockType)
	protocol.Single(io, &pk.Settings)
	io.Varint32(&pk.RedstoneSaveMode)
	io.Bool(&pk.ShouldTrigger)
	io.Bool(&pk.Waterlogged)
}
