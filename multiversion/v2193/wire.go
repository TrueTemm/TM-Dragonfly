package v2193

import (
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type skin struct{ protocol.Skin }

func (x *skin) Marshal(r protocol.IO) {
	playFabID := ""
	r.String(&x.SkinID)
	r.String(&playFabID)
	r.ByteSlice(&x.SkinResourcePatch)
	r.Uint32(&x.SkinImageWidth)
	r.Uint32(&x.SkinImageHeight)
	r.ByteSlice(&x.SkinData)
	protocol.Slice(r, &x.Animations)
	r.Uint32(&x.CapeImageWidth)
	r.Uint32(&x.CapeImageHeight)
	r.ByteSlice(&x.CapeData)
	r.ByteSlice(&x.SkinGeometry)
	r.ByteSlice(&x.GeometryDataEngineVersion)
	r.ByteSlice(&x.AnimationData)
	r.String(&x.CapeID)
	r.String(&x.FullID)
	r.Uint8(&x.ArmSize)
	r.BEARGB(&x.SkinColour)
	protocol.Slice(r, &x.PersonaPieces)
	protocol.Slice(r, &x.PieceTintColours)
	r.Bool(&x.PremiumSkin)
	r.Bool(&x.PersonaSkin)
	r.Bool(&x.PersonaCapeOnClassicSkin)
	r.Bool(&x.PrimaryUser)
	r.Bool(&x.OverrideAppearance)
	trusted := "false"
	if x.Trusted {
		trusted = "true"
	}
	r.String(&trusted)
	x.Trusted = strings.EqualFold(trusted, "true")
	r.String(&x.ProfileHash)
}

type PlayerSkin struct{ packet.PlayerSkin }

func (pk *PlayerSkin) Marshal(io protocol.IO) {
	io.UUID(&pk.UUID)
	s := skin{pk.Skin}
	protocol.Single(io, &s)
	pk.Skin = s.Skin
	io.String(&pk.NewSkinName)
	io.String(&pk.OldSkinName)
}

type listEntry struct{ protocol.PlayerListEntry }

func (x *listEntry) Marshal(r protocol.IO) {
	variant := uint32(0)
	if x.ActionType == protocol.PlayerListActionAdd {
		variant = 1
	}
	r.Varuint32(&variant)
	legacyAction := x.ActionType
	r.Uint8(&legacyAction)
	x.ActionType = protocol.PlayerListActionRemove
	if variant == 1 {
		x.ActionType = protocol.PlayerListActionAdd
	} else if variant != 0 {
		r.UnknownEnumOption(variant, "player list entry variant")
	}
	r.UUID(&x.UUID)
	if x.ActionType == protocol.PlayerListActionRemove {
		return
	}
	r.ActorUniqueID(&x.EntityUniqueID)
	r.String(&x.Username)
	r.String(&x.XUID)
	r.String(&x.PlatformChatID)
	r.Int32(&x.BuildPlatform)
	s := skin{x.Skin}
	protocol.Single(r, &s)
	x.Skin = s.Skin
	r.Bool(&x.Teacher)
	r.Bool(&x.Host)
	r.Bool(&x.SubClient)
	r.BEARGB(&x.PlayerColour)
}

type PlayerList struct{ packet.PlayerList }

func (pk *PlayerList) Marshal(io protocol.IO) {
	entries := make([]listEntry, len(pk.Entries))
	for i, e := range pk.Entries {
		entries[i] = listEntry{e}
	}
	protocol.Slice(io, &entries)
	pk.Entries = make([]protocol.PlayerListEntry, len(entries))
	for i, e := range entries {
		pk.Entries[i] = e.PlayerListEntry
	}
}

type dimension struct{ protocol.DimensionDefinition }

func (x *dimension) Marshal(r protocol.IO) {
	r.String(&x.Name)
	r.Varint32(&x.MinimumY)
	r.Varint32(&x.HeightRange)
	r.Varint32(&x.Generator)
	r.Varint32(&x.DimensionType)
	r.UUID(&x.PackID)
	r.String(&x.DefaultBiome)
}

type DimensionData struct{ packet.DimensionData }

func (pk *DimensionData) Marshal(io protocol.IO) {
	defs := make([]dimension, len(pk.Definitions))
	for i, d := range pk.Definitions {
		defs[i] = dimension{d}
	}
	protocol.Slice(io, &defs)
	pk.Definitions = make([]protocol.DimensionDefinition, len(defs))
	for i, d := range defs {
		pk.Definitions[i] = d.DimensionDefinition
	}
}

type InventoryTransaction struct{ packet.InventoryTransaction }

func (pk *InventoryTransaction) Marshal(io protocol.IO) {
	io.Varint32(&pk.LegacyRequestID)
	hasLegacy := pk.LegacyRequestID < -1 && (pk.LegacyRequestID&1) == 0
	io.Bool(&hasLegacy)
	if hasLegacy {
		protocol.Slice(io, &pk.LegacySetItemSlots)
	}
	io.TransactionDataType(&pk.TransactionData)
	protocol.Slice(io, &pk.Actions)
	TransactionBody(io, pk.TransactionData)
}

func TransactionBody(io protocol.IO, data protocol.InventoryTransactionData) {
	switch d := data.(type) {
	case *protocol.UseItemOnEntityTransactionData:
		io.ActorRuntimeID(&d.TargetEntityRuntimeID)
		io.Varint32(&d.ActionType)
		io.Varint32(&d.HotBarSlot)
		io.ItemInstance(&d.HeldItem)
		io.Vec3(&d.Position)
		io.Vec3(&d.ClickedPosition)
	case *protocol.ReleaseItemTransactionData:
		io.Varint32(&d.ActionType)
		io.Varint32(&d.HotBarSlot)
		io.ItemInstance(&d.HeldItem)
		io.Vec3(&d.HeadPosition)
	default:
		data.Marshal(io)
	}
}
