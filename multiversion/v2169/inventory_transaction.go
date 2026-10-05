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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryTransaction struct {
	packet.InventoryTransaction
}

func (pk *InventoryTransaction) Marshal(io protocol.IO) {
	io.Varint32(&pk.LegacyRequestID)
	hasLegacy := pk.LegacyRequestID < -1 && (pk.LegacyRequestID&1) == 0
	io.Bool(&hasLegacy)
	if hasLegacy {
		protocol.Slice(io, &pk.LegacySetItemSlots)
	}

	hasType := true
	io.Bool(&hasType)
	if !hasType {
		io.InvalidValue(hasType, "InventoryTransaction transaction type", "expected presence bool to be true")
	}
	io.TransactionDataType(&pk.TransactionData)
	hasActions := true
	io.Bool(&hasActions)
	if !hasActions {
		io.InvalidValue(hasActions, "InventoryTransaction actions", "expected presence bool to be true")
	}
	actions := inventoryActions(pk.Actions)
	protocol.Slice(io, &actions)
	pk.Actions = latestInventoryActions(actions)

	if use, ok := pk.TransactionData.(*protocol.UseItemTransactionData); ok {
		useItemTransactionData(io, use)
		return
	}
	pk.TransactionData.Marshal(io)
}

func toLatestInventoryTransaction(pk *InventoryTransaction) *packet.InventoryTransaction {
	out := pk.InventoryTransaction
	return &out
}

func fromLatestInventoryTransaction(pk *packet.InventoryTransaction) *InventoryTransaction {
	return &InventoryTransaction{InventoryTransaction: *pk}
}
