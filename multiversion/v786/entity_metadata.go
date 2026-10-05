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
	"reflect"
	"sort"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func readEntityMetadata786(r *protocol.Reader, x *protocol.EntityMetadata) {
	*x = protocol.EntityMetadata{}

	var count uint32
	r.Varuint32(&count)
	for i := uint32(0); i < count; i++ {
		var key, dataType uint32
		r.Varuint32(&key)
		r.Varuint32(&dataType)
		switch dataType {
		case protocol.EntityDataTypeByte:
			var v byte
			r.Uint8(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeInt16:
			var v int16
			r.Int16(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeInt32:
			var v int32
			r.Varint32(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeFloat32:
			var v float32
			r.Float32(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeString:
			var v string
			r.String(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeCompoundTag:
			var v map[string]any
			r.NBT(&v, nbt.NetworkLittleEndian)
			(*x)[key] = v
		case protocol.EntityDataTypeBlockPos:
			var v protocol.BlockPos
			r.BlockPos(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeInt64:
			var v int64
			r.Varint64(&v)
			(*x)[key] = v
		case protocol.EntityDataTypeVec3:
			var v mgl32.Vec3
			r.Vec3(&v)
			(*x)[key] = v
		default:
			r.UnknownEnumOption(dataType, "entity metadata")
		}
	}
}

func writeEntityMetadata786(w *protocol.Writer, x *protocol.EntityMetadata) {
	l := uint32(len(*x))
	w.Varuint32(&l)

	keys := make([]int, 0, l)
	for k := range *x {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		key := uint32(k)
		value := (*x)[uint32(k)]
		w.Varuint32(&key)
		switch v := value.(type) {
		case byte:
			dataType := protocol.EntityDataTypeByte
			w.Varuint32(&dataType)
			w.Uint8(&v)
		case int16:
			dataType := protocol.EntityDataTypeInt16
			w.Varuint32(&dataType)
			w.Int16(&v)
		case int32:
			dataType := protocol.EntityDataTypeInt32
			w.Varuint32(&dataType)
			w.Varint32(&v)
		case float32:
			dataType := protocol.EntityDataTypeFloat32
			w.Varuint32(&dataType)
			w.Float32(&v)
		case string:
			dataType := protocol.EntityDataTypeString
			w.Varuint32(&dataType)
			w.String(&v)
		case map[string]any:
			dataType := protocol.EntityDataTypeCompoundTag
			w.Varuint32(&dataType)
			w.NBT(&v, nbt.NetworkLittleEndian)
		case protocol.BlockPos:
			dataType := protocol.EntityDataTypeBlockPos
			w.Varuint32(&dataType)
			w.BlockPos(&v)
		case int64:
			dataType := protocol.EntityDataTypeInt64
			w.Varuint32(&dataType)
			w.Varint64(&v)
		case mgl32.Vec3:
			dataType := protocol.EntityDataTypeVec3
			w.Varuint32(&dataType)
			w.Vec3(&v)
		default:
			w.UnknownEnumOption(reflect.TypeOf(value), "entity metadata")
		}
	}
}
