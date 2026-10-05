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

package playerdb

import (
	"bytes"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type InventoryData struct {
	Items []item.Stack

	Boots      item.Stack
	Leggings   item.Stack
	Chestplate item.Stack
	Helmet     item.Stack

	OffHand item.Stack

	MainHandSlot uint32
}

func invToData(data InventoryData) jsonInventoryData {
	d := jsonInventoryData{
		MainHandSlot: data.MainHandSlot,
		OffHand:      encodeItem(data.OffHand),
	}
	d.Items = encodeItems(data.Items)
	d.Boots = encodeItem(data.Boots)
	d.Leggings = encodeItem(data.Leggings)
	d.Chestplate = encodeItem(data.Chestplate)
	d.Helmet = encodeItem(data.Helmet)
	return d
}

func dataToInv(data jsonInventoryData) InventoryData {
	d := InventoryData{
		MainHandSlot: data.MainHandSlot,
		OffHand:      decodeItem(data.OffHand),
		Items:        make([]item.Stack, 36),
	}
	decodeItems(data.Items, d.Items)
	d.Boots = decodeItem(data.Boots)
	d.Leggings = decodeItem(data.Leggings)
	d.Chestplate = decodeItem(data.Chestplate)
	d.Helmet = decodeItem(data.Helmet)
	return d
}

func encodeItems(items []item.Stack) (encoded []jsonSlot) {
	encoded = make([]jsonSlot, 0, len(items))
	for slot, i := range items {
		data := encodeItem(i)
		if data == nil {
			continue
		}
		encoded = append(encoded, jsonSlot{Slot: slot, Item: data})
	}
	return
}

func decodeItems(encoded []jsonSlot, items []item.Stack) {
	for _, i := range encoded {
		items[i.Slot] = decodeItem(i.Item)
	}
}

func encodeItem(stack item.Stack) []byte {
	if stack.Empty() {
		return nil
	}

	var b bytes.Buffer
	itemNBT := item.WriteNBT(stack, true)
	encoder := nbt.NewEncoderWithEncoding(&b, nbt.LittleEndian)
	err := encoder.Encode(itemNBT)
	if err != nil {
		return nil
	}
	return b.Bytes()
}

func decodeItem(data []byte) item.Stack {
	var itemNBT map[string]any
	decoder := nbt.NewDecoderWithEncoding(bytes.NewBuffer(data), nbt.LittleEndian)
	err := decoder.Decode(&itemNBT)
	if err != nil {
		return item.Stack{}
	}
	return item.ReadNBT(itemNBT, nil)
}
