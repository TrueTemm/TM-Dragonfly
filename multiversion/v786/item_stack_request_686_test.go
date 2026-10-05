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
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestItemStackRequestCreativeTake686(t *testing.T) {
	req := protocol.ItemStackRequest{RequestID: -3, Actions: []protocol.StackRequestAction{
		&protocol.CraftCreativeStackRequestAction{CreativeItemNetworkID: 391, NumberOfCrafts: 1}, nil,
	}, FilterStrings: nil, FilterCause: -1}
	take := &protocol.TakeStackRequestAction{}
	take.Count = 64
	take.Source = protocol.StackRequestSlotInfo{Container: protocol.FullContainerName{ContainerID: protocol.ContainerCreatedOutput}, Slot: 50, StackNetworkID: -3}
	take.Destination = protocol.StackRequestSlotInfo{Container: protocol.FullContainerName{ContainerID: protocol.ContainerHotBar}, Slot: 3}
	req.Actions[1] = take
	for _, proto := range []uint32{686, 712, 786} {
		buf := bytes.NewBuffer(nil)
		w := Writer786{Writer: protocol.NewWriter(buf, 0), Proto: proto}
		pk := &ItemStackRequest{ItemStackRequest: packet.ItemStackRequest{Requests: []protocol.ItemStackRequest{req}}}
		pk.Marshal(w)
		raw := buf.Bytes()
		r := Reader786{Reader: protocol.NewReader(bytes.NewBuffer(raw), 0, false), Proto: proto}
		var back ItemStackRequest
		back.Marshal(r)
		if len(back.Requests) != 1 || len(back.Requests[0].Actions) != 2 {
			t.Fatalf("%d: decoded %+v", proto, back)
		}
		cc := back.Requests[0].Actions[0].(*protocol.CraftCreativeStackRequestAction)
		if cc.CreativeItemNetworkID != 391 || (proto >= 712 && cc.NumberOfCrafts != 1) {
			t.Fatalf("%d: craft creative decoded as %+v", proto, cc)
		}
		if proto == 686 && len(raw) != len(bytesFor(t, 712, req))-9 {
			t.Errorf("686 should be 9 bytes shorter than 712 (NumberOfCrafts + two Uint32 dynamic container ids): 686=%d 712=%d", len(raw), len(bytesFor(t, 712, req)))
		}
		if dir := os.Getenv("BLINDVERIFY_DUMP"); dir != "" && proto == 686 {
			line := fmt.Sprintf("ItemStackRequest\t%d\tok\t%s\n", packet.IDItemStackRequest, hex.EncodeToString(raw))
			_ = os.WriteFile(dir+"/send_686_isr.txt", []byte(line), 0o644)
		}
	}
}

func bytesFor(t *testing.T, proto uint32, req protocol.ItemStackRequest) []byte {
	buf := bytes.NewBuffer(nil)
	(&ItemStackRequest{ItemStackRequest: packet.ItemStackRequest{Requests: []protocol.ItemStackRequest{req}}}).Marshal(Writer786{Writer: protocol.NewWriter(buf, 0), Proto: proto})
	return buf.Bytes()
}
