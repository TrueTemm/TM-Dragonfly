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
	"encoding/hex"
	"testing"

	"github.com/df-mc/dragonfly/multiversion/v819"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestReportedCreative819(t *testing.T) {
	for _, tc := range []struct {
		name  string
		count uint16
		wire  string
	}{
		{"minecraft:white_wool", 64, "0105030e870301130146400000000a000000000000000000000101403c0032050c00000000ffffffff"},
		{"minecraft:snowball", 16, "0105030e87030113019e06100000000a000000000000000000000101103c0032050c00000000ffffffff"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := hex.DecodeString(tc.wire)
			if err != nil {
				t.Fatal(err)
			}
			p := v819.Protocol{}
			buf := bytes.NewBuffer(raw)
			pk := p.Packets(true)[packet.IDItemStackRequest]()
			pk.Marshal(p.NewReader(buf, 0, true))
			if buf.Len() != 0 {
				t.Fatalf("%d unread bytes %x", buf.Len(), buf.Bytes())
			}
			req := p.ConvertToLatest(pk, nil)[0].(*packet.ItemStackRequest).Requests[0]
			result := req.Actions[1].(*protocol.CraftResultsDeprecatedStackRequestAction)
			place := req.Actions[2].(*protocol.PlaceStackRequestAction)
			if len(result.ResultItems) != 1 || result.ResultItems[0].Identifier != tc.name || result.ResultItems[0].Count != tc.count || place.Count != byte(tc.count) {
				t.Fatalf("creative stack changed %+v %+v", result, place)
			}
			if req.RequestID != -3 || req.FilterCause != -1 || place.Destination.Container.ContainerID != 12 {
				t.Fatal("request fields changed")
			}
			encoded := new(bytes.Buffer)
			pk.Marshal(p.NewWriter(encoded, 0))
			if !bytes.Equal(encoded.Bytes(), raw) {
				t.Fatalf("wire changed %x", encoded.Bytes())
			}
		})
	}
}
