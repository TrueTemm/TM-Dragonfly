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

type SubChunkEntry786 struct {
	Offset protocol.SubChunkOffset

	Result byte

	RawPayload []byte

	HeightMapType byte

	HeightMapData []int8

	BlobHash uint64
}

func marshalSubChunkEntry786NoCache(io protocol.IO, x *SubChunkEntry786) {
	protocol.Single(io, &x.Offset)
	io.Uint8(&x.Result)
	io.ByteSlice(&x.RawPayload)
	io.Uint8(&x.HeightMapType)
	if x.HeightMapType == protocol.HeightMapDataHasData {
		protocol.FuncSliceOfLen(io, 256, &x.HeightMapData, io.Int8)
	}
}

func subChunkEntriesUint32Length786(io protocol.IO, x *[]SubChunkEntry786) {
	count := uint32(len(*x))
	io.Uint32(&count)
	protocol.FuncIOSliceOfLen(io, count, x, marshalSubChunkEntry786NoCache)
}

type SubChunk struct {
	CacheEnabled bool

	Dimension int32

	Position protocol.SubChunkPos

	SubChunkEntries []SubChunkEntry786
}

func (*SubChunk) ID() uint32 {
	return IDSubChunk
}

func (pk *SubChunk) Marshal(io protocol.IO) {
	io.Bool(&pk.CacheEnabled)
	io.Varint32(&pk.Dimension)
	io.SubChunkPos(&pk.Position)

	subChunkEntriesUint32Length786(io, &pk.SubChunkEntries)
}
