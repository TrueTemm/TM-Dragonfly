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
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type AddVolumeEntity struct {
	EntityRuntimeID uint64

	EntityMetadata map[string]any

	EncodingIdentifier string

	InstanceIdentifier string

	Bounds [2]protocol.BlockPos

	Dimension int32

	EngineVersion string
}

func (*AddVolumeEntity) ID() uint32 {
	return IDAddVolumeEntity
}

func (pk *AddVolumeEntity) Marshal(io protocol.IO) {
	io.Uint64(&pk.EntityRuntimeID)
	io.NBT(&pk.EntityMetadata, nbt.NetworkLittleEndian)
	io.String(&pk.EncodingIdentifier)
	io.String(&pk.InstanceIdentifier)
	UBlockPos786(io, &pk.Bounds[0])
	UBlockPos786(io, &pk.Bounds[1])
	io.Varint32(&pk.Dimension)
	io.String(&pk.EngineVersion)
}
