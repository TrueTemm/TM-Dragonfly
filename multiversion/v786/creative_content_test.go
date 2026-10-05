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

package v786

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestCreativeGroupCategoryIsInt32(t *testing.T) {
	pk := &CreativeContent{
		Groups: []CreativeGroup786{{Category: 1, Name: "itemGroup.name.planks"}},
	}
	buf := new(bytes.Buffer)
	w := Writer786{Writer: protocol.NewWriter(buf, 0)}
	pk.Marshal(w)

	b := buf.Bytes()

	if len(b) < 6 {
		t.Fatalf("encoded too short: %x", b)
	}
	if b[0] != 0x01 {
		t.Fatalf("expected group count 1, got %#x (%x)", b[0], b)
	}

	if !bytes.Equal(b[1:5], []byte{0x01, 0x00, 0x00, 0x00}) {
		t.Fatalf("Category not written as 4-byte Int32; got % x (full: %x)", b[1:5], b)
	}
	if b[5] != 0x15 {
		t.Fatalf("expected name length 0x15 right after 4-byte Category, got %#x (%x)", b[5], b)
	}
}

func TestCreativeContentRoundTrips(t *testing.T) {
	in := &CreativeContent{
		Groups: []CreativeGroup786{
			{Category: 1, Name: "itemGroup.name.planks", Icon: protocol.ItemStack{}},
			{Category: 2, Name: "", Icon: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: 5}, Count: 1}},
		},
		Items: []protocol.CreativeItem{
			{CreativeItemNetworkID: 0, Item: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: 5}, Count: 1}, GroupIndex: 0},
			{CreativeItemNetworkID: 1, Item: protocol.ItemStack{}, GroupIndex: 1},
		},
	}
	buf := new(bytes.Buffer)
	w := Writer786{Writer: protocol.NewWriter(buf, 0)}
	in.Marshal(w)

	out := &CreativeContent{}
	r := Reader786{Reader: protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)}
	out.Marshal(r)

	if len(out.Groups) != 2 || len(out.Items) != 2 {
		t.Fatalf("counts mismatch: groups=%d items=%d", len(out.Groups), len(out.Items))
	}
	if out.Groups[0].Category != 1 || out.Groups[0].Name != "itemGroup.name.planks" {
		t.Fatalf("group[0] mismatch: %+v", out.Groups[0])
	}
	if out.Groups[1].Category != 2 || out.Groups[1].Icon.NetworkID != 5 {
		t.Fatalf("group[1] mismatch: %+v", out.Groups[1])
	}
	if out.Items[0].GroupIndex != 0 || out.Items[0].Item.NetworkID != 5 {
		t.Fatalf("item[0] mismatch: %+v", out.Items[0])
	}
	if out.Items[1].GroupIndex != 1 {
		t.Fatalf("item[1] mismatch: %+v", out.Items[1])
	}
}
