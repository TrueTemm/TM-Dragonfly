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
	"math"
	"slices"
)

type paletteSize byte

type Palette struct {
	last      uint32
	lastIndex int16
	size      paletteSize

	values []uint32
}

func newPalette(size paletteSize, values []uint32) *Palette {
	return &Palette{size: size, values: values, last: math.MaxUint32}
}

func (palette *Palette) Clone() *Palette {
	return &Palette{
		last:      palette.last,
		lastIndex: palette.lastIndex,
		size:      palette.size,
		values:    slices.Clone(palette.values),
	}
}

func (palette *Palette) Len() int {
	return len(palette.values)
}

func (palette *Palette) Add(v uint32) (index int16, resize bool) {
	i := int16(len(palette.values))
	palette.values = append(palette.values, v)

	if palette.needsResize() {
		palette.increaseSize()
		return i, true
	}
	return i, false
}

func (palette *Palette) Replace(f func(v uint32) uint32) {

	palette.last = math.MaxUint32
	for index, v := range palette.values {
		palette.values[index] = f(v)
	}
}

func (palette *Palette) Index(runtimeID uint32) int16 {
	if runtimeID == palette.last {

		return palette.lastIndex
	}

	return palette.indexSlow(runtimeID)
}

func (palette *Palette) indexSlow(runtimeID uint32) int16 {
	l := len(palette.values)
	for i := 0; i < l; i++ {
		if palette.values[i] == runtimeID {
			palette.last = runtimeID
			v := int16(i)
			palette.lastIndex = v
			return v
		}
	}
	return -1
}

func (palette *Palette) Value(i uint16) uint32 {
	return palette.values[i]
}

func (palette *Palette) needsResize() bool {
	return len(palette.values) > (1 << palette.size)
}

var sizes = [...]paletteSize{0, 1, 2, 3, 4, 5, 6, 8, 16}
var offsets = [...]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 8: 7, 16: 8}

func (palette *Palette) increaseSize() {
	palette.size = sizes[offsets[palette.size]+1]
}

func (p paletteSize) padded() bool {
	return p == 3 || p == 5 || p == 6
}

func paletteSizeFor(n int) paletteSize {
	for _, size := range sizes {
		if n <= (1 << size) {
			return size
		}
	}

	return 0
}

func (p paletteSize) uint32s() (n int) {
	uint32Count := 0
	if p != 0 {

		indicesPerUint32 := 32 / int(p)

		uint32Count = 4096 / indicesPerUint32
	}
	if p.padded() {

		uint32Count++
	}
	return uint32Count
}
