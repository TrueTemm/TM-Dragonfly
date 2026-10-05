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

package item

import (
	"bytes"
	"encoding/gob"
	"sort"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

func WriteNBT(s Stack, disk bool) map[string]any {
	tag := make(map[string]any)
	if s.Empty() {
		return tag
	}
	if nbt, ok := s.Item().(world.NBTer); ok {
		for k, v := range nbt.EncodeNBT() {
			tag[k] = v
		}
	}
	writeAnvilCost(tag, s)
	writeDamage(tag, s, disk)
	writeDisplay(tag, s)
	writeDragonflyData(tag, s)
	writeEnchantments(tag, s)
	writeUnbreakable(tag, s)

	data := make(map[string]any)
	if disk {
		writeItemStack(data, tag, s)
	} else {
		for k, v := range tag {
			data[k] = v
		}
	}
	return data
}

func ReadNBT(data map[string]any, s *Stack) Stack {
	disk, tag := s == nil, data
	if disk {
		t, ok := data["tag"].(map[string]any)
		if !ok {
			t = map[string]any{}
		}
		tag = t

		a := readItemStack(data, tag)
		s = &a
	}

	readAnvilCost(tag, s)
	readDamage(tag, s, disk)
	readDisplay(tag, s)
	readDragonflyData(tag, s)
	readEnchantments(tag, s)
	readUnbreakable(tag, s)
	return *s
}

func MapNBT(x map[string]any, k string) Stack {
	if m, ok := x[k].(map[string]any); ok {
		return ReadNBT(m, nil)
	}
	return Stack{}
}

func writeItemStack(m, t map[string]any, s Stack) {
	m["Name"], m["Damage"] = s.Item().EncodeItem()
	if b, ok := s.Item().(world.Block); ok {
		v := map[string]any{}
		writeBlock(v, b)
		m["Block"] = v
	}
	m["Count"] = byte(s.Count())
	if len(t) > 0 {
		m["tag"] = t
	}
}

func writeBlock(m map[string]any, b world.Block) {
	m["name"], m["states"] = b.EncodeBlock()
	m["version"] = chunk.CurrentBlockVersion
}

func writeDragonflyData(m map[string]any, s Stack) {
	if v := s.Values(); len(v) != 0 {
		buf := new(bytes.Buffer)
		if err := gob.NewEncoder(buf).Encode(mapToSlice(v)); err != nil {
			panic("error encoding item user data: " + err.Error())
		}
		m["tmData"] = buf.Bytes()
	}
}

func mapToSlice(m map[string]any) []mapValue {
	values := make([]mapValue, 0, len(m))
	for k, v := range m {
		values = append(values, mapValue{K: k, V: v})
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].K < values[j].K
	})
	return values
}

type mapValue struct {
	K string
	V any
}

func writeEnchantments(m map[string]any, s Stack) {
	if len(s.Enchantments()) != 0 {
		var enchantments []map[string]any
		for _, e := range s.Enchantments() {
			if eType, ok := EnchantmentID(e.Type()); ok {
				enchantments = append(enchantments, map[string]any{
					"id":  int16(eType),
					"lvl": int16(e.Level()),
				})
			}
		}
		m["ench"] = enchantments
	}
}

func writeDisplay(m map[string]any, s Stack) {
	name, lore := s.CustomName(), s.Lore()
	v := map[string]any{}
	if name != "" {
		v["Name"] = name
	}
	if len(lore) != 0 {
		v["Lore"] = lore
	}
	if len(v) != 0 {
		m["display"] = v
	}
}

func writeDamage(m map[string]any, s Stack, disk bool) {
	if v, ok := m["Damage"]; !ok || v.(int16) == 0 {
		if _, ok := s.Item().(Durable); ok {
			if disk {
				m["Damage"] = int16(s.MaxDurability() - s.Durability())
			} else {
				m["Damage"] = int32(s.MaxDurability() - s.Durability())
			}
		}
	}
}

func writeAnvilCost(m map[string]any, s Stack) {
	if cost := s.AnvilCost(); cost > 0 {
		m["RepairCost"] = int32(cost)
	}
}

func writeUnbreakable(m map[string]any, s Stack) {
	if s.Unbreakable() {
		m["Unbreakable"] = byte(1)
	}
}

func readItemStack(m, t map[string]any) Stack {
	var it world.Item
	if blockItem, ok := nbtBlock(m, "Block").(world.Item); ok {
		it = blockItem
	}
	if v, ok := world.ItemByName(nbtString(m, "Name"), nbtInt16(m, "Damage")); ok {
		it = v
	}
	if it == nil {
		return Stack{}
	}
	if n, ok := it.(world.NBTer); ok {
		it = n.DecodeNBT(t).(world.Item)
	}
	return NewStack(it, int(nbtUint8(m, "Count")))
}

func readDamage(m map[string]any, s *Stack, disk bool) {
	if disk {
		*s = s.Damage(int(nbtInt16(m, "Damage")))
		return
	}
	*s = s.Damage(int(nbtInt32(m, "Damage")))
}

func readAnvilCost(m map[string]any, s *Stack) {
	*s = s.WithAnvilCost(int(nbtInt32(m, "RepairCost")))
}

func readEnchantments(m map[string]any, s *Stack) {
	enchantments, ok := m["ench"].([]map[string]any)
	if !ok {
		for _, e := range nbtSlice(m, "ench") {
			if v, ok := e.(map[string]any); ok {
				enchantments = append(enchantments, v)
			}
		}
	}
	for _, ench := range enchantments {
		if t, ok := EnchantmentByID(int(nbtInt16(ench, "id"))); ok {

			*s = s.WithForcedEnchantments(NewEnchantment(t, int(max(nbtInt16(ench, "lvl"), 1))))
		}
	}
}

func readDisplay(m map[string]any, s *Stack) {
	if display, ok := m["display"].(map[string]any); ok {
		if name, ok := display["Name"].(string); ok {

			*s = s.WithCustomName(name)
		}
		if lore, ok := display["Lore"].([]string); ok {
			*s = s.WithLore(lore...)
		} else if lore, ok := display["Lore"].([]any); ok {
			loreLines := make([]string, 0, len(lore))
			for _, l := range lore {
				loreLines = append(loreLines, l.(string))
			}
			*s = s.WithLore(loreLines...)
		}
	}
}

func readDragonflyData(m map[string]any, s *Stack) {
	if customData, ok := m["tmData"]; ok {
		d, ok := customData.([]byte)
		if !ok {
			if itf, ok := customData.([]any); ok {
				for _, v := range itf {
					b, _ := v.(byte)
					d = append(d, b)
				}
			}
		}
		var values []mapValue
		if err := gob.NewDecoder(bytes.NewBuffer(d)).Decode(&values); err != nil {
			panic("error decoding item user data: " + err.Error())
		}
		for _, val := range values {
			*s = s.WithValue(val.K, val.V)
		}
	}
}

func readUnbreakable(m map[string]any, s *Stack) {
	if nbtBool(m, "Unbreakable") {
		*s = s.AsUnbreakable()
	}
}

func nbtBool(m map[string]any, k string) bool { return nbtUint8(m, k) == 1 }

func nbtUint8(m map[string]any, k string) uint8 {
	v, _ := m[k].(uint8)
	return v
}

func nbtString(m map[string]any, k string) string {
	v, _ := m[k].(string)
	return v
}

func nbtInt16(m map[string]any, k string) int16 {
	v, _ := m[k].(int16)
	return v
}

func nbtInt32(m map[string]any, k string) int32 {
	v, _ := m[k].(int32)
	return v
}

func nbtSlice(m map[string]any, k string) []any {
	v, _ := m[k].([]any)
	return v
}

func nbtBlock(m map[string]any, k string) world.Block {
	if mk, ok := m[k].(map[string]any); ok {
		name, _ := mk["name"].(string)
		properties, _ := mk["states"].(map[string]any)
		b, _ := world.BlockByName(name, properties)
		return b
	}
	return nil
}
