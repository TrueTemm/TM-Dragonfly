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

package v859

import (
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{}

func (Protocol) ID() int32 { return 859 }

func (Protocol) Ver() string { return "1.21.120" }

func (Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return NewClientPool()
	}
	return NewServerPool()
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return v786.Reader786{Reader: protocol.NewReader(r, shieldID, enableLimits), Proto: 859}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return v786.Writer786{Writer: protocol.NewWriter(w, shieldID), Proto: 859}
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if out, ok := convertToLatest(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	if out, ok := convertFromLatest(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}
