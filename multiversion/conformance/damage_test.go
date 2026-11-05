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
	"math"
	"os"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestHistoricalDamagePackets(t *testing.T) {
	data, err := os.ReadFile("testdata/damage.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Revision string            `json:"revision"`
		Packets  map[string][]byte `json:"packets"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, p := range protocols() {
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			f, ok := fixtures[fmt.Sprint(p.ID())]
			if !ok || len(f.Revision) != 40 {
				t.Fatal("historical damage fixtures missing")
			}
			sources := map[string]packet.Packet{
				"hurt":  &packet.ActorEvent{EntityRuntimeID: 1, EventType: packet.ActorEventHurt},
				"death": &packet.ActorEvent{EntityRuntimeID: 1, EventType: packet.ActorEventDeath},
			}
			for name, value := range map[string]float32{"health_full": 20, "health_hurt": 16, "health_dead": 0} {
				sources[name] = &packet.UpdateAttributes{EntityRuntimeID: 1, Attributes: []protocol.Attribute{
					{AttributeValue: protocol.AttributeValue{Name: "minecraft:health", Value: value, Max: 20}, DefaultMax: 20, Default: 20},
					{AttributeValue: protocol.AttributeValue{Name: "minecraft:absorption", Max: math.MaxFloat32}, DefaultMax: math.MaxFloat32},
				}}
			}
			for name, source := range sources {
				t.Run(name, func(t *testing.T) {
					out := p.ConvertFromLatest(source, nil)
					if len(out) != 1 {
						t.Fatal("damage packet dropped")
					}
					encoded := new(bytes.Buffer)
					out[0].Marshal(p.NewWriter(encoded, 0))
					if !bytes.Equal(encoded.Bytes(), f.Packets[name]) {
						t.Fatalf("damage wire changed: got %x want %x", encoded.Bytes(), f.Packets[name])
					}
					decoded := p.Packets(false)[source.ID()]()
					decoded.Marshal(p.NewReader(encoded, 0, true))
					if encoded.Len() != 0 {
						t.Fatalf("%d damage bytes unread", encoded.Len())
					}
					back := p.ConvertToLatest(decoded, nil)
					if len(back) != 1 {
						t.Fatal("damage packet not restored")
					}
					switch got := back[0].(type) {
					case *packet.ActorEvent:
						if got.EventType != source.(*packet.ActorEvent).EventType || got.EntityRuntimeID != 1 {
							t.Fatal("damage event changed")
						}
					case *packet.UpdateAttributes:
						if len(got.Attributes) != 2 || got.Attributes[0].Name != "minecraft:health" || got.Attributes[0].Value != source.(*packet.UpdateAttributes).Attributes[0].Value || got.Attributes[0].Max != 20 {
							t.Fatalf("health attribute changed: %+v", got)
						}
					default:
						t.Fatalf("unexpected damage packet %T", got)
					}
				})
			}
		})
	}
}
