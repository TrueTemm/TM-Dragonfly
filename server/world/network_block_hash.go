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

package world

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/segmentio/fasthash/fnv1a"
)

func networkBlockHash(name string, properties map[string]any, scratch []byte) (uint32, []byte) {
	if name == "minecraft:unknown" {
		return 0xfffffffe, scratch
	}

	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	data := scratch[:0]
	writeString := func(str string) {
		data = binary.LittleEndian.AppendUint16(data, uint16(len(str)))
		data = append(data, []byte(str)...)
	}

	data = append(data, 10)
	data = append(data, 0)
	data = append(data, 0)

	data = append(data, 8)
	writeString("name")
	writeString(name)

	data = append(data, 10)
	writeString("states")
	for _, k := range keys {
		v := properties[k]
		switch v := v.(type) {
		case string:
			data = append(data, 8)
			writeString(k)
			writeString(v)

		case uint8:
			data = append(data, 1)
			writeString(k)
			data = append(data, byte(v))
		case int8:
			data = append(data, 1)
			writeString(k)
			data = append(data, byte(v))
		case bool:
			b := 0
			if v {
				b = 1
			}
			data = append(data, 1)
			writeString(k)
			data = append(data, byte(b))

		case uint16:
			data = append(data, 2)
			writeString(k)
			data = binary.LittleEndian.AppendUint16(data, uint16(v))
		case int16:
			data = append(data, 2)
			writeString(k)
			data = binary.LittleEndian.AppendUint16(data, uint16(v))

		case uint32:
			data = append(data, 3)
			writeString(k)
			data = binary.LittleEndian.AppendUint32(data, uint32(v))
		case int32:
			data = append(data, 3)
			writeString(k)
			data = binary.LittleEndian.AppendUint32(data, uint32(v))
		default:
			panic(fmt.Sprintf("unhandled nbt type: %T", v))
		}
	}
	data = append(data, 0)
	data = append(data, 0)

	return fnv1a.HashBytes32(data), data
}

func NetworkBlockHash(name string, properties map[string]any) uint32 {
	h, _ := networkBlockHash(name, properties, nil)
	return h
}
