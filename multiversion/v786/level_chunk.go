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

type LevelChunk struct {
	Position protocol.ChunkPos

	Dimension int32

	HighestSubChunk uint16

	SubChunkCount uint32

	CacheEnabled bool

	BlobHashes []uint64

	RawPayload []byte
}

func (*LevelChunk) ID() uint32 {
	return IDLevelChunk
}

func (pk *LevelChunk) Marshal(io protocol.IO) {
	io.ChunkPos(&pk.Position)
	io.Varint32(&pk.Dimension)
	io.Varuint32(&pk.SubChunkCount)
	if pk.SubChunkCount == SubChunkRequestModeLimited786 {
		io.Uint16(&pk.HighestSubChunk)
	}
	io.Bool(&pk.CacheEnabled)
	if pk.CacheEnabled {
		protocol.FuncSlice(io, &pk.BlobHashes, io.Uint64)
	}
	io.ByteSlice(&pk.RawPayload)
}
