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

package nbtconv

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
)

func InvFromNBT(inv *inventory.Inventory, items []any) {
	for _, itemData := range items {
		data, _ := itemData.(map[string]any)
		it := item.ReadNBT(data, nil)
		if it.Empty() {
			continue
		}
		_ = inv.SetItem(int(Uint8(data, "Slot")), it)
	}
}

func InvToNBT(inv *inventory.Inventory) []any {
	var items []any
	for index, i := range inv.Slots() {
		if i.Empty() {
			continue
		}
		data := item.WriteNBT(i, true)
		data["Slot"] = byte(index)
		items = append(items, data)
	}
	return items
}
