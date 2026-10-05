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
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ItemStackResponse struct {
	packet.ItemStackResponse
}

func (*ItemStackResponse) ID() uint32 { return packet.IDItemStackResponse }

func (pk *ItemStackResponse) Marshal(io protocol.IO) {
	protocol.FuncSlice(io, &pk.Responses, func(r *protocol.ItemStackResponse) {
		itemStackResponse786(io, r)
	})
}

func itemStackResponse786(r protocol.IO, x *protocol.ItemStackResponse) {
	r.Uint8(&x.Status)
	r.Varint32(&x.RequestID)
	if x.Status == protocol.ItemStackResponseStatusOK {
		protocol.FuncSlice(r, &x.ContainerInfo, func(c *protocol.StackResponseContainerInfo) {
			stackRespContainer786(r, c)
		})
	}
}

func stackRespContainer786(r protocol.IO, x *protocol.StackResponseContainerInfo) {

	if p := ProtoOf(r); p != 0 && p < 729 {
		r.Uint8(&x.Container.ContainerID)
	} else {
		protocol.Single(r, &x.Container)
	}
	protocol.FuncSlice(r, &x.SlotInfo, func(s *protocol.StackResponseSlotInfo) {
		stackRespSlot786(r, s)
	})
}

func stackRespSlot786(r protocol.IO, x *protocol.StackResponseSlotInfo) {
	r.Uint8(&x.Slot)
	r.Uint8(&x.HotbarSlot)
	r.Uint8(&x.Count)
	r.Varint32(&x.StackNetworkID)
	r.String(&x.CustomName)

	if p := ProtoOf(r); p == 0 || p >= 766 {

		filtered, _ := x.FilteredCustomName.Value()
		r.String(&filtered)
		x.FilteredCustomName = protocol.Option(filtered)
	}
	r.Varint32(&x.DurabilityCorrection)
}

func ToLatestItemStackResponse(pk *ItemStackResponse) *packet.ItemStackResponse {
	out := pk.ItemStackResponse
	return &out
}
func FromLatestItemStackResponse(pk *packet.ItemStackResponse) *ItemStackResponse {
	return &ItemStackResponse{ItemStackResponse: *pk}
}
