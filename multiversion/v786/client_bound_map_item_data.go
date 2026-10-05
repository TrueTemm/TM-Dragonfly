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
	"image/color"
)

const (
	MapUpdateFlagTexture = 1 << (iota + 1)
	MapUpdateFlagDecoration
	MapUpdateFlagInitialisation
)

type ClientBoundMapItemData struct {
	MapID int64

	UpdateFlags uint32

	Dimension byte

	LockedMap bool

	Origin protocol.BlockPos

	Scale byte

	MapsIncludedIn []int64

	TrackedObjects []protocol.MapTrackedObject

	Decorations []protocol.MapDecoration

	Height int32

	Width int32

	XOffset int32

	YOffset int32

	Pixels []color.RGBA
}

func (*ClientBoundMapItemData) ID() uint32 {
	return IDClientBoundMapItemData
}

func (pk *ClientBoundMapItemData) Marshal(io protocol.IO) {
	io.Varint64(&pk.MapID)
	io.Varuint32(&pk.UpdateFlags)
	io.Uint8(&pk.Dimension)
	io.Bool(&pk.LockedMap)
	io.BlockPos(&pk.Origin)

	if pk.UpdateFlags&MapUpdateFlagInitialisation != 0 {
		protocol.FuncSlice(io, &pk.MapsIncludedIn, io.Varint64)
	}
	if pk.UpdateFlags&(MapUpdateFlagInitialisation|MapUpdateFlagDecoration|MapUpdateFlagTexture) != 0 {
		io.Uint8(&pk.Scale)
	}
	if pk.UpdateFlags&MapUpdateFlagDecoration != 0 {
		protocol.Slice(io, &pk.TrackedObjects)
		protocol.Slice(io, &pk.Decorations)
	}
	if pk.UpdateFlags&MapUpdateFlagTexture != 0 {
		io.Varint32(&pk.Width)
		io.Varint32(&pk.Height)
		io.Varint32(&pk.XOffset)
		io.Varint32(&pk.YOffset)
		protocol.FuncIOSlice(io, &pk.Pixels, VarRGBA786)
	}
}
