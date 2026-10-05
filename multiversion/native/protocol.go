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

package native

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const proto = uint32(protocol.CurrentProtocol)

type Protocol struct{}

func (Protocol) ID() int32 { return protocol.CurrentProtocol }

func (Protocol) Ver() string { return protocol.CurrentVersion }

func (Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return packet.NewClientPool()
	}
	return packet.NewServerPool()
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return Reader{Reader: protocol.NewReader(r, shieldID, enableLimits)}
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return Writer{Writer: protocol.NewWriter(w, shieldID)}
}

func (Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return convertToLatest(pk)
}

func (Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return convertFromLatest(pk)
}
