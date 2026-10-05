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

package v729

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

type Protocol struct{}

func (Protocol) ID() int32 { return 729 }

func (Protocol) Ver() string { return "1.21.30" }

func (Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return NewClientPool()
	}
	return NewServerPool()
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return v786.Reader786{Reader: protocol.NewReader(r, shieldID, enableLimits), Proto: 729}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return v786.Writer786{Writer: protocol.NewWriter(w, shieldID), Proto: 729}
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if out, ok := ToLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if out, ok := FromLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}
