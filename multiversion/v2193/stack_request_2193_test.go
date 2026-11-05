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

package v2193

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestItemStackRequest2193RoundTrip(t *testing.T) {
	src := &protocol.TakeStackRequestAction{}
	src.Count = 1
	src.Source.Container.ContainerID = protocol.ContainerCreatedOutput
	src.Source.Container.DynamicContainerID = protocol.Option(uint32(0x01020304))
	src.Source.Slot = 7
	src.Source.StackNetworkID = -3
	src.Destination.Container.ContainerID = protocol.ContainerHotBar
	src.Destination.Slot = 3
	src.Destination.StackNetworkID = 22
	original := &packet.ItemStackRequest{Requests: []protocol.ItemStackRequest{{RequestID: -1, Actions: []protocol.StackRequestAction{src}}}}
	var encoded bytes.Buffer
	original.Marshal((Protocol{}).NewWriter(&encoded, 0))
	raw := append([]byte(nil), encoded.Bytes()...)
	if len(raw) < 15 {
		t.Fatalf("request too short: %x", raw)
	}
	var decoded packet.ItemStackRequest
	decoded.Marshal((Protocol{}).NewReader(bytes.NewBuffer(raw), 0, true))
	if len(decoded.Requests) != 1 || len(decoded.Requests[0].Actions) != 1 {
		t.Fatalf("roundtrip actions: %+v", decoded.Requests)
	}
	got, ok := decoded.Requests[0].Actions[0].(*protocol.TakeStackRequestAction)
	if !ok {
		t.Fatalf("action type %T", decoded.Requests[0].Actions[0])
	}
	if got.Source.Container.ContainerID != src.Source.Container.ContainerID ||
		got.Source.StackNetworkID != -3 || got.Destination.StackNetworkID != 22 {
		t.Fatalf("slot mismatch: %+v", got)
	}
	dyn, _ := got.Source.Container.DynamicContainerID.Value()
	if dyn != 0x01020304 {
		t.Fatalf("dynamic id = %x", dyn)
	}
}

func TestItemStackActionIDs2193(t *testing.T) {
	for id := uint8(0); id < 18; id++ {
		a, ok := makeAction2193(id)
		if !ok {
			t.Fatalf("missing action %d", id)
		}
		roundtrip, ok := legacyActionIndex2193(a)
		if !ok || roundtrip != id {
			t.Fatalf("action %d => %d", id, roundtrip)
		}
	}
	if _, ok := makeAction2193(18); ok {
		t.Fatal("reserved action cannot appear in 2193")
	}
}

func TestItemStackResponse2193RoundTrip(t *testing.T) {
	container := protocol.FullContainerName{ContainerID: protocol.ContainerHotBar, DynamicContainerID: protocol.Option(uint32(0x55667788))}
	orig := &ItemStackResponse2193{packet.ItemStackResponse{Responses: []protocol.ItemStackResponse{{
		Status: protocol.ItemStackResponseStatusOK, RequestID: 6,
		ContainerInfo: []protocol.StackResponseContainerInfo{{Container: container,
			SlotInfo: []protocol.StackResponseSlotInfo{{Slot: 3, HotbarSlot: 3, Count: 1, StackNetworkID: 7}}}},
	}}}}
	var buf bytes.Buffer
	orig.Marshal((Protocol{}).NewWriter(&buf, 0))
	var got ItemStackResponse2193
	got.Marshal((Protocol{}).NewReader(bytes.NewBuffer(buf.Bytes()), 0, true))
	if len(got.Responses) != 1 || len(got.Responses[0].ContainerInfo) != 1 {
		t.Fatalf("response %+v", got.Responses)
	}
	dynamic, _ := got.Responses[0].ContainerInfo[0].Container.DynamicContainerID.Value()
	if dynamic != 0x55667788 {
		t.Fatalf("response dynamic ID %x", dynamic)
	}
}
