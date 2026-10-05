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
	"bytes"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func ReadItem786(r *protocol.Reader, x *protocol.ItemStack) { readItem786(r, r.ShieldID(), x) }

func readItem786(r *protocol.Reader, shield int32, x *protocol.ItemStack) {
	x.NBTData = make(map[string]any)
	r.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		x.MetadataValue, x.Count, x.CanBePlacedOn, x.CanBreak = 0, 0, nil, nil
		return
	}

	r.Uint16(&x.Count)
	r.Varuint32(&x.MetadataValue)
	r.Varint32(&x.BlockRuntimeID)

	var extraData []byte
	r.ByteSlice(&extraData)
	if len(extraData) == 0 {
		return
	}
	buf := bytes.NewBuffer(extraData)
	br := protocol.NewReader(buf, r.ShieldID(), false)

	var length int16
	br.Int16(&length)
	switch length {
	case 0:
	case -1:
		var version uint8
		br.Uint8(&version)
		if version != 1 {
			br.UnknownEnumOption(version, "item user data version")
			return
		}
		br.NBT(&x.NBTData, nbt.LittleEndian)
	default:
		br.NBT(&x.NBTData, nbt.LittleEndian)
	}
	protocol.FuncSliceUint32Length(br, &x.CanBePlacedOn, br.StringUTF)
	protocol.FuncSliceUint32Length(br, &x.CanBreak, br.StringUTF)
	if x.NetworkID == shield {
		br.Int64(&x.BlockingTick)
	}
}

func WriteItem786(w *protocol.Writer, x *protocol.ItemStack) { writeItem786(w, w.ShieldID(), x) }

func writeItem786(w *protocol.Writer, shield int32, x *protocol.ItemStack) {
	w.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		return
	}

	w.Uint16(&x.Count)
	w.Varuint32(&x.MetadataValue)
	w.Varint32(&x.BlockRuntimeID)

	buf := new(bytes.Buffer)
	bw := protocol.NewWriter(buf, w.ShieldID())

	var length int16
	if len(x.NBTData) != 0 {
		length = -1
		version := uint8(1)
		bw.Int16(&length)
		bw.Uint8(&version)
		bw.NBT(&x.NBTData, nbt.LittleEndian)
	} else {
		bw.Int16(&length)
	}
	protocol.FuncSliceUint32Length(bw, &x.CanBePlacedOn, bw.StringUTF)
	protocol.FuncSliceUint32Length(bw, &x.CanBreak, bw.StringUTF)
	if x.NetworkID == shield {
		blockingTick := x.BlockingTick
		bw.Int64(&blockingTick)
	}
	extraData := buf.Bytes()
	w.ByteSlice(&extraData)
}

func ReadItemInstance786(r *protocol.Reader, i *protocol.ItemInstance) {
	readItemInstance786(r, r.ShieldID(), i)
}

func readItemInstance786(r *protocol.Reader, shield int32, i *protocol.ItemInstance) {
	x := &i.Stack
	x.NBTData = make(map[string]any)
	r.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		x.MetadataValue, x.Count, x.CanBePlacedOn, x.CanBreak = 0, 0, nil, nil
		return
	}

	r.Uint16(&x.Count)
	r.Varuint32(&x.MetadataValue)

	var hasNetID bool
	r.Bool(&hasNetID)
	if hasNetID {
		r.Varint32(&i.StackNetworkID)
	}

	r.Varint32(&x.BlockRuntimeID)

	var extraData []byte
	r.ByteSlice(&extraData)
	if len(extraData) == 0 {
		return
	}
	buf := bytes.NewBuffer(extraData)
	br := protocol.NewReader(buf, r.ShieldID(), false)

	var length int16
	br.Int16(&length)
	switch length {
	case 0:
	case -1:
		var version uint8
		br.Uint8(&version)
		if version != 1 {
			br.UnknownEnumOption(version, "item user data version")
			return
		}
		br.NBT(&x.NBTData, nbt.LittleEndian)
	default:
		br.NBT(&x.NBTData, nbt.LittleEndian)
	}
	protocol.FuncSliceUint32Length(br, &x.CanBePlacedOn, br.StringUTF)
	protocol.FuncSliceUint32Length(br, &x.CanBreak, br.StringUTF)
	if x.NetworkID == shield {
		br.Int64(&x.BlockingTick)
	}
}

func WriteItemInstance786(w *protocol.Writer, i *protocol.ItemInstance) {
	writeItemInstance786(w, w.ShieldID(), i)
}

func writeItemInstance786(w *protocol.Writer, shield int32, i *protocol.ItemInstance) {
	x := &i.Stack
	w.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		return
	}

	w.Uint16(&x.Count)
	w.Varuint32(&x.MetadataValue)

	hasNetID := i.StackNetworkID != 0
	w.Bool(&hasNetID)
	if hasNetID {
		w.Varint32(&i.StackNetworkID)
	}

	w.Varint32(&x.BlockRuntimeID)

	buf := new(bytes.Buffer)
	bw := protocol.NewWriter(buf, w.ShieldID())

	var length int16
	if len(x.NBTData) != 0 {
		length = -1
		version := uint8(1)
		bw.Int16(&length)
		bw.Uint8(&version)
		bw.NBT(&x.NBTData, nbt.LittleEndian)
	} else {
		bw.Int16(&length)
	}
	protocol.FuncSliceUint32Length(bw, &x.CanBePlacedOn, bw.StringUTF)
	protocol.FuncSliceUint32Length(bw, &x.CanBreak, bw.StringUTF)
	if x.NetworkID == shield {
		blockingTick := x.BlockingTick
		bw.Int64(&blockingTick)
	}
	extraData := buf.Bytes()
	w.ByteSlice(&extraData)
}

type Reader786 struct {
	*protocol.Reader

	Proto uint32
}

func (r Reader786) Item(x *protocol.ItemStack) {
	readItem786(r.Reader, itemdata.ShieldID(r.Proto), x)
	untranslateItemIDs(r.Proto, x)
}
func (r Reader786) ItemInstance(i *protocol.ItemInstance) {
	readItemInstance786(r.Reader, itemdata.ShieldID(r.Proto), i)
	untranslateItemIDs(r.Proto, &i.Stack)
}
func (r Reader786) EntityMetadata(x *protocol.EntityMetadata) { readEntityMetadata786(r.Reader, x) }

type Writer786 struct {
	*protocol.Writer

	Proto uint32
}

func (w Writer786) Item(x *protocol.ItemStack) {
	cp := translateItemIDs(w.Proto, *x)
	writeItem786(w.Writer, itemdata.ShieldID(w.Proto), &cp)
}
func (w Writer786) ItemInstance(i *protocol.ItemInstance) {
	cp := *i
	cp.Stack = translateItemIDs(w.Proto, i.Stack)
	writeItemInstance786(w.Writer, itemdata.ShieldID(w.Proto), &cp)
}
func (w Writer786) EntityMetadata(x *protocol.EntityMetadata) { writeEntityMetadata786(w.Writer, x) }

func translateItemIDs(proto uint32, x protocol.ItemStack) protocol.ItemStack {
	x.NetworkID = itemdata.Translate(proto, x.NetworkID)
	if x.BlockRuntimeID != 0 {
		x.BlockRuntimeID = int32(blockpalette.HashFor(proto, uint32(x.BlockRuntimeID)))
	}
	return x
}

func untranslateItemIDs(proto uint32, x *protocol.ItemStack) {
	x.NetworkID = itemdata.ReverseTranslate(proto, x.NetworkID)
	if x.BlockRuntimeID != 0 {
		x.BlockRuntimeID = int32(blockpalette.RuntimeIDFor(proto, uint32(x.BlockRuntimeID)))
	}
}
