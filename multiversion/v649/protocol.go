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

package v649

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v686"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

type Protocol struct{}

func (Protocol) ID() int32 { return 649 }

func (Protocol) Ver() string { return "1.20.60" }

func (Protocol) Packets(listener bool) packet.Pool {
	var p packet.Pool
	if listener {
		p = v686.NewClientPool()
	} else {
		p = v686.NewServerPool()
	}
	delete(p, v786.IDClientBoundCloseForm)
	delete(p, v786.IDAwardAchievement)
	mapLocked, onScreen := p[v786.IDMapCreateLockedCopy], p[v786.IDOnScreenTextureAnimation]
	delete(p, v786.IDMapCreateLockedCopy)
	delete(p, v786.IDOnScreenTextureAnimation)
	if mapLocked != nil {
		p[v786.IDOnScreenTextureAnimation] = mapLocked // swapped below 685
	}
	if onScreen != nil {
		p[v786.IDMapCreateLockedCopy] = onScreen
	}
	p[idItemFrameDropItem] = func() packet.Packet { return &itemFrameDropItem{} } // gone from 662
	return p
}

const idItemFrameDropItem = 71 // 649's slot, reserved from 662

type itemFrameDropItem struct{ Position protocol.BlockPos }

func (*itemFrameDropItem) ID() uint32                { return idItemFrameDropItem }
func (pk *itemFrameDropItem) Marshal(io protocol.IO) { v786.UBlockPos786(io, &pk.Position) }

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return v786.Reader786{Reader: protocol.NewReader(r, shieldID, enableLimits), Proto: 649}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return v786.Writer786{Writer: protocol.NewWriter(w, shieldID), Proto: 649}
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if _, ok := pk.(*itemFrameDropItem); ok {
		return nil // no latest equivalent
	}
	if out, ok := v686.ToLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}

type mapLocked649 struct{ *packet.MapCreateLockedCopy }

func (mapLocked649) ID() uint32 { return v786.IDOnScreenTextureAnimation }

type onScreen649 struct {
	*packet.OnScreenTextureAnimation
}

func (onScreen649) ID() uint32 { return v786.IDMapCreateLockedCopy }

var argType649 = map[uint32]uint32{47: 43, 56: 44, 64: 52, 65: 53, 67: 55, 70: 58, 74: 62, 84: 71, 87: 74} // native arg type to 1.20.60

func remapCommands649(pk *packet.AvailableCommands) *packet.AvailableCommands {
	cp := *pk
	cp.Commands = make([]protocol.Command, len(pk.Commands))
	copy(cp.Commands, pk.Commands)
	for i := range cp.Commands {
		ov := append([]protocol.CommandOverload(nil), cp.Commands[i].Overloads...)
		for j := range ov {
			ps := append([]protocol.CommandParameter(nil), ov[j].Parameters...)
			for k := range ps {
				t := ps[k].Type
				if t&protocol.CommandArgEnum != 0 || t&protocol.CommandArgSoftEnum != 0 || t&protocol.CommandArgValid == 0 {
					continue // enum/soft-enum index, not an arg type
				}
				if m, ok := argType649[t&0xffff]; ok {
					ps[k].Type = (t &^ 0xffff) | m
				}
			}
			ov[j].Parameters = ps
		}
		cp.Commands[i].Overloads = ov
	}
	return &cp
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	switch pk := pk.(type) {
	case *packet.ClientBoundCloseForm:
		return nil
	case *packet.AwardAchievement:
		return nil
	case *packet.MapCreateLockedCopy:
		return []packet.Packet{&mapLocked649{pk}}
	case *packet.OnScreenTextureAnimation:
		return []packet.Packet{&onScreen649{pk}}
	case *packet.AvailableCommands:
		if out, ok := v686.FromLatestShared(uint32(p.ID()), remapCommands649(pk)); ok {
			return out
		}
		return nil
	case *packet.StartGame:
		out, _ := v686.FromLatestShared(uint32(p.ID()), pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.20.60", "1.20.60"
		items := itemdata.Items649()
		sg.Items = make([]v766.ItemEntry, len(items))
		for i, it := range items {
			sg.Items[i] = v766.ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID, ComponentBased: it.ComponentBased}
		}
		return out
	}
	if out, ok := v686.FromLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}
