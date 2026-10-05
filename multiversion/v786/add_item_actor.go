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

type AddItemActor struct {
	EntityUniqueID int64

	EntityRuntimeID uint64

	Item protocol.ItemInstance

	Position mgl32.Vec3

	Velocity mgl32.Vec3

	EntityMetadata protocol.EntityMetadata

	FromFishing bool
}

func (*AddItemActor) ID() uint32 {
	return IDAddItemActor
}

func (pk *AddItemActor) Marshal(io protocol.IO) {
	io.Varint64(&pk.EntityUniqueID)
	io.Varuint64(&pk.EntityRuntimeID)
	io.ItemInstance(&pk.Item)
	io.Vec3(&pk.Position)
	io.Vec3(&pk.Velocity)
	io.EntityMetadata(&pk.EntityMetadata)
	io.Bool(&pk.FromFishing)
}
