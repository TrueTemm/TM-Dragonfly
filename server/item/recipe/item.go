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

package recipe

import (
	"math"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Item interface {
	Count() int

	Empty() bool
}

type blockState struct {
	Name       string         `nbt:"name"`
	Properties map[string]any `nbt:"states"`
	Version    int32          `nbt:"version"`
}

type inputItem struct {
	Name string `nbt:"name"`

	Meta int32 `nbt:"meta"`

	Count int32 `nbt:"count"`

	State blockState `nbt:"block"`

	Tag string `nbt:"tag"`
}

func (i inputItem) Item() (Item, bool) {
	if i.Tag != "" {
		return NewItemTag(i.Tag, int(i.Count)), true
	}

	it, ok := world.ItemByName(i.Name, int16(i.Meta))
	if !ok {
		return nil, false
	}
	st := item.NewStack(it, int(i.Count))
	if i.Meta == math.MaxInt16 {
		st = st.WithValue("variants", true)
	}

	return st, true
}

type inputItems []inputItem

func (d inputItems) Items() ([]Item, bool) {
	s := make([]Item, 0, len(d))
	for _, i := range d {
		itemInput, ok := i.Item()
		if !ok {
			return nil, false
		}
		s = append(s, itemInput)
	}
	return s, true
}

type outputItem struct {
	Name string `nbt:"name"`

	Meta int32 `nbt:"meta"`

	Count int16 `nbt:"count"`

	State blockState `nbt:"block"`

	NBTData map[string]any `nbt:"data"`
}

func (o outputItem) Stack() (item.Stack, bool) {
	it, ok := o.item()
	if !ok {
		return item.Stack{}, false
	}
	if n, ok := it.(world.NBTer); ok && len(o.NBTData) > 0 {
		it = n.DecodeNBT(o.NBTData).(world.Item)
	}

	return item.NewStack(it, int(o.Count)), true
}

func (o outputItem) item() (world.Item, bool) {
	if o.State.Name != "" {
		if b, ok := world.BlockByName(o.State.Name, o.State.Properties); ok {
			if it, ok := b.(world.Item); ok {
				return it, true
			}
		}
	}
	return world.ItemByName(o.Name, int16(o.Meta))
}

type outputItems []outputItem

func (d outputItems) Stacks() ([]item.Stack, bool) {
	s := make([]item.Stack, 0, len(d))
	for _, o := range d {
		itemOutput, ok := o.Stack()
		if !ok {
			return nil, false
		}
		s = append(s, itemOutput)
	}
	return s, true
}
