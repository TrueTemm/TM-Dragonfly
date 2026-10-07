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

package v671

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

func (Protocol) ID() int32 { return 671 }

func (Protocol) Ver() string { return "1.20.81" }

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
		p[v786.IDOnScreenTextureAnimation] = mapLocked // swapped at 671
	}
	if onScreen != nil {
		p[v786.IDMapCreateLockedCopy] = onScreen
	}
	return p
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return v786.Reader786{Reader: protocol.NewReader(r, shieldID, enableLimits), Proto: 671}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return v786.Writer786{Writer: protocol.NewWriter(w, shieldID), Proto: 671}
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if out, ok := v686.ToLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}

type mapLocked671 struct{ *packet.MapCreateLockedCopy }

func (mapLocked671) ID() uint32 { return v786.IDOnScreenTextureAnimation }

type onScreen671 struct{ *packet.OnScreenTextureAnimation }

func (onScreen671) ID() uint32 { return v786.IDMapCreateLockedCopy }

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	switch pk := pk.(type) {
	case *packet.ClientBoundCloseForm:
		return nil
	case *packet.AwardAchievement:
		return nil
	case *packet.MapCreateLockedCopy:
		return []packet.Packet{&mapLocked671{pk}}
	case *packet.OnScreenTextureAnimation:
		return []packet.Packet{&onScreen671{pk}}
	case *packet.StartGame:
		out, _ := v686.FromLatestShared(uint32(p.ID()), pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.20.81", "1.20.81"
		items := itemdata.Items671()
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
