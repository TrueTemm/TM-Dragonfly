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
	"slices"
	"unsafe"
)

const (
	uint32ByteSize = 4

	uint32BitSize = uint32ByteSize * 8
)

type PalettedStorage struct {
	bitsPerIndex uint16

	filledBitsPerIndex uint16

	indexMask uint32

	indicesStart unsafe.Pointer

	palette *Palette

	indices []uint32
}

func newPalettedStorage(indices []uint32, palette *Palette) *PalettedStorage {
	var (
		bitsPerIndex       = uint16(len(indices) / uint32BitSize / uint32ByteSize)
		indexMask          = (uint32(1) << bitsPerIndex) - 1
		indicesStart       = (unsafe.Pointer)(unsafe.SliceData(indices))
		filledBitsPerIndex uint16
	)
	if bitsPerIndex != 0 {
		filledBitsPerIndex = uint32BitSize / bitsPerIndex * bitsPerIndex
	}
	return &PalettedStorage{filledBitsPerIndex: filledBitsPerIndex, indexMask: indexMask, indicesStart: indicesStart, bitsPerIndex: bitsPerIndex, indices: indices, palette: palette}
}

func emptyStorage(v uint32) *PalettedStorage {
	return newPalettedStorage([]uint32{}, newPalette(0, []uint32{v}))
}

func (storage *PalettedStorage) Clone() *PalettedStorage {
	return newPalettedStorage(slices.Clone(storage.indices), storage.palette.Clone())
}

func (storage *PalettedStorage) Palette() *Palette {
	return storage.palette
}

func (storage *PalettedStorage) At(x, y, z byte) uint32 {
	return storage.palette.Value(storage.paletteIndex(x&15, y&15, z&15))
}

func (storage *PalettedStorage) PaletteIndex(x, y, z byte) uint16 {
	return storage.paletteIndex(x&15, y&15, z&15)
}

func (storage *PalettedStorage) Set(x, y, z byte, v uint32) {
	index := storage.palette.Index(v)
	if index == -1 {

		index = storage.addNew(v)
	}
	storage.setPaletteIndex(x&15, y&15, z&15, uint16(index))
}

func (storage *PalettedStorage) Equal(other *PalettedStorage) bool {
	if storage == nil || other == nil {
		return false
	}
	if len(storage.indices) == 0 || len(other.indices) == 0 || storage.palette.values[0] == 0 || other.palette.values[0] == 0 {
		return false
	}
	indicesA := unsafe.Slice((*byte)(unsafe.Pointer(&storage.indices[0])), len(storage.indices)*4)
	indicesB := unsafe.Slice((*byte)(unsafe.Pointer(&other.indices[0])), len(other.indices)*4)
	if !bytes.Equal(indicesA, indicesB) {
		return false
	}
	paletteA := unsafe.Slice((*byte)(unsafe.Pointer(&storage.palette.values[0])), len(storage.palette.values)*4)
	paletteB := unsafe.Slice((*byte)(unsafe.Pointer(&other.palette.values[0])), len(other.palette.values)*4)
	return bytes.Equal(paletteA, paletteB)
}

func (storage *PalettedStorage) addNew(v uint32) int16 {
	index, resize := storage.palette.Add(v)
	if resize {
		storage.resize(storage.palette.size)
	}
	return index
}

func (storage *PalettedStorage) paletteIndex(x, y, z byte) uint16 {
	if storage.bitsPerIndex == 0 {

		return 0
	}
	offset := ((uint16(x) << 8) | (uint16(z) << 4) | uint16(y)) * storage.bitsPerIndex
	uint32Offset, bitOffset := offset/storage.filledBitsPerIndex, offset%storage.filledBitsPerIndex

	w := *(*uint32)(unsafe.Pointer(uintptr(storage.indicesStart) + uintptr(uint32Offset<<2)))
	return uint16((w >> bitOffset) & storage.indexMask)
}

func (storage *PalettedStorage) setPaletteIndex(x, y, z byte, i uint16) {
	if storage.bitsPerIndex == 0 {
		return
	}
	offset := ((uint16(x) << 8) | (uint16(z) << 4) | uint16(y)) * storage.bitsPerIndex
	uint32Offset, bitOffset := offset/storage.filledBitsPerIndex, offset%storage.filledBitsPerIndex

	ptr := (*uint32)(unsafe.Pointer(uintptr(storage.indicesStart) + uintptr(uint32Offset<<2)))
	*ptr = (*ptr &^ (storage.indexMask << bitOffset)) | (uint32(i) << bitOffset)
}

func (storage *PalettedStorage) resize(newPaletteSize paletteSize) {
	if newPaletteSize == paletteSize(storage.bitsPerIndex) {
		return
	}

	newStorage := newPalettedStorage(make([]uint32, newPaletteSize.uint32s()), storage.palette)
	for x := byte(0); x < 16; x++ {
		for y := byte(0); y < 16; y++ {
			for z := byte(0); z < 16; z++ {
				newStorage.setPaletteIndex(x, y, z, storage.paletteIndex(x, y, z))
			}
		}
	}

	*storage = *newStorage
}

func (storage *PalettedStorage) compact() {
	if storage.palette.Len() == 0 {
		return
	}
	if storage.palette.Len() == 1 {

		storage.bitsPerIndex = 0
		storage.filledBitsPerIndex = 0
		storage.indexMask = 0
		storage.indicesStart = nil
		storage.indices = nil
		storage.palette.size = 0
		return
	}

	usedIndices := make([]bool, storage.palette.Len())
	for x := byte(0); x < 16; x++ {
		for y := byte(0); y < 16; y++ {
			for z := byte(0); z < 16; z++ {
				usedIndices[storage.paletteIndex(x, y, z)] = true
			}
		}
	}

	usedCount := 0
	allUsed := true
	for _, used := range usedIndices {
		if used {
			usedCount++
		} else {
			allUsed = false
		}
	}

	size := paletteSizeFor(usedCount)
	if allUsed && size == storage.palette.size {
		return
	}

	newRuntimeIDs := make([]uint32, 0, usedCount)
	conversion := make([]uint16, len(usedIndices))
	for index, used := range usedIndices {
		if used {
			conversion[index] = uint16(len(newRuntimeIDs))
			newRuntimeIDs = append(newRuntimeIDs, storage.palette.values[index])
		}
	}

	newStorage := newPalettedStorage(make([]uint32, size.uint32s()), newPalette(size, newRuntimeIDs))

	for x := byte(0); x < 16; x++ {
		for y := byte(0); y < 16; y++ {
			for z := byte(0); z < 16; z++ {

				newStorage.setPaletteIndex(x, y, z, conversion[storage.paletteIndex(x, y, z)])
			}
		}
	}
	*storage = *newStorage
}
