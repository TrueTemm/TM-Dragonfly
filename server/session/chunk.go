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

package session

import (
	"bytes"

	"github.com/cespare/xxhash/v2"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const subChunkRequests = false

func (s *Session) ViewChunk(pos world.ChunkPos, dim world.Dimension, c *world.Column) {
	if !s.conn.ClientCacheEnabled() {
		s.sendNetworkChunk(pos, dim, c)
		return
	}
	s.sendBlobHashes(pos, dim, c)
}

func (s *Session) ViewSubChunks(centre world.SubChunkPos, offsets []protocol.SubChunkOffset, tx *world.Tx) {
	r := tx.Range()

	entries := make([]protocol.SubChunkEntry, 0, len(offsets))
	transaction := make(map[uint64]struct{})
	for _, offset := range offsets {
		ind := int16(centre.Y()) + int16(offset[1]) - int16(r[0])>>4
		if ind < 0 || ind > int16(r.Height()>>4) {
			entries = append(entries, protocol.SubChunkEntry{Result: protocol.SubChunkResultIndexOutOfBounds, Offset: offset})
			continue
		}
		col, ok := s.chunkLoader.Chunk(world.ChunkPos{
			centre.X() + int32(offset[0]),
			centre.Z() + int32(offset[2]),
		})
		if !ok {
			entries = append(entries, protocol.SubChunkEntry{Result: protocol.SubChunkResultChunkNotFound, Offset: offset})
			continue
		}
		entries = append(entries, s.subChunkEntry(offset, ind, col, transaction))
	}
	if s.conn.ClientCacheEnabled() && len(transaction) > 0 {
		s.blobMu.Lock()
		s.openChunkTransactions = append(s.openChunkTransactions, transaction)
		s.blobMu.Unlock()
	}
	dim, _ := world.DimensionID(tx.World().Dimension())
	s.writePacket(&packet.SubChunk{
		Dimension:       int32(dim),
		Position:        protocol.SubChunkPos(centre),
		CacheEnabled:    s.conn.ClientCacheEnabled(),
		SubChunkEntries: entries,
	})
}

func (s *Session) subChunkEntry(offset protocol.SubChunkOffset, ind int16, col *world.Column, transaction map[uint64]struct{}) protocol.SubChunkEntry {
	chunkMap := col.HeightMap()

	subMapType, subMap, hasData := byte(protocol.HeightMapDataHasData), protocol.HeightMap{}, true
	higher, lower := true, true
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			y := chunkMap.At(x, z)
			otherInd := col.SubIndex(y)
			switch {
			case otherInd > ind:
				subMap[z][x], lower = 16, false
			case otherInd < ind:
				subMap[z][x], higher = -1, false
			default:
				subMap[z][x], lower, higher = int8(y-col.SubY(otherInd)), false, false
			}
		}
	}
	if higher {
		subMapType, hasData = protocol.HeightMapDataTooHigh, false
	} else if lower {
		subMapType, hasData = protocol.HeightMapDataTooLow, false
	}
	var subMapData protocol.Optional[protocol.HeightMap]
	if hasData {
		subMapData = protocol.Option(subMap)
	}

	sub := col.Sub()[ind]
	if sub.Empty() {
		return protocol.SubChunkEntry{
			Result:              protocol.SubChunkResultSuccessAllAir,
			HeightMapType:       subMapType,
			HeightMapData:       subMapData,
			RenderHeightMapType: subMapType,
			RenderHeightMapData: subMapData,
			Offset:              offset,
		}
	}

	serialisedSubChunk := chunk.EncodeSubChunk(col.Chunk, chunk.NetworkEncoding, int(ind))

	blockEntityBuf := bytes.NewBuffer(nil)
	enc := nbt.NewEncoderWithEncoding(blockEntityBuf, nbt.NetworkLittleEndian)
	for pos, b := range col.BlockEntities {
		if n, ok := b.(world.NBTer); ok && col.SubIndex(int16(pos.Y())) == ind {
			d := n.EncodeNBT()
			d["x"], d["y"], d["z"] = int32(pos[0]), int32(pos[1]), int32(pos[2])
			_ = enc.Encode(d)
		}
	}

	entry := protocol.SubChunkEntry{
		Result:              protocol.SubChunkResultSuccess,
		RawPayload:          protocol.Option(append(serialisedSubChunk, blockEntityBuf.Bytes()...)),
		HeightMapType:       subMapType,
		HeightMapData:       subMapData,
		RenderHeightMapType: subMapType,
		RenderHeightMapData: subMapData,
		Offset:              offset,
	}
	if s.conn.ClientCacheEnabled() {
		if hash := xxhash.Sum64(serialisedSubChunk); s.trackBlob(hash, serialisedSubChunk) {
			transaction[hash] = struct{}{}

			entry.BlobHash = protocol.Option(hash)
			entry.RawPayload = protocol.Option(blockEntityBuf.Bytes())
		}
	}
	return entry
}

func (s *Session) dimensionID(dim world.Dimension) int32 {
	d, _ := world.DimensionID(dim)
	return int32(d)
}

func (s *Session) sendBlobHashes(pos world.ChunkPos, dim world.Dimension, c *world.Column) {
	if subChunkRequests {
		biomes := chunk.EncodeBiomes(c.Chunk, chunk.NetworkEncoding)
		if hash := xxhash.Sum64(biomes); s.trackBlob(hash, biomes) {
			s.writePacket(&packet.LevelChunk{
				Dimension:     s.dimensionID(dim),
				SubChunkCount: 0,
				Position:      protocol.ChunkPos(pos),
				SubChunkLimit: protocol.Option(int32(c.HighestFilledSubChunk())),
				BlobHashes:    []uint64{hash},
				RawPayload:    []byte{0},
				CacheEnabled:  true,
			})
			return
		}
	}

	net := c.NetworkChunk()
	blobs, hashes := net.Blobs()
	count := uint32(len(net.SubChunks))
	m := make(map[uint64]struct{}, len(blobs))
	for _, h := range hashes {
		m[h] = struct{}{}
	}

	s.blobMu.Lock()
	s.openChunkTransactions = append(s.openChunkTransactions, m)
	if l := len(s.blobs); l > 4096 {
		s.blobMu.Unlock()
		s.conf.Log.Error("too many blobs pending", "n", l)
		return
	}
	for i := range hashes {
		s.blobs[hashes[i]] = blobs[i]
	}
	s.blobMu.Unlock()

	s.writePacket(&packet.LevelChunk{
		Dimension:     s.dimensionID(dim),
		Position:      protocol.ChunkPos{pos.X(), pos.Z()},
		SubChunkCount: count,
		CacheEnabled:  true,
		BlobHashes:    hashes,
		RawPayload:    net.Tail,
	})
}

func (s *Session) sendNetworkChunk(pos world.ChunkPos, dim world.Dimension, c *world.Column) {
	if subChunkRequests {
		s.writePacket(&packet.LevelChunk{
			Dimension:     s.dimensionID(dim),
			SubChunkCount: 0,
			Position:      protocol.ChunkPos(pos),
			SubChunkLimit: protocol.Option(int32(c.HighestFilledSubChunk())),
			RawPayload:    append(chunk.EncodeBiomes(c.Chunk, chunk.NetworkEncoding), 0),
		})
		return
	}

	net := c.NetworkChunk()
	s.writePacket(&packet.LevelChunk{
		Dimension:     s.dimensionID(dim),
		Position:      protocol.ChunkPos{pos.X(), pos.Z()},
		SubChunkCount: uint32(len(net.SubChunks)),
		RawPayload:    net.Payload(),
	})
}

func (s *Session) trackBlob(hash uint64, blob []byte) bool {
	s.blobMu.Lock()
	if l := len(s.blobs); l > 4096 {
		s.blobMu.Unlock()
		s.conf.Log.Error("too many blobs pending", "n", l)
		return false
	}
	s.blobs[hash] = blob
	s.blobMu.Unlock()
	return true
}
