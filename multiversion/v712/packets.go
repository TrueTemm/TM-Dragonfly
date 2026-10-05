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

package v712

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

type Attribute struct {
	Min, Max, Value, Default float32
	Name                     string
	Modifiers                []protocol.AttributeModifier
}

func (x *Attribute) Marshal(r protocol.IO) {
	r.Float32(&x.Min)
	r.Float32(&x.Max)
	r.Float32(&x.Value)
	r.Float32(&x.Default)
	r.String(&x.Name)
	protocol.Slice(r, &x.Modifiers)
}

type UpdateAttributes struct {
	EntityRuntimeID uint64
	Attributes      []Attribute
	Tick            uint64
}

func (*UpdateAttributes) ID() uint32 { return v786.IDUpdateAttributes }

func (pk *UpdateAttributes) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	protocol.Slice(io, &pk.Attributes)
	io.Varuint64(&pk.Tick)
}

func fromLatestUpdateAttributes(pk *packet.UpdateAttributes) *UpdateAttributes {
	out := &UpdateAttributes{EntityRuntimeID: pk.EntityRuntimeID, Attributes: make([]Attribute, len(pk.Attributes)), Tick: pk.Tick}
	for i, a := range pk.Attributes {
		out.Attributes[i] = Attribute{Min: a.Min, Max: a.Max, Value: a.Value, Default: a.Default, Name: a.Name, Modifiers: a.Modifiers}
	}
	return out
}

func toLatestUpdateAttributes(pk *UpdateAttributes) *packet.UpdateAttributes {
	out := &packet.UpdateAttributes{EntityRuntimeID: pk.EntityRuntimeID, Attributes: make([]protocol.Attribute, len(pk.Attributes)), Tick: pk.Tick}
	for i, a := range pk.Attributes {
		out.Attributes[i] = protocol.Attribute{AttributeValue: protocol.AttributeValue{Name: a.Name, Value: a.Value, Min: a.Min, Max: a.Max},
			DefaultMin: a.Min, DefaultMax: a.Max, Default: a.Default, Modifiers: a.Modifiers}
	}
	return out
}

type InventoryContent struct {
	WindowID        uint32
	Content         []protocol.ItemInstance
	DynamicWindowID uint32
}

func (*InventoryContent) ID() uint32 { return v786.IDInventoryContent }

func (pk *InventoryContent) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	protocol.FuncSlice(io, &pk.Content, io.ItemInstance)
	io.Varuint32(&pk.DynamicWindowID)
}

func fromLatestInventoryContent(pk *packet.InventoryContent) *InventoryContent {
	return &InventoryContent{WindowID: pk.WindowID, Content: pk.Content}
}

func toLatestInventoryContent(pk *InventoryContent) *packet.InventoryContent {
	return &packet.InventoryContent{WindowID: pk.WindowID, Content: pk.Content,
		Container: protocol.FullContainerName{ContainerID: byte(pk.WindowID)}}
}

type InventorySlot struct {
	WindowID        uint32
	Slot            uint32
	DynamicWindowID uint32
	NewItem         protocol.ItemInstance
}

func (*InventorySlot) ID() uint32 { return v786.IDInventorySlot }

func (pk *InventorySlot) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	io.Varuint32(&pk.Slot)
	io.Varuint32(&pk.DynamicWindowID)
	io.ItemInstance(&pk.NewItem)
}

func fromLatestInventorySlot(pk *packet.InventorySlot) *InventorySlot {
	return &InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, NewItem: pk.NewItem}
}

func toLatestInventorySlot(pk *InventorySlot) *packet.InventorySlot {
	return v786.ToLatestInventorySlot(&v786.InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot,
		Container: protocol.FullContainerName{ContainerID: byte(pk.WindowID)}, NewItem: pk.NewItem})
}

type Emote struct {
	EntityRuntimeID uint64
	EmoteID         string
	XUID            string
	PlatformID      string
	Flags           byte
}

func (*Emote) ID() uint32 { return v786.IDEmote }

func (pk *Emote) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.String(&pk.EmoteID)
	io.String(&pk.XUID)
	io.String(&pk.PlatformID)
	io.Uint8(&pk.Flags)
}

func fromLatestEmote(pk *packet.Emote) *Emote {
	return &Emote{EntityRuntimeID: pk.EntityRuntimeID, EmoteID: pk.EmoteID, XUID: pk.XUID, PlatformID: pk.PlatformID, Flags: pk.Flags}
}

func toLatestEmote(pk *Emote) *packet.Emote {
	return &packet.Emote{EntityRuntimeID: pk.EntityRuntimeID, EmoteID: pk.EmoteID, XUID: pk.XUID, PlatformID: pk.PlatformID, Flags: pk.Flags}
}

type Transfer struct {
	Address string
	Port    uint16
}

func (*Transfer) ID() uint32 { return v786.IDTransfer }

func (pk *Transfer) Marshal(io protocol.IO) {
	io.String(&pk.Address)
	io.Uint16(&pk.Port)
}

func fromLatestTransfer(pk *packet.Transfer) *Transfer {
	x := v786.FromLatestTransfer(pk)
	return &Transfer{Address: x.Address, Port: x.Port}
}

type BehaviourPackInfo struct {
	UUID            string
	Version         string
	Size            uint64
	ContentKey      string
	SubPackName     string
	ContentIdentity string
	HasScripts      bool
	AddonPack       bool
}

func (x *BehaviourPackInfo) Marshal(r protocol.IO) {
	r.String(&x.UUID)
	r.String(&x.Version)
	r.Uint64(&x.Size)
	r.String(&x.ContentKey)
	r.String(&x.SubPackName)
	r.String(&x.ContentIdentity)
	r.Bool(&x.HasScripts)
	r.Bool(&x.AddonPack)
}

type ResourcePacksInfo struct {
	TexturePackRequired bool
	HasAddons           bool
	HasScripts          bool
	ForcingServerPacks  bool
	BehaviourPacks      []BehaviourPackInfo
	TexturePacks        []v729.TexturePackInfo
	PackURLs            []v729.PackURL
}

func (*ResourcePacksInfo) ID() uint32 { return v786.IDResourcePacksInfo }

func (pk *ResourcePacksInfo) Marshal(io protocol.IO) {
	io.Bool(&pk.TexturePackRequired)
	io.Bool(&pk.HasAddons)
	io.Bool(&pk.HasScripts)
	io.Bool(&pk.ForcingServerPacks)
	v786.SliceUint16Length786(io, &pk.BehaviourPacks)
	v786.SliceUint16Length786(io, &pk.TexturePacks)
	protocol.Slice(io, &pk.PackURLs)
}

func fromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	x := v729.FromLatestResourcePacksInfo(pk)

	return &ResourcePacksInfo{TexturePackRequired: x.TexturePackRequired, HasAddons: x.HasAddons, HasScripts: x.HasScripts,
		BehaviourPacks: []BehaviourPackInfo{}, TexturePacks: x.TexturePacks, PackURLs: x.PackURLs}
}
