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
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryContent2193 struct{ packet.InventoryContent }

func (*InventoryContent2193) ID() uint32 { return packet.IDInventoryContent }
func (p *InventoryContent2193) Marshal(io protocol.IO) {
	io.Varuint32(&p.WindowID)
	protocol.FuncSlice(io, &p.Content, io.ItemInstance)
	fullContainer2193(io, &p.Container)
	io.ItemInstance(&p.StorageItem)
}

type InventorySlot2193 struct{ packet.InventorySlot }

func (*InventorySlot2193) ID() uint32 { return packet.IDInventorySlot }
func (p *InventorySlot2193) Marshal(io protocol.IO) {
	io.Varuint32(&p.WindowID)
	io.Varuint32(&p.Slot)
	protocol.OptionalFunc(io, &p.Container, func(c *protocol.FullContainerName) { fullContainer2193(io, c) })
	protocol.OptionalFunc(io, &p.StorageItem, io.ItemInstance)
	io.ItemInstance(&p.NewItem)
}

type ItemStackResponse2193 struct{ packet.ItemStackResponse }

func (*ItemStackResponse2193) ID() uint32 { return packet.IDItemStackResponse }
func (p *ItemStackResponse2193) Marshal(io protocol.IO) {
	protocol.FuncSlice(io, &p.Responses, func(x *protocol.ItemStackResponse) {
		io.Uint8(&x.Status)
		io.Varint32(&x.RequestID)
		var containers protocol.Optional[[]protocol.StackResponseContainerInfo]
		if len(x.ContainerInfo) != 0 {
			containers = protocol.Option(x.ContainerInfo)
		}
		protocol.OptionalFunc(io, &containers, func(cs *[]protocol.StackResponseContainerInfo) {
			protocol.FuncSlice(io, cs, func(c *protocol.StackResponseContainerInfo) {
				fullContainer2193(io, &c.Container)
				protocol.Slice(io, &c.SlotInfo)
			})
		})
		if v, ok := containers.Value(); ok {
			x.ContainerInfo = v
		}
	})
}
