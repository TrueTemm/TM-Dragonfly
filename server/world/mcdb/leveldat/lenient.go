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

package leveldat

import (
	"reflect"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

func unmarshalLenient(data []byte, dst any) error {
	var all map[string]any
	if err := nbt.UnmarshalEncoding(data, &all, nbt.LittleEndian); err != nil {
		return err
	}
	filtered := filterCompound(all, reflect.TypeOf(dst).Elem())
	clean, err := nbt.MarshalEncoding(filtered, nbt.LittleEndian)
	if err != nil {
		return err
	}
	return nbt.UnmarshalEncoding(clean, dst, nbt.LittleEndian)
}

func filterCompound(m map[string]any, t reflect.Type) map[string]any {
	if t.Kind() != reflect.Struct {

		return m
	}
	fields := make(map[string]reflect.Type, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		name := f.Name
		if tag, ok := f.Tag.Lookup("nbt"); ok {
			if tag == "-" {
				continue
			}
			if before, _, found := strings.Cut(tag, ","); found {
				name = before
			} else {
				name = tag
			}
		}
		if name != "" {
			fields[name] = f.Type
		}
	}

	out := make(map[string]any, len(m))
	for name, v := range m {
		ft, ok := fields[name]
		if !ok {
			continue
		}
		if nested, ok := v.(map[string]any); ok {
			out[name] = filterCompound(nested, ft)
			continue
		}
		out[name] = v
	}
	return out
}
