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

package creative

import (
	_ "embed"
	"fmt"

	"github.com/df-mc/dragonfly/server/internal/nbtconv"

	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type Item struct {
	Stack item.Stack

	Group string
}

type Group struct {
	Category Category

	Name string

	Icon item.Stack
}

func Groups() []Group {
	return creativeGroups
}

func RegisterGroup(group Group) {
	creativeGroups = append(creativeGroups, group)
}

func Items() []Item {
	return creativeItemStacks
}

func RegisterItem(item Item) {
	creativeItemStacks = append(creativeItemStacks, item)
}

var (
	//go:embed creative_items.nbt
	creativeItemData []byte

	creativeGroups []Group

	creativeItemStacks []Item
)

type creativeGroupEntry struct {
	Category int32             `nbt:"category"`
	Name     string            `nbt:"name"`
	Icon     creativeItemEntry `nbt:"icon"`
}

type creativeItemEntry struct {
	Name            string         `nbt:"name"`
	Meta            int16          `nbt:"meta"`
	NBT             map[string]any `nbt:"nbt,omitempty"`
	BlockProperties map[string]any `nbt:"block_properties,omitempty"`
	GroupIndex      int32          `nbt:"group_index,omitempty"`
}

func registerCreativeItems() {
	var m struct {
		Groups []creativeGroupEntry `nbt:"groups"`
		Items  []creativeItemEntry  `nbt:"items"`
	}
	if err := nbt.Unmarshal(creativeItemData, &m); err != nil {
		panic(err)
	}
	for i, group := range m.Groups {
		name := group.Name
		if name == "" {
			name = fmt.Sprint("anon", i)
		}
		st, _ := itemStackFromEntry(group.Icon)
		c := Category{category(group.Category)}
		RegisterGroup(Group{Category: c, Name: name, Icon: st})
	}
	for _, data := range m.Items {
		if data.GroupIndex >= int32(len(creativeGroups)) {
			panic(fmt.Errorf("invalid group index %v for item %v", data.GroupIndex, data.Name))
		}
		st, ok := itemStackFromEntry(data)
		if !ok {
			continue
		}
		RegisterItem(Item{st, creativeGroups[data.GroupIndex].Name})
	}
}

func itemStackFromEntry(data creativeItemEntry) (item.Stack, bool) {
	var (
		it world.Item
		ok bool
	)
	if len(data.BlockProperties) > 0 {

		if b, ok := world.BlockByName(data.Name, data.BlockProperties); ok {
			if it, ok = b.(world.Item); !ok {
				return item.Stack{}, false
			}
		}
	} else {
		if it, ok = world.ItemByName(data.Name, data.Meta); !ok {

			return item.Stack{}, false
		}
		if _, resultingMeta := it.EncodeItem(); resultingMeta != data.Meta {

			return item.Stack{}, false
		}
	}

	if n, ok := it.(world.NBTer); ok {
		if len(data.NBT) > 0 {
			it = n.DecodeNBT(data.NBT).(world.Item)
		}
	}

	st := item.NewStack(it, 1)
	if len(data.NBT) > 0 {
		var invalid bool
		for _, e := range nbtconv.Slice(data.NBT, "ench") {
			if v, ok := e.(map[string]any); ok {
				t, ok := item.EnchantmentByID(int(nbtconv.Int16(v, "id")))
				if !ok {
					invalid = true
					break
				}
				st = st.WithEnchantments(item.NewEnchantment(t, int(nbtconv.Int16(v, "lvl"))))
			}
		}
		if invalid {

			return item.Stack{}, false
		}
	}
	return st, true
}
