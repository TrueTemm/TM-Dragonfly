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

package v685

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v686"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

type Protocol struct{}

func (Protocol) ID() int32 { return 685 }

func (Protocol) Ver() string { return "1.21.1" }

func (Protocol) Packets(listener bool) packet.Pool {
	var p packet.Pool
	if listener {
		p = v686.NewClientPool()
	} else {
		p = v686.NewServerPool()
	}
	delete(p, v786.IDClientBoundCloseForm)
	return p
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return v786.Reader786{Reader: protocol.NewReader(r, shieldID, enableLimits), Proto: 685}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return v786.Writer786{Writer: protocol.NewWriter(w, shieldID), Proto: 685}
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if out, ok := v686.ToLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	switch pk.(type) {
	case *packet.ClientBoundCloseForm:
		return nil
	case *packet.StartGame:
		out, _ := v686.FromLatestShared(uint32(p.ID()), pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.21.1", "1.21.1"
		return out
	}
	if out, ok := v686.FromLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}
