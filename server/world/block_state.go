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
	"bytes"
	_ "embed"
	"fmt"
	"maps"
	"math"
	"slices"
	"strings"
	"unsafe"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte
)

func init() {
	dec := nbt.NewDecoder(bytes.NewBuffer(blockStateData))

	for {
		var s BlockState
		if err := dec.Decode(&s); err != nil {
			break
		}
		DefaultBlockRegistry.RegisterBlockState(s)
	}
}

type BlockState struct {
	Name       string         `nbt:"name"`
	Properties map[string]any `nbt:"states"`
	Version    int32          `nbt:"version"`
}

type unknownBlock struct {
	BlockState
	data map[string]any
}

func (b unknownBlock) EncodeBlock() (string, map[string]any) {
	return b.Name, b.Properties
}

func (unknownBlock) Model() BlockModel {
	return unknownModel{}
}

func (b unknownBlock) Hash() (uint64, uint64) {
	return 0, math.MaxUint64
}

func (b unknownBlock) EncodeNBT() map[string]any {
	return maps.Clone(b.data)
}

func (b unknownBlock) DecodeNBT(data map[string]any) any {
	b.data = maps.Clone(data)
	return b
}

type stateHash struct {
	name, properties string
}

func hashProperties(properties map[string]any) string {
	if properties == nil {
		return ""
	}
	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var b strings.Builder
	for _, k := range keys {
		switch v := properties[k].(type) {
		case bool:
			if v {
				b.WriteByte(1)
			} else {
				b.WriteByte(0)
			}
		case uint8:
			b.WriteByte(v)
		case int32:
			a := *(*[4]byte)(unsafe.Pointer(&v))
			b.Write(a[:])
		case string:
			b.WriteString(v)
		default:

			panic(fmt.Sprintf("invalid block property type %T for property %v", v, k))
		}
	}

	return b.String()
}
