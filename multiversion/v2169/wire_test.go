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
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const realPlayerAuthInput = "000048410080a0c200004841713d83420000ecc1000000000000803f0080a0c20107280c104446485a020102000048410080a0c2d525cdcccc3d2e90a0bd0ad7a33c01010301011c0204050101030001017c0100050000000000000000000000000000000002010001010100000000000000000000000000000000000301000100090000000000000000000000000000000000011880013b0406000000000000000000004841713d83420000ecc10000003f0000803f0000803e9208010101010e010404011c00035b00000001046d656f7702000000010102001880013b08021880013b0801010000c03f000020c00101d879000000000000803f00000000000000000000803f000000000000803f"

const realInventoryTransaction = "0301011c020405010201020001017c01000500000000000000000000000000000000020100010101000000000000000000000000000000000000011880013b0406000000000000000000004841713d83420000ecc10000003f0000803f0000803e92080101"

func decode(t *testing.T, hx string, pk packet.Packet) {
	t.Helper()
	b, err := hex.DecodeString(hx)
	if err != nil {
		t.Fatalf("test data is not hex: %v", err)
	}
	buf := bytes.NewBuffer(b)
	pk.Marshal(Protocol{}.NewReader(buf, 0, false))
	if buf.Len() != 0 {
		t.Fatalf("%T: %v bytes left over after decoding — our reader and the real 2169 writer disagree "+
			"about this packet's shape", pk, buf.Len())
	}
	out := bytes.NewBuffer(nil)
	pk.Marshal(Protocol{}.NewWriter(out, 0))
	if got := hex.EncodeToString(out.Bytes()); got != hx {
		t.Fatalf("%T: re-encoded to different bytes:"+"\n want %v"+"\n  got %v", pk, hx, got)
	}
}

func TestPlayerAuthInputAgainstReal2169(t *testing.T) {
	pk := &PlayerAuthInput{}
	decode(t, realPlayerAuthInput, pk)

	if !pk.InputData.Present() {
		t.Fatal("input flags were read as absent")
	}
	for _, flag := range []int{packet.InputFlagSprinting, packet.InputFlagJumping, packet.InputFlagSneaking,
		packet.InputFlagPerformItemInteraction, packet.InputFlagPerformBlockActions,
		packet.InputFlagPerformItemStackRequest, packet.InputFlagClientPredictedVehicle} {
		if !pk.InputData.Load(flag) {
			t.Errorf("input flag %v was not read", flag)
		}
	}
	if pk.Tick != 4821 {
		t.Errorf("tick: want 4821, got %v", pk.Tick)
	}
	use, ok := pk.ItemInteractionData.Value()
	if !ok {
		t.Fatal("item interaction data was read as absent")
	}
	if len(use.Actions) != 3 {
		t.Fatalf("item interaction actions: want 3, got %v", len(use.Actions))
	}
	if window, ok := use.Actions[0].WindowID.Value(); !ok || window != 124 {
		t.Errorf("first action's window ID: want 124, got %v (present: %v)", window, ok)
	}
	if flags, ok := use.Actions[1].SourceFlags.Value(); !ok || flags != 1 {
		t.Errorf("second action's source flags: want 1, got %v (present: %v)", flags, ok)
	}
	if use.BlockRuntimeID != 1042 {
		t.Errorf("block runtime ID: want 1042, got %v", use.BlockRuntimeID)
	}

	if use.Hand != protocol.HandSlotMainHand {
		t.Errorf("hand: want the main hand, got %v", use.Hand)
	}
	if actions, ok := pk.BlockActions.Value(); !ok || len(actions) != 2 {
		t.Errorf("block actions: want 2, got %v (present: %v)", len(actions), ok)
	}
	if vehicle, ok := pk.ClientPredictedVehicle.Value(); !ok || vehicle != 7788 {
		t.Errorf("predicted vehicle: want 7788, got %v (present: %v)", vehicle, ok)
	}
	if request, ok := pk.ItemStackRequest.Value(); !ok || request.RequestID != 7 {
		t.Errorf("item stack request ID: want 7, got %v (present: %v)", request.RequestID, ok)
	}
}

func TestInventoryTransactionAgainstReal2169(t *testing.T) {
	pk := &InventoryTransaction{}
	decode(t, realInventoryTransaction, pk)

	if len(pk.LegacySetItemSlots) != 1 {
		t.Fatalf("legacy slots: want 1, got %v", len(pk.LegacySetItemSlots))
	}
	if len(pk.Actions) != 2 {
		t.Fatalf("actions: want 2, got %v", len(pk.Actions))
	}
	if window, ok := pk.Actions[0].WindowID.Value(); !ok || window != 124 {
		t.Errorf("first action's window ID: want 124, got %v (present: %v)", window, ok)
	}
	use, ok := pk.TransactionData.(*protocol.UseItemTransactionData)
	if !ok {
		t.Fatalf("transaction data: want a UseItem body, got %T", pk.TransactionData)
	}
	if use.BlockPosition != (protocol.BlockPos{12, 64, -30}) {
		t.Errorf("block position: want 12/64/-30, got %v", use.BlockPosition)
	}
	if use.HotBarSlot != 3 {
		t.Errorf("hot bar slot: want 3, got %v", use.HotBarSlot)
	}
	if use.Hand != protocol.HandSlotMainHand {
		t.Errorf("hand: want the main hand, got %v", use.Hand)
	}
	if use.ClientCooldownState != 1 {
		t.Errorf("client cooldown state: want 1, got %v", use.ClientCooldownState)
	}
}
