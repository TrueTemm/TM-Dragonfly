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

package conformance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type stackFixtureSet struct {
	Revision   string            `json:"revision"`
	StoneSword int32             `json:"stone_sword"`
	Fixtures   map[string][]byte `json:"fixtures"`
	Inventory  map[string][]byte `json:"inventory"`
}

func TestHistoricalItemStackActions(t *testing.T) {
	data, err := os.ReadFile("testdata/item_stack.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]stackFixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, p := range protocols() {
		f, ok := fixtures[fmt.Sprint(p.ID())]
		if !ok {
			t.Fatalf("protocol %d has no historical fixtures", p.ID())
		}
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			if len(f.Revision) != 40 || f.StoneSword == 0 {
				t.Fatal("historical provenance missing")
			}
			for name, raw := range f.Fixtures {
				t.Run(name, func(t *testing.T) {
					defer func() {
						if err := recover(); err != nil {
							t.Fatalf("historical decoder panicked: %v", err)
						}
					}()
					buf := bytes.NewBuffer(raw)
					pk := p.Packets(true)[packet.IDItemStackRequest]()
					pk.Marshal(p.NewReader(buf, 0, true))
					if buf.Len() != 0 {
						t.Fatalf("%d unread bytes: %x", buf.Len(), buf.Bytes())
					}
					out := p.ConvertToLatest(pk, nil)
					req := out[0].(*packet.ItemStackRequest).Requests[0]
					if req.RequestID != -3 || req.FilterCause != -1 {
						t.Fatalf("request changed: %+v", req)
					}
					if name == "results" || name == "creative_take" {
						i := 0
						if name == "creative_take" {
							i = 1
							if len(req.Actions) != 3 {
								t.Fatal("creative action count changed")
							}
						}
						result, ok := req.Actions[i].(*protocol.CraftResultsDeprecatedStackRequestAction)
						if !ok || len(result.ResultItems) != 1 || result.ResultItems[0].Identifier != "minecraft:stone_sword" || result.ResultItems[0].Count != 1 || result.TimesCrafted != 1 {
							t.Fatalf("stone sword result changed: %#v", req.Actions[i])
						}
					}
					if name == "creative_take" {
						take, ok := req.Actions[2].(*protocol.TakeStackRequestAction)
						if !ok || take.Count != 1 || take.Source.StackNetworkID != -3 || take.Source.Slot != 50 || take.Destination.Slot != 0 {
							t.Fatalf("creative transfer changed: %#v", req.Actions[2])
						}
					}
					encoded := new(bytes.Buffer)
					pk.Marshal(p.NewWriter(encoded, 0))
					if !bytes.Equal(encoded.Bytes(), raw) {
						t.Fatalf("historical bytes changed: got %x want %x", encoded.Bytes(), raw)
					}
					if strings.HasPrefix(name, "creative") && p.ID() < 712 {
						if req.Actions[0].(*protocol.CraftCreativeStackRequestAction).NumberOfCrafts != 1 {
							t.Fatal("legacy creative count must be one")
						}
					}
				})
			}
		})
	}
}

func TestHistoricalInventoryContainers(t *testing.T) {
	data, err := os.ReadFile("testdata/item_stack.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]stackFixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, p := range protocols() {
		if p.ID() < 2168 {
			continue
		}
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			for name, raw := range fixtures[fmt.Sprint(p.ID())].Inventory {
				t.Run(name, func(t *testing.T) {
					defer func() {
						if err := recover(); err != nil {
							t.Fatalf("inventory decoder panicked: %v", err)
						}
					}()
					c := protocol.FullContainerName{ContainerID: protocol.ContainerHotBar}
					if strings.HasSuffix(name, "dynamic") {
						c.DynamicContainerID = protocol.Option(uint32(0x01020304))
					}
					var source packet.Packet
					switch {
					case strings.HasPrefix(name, "content_"):
						source = &packet.InventoryContent{WindowID: 0, Container: c}
					case strings.HasPrefix(name, "slot_"):
						source = &packet.InventorySlot{WindowID: 0, Slot: 3, Container: protocol.Option(c)}
					case strings.HasPrefix(name, "response_"):
						source = &packet.ItemStackResponse{Responses: []protocol.ItemStackResponse{{Status: protocol.ItemStackResponseStatusOK, RequestID: -3, ContainerInfo: []protocol.StackResponseContainerInfo{{Container: c, SlotInfo: []protocol.StackResponseSlotInfo{{Slot: 3, HotbarSlot: 3, Count: 1, StackNetworkID: 27, CustomName: "blade", DurabilityCorrection: 2}}}}}}}
					default:
						t.Fatal("unknown inventory fixture")
					}
					out := p.ConvertFromLatest(source, nil)
					if len(out) != 1 {
						t.Fatal("inventory packet dropped")
					}
					encoded := new(bytes.Buffer)
					out[0].Marshal(p.NewWriter(encoded, 0))
					if !bytes.Equal(encoded.Bytes(), raw) {
						t.Fatalf("inventory bytes changed: got %x want %x", encoded.Bytes(), raw)
					}
					buf := bytes.NewBuffer(raw)
					decoded := p.Packets(false)[source.ID()]()
					decoded.Marshal(p.NewReader(buf, 0, true))
					if buf.Len() != 0 {
						t.Fatalf("%d inventory bytes unread", buf.Len())
					}
				})
			}
		})
	}
}
