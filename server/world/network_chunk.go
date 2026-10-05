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

package world

import (
	"bytes"

	"github.com/cespare/xxhash/v2"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type NetworkChunk struct {
	SubChunks [][]byte

	Biomes []byte

	Tail []byte

	full   []byte
	hashes []uint64
}

func encodeNetworkChunk(c *Column) *NetworkChunk {
	data := chunk.Encode(c.Chunk, chunk.NetworkEncoding)

	tail := bytes.NewBuffer(make([]byte, 1, 32))
	if len(c.BlockEntities) > 0 {
		enc := nbt.NewEncoderWithEncoding(tail, nbt.NetworkLittleEndian)
		for bp, b := range c.BlockEntities {
			if n, ok := b.(NBTer); ok {
				d := n.EncodeNBT()
				d["x"], d["y"], d["z"] = int32(bp[0]), int32(bp[1]), int32(bp[2])
				_ = enc.Encode(d)
			}
		}
	}
	return &NetworkChunk{SubChunks: data.SubChunks, Biomes: data.Biomes, Tail: tail.Bytes()}
}

func (n *NetworkChunk) Payload() []byte {
	if n.full == nil {
		size := len(n.Biomes) + len(n.Tail)
		for _, s := range n.SubChunks {
			size += len(s)
		}
		full := make([]byte, 0, size)
		for _, s := range n.SubChunks {
			full = append(full, s...)
		}
		full = append(full, n.Biomes...)
		n.full = append(full, n.Tail...)
	}
	return n.full
}

func (n *NetworkChunk) Blobs() ([][]byte, []uint64) {
	blobs := make([][]byte, 0, len(n.SubChunks)+1)
	blobs = append(blobs, n.SubChunks...)
	blobs = append(blobs, n.Biomes)
	if n.hashes == nil {
		hashes := make([]uint64, len(blobs))
		for i, blob := range blobs {
			hashes[i] = xxhash.Sum64(blob)
		}
		n.hashes = hashes
	}
	return blobs, n.hashes
}

func (c *Column) NetworkChunk() *NetworkChunk {
	if c.net == nil {
		c.net = encodeNetworkChunk(c)
	}
	return c.net
}

func (c *Column) markModified() {
	c.modified = true
	c.net = nil
}
