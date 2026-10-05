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

package v2168

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v2169"
)

type StackResponseSlotInfo struct {
	Slot, HotbarSlot     byte
	Count                byte
	StackNetworkID       int32
	CustomName           string
	FilteredCustomName   string
	DurabilityCorrection int32
}

func (x *StackResponseSlotInfo) Marshal(r protocol.IO) {
	r.Uint8(&x.Slot)
	r.Uint8(&x.HotbarSlot)
	r.Uint8(&x.Count)
	var stackNetworkID protocol.Optional[int32]
	if x.StackNetworkID > 0 {
		stackNetworkID = protocol.Option(x.StackNetworkID)
	}
	v2169.DoubleOptionalFunc(r, &stackNetworkID, r.Varint32)
	if value, ok := stackNetworkID.Value(); ok {
		x.StackNetworkID = value
	}
	r.String(&x.CustomName)
	r.String(&x.FilteredCustomName)
	r.Varint32(&x.DurabilityCorrection)
	if x.DurabilityCorrection < -32768 || x.DurabilityCorrection > 32767 {
		r.InvalidValue(x.DurabilityCorrection, "durability correction", "must fit in an int16")
	}
}

type StackResponseContainerInfo struct {
	Container protocol.FullContainerName
	SlotInfo  []StackResponseSlotInfo
}

func (x *StackResponseContainerInfo) Marshal(r protocol.IO) {
	protocol.Single(r, &x.Container)
	protocol.Slice(r, &x.SlotInfo)
}

type ItemStackResponseEntry struct {
	Status        uint8
	RequestID     int32
	ContainerInfo []StackResponseContainerInfo
}

func (x *ItemStackResponseEntry) Marshal(r protocol.IO) {
	r.Uint8(&x.Status)
	r.Varint32(&x.RequestID)
	var containerInfo protocol.Optional[[]StackResponseContainerInfo]
	if len(x.ContainerInfo) != 0 {
		containerInfo = protocol.Option(x.ContainerInfo)
	}
	v2169.DoubleOptionalFunc(r, &containerInfo, func(containerInfo *[]StackResponseContainerInfo) {
		protocol.Slice(r, containerInfo)
	})
	if value, ok := containerInfo.Value(); ok {
		x.ContainerInfo = value
	}
}

type ItemStackResponse struct {
	Responses []ItemStackResponseEntry
}

func (*ItemStackResponse) ID() uint32 { return packet.IDItemStackResponse }

func (pk *ItemStackResponse) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Responses)
}

func fromLatestItemStackResponse(pk *packet.ItemStackResponse) *ItemStackResponse {
	out := &ItemStackResponse{Responses: make([]ItemStackResponseEntry, len(pk.Responses))}
	for i, resp := range pk.Responses {
		entry := ItemStackResponseEntry{
			Status:        resp.Status,
			RequestID:     resp.RequestID,
			ContainerInfo: make([]StackResponseContainerInfo, len(resp.ContainerInfo)),
		}
		for j, info := range resp.ContainerInfo {
			container := StackResponseContainerInfo{
				Container: info.Container,
				SlotInfo:  make([]StackResponseSlotInfo, len(info.SlotInfo)),
			}
			for k, slot := range info.SlotInfo {
				filtered, _ := slot.FilteredCustomName.Value()
				container.SlotInfo[k] = StackResponseSlotInfo{
					Slot: slot.Slot, HotbarSlot: slot.HotbarSlot, Count: slot.Count,
					StackNetworkID: slot.StackNetworkID, CustomName: slot.CustomName,
					FilteredCustomName: filtered, DurabilityCorrection: slot.DurabilityCorrection,
				}
			}
			entry.ContainerInfo[j] = container
		}
		out.Responses[i] = entry
	}
	return out
}

func toLatestItemStackResponse(pk *ItemStackResponse) *packet.ItemStackResponse {
	out := &packet.ItemStackResponse{Responses: make([]protocol.ItemStackResponse, len(pk.Responses))}
	for i, resp := range pk.Responses {
		entry := protocol.ItemStackResponse{
			Status:        resp.Status,
			RequestID:     resp.RequestID,
			ContainerInfo: make([]protocol.StackResponseContainerInfo, len(resp.ContainerInfo)),
		}
		for j, info := range resp.ContainerInfo {
			container := protocol.StackResponseContainerInfo{
				Container: info.Container,
				SlotInfo:  make([]protocol.StackResponseSlotInfo, len(info.SlotInfo)),
			}
			for k, slot := range info.SlotInfo {
				var filtered protocol.Optional[string]
				if slot.FilteredCustomName != "" {
					filtered = protocol.Option(slot.FilteredCustomName)
				}
				container.SlotInfo[k] = protocol.StackResponseSlotInfo{
					Slot: slot.Slot, HotbarSlot: slot.HotbarSlot, Count: slot.Count,
					StackNetworkID: slot.StackNetworkID, CustomName: slot.CustomName,
					FilteredCustomName: filtered, DurabilityCorrection: slot.DurabilityCorrection,
				}
			}
			entry.ContainerInfo[j] = container
		}
		out.Responses[i] = entry
	}
	return out
}
