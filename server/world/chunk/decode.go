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

package chunk

import (
	"bytes"
	"fmt"

	"github.com/df-mc/dragonfly/server/block/cube"
)

func NetworkDecode(br BlockRegistry, data []byte, count int, r cube.Range) (*Chunk, error) {
	return NetworkDecodeBuffer(br, bytes.NewBuffer(data), count, r)
}

func NetworkDecodeBuffer(br BlockRegistry, buf *bytes.Buffer, count int, r cube.Range) (*Chunk, error) {
	c := New(br, r)

	if count < 0 || count > len(c.sub) {
		return nil, fmt.Errorf("invalid sub-chunk count %d: chunk range has %d sub-chunks", count, len(c.sub))
	}
	for i := 0; i < count; i++ {
		index := uint8(i)
		sub, err := decodeSubChunk(buf, c, &index, NetworkEncoding)
		if err != nil {
			return nil, err
		}

		if int(index) >= len(c.sub) {
			return nil, fmt.Errorf("invalid sub-chunk index %d: chunk range has %d sub-chunks", index, len(c.sub))
		}
		c.sub[index] = sub
	}
	var last *PalettedStorage
	for i := 0; i < len(c.sub); i++ {
		b, err := decodePalettedStorage(buf, NetworkEncoding, BiomePaletteEncoding)
		if err != nil {
			return nil, err
		}
		if b == nil {

			if i == 0 {

				return nil, fmt.Errorf("first biome storage pointed to previous one")
			}
			b = last
		} else {
			last = b
		}
		c.biomes[i] = b
	}
	return c, nil
}

func DiskDecode(br BlockRegistry, data SerialisedData, r cube.Range) (*Chunk, error) {
	c := New(br, r)

	err := decodeBiomes(bytes.NewBuffer(data.Biomes), c, DiskEncoding)
	if err != nil {
		return nil, err
	}
	for i, sub := range data.SubChunks {
		if len(sub) == 0 {

			continue
		}
		index := uint8(i)
		if c.sub[index], err = decodeSubChunk(bytes.NewBuffer(sub), c, &index, DiskEncoding); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func decodeSubChunk(buf *bytes.Buffer, c *Chunk, index *byte, e Encoding) (*SubChunk, error) {
	ver, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}
	sub := NewSubChunk(c.air)
	switch ver {
	default:
		return nil, fmt.Errorf("unknown sub chunk version %v: can't decode", ver)
	case 1:

		storage, err := decodePalettedStorage(buf, e, BlockPaletteEncoding{Blocks: c.br})
		if err != nil {
			return nil, err
		}
		sub.storages = append(sub.storages, storage)
	case 8, 9:

		storageCount, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error reading storage count: %w", err)
		}
		if ver == 9 {
			uIndex, err := buf.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("error reading sub-chunk index: %w", err)
			}

			*index = uint8(int8(uIndex) - int8(c.r[0]>>4))
		}
		sub.storages = make([]*PalettedStorage, storageCount)

		for i := byte(0); i < storageCount; i++ {
			sub.storages[i], err = decodePalettedStorage(buf, e, BlockPaletteEncoding{Blocks: c.br})
			if err != nil {
				return nil, err
			}
		}
	}
	return sub, nil
}

func decodeBiomes(buf *bytes.Buffer, c *Chunk, e Encoding) error {
	var last *PalettedStorage
	if buf.Len() != 0 {
		for i := 0; i < len(c.sub); i++ {
			b, err := decodePalettedStorage(buf, e, BiomePaletteEncoding)
			if err != nil {
				return err
			}

			if i == 0 && b == nil {

				return fmt.Errorf("first biome storage pointed to previous one")
			}
			if b == nil {

				b = last
			} else {
				last = b
			}
			c.biomes[i] = b
		}
	}
	return nil
}

func decodePalettedStorage(buf *bytes.Buffer, e Encoding, pe paletteEncoding) (*PalettedStorage, error) {
	blockSize, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading block size: %w", err)
	}
	blockSize >>= 1
	if blockSize == 0x7f {
		return nil, nil
	}

	size := paletteSize(blockSize)
	if size > 32 {
		return nil, fmt.Errorf("cannot read paletted storage (size=%v) %T: size too large", blockSize, pe)
	}
	uint32Count := size.uint32s()

	uint32s := make([]uint32, uint32Count)
	byteCount := uint32Count * 4

	data := buf.Next(byteCount)
	if len(data) != byteCount {
		return nil, fmt.Errorf("cannot read paletted storage (size=%v) %T: not enough block data present: expected %v bytes, got %v", blockSize, pe, byteCount, len(data))
	}
	for i := 0; i < uint32Count; i++ {

		uint32s[i] = uint32(data[i*4]) | uint32(data[i*4+1])<<8 | uint32(data[i*4+2])<<16 | uint32(data[i*4+3])<<24
	}
	p, err := e.decodePalette(buf, paletteSize(blockSize), pe)
	return newPalettedStorage(uint32s, p), err
}
