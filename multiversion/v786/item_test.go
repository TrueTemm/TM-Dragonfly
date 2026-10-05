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

func TestItem786AirIsOneByte(t *testing.T) {
	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	x := &protocol.ItemStack{}
	WriteItem786(w, x)
	if buf.Len() != 1 {
		t.Fatalf("expected air item to be exactly 1 byte on the wire, got %d: %x", buf.Len(), buf.Bytes())
	}

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	out := &protocol.ItemStack{}
	ReadItem786(r, out)
	if out.NetworkID != 0 {
		t.Fatalf("expected NetworkID=0 after round-trip, got %d", out.NetworkID)
	}
}

func TestItem786NonAirRoundTrips(t *testing.T) {
	in := &protocol.ItemStack{}
	in.NetworkID = 42
	in.Count = 3
	in.MetadataValue = 7
	in.BlockRuntimeID = 99
	in.NBTData = map[string]any{"foo": "bar"}
	in.CanBePlacedOn = []string{"minecraft:stone"}
	in.CanBreak = []string{"minecraft:dirt"}

	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	WriteItem786(w, in)

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	out := &protocol.ItemStack{}
	ReadItem786(r, out)

	if out.NetworkID != in.NetworkID || out.Count != in.Count || out.MetadataValue != in.MetadataValue || out.BlockRuntimeID != in.BlockRuntimeID {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", out, in)
	}
	if out.NBTData["foo"] != "bar" {
		t.Fatalf("NBT round-trip mismatch: %+v", out.NBTData)
	}
	if len(out.CanBePlacedOn) != 1 || out.CanBePlacedOn[0] != "minecraft:stone" {
		t.Fatalf("CanBePlacedOn round-trip mismatch: %+v", out.CanBePlacedOn)
	}
	if len(out.CanBreak) != 1 || out.CanBreak[0] != "minecraft:dirt" {
		t.Fatalf("CanBreak round-trip mismatch: %+v", out.CanBreak)
	}
}

func TestItemInstance786AirIsOneByte(t *testing.T) {
	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	x := &protocol.ItemInstance{}
	WriteItemInstance786(w, x)
	if buf.Len() != 1 {
		t.Fatalf("expected air ItemInstance to be exactly 1 byte on the wire, got %d: %x", buf.Len(), buf.Bytes())
	}

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	out := &protocol.ItemInstance{}
	ReadItemInstance786(r, out)
	if out.Stack.NetworkID != 0 {
		t.Fatalf("expected NetworkID=0 after round-trip, got %d", out.Stack.NetworkID)
	}
}

func TestItemInstance786NonAirRoundTrips(t *testing.T) {
	in := &protocol.ItemInstance{StackNetworkID: 5}
	in.Stack.NetworkID = 42
	in.Stack.Count = 3
	in.Stack.MetadataValue = 7
	in.Stack.BlockRuntimeID = 99
	in.Stack.NBTData = map[string]any{"foo": "bar"}

	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	WriteItemInstance786(w, in)

	var latestBuf bytes.Buffer
	latestW := protocol.NewWriter(&latestBuf, 0)
	latestIn := *in
	latestW.ItemInstance(&latestIn)
	if bytes.Equal(buf.Bytes(), latestBuf.Bytes()) {
		t.Fatalf("expected 786 encoding to differ from latest's, both were %x", buf.Bytes())
	}

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	out := &protocol.ItemInstance{}
	ReadItemInstance786(r, out)

	if out.Stack.NetworkID != in.Stack.NetworkID || out.Stack.Count != in.Stack.Count ||
		out.Stack.MetadataValue != in.Stack.MetadataValue || out.Stack.BlockRuntimeID != in.Stack.BlockRuntimeID ||
		out.StackNetworkID != in.StackNetworkID {
		t.Fatalf("round-trip mismatch: got %+v (stack %+v), want %+v (stack %+v)", out, out.Stack, in, in.Stack)
	}
	if out.Stack.NBTData["foo"] != "bar" {
		t.Fatalf("NBT round-trip mismatch: %+v", out.Stack.NBTData)
	}
}

func TestReaderWriter786SatisfyIO(t *testing.T) {
	var _ protocol.IO = Reader786{}
	var _ protocol.IO = Writer786{}
}

func TestEntityMetadata786RoundTripsWithoutLegacyByte(t *testing.T) {
	in := protocol.EntityMetadata{
		protocol.EntityDataKeyFlags: int64(1 << protocol.EntityDataFlagOnFire),
		protocol.EntityDataKeyName:  "hi",
		protocol.EntityDataKeyScale: float32(1.5),
	}

	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	writeEntityMetadata786(w, &in)

	var latestBuf bytes.Buffer
	latestW := protocol.NewWriter(&latestBuf, 0)
	latestW.EntityMetadata(&in)
	if buf.Len() >= latestBuf.Len() {
		t.Fatalf("expected 786 encoding (%d bytes) to be smaller than latest's (%d bytes, includes a legacy-type byte per entry)", buf.Len(), latestBuf.Len())
	}

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	var out protocol.EntityMetadata
	readEntityMetadata786(r, &out)

	if len(out) != len(in) {
		t.Fatalf("round-trip entry count mismatch: got %d, want %d (%+v)", len(out), len(in), out)
	}
	if out[protocol.EntityDataKeyFlags] != in[protocol.EntityDataKeyFlags] {
		t.Fatalf("Flags round-trip mismatch: got %v, want %v", out[protocol.EntityDataKeyFlags], in[protocol.EntityDataKeyFlags])
	}
	if out[protocol.EntityDataKeyName] != in[protocol.EntityDataKeyName] {
		t.Fatalf("Name round-trip mismatch: got %v, want %v", out[protocol.EntityDataKeyName], in[protocol.EntityDataKeyName])
	}
	if out[protocol.EntityDataKeyScale] != in[protocol.EntityDataKeyScale] {
		t.Fatalf("Scale round-trip mismatch: got %v, want %v", out[protocol.EntityDataKeyScale], in[protocol.EntityDataKeyScale])
	}
}
