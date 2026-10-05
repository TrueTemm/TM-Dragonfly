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

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type UpdateSubChunkBlocks struct {
	Position protocol.SubChunkPos

	Blocks []protocol.BlockChangeEntry

	Extra []protocol.BlockChangeEntry
}

func (*UpdateSubChunkBlocks) ID() uint32 {
	return IDUpdateSubChunkBlocks
}

func (pk *UpdateSubChunkBlocks) Marshal(io protocol.IO) {
	io.SubChunkPos(&pk.Position)
	protocol.FuncIOSlice(io, &pk.Blocks, blockChangeEntry786)
	protocol.FuncIOSlice(io, &pk.Extra, blockChangeEntry786)
}

func blockChangeEntry786(io protocol.IO, x *protocol.BlockChangeEntry) {
	if p := ProtoOf(io); p != 0 && p < 712 {
		io.BlockPos(&x.BlockPos)
	} else {
		UBlockPos786(io, &x.BlockPos)
	}
	io.Varuint32(&x.BlockRuntimeID)
	io.Varuint32(&x.Flags)
	io.Varuint64(&x.SyncedUpdateEntityUniqueID)
	io.Varuint32(&x.SyncedUpdateType)
}
