package v2193

import (
	"github.com/df-mc/dragonfly/multiversion/native"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const Proto = 2193

type Protocol struct{}

func (Protocol) ID() int32 { return Proto }

func (Protocol) Ver() string { return "1.26.50" }

func (Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return NewClientPool()
	}
	return NewServerPool()
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return Reader2193{Reader: protocol.NewReader(r, shieldID, enableLimits)}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return Writer2193{Writer: protocol.NewWriter(w, shieldID)}
}

func (Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return ToLatestShared(Proto, pk)
}

func (Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return FromLatestShared(Proto, pk)
}

func FromLatestShared(proto uint32, pk packet.Packet) []packet.Packet {
	pk = native.FromLatestFor(proto, pk)[0]
	if sg, ok := pk.(*packet.StartGame); ok {
		out := *sg
		out.BaseGameVersion, out.GameVersion = "1.26.50", "1.26.50"
		pk = &out
	}
	if out, ok := FromLatest(pk); ok {
		return []packet.Packet{out}
	}
	return []packet.Packet{pk}
}

func ToLatestShared(proto uint32, pk packet.Packet) []packet.Packet {
	if out, ok := ToLatest(pk); ok {
		pk = out
	}
	return native.ToLatestFor(proto, pk)
}

func NewClientPool() packet.Pool {
	p := packet.NewClientPool()
	ApplyDeltas(p)
	return p
}

func NewServerPool() packet.Pool {
	p := packet.NewServerPool()
	ApplyDeltas(p)
	return p
}

func ApplyDeltas(p packet.Pool) {
	p[packet.IDAddActor] = func() packet.Packet { return &AddActor{} }
	p[packet.IDAddPlayer] = func() packet.Packet { return &AddPlayer{} }
	p[packet.IDAnimate] = func() packet.Packet { return &Animate{} }
	p[packet.IDLevelChunk] = func() packet.Packet { return &LevelChunk{} }
	p[packet.IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[packet.IDPlayerList] = func() packet.Packet { return &PlayerList{} }
	p[packet.IDPlayerSkin] = func() packet.Packet { return &PlayerSkin{} }
	p[packet.IDDimensionData] = func() packet.Packet { return &DimensionData{} }
	p[packet.IDInventoryTransaction] = func() packet.Packet { return &InventoryTransaction{} }
	p[packet.IDInventoryContent] = func() packet.Packet { return &InventoryContent2193{} }
	p[packet.IDInventorySlot] = func() packet.Packet { return &InventorySlot2193{} }
	p[packet.IDItemStackResponse] = func() packet.Packet { return &ItemStackResponse2193{} }
	for _, id := range []uint32{
		packet.IDClientboundMatchmakingState, packet.IDServerboundStonecutterSetRecipe,
		packet.IDClientboundStonecutterSetRecipe, packet.IDServerboundMatchmakingCancel,
		packet.IDSetPassengerOfBlock, packet.IDServerboundCursorItemDrag,
		packet.IDClientboundPlayAudioContent, packet.IDServerboundRegisterAudioContent,
	} {
		delete(p, id)
	}
}

func FromLatest(pk packet.Packet) (packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.AddActor:
		return &AddActor{*pk}, true
	case *packet.AddPlayer:
		return &AddPlayer{*pk}, true
	case *packet.Animate:
		return &Animate{*pk}, true
	case *packet.LevelChunk:
		return &LevelChunk{*pk}, true
	case *packet.StartGame:
		return &StartGame{*pk}, true
	case *packet.PlayerList:
		return &PlayerList{*pk}, true
	case *packet.PlayerSkin:
		return &PlayerSkin{*pk}, true
	case *packet.DimensionData:
		return &DimensionData{*pk}, true
	case *packet.InventoryTransaction:
		return &InventoryTransaction{*pk}, true
	case *packet.InventoryContent:
		return &InventoryContent2193{*pk}, true
	case *packet.InventorySlot:
		return &InventorySlot2193{*pk}, true
	case *packet.ItemStackResponse:
		return &ItemStackResponse2193{*pk}, true
	}
	return nil, false
}

func ToLatest(pk packet.Packet) (packet.Packet, bool) {
	switch pk := pk.(type) {
	case *AddActor:
		return &pk.AddActor, true
	case *AddPlayer:
		return &pk.AddPlayer, true
	case *Animate:
		return &pk.Animate, true
	case *LevelChunk:
		return &pk.LevelChunk, true
	case *StartGame:
		return &pk.StartGame, true
	case *PlayerList:
		return &pk.PlayerList, true
	case *PlayerSkin:
		return &pk.PlayerSkin, true
	case *DimensionData:
		return &pk.DimensionData, true
	case *InventoryTransaction:
		return &pk.InventoryTransaction, true
	case *InventoryContent2193:
		return &pk.InventoryContent, true
	case *InventorySlot2193:
		return &pk.InventorySlot, true
	case *ItemStackResponse2193:
		return &pk.ItemStackResponse, true
	}
	return nil, false
}
