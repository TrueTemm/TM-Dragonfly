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
	"image/color"
	"math"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	SubChunkRequestModeLimitless786 = math.MaxUint32 - iota
	SubChunkRequestModeLimited786
)

func SliceVarint32Length786[T any, S ~*[]T, A protocol.PtrMarshaler[T]](r protocol.IO, x S) {
	count := int32(len(*x))
	r.Varint32(&count)
	protocol.SliceOfLen[T, S, A](r, uint32(count), x)
}

func SliceUint16Length786[T any, S ~*[]T, A protocol.PtrMarshaler[T]](r protocol.IO, x S) {
	count := uint16(len(*x))
	r.Uint16(&count)
	protocol.SliceOfLen[T, S, A](r, uint32(count), x)
}

func FuncSliceUint16Length786[T any, S ~*[]T](r protocol.IO, x S, f func(*T)) {
	count := uint16(len(*x))
	r.Uint16(&count)
	protocol.FuncSliceOfLen(r, uint32(count), x, f)
}

func PlayerListRemoveEntry786(r protocol.IO, x *PlayerListEntry786) {
	r.UUID(&x.UUID)
}

func CommandOriginData786(r protocol.IO, x *protocol.CommandOrigin) {
	r.Varuint32(&x.Origin)
	r.UUID(&x.UUID)
	r.String(&x.RequestID)
	if x.Origin == uint32(protocol.CommandOriginDevConsole) || x.Origin == uint32(protocol.CommandOriginTest) {
		r.Varint64(&x.PlayerUniqueID)
	}
}

func CommandOutputMessage786(r protocol.IO, x *protocol.CommandOutputMessage) {
	r.Bool(&x.Success)
	r.String(&x.Message)
	protocol.FuncSlice(r, &x.Parameters, r.String)
}

func UBlockPos786(io protocol.IO, x *protocol.BlockPos) {
	if p := ProtoOf(io); p >= 944 {
		io.BlockPos(x)
		return
	}
	io.Varint32(&x[0])
	y := uint32(x[1])
	io.Varuint32(&y)
	x[1] = int32(y)
	io.Varint32(&x[2])
}

func VarRGBA786(io protocol.IO, x *color.RGBA) {
	val := uint32(x.R) | uint32(x.G)<<8 | uint32(x.B)<<16 | uint32(x.A)<<24
	io.Varuint32(&val)
	*x = color.RGBA{R: byte(val), G: byte(val >> 8), B: byte(val >> 16), A: byte(val >> 24)}
}

type PlayerMovementSettings786 struct {
	MovementType int32

	RewindHistorySize int32

	ServerAuthoritativeBlockBreaking bool
}

const PlayerMovementModeServer786 = 1

func PlayerMoveSettings786(io protocol.IO, x *PlayerMovementSettings786) {
	io.Varint32(&x.MovementType)
	io.Varint32(&x.RewindHistorySize)
	io.Bool(&x.ServerAuthoritativeBlockBreaking)
}

func ProtoOf(r protocol.IO) uint32 {
	switch io := r.(type) {
	case Reader786:
		return io.Proto
	case Writer786:
		return io.Proto
	}
	return 0
}
