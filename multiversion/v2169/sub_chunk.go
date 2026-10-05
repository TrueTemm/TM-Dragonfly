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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SubChunkEntry struct {
	Offset              protocol.SubChunkOffset
	Result              byte
	RawPayload          protocol.Optional[[]byte]
	HeightMapType       byte
	HeightMapData       protocol.Optional[[]int8]
	RenderHeightMapType byte
	RenderHeightMapData protocol.Optional[[]int8]
	BlobHash            protocol.Optional[uint64]
}

func (x *SubChunkEntry) Marshal(r protocol.IO) {
	protocol.Single(r, &x.Offset)
	r.Uint8(&x.Result)
	protocol.OptionalFunc(r, &x.RawPayload, r.ByteSlice)
	r.Uint8(&x.HeightMapType)
	protocol.OptionalFunc(r, &x.HeightMapData, func(data *[]int8) {
		protocol.FuncSliceOfLen(r, 256, data, r.Int8)
	})
	r.Uint8(&x.RenderHeightMapType)
	protocol.OptionalFunc(r, &x.RenderHeightMapData, func(data *[]int8) {
		protocol.FuncSliceOfLen(r, 256, data, r.Int8)
	})
	protocol.OptionalFunc(r, &x.BlobHash, r.Uint64)
}

type SubChunk struct {
	CacheEnabled    bool
	Dimension       int32
	Position        protocol.SubChunkPos
	SubChunkEntries []SubChunkEntry
}

func (*SubChunk) ID() uint32 { return packet.IDSubChunk }

func (pk *SubChunk) Marshal(io protocol.IO) {
	io.Bool(&pk.CacheEnabled)
	io.Varint32(&pk.Dimension)
	io.SubChunkPos(&pk.Position)
	protocol.Slice(io, &pk.SubChunkEntries)
}

func flatHeightMap(m protocol.HeightMap) []int8 {
	flat := make([]int8, 256)
	for z := range m {
		for x := range m[z] {
			flat[(z<<4)|x] = m[z][x]
		}
	}
	return flat
}

func gridHeightMap(flat []int8) protocol.HeightMap {
	var m protocol.HeightMap
	if len(flat) != 256 {
		return m
	}
	for z := range m {
		for x := range m[z] {
			m[z][x] = flat[(z<<4)|x]
		}
	}
	return m
}

func fromLatestSubChunk(pk *packet.SubChunk) *SubChunk {
	out := &SubChunk{
		CacheEnabled:    pk.CacheEnabled,
		Dimension:       pk.Dimension,
		Position:        pk.Position,
		SubChunkEntries: make([]SubChunkEntry, len(pk.SubChunkEntries)),
	}
	for i, e := range pk.SubChunkEntries {
		entry := SubChunkEntry{
			Offset:              e.Offset,
			Result:              e.Result,
			RawPayload:          e.RawPayload,
			HeightMapType:       e.HeightMapType,
			RenderHeightMapType: e.RenderHeightMapType,
			BlobHash:            e.BlobHash,
		}
		if m, ok := e.HeightMapData.Value(); ok {
			entry.HeightMapData = protocol.Option(flatHeightMap(m))
		}
		if m, ok := e.RenderHeightMapData.Value(); ok {
			entry.RenderHeightMapData = protocol.Option(flatHeightMap(m))
		}
		out.SubChunkEntries[i] = entry
	}
	return out
}

func toLatestSubChunk(pk *SubChunk) *packet.SubChunk {
	out := &packet.SubChunk{
		CacheEnabled:    pk.CacheEnabled,
		Dimension:       pk.Dimension,
		Position:        pk.Position,
		SubChunkEntries: make([]protocol.SubChunkEntry, len(pk.SubChunkEntries)),
	}
	for i, e := range pk.SubChunkEntries {
		entry := protocol.SubChunkEntry{
			Offset:              e.Offset,
			Result:              e.Result,
			RawPayload:          e.RawPayload,
			HeightMapType:       e.HeightMapType,
			RenderHeightMapType: e.RenderHeightMapType,
			BlobHash:            e.BlobHash,
		}
		if flat, ok := e.HeightMapData.Value(); ok {
			entry.HeightMapData = protocol.Option(gridHeightMap(flat))
		}
		if flat, ok := e.RenderHeightMapData.Value(); ok {
			entry.RenderHeightMapData = protocol.Option(gridHeightMap(flat))
		}
		out.SubChunkEntries[i] = entry
	}
	return out
}
