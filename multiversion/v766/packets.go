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

package v766

import (
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

type AbilityLayer struct {
	Type      uint16
	Abilities uint32
	Values    uint32
	FlySpeed  float32
	WalkSpeed float32
}

func (x *AbilityLayer) Marshal(r protocol.IO) {
	r.Uint16(&x.Type)
	r.Uint32(&x.Abilities)
	r.Uint32(&x.Values)
	r.Float32(&x.FlySpeed)
	r.Float32(&x.WalkSpeed)
}

type AbilityData struct {
	EntityUniqueID     int64
	PlayerPermissions  byte
	CommandPermissions byte
	Layers             []AbilityLayer
}

func (x *AbilityData) Marshal(r protocol.IO) {
	r.Int64(&x.EntityUniqueID)
	r.Uint8(&x.PlayerPermissions)
	r.Uint8(&x.CommandPermissions)
	protocol.SliceUint8Length(r, &x.Layers)
}

type UpdateAbilities struct {
	AbilityData AbilityData
}

func (*UpdateAbilities) ID() uint32 { return v786.IDUpdateAbilities }

func (pk *UpdateAbilities) Marshal(io protocol.IO) { protocol.Single(io, &pk.AbilityData) }

func fromLatestUpdateAbilities(pk *packet.UpdateAbilities) *UpdateAbilities {
	d := pk.AbilityData
	out := &UpdateAbilities{AbilityData: AbilityData{EntityUniqueID: d.EntityUniqueID,
		PlayerPermissions: d.PlayerPermissions, CommandPermissions: d.CommandPermissions,
		Layers: make([]AbilityLayer, len(d.Layers))}}
	for i, l := range d.Layers {
		out.AbilityData.Layers[i] = AbilityLayer{Type: l.Type, Abilities: l.Abilities, Values: l.Values,
			FlySpeed: l.FlySpeed, WalkSpeed: l.WalkSpeed}
	}
	return out
}

func toLatestUpdateAbilities(pk *UpdateAbilities) *packet.UpdateAbilities {
	d := pk.AbilityData
	out := &packet.UpdateAbilities{AbilityData: protocol.AbilityData{EntityUniqueID: d.EntityUniqueID,
		PlayerPermissions: d.PlayerPermissions, CommandPermissions: d.CommandPermissions,
		Layers: make([]protocol.AbilityLayer, len(d.Layers))}}
	for i, l := range d.Layers {

		out.AbilityData.Layers[i] = protocol.AbilityLayer{Type: l.Type, Abilities: l.Abilities, Values: l.Values,
			FlySpeed: l.FlySpeed, VerticalFlySpeed: protocol.AbilityBaseVerticalFlySpeed, WalkSpeed: l.WalkSpeed}
	}
	return out
}

type BossEvent struct {
	BossEntityUniqueID int64
	EventType          uint32
	PlayerUniqueID     int64
	BossBarTitle       string
	HealthPercentage   float32
	ScreenDarkening    uint16
	Colour             uint32
	Overlay            uint32
}

func (*BossEvent) ID() uint32 { return v786.IDBossEvent }

func (pk *BossEvent) Marshal(io protocol.IO) {
	io.Varint64(&pk.BossEntityUniqueID)
	io.Varuint32(&pk.EventType)
	switch pk.EventType {
	case packet.BossEventShow:
		io.String(&pk.BossBarTitle)
		io.Float32(&pk.HealthPercentage)
		io.Uint16(&pk.ScreenDarkening)
		io.Varuint32(&pk.Colour)
		io.Varuint32(&pk.Overlay)
	case packet.BossEventRegisterPlayer, packet.BossEventUnregisterPlayer, packet.BossEventRequest:
		io.Varint64(&pk.PlayerUniqueID)
	case packet.BossEventHide:
	case packet.BossEventHealthPercentage:
		io.Float32(&pk.HealthPercentage)
	case packet.BossEventTitle:
		io.String(&pk.BossBarTitle)
	case packet.BossEventAppearanceProperties:
		io.Uint16(&pk.ScreenDarkening)
		io.Varuint32(&pk.Colour)
		io.Varuint32(&pk.Overlay)
	case packet.BossEventTexture:
		io.Varuint32(&pk.Colour)
		io.Varuint32(&pk.Overlay)
	default:
		io.UnknownEnumOption(pk.EventType, "boss event type")
	}
}

func fromLatestBossEvent(pk *packet.BossEvent) *BossEvent {
	x := v786.FromLatestBossEvent(pk)
	return &BossEvent{BossEntityUniqueID: x.BossEntityUniqueID, EventType: x.EventType, PlayerUniqueID: x.PlayerUniqueID,
		BossBarTitle: x.BossBarTitle, HealthPercentage: x.HealthPercentage, ScreenDarkening: x.ScreenDarkening,
		Colour: x.Colour, Overlay: x.Overlay}
}

func toLatestBossEvent(pk *BossEvent) *packet.BossEvent {
	return v786.ToLatestBossEvent(&v786.BossEvent{BossEntityUniqueID: pk.BossEntityUniqueID, EventType: pk.EventType,
		PlayerUniqueID: pk.PlayerUniqueID, BossBarTitle: pk.BossBarTitle, HealthPercentage: pk.HealthPercentage,
		ScreenDarkening: pk.ScreenDarkening, Colour: pk.Colour, Overlay: pk.Overlay})
}

type CommandBlockUpdate struct {
	Block                   bool
	Position                protocol.BlockPos
	Mode                    uint32
	NeedsRedstone           bool
	Conditional             bool
	MinecartEntityRuntimeID uint64
	Command                 string
	LastOutput              string
	Name                    string
	ShouldTrackOutput       bool
	TickDelay               int32
	ExecuteOnFirstTick      bool
}

func (*CommandBlockUpdate) ID() uint32 { return v786.IDCommandBlockUpdate }

func (pk *CommandBlockUpdate) Marshal(io protocol.IO) {
	io.Bool(&pk.Block)
	if pk.Block {
		v786.UBlockPos786(io, &pk.Position)
		io.Varuint32(&pk.Mode)
		io.Bool(&pk.NeedsRedstone)
		io.Bool(&pk.Conditional)
	} else {
		io.Varuint64(&pk.MinecartEntityRuntimeID)
	}
	io.String(&pk.Command)
	io.String(&pk.LastOutput)
	io.String(&pk.Name)
	io.Bool(&pk.ShouldTrackOutput)
	io.Int32(&pk.TickDelay)
	io.Bool(&pk.ExecuteOnFirstTick)
}

func toLatestCommandBlockUpdate(pk *CommandBlockUpdate) *packet.CommandBlockUpdate {
	return v786.ToLatestCommandBlockUpdate(&v786.CommandBlockUpdate{Block: pk.Block, Position: pk.Position, Mode: pk.Mode,
		NeedsRedstone: pk.NeedsRedstone, Conditional: pk.Conditional, MinecartEntityRuntimeID: pk.MinecartEntityRuntimeID,
		Command: pk.Command, LastOutput: pk.LastOutput, Name: pk.Name, ShouldTrackOutput: pk.ShouldTrackOutput,
		TickDelay: pk.TickDelay, ExecuteOnFirstTick: pk.ExecuteOnFirstTick})
}

func fromLatestCommandBlockUpdate(pk *packet.CommandBlockUpdate) *CommandBlockUpdate {
	x := v786.FromLatestCommandBlockUpdate(pk)
	return &CommandBlockUpdate{Block: x.Block, Position: x.Position, Mode: x.Mode, NeedsRedstone: x.NeedsRedstone,
		Conditional: x.Conditional, MinecartEntityRuntimeID: x.MinecartEntityRuntimeID, Command: x.Command,
		LastOutput: x.LastOutput, Name: x.Name, ShouldTrackOutput: x.ShouldTrackOutput, TickDelay: x.TickDelay,
		ExecuteOnFirstTick: x.ExecuteOnFirstTick}
}

type StructureBlockUpdate struct {
	Position           protocol.BlockPos
	StructureName      string
	DataField          string
	IncludePlayers     bool
	ShowBoundingBox    bool
	StructureBlockType int32
	Settings           protocol.StructureSettings
	RedstoneSaveMode   int32
	ShouldTrigger      bool
	Waterlogged        bool
}

func (*StructureBlockUpdate) ID() uint32 { return v786.IDStructureBlockUpdate }

func (pk *StructureBlockUpdate) Marshal(io protocol.IO) {
	v786.UBlockPos786(io, &pk.Position)
	io.String(&pk.StructureName)
	io.String(&pk.DataField)
	io.Bool(&pk.IncludePlayers)
	io.Bool(&pk.ShowBoundingBox)
	io.Varint32(&pk.StructureBlockType)
	protocol.Single(io, &pk.Settings)
	io.Varint32(&pk.RedstoneSaveMode)
	io.Bool(&pk.ShouldTrigger)
	io.Bool(&pk.Waterlogged)
}

func toLatestStructureBlockUpdate(pk *StructureBlockUpdate) *packet.StructureBlockUpdate {
	return v786.ToLatestStructureBlockUpdate(&v786.StructureBlockUpdate{Position: pk.Position, StructureName: pk.StructureName,
		DataField: pk.DataField, IncludePlayers: pk.IncludePlayers, ShowBoundingBox: pk.ShowBoundingBox,
		StructureBlockType: pk.StructureBlockType, Settings: pk.Settings, RedstoneSaveMode: pk.RedstoneSaveMode,
		ShouldTrigger: pk.ShouldTrigger, Waterlogged: pk.Waterlogged})
}

func fromLatestStructureBlockUpdate(pk *packet.StructureBlockUpdate) *StructureBlockUpdate {
	x := v786.FromLatestStructureBlockUpdate(pk)
	return &StructureBlockUpdate{Position: x.Position, StructureName: x.StructureName, DataField: x.DataField,
		IncludePlayers: x.IncludePlayers, ShowBoundingBox: x.ShowBoundingBox, StructureBlockType: x.StructureBlockType,
		Settings: x.Settings, RedstoneSaveMode: x.RedstoneSaveMode, ShouldTrigger: x.ShouldTrigger, Waterlogged: x.Waterlogged}
}

type CreativeItem struct {
	CreativeItemNetworkID uint32
	Item                  protocol.ItemStack
}

func (x *CreativeItem) Marshal(r protocol.IO) {
	r.Varuint32(&x.CreativeItemNetworkID)
	r.Item(&x.Item)
}

type CreativeContent struct {
	Items []CreativeItem
}

func (*CreativeContent) ID() uint32 { return v786.IDCreativeContent }

func (pk *CreativeContent) Marshal(io protocol.IO) { protocol.Slice(io, &pk.Items) }

func fromLatestCreativeContent(pk *packet.CreativeContent) *CreativeContent {
	x := v786.FromLatestCreativeContent786(pk)
	out := &CreativeContent{Items: make([]CreativeItem, len(x.Items))}
	for i, it := range x.Items {
		out.Items[i] = CreativeItem{CreativeItemNetworkID: it.CreativeItemNetworkID, Item: it.Item}
	}
	return out
}

func toLatestCreativeContent(pk *CreativeContent) *packet.CreativeContent {
	x := &v786.CreativeContent{Items: make([]protocol.CreativeItem, len(pk.Items))}
	for i, it := range pk.Items {
		x.Items[i] = protocol.CreativeItem{CreativeItemNetworkID: it.CreativeItemNetworkID, Item: it.Item}
	}
	return v786.ToLatestCreativeContent786(x)
}

type ItemComponentEntry struct {
	Name string
	Data map[string]any
}

func (x *ItemComponentEntry) Marshal(r protocol.IO) {
	r.String(&x.Name)
	r.NBT(&x.Data, nbt.NetworkLittleEndian)
}

type ItemComponent struct {
	Items []ItemComponentEntry
}

func (*ItemComponent) ID() uint32 { return v786.IDItemRegistry }

func (pk *ItemComponent) Marshal(io protocol.IO) { protocol.Slice(io, &pk.Items) }
