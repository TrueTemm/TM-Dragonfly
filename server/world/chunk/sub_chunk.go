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

import "slices"

type SubChunk struct {
	air        uint32
	storages   []*PalettedStorage
	blockLight []uint8
	skyLight   []uint8
}

func (sub *SubChunk) Equals(s *SubChunk) bool {
	if s.air != sub.air || len(s.storages) != len(sub.storages) {
		return false
	}

	for i, st := range s.storages {
		if !st.Equal(sub.storages[i]) {
			return false
		}
	}

	return true
}

func NewSubChunk(air uint32) *SubChunk {
	return &SubChunk{air: air}
}

func (sub *SubChunk) Clone() *SubChunk {
	clone := &SubChunk{
		air:        sub.air,
		storages:   make([]*PalettedStorage, len(sub.storages)),
		blockLight: cloneLight(sub.blockLight),
		skyLight:   cloneLight(sub.skyLight),
	}
	for i, storage := range sub.storages {
		clone.storages[i] = storage.Clone()
	}
	return clone
}

func cloneLight(light []uint8) []uint8 {
	if len(light) == 0 {
		return slices.Clone(light)
	}
	switch &light[0] {
	case noLightPtr:
		return noLight
	case fullLightPtr:
		return fullLight
	default:
		return slices.Clone(light)
	}
}

func (sub *SubChunk) Empty() bool {
	return len(sub.storages) == 0 || (len(sub.storages) == 1 && len(sub.storages[0].palette.values) == 1 && sub.storages[0].palette.values[0] == sub.air)
}

func (sub *SubChunk) Layer(layer uint8) *PalettedStorage {
	for uint8(len(sub.storages)) <= layer {

		sub.storages = append(sub.storages, emptyStorage(sub.air))
	}
	return sub.storages[layer]
}

func (sub *SubChunk) Layers() []*PalettedStorage {
	return sub.storages
}

func (sub *SubChunk) Block(x, y, z byte, layer uint8) uint32 {
	if uint8(len(sub.storages)) <= layer {
		return sub.air
	}
	return sub.storages[layer].At(x, y, z)
}

func (sub *SubChunk) SetBlock(x, y, z byte, layer uint8, block uint32) {
	sub.Layer(layer).Set(x, y, z, block)
}

func (sub *SubChunk) SetBlockLight(x, y, z byte, level uint8) {
	if ptr := &sub.blockLight[0]; ptr == noLightPtr {

		sub.blockLight = append([]byte(nil), sub.blockLight...)
	}
	index := (uint16(x) << 8) | (uint16(z) << 4) | uint16(y)

	i := index >> 1
	bit := (index & 1) << 2
	sub.blockLight[i] = (sub.blockLight[i] & (0xf0 >> bit)) | (level << bit)
}

func (sub *SubChunk) BlockLight(x, y, z byte) uint8 {
	index := (uint16(x) << 8) | (uint16(z) << 4) | uint16(y)
	return (sub.blockLight[index>>1] >> ((index & 1) << 2)) & 0xf
}

func (sub *SubChunk) SetSkyLight(x, y, z byte, level uint8) {
	if ptr := &sub.skyLight[0]; ptr == fullLightPtr || ptr == noLightPtr {

		sub.skyLight = append([]byte(nil), sub.skyLight...)
	}
	index := (uint16(x) << 8) | (uint16(z) << 4) | uint16(y)

	i := index >> 1
	bit := (index & 1) << 2
	sub.skyLight[i] = (sub.skyLight[i] & (0xf0 >> bit)) | (level << bit)
}

func (sub *SubChunk) SkyLight(x, y, z byte) uint8 {
	index := (uint16(x) << 8) | (uint16(z) << 4) | uint16(y)
	return (sub.skyLight[index>>1] >> ((index & 1) << 2)) & 0xf
}

func (sub *SubChunk) compact() {
	newStorages := make([]*PalettedStorage, 0, len(sub.storages))
	for _, storage := range sub.storages {
		storage.compact()
		if len(storage.palette.values) == 1 && storage.palette.values[0] == sub.air {

			continue
		}
		newStorages = append(newStorages, storage)
	}
	sub.storages = newStorages
}
