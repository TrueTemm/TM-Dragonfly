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

	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
)

func ItemInstanceNew(io protocol.IO, i *protocol.ItemInstance) {
	switch io := io.(type) {
	case Reader786:
		readItemInstanceNew(io.Reader, itemdata.ShieldID(io.Proto), i)
		untranslateItemIDs(io.Proto, &i.Stack)
	case Writer786:
		cp := *i
		cp.Stack = translateItemIDs(io.Proto, i.Stack)
		writeItemInstanceNew(io.Writer, itemdata.ShieldID(io.Proto), &cp)
	default:

		io.ItemInstance(i)
	}
}

func readItemInstanceNew(r *protocol.Reader, shield int32, i *protocol.ItemInstance) {
	x := &i.Stack
	var id int16
	r.Int16(&id)
	x.NetworkID = int32(id)

	r.Uint16(&x.Count)
	r.Varuint32(&x.MetadataValue)

	var hasNetID bool
	r.Bool(&hasNetID)
	if hasNetID {
		var empty uint32
		r.Varuint32(&empty)
		r.Varint32(&i.StackNetworkID)
	} else {
		i.StackNetworkID = 0
	}

	var runtimeID uint32
	r.Varuint32(&runtimeID)
	x.BlockRuntimeID = int32(runtimeID)

	var extraData []byte
	r.ByteSlice(&extraData)
	if len(extraData) == 0 {
		x.NBTData, x.CanBePlacedOn, x.CanBreak = nil, nil, nil
		x.BlockingTick = 0
		return
	}
	buf := bytes.NewBuffer(extraData)
	br := protocol.NewReader(buf, r.ShieldID(), false)

	var length int16
	br.Int16(&length)
	x.NBTData = make(map[string]any)
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

func writeItemInstanceNew(w *protocol.Writer, shield int32, i *protocol.ItemInstance) {
	x := &i.Stack
	id := int16(x.NetworkID)
	w.Int16(&id)

	w.Uint16(&x.Count)
	w.Varuint32(&x.MetadataValue)

	hasNetID := i.StackNetworkID != 0
	w.Bool(&hasNetID)
	if hasNetID {
		var zero uint32
		w.Varuint32(&zero)
		w.Varint32(&i.StackNetworkID)
	}

	runtimeID := uint32(x.BlockRuntimeID)
	w.Varuint32(&runtimeID)

	if x.NetworkID == 0 {
		var zero uint32
		w.Varuint32(&zero)
		return
	}

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
