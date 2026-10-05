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

package blockpalette

import (
	"bytes"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func RemapChunkBlockPalette(raw []byte, subChunkCount, protocolID uint32) []byte {
	if subChunkCount == 0 || len(raw) == 0 || !Supported(protocolID) {
		return raw
	}
	in := bytes.NewBuffer(raw)
	out := new(bytes.Buffer)
	out.Grow(len(raw))

	for sc := uint32(0); sc < subChunkCount; sc++ {
		ver, err := in.ReadByte()
		if err != nil {
			return raw
		}
		out.WriteByte(ver)
		if ver != 8 && ver != 9 {
			return raw
		}
		storageCount, err := in.ReadByte()
		if err != nil {
			return raw
		}
		out.WriteByte(storageCount)
		if ver == 9 {
			b, err := in.ReadByte()
			if err != nil {
				return raw
			}
			out.WriteByte(b)
		}
		for st := byte(0); st < storageCount; st++ {
			header, err := in.ReadByte()
			if err != nil {
				return raw
			}
			out.WriteByte(header)
			blockSize := header >> 1
			if blockSize == 0x7f {
				continue
			}
			words := paletteWordCount(blockSize)
			idx := in.Next(words * 4)
			if len(idx) != words*4 {
				return raw
			}
			out.Write(idx)

			var count int32 = 1
			if blockSize != 0 {
				if err := protocol.Varint32(in, &count); err != nil || count <= 0 {
					return raw
				}
				_ = protocol.WriteVarint32(out, count)
			}
			for i := int32(0); i < count; i++ {
				var v int32
				if err := protocol.Varint32(in, &v); err != nil {
					return raw
				}
				_ = protocol.WriteVarint32(out, int32(HashFor(protocolID, uint32(v))))
			}
		}
	}
	out.Write(in.Bytes())
	return out.Bytes()
}

func paletteWordCount(bitsPerIndex byte) int {
	if bitsPerIndex == 0 {
		return 0
	}
	indicesPerUint32 := 32 / int(bitsPerIndex)
	n := 4096 / indicesPerUint32
	if bitsPerIndex == 3 || bitsPerIndex == 5 || bitsPerIndex == 6 {
		n++
	}
	return n
}
