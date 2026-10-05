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

package session

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SubChunkRequestHandler struct{}

func (*SubChunkRequestHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, _ Controllable) error {
	pk := p.(*packet.SubChunkRequest)
	if dimID, _ := world.DimensionID(tx.World().Dimension()); pk.Dimension != int32(dimID) {

		s.writePacket(&packet.SubChunk{
			Dimension:       pk.Dimension,
			Position:        pk.Position,
			CacheEnabled:    s.conn.ClientCacheEnabled(),
			SubChunkEntries: []protocol.SubChunkEntry{},
		})
		return nil
	}
	s.ViewSubChunks(world.SubChunkPos(pk.Position), pk.Offsets, tx)
	return nil
}
