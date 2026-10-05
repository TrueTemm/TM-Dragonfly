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

package native

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
)

type Reader struct {
	*protocol.Reader
}

type Writer struct {
	*protocol.Writer
}

func hash(rid int32) int32 {
	if rid == 0 {
		return 0
	}
	return int32(blockpalette.HashFor(proto, uint32(rid)))
}

func unhash(h int32) int32 {
	if h == 0 {
		return 0
	}
	return int32(blockpalette.RuntimeIDFor(proto, uint32(h)))
}

func (r Reader) Item(x *protocol.ItemStack) {
	r.Reader.Item(x)
	x.BlockRuntimeID = unhash(x.BlockRuntimeID)
}

func (w Writer) Item(x *protocol.ItemStack) {
	y := *x
	y.BlockRuntimeID = hash(y.BlockRuntimeID)
	w.Writer.Item(&y)
}

func (r Reader) ItemInstance(x *protocol.ItemInstance) {
	r.Reader.ItemInstance(x)
	x.Stack.BlockRuntimeID = unhash(x.Stack.BlockRuntimeID)
}

func (w Writer) ItemInstance(x *protocol.ItemInstance) {
	y := *x
	y.Stack.BlockRuntimeID = hash(y.Stack.BlockRuntimeID)
	w.Writer.ItemInstance(&y)
}
