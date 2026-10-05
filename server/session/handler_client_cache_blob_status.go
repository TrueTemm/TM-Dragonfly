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

type ClientCacheBlobStatusHandler struct {
}

func (c *ClientCacheBlobStatusHandler) Handle(p packet.Packet, s *Session, _ *world.Tx, _ Controllable) error {
	pk := p.(*packet.ClientCacheBlobStatus)

	resp := &packet.ClientCacheMissResponse{Blobs: make([]protocol.CacheBlob, 0, len(pk.MissHashes))}

	s.blobMu.Lock()
	for _, hit := range pk.HitHashes {
		delete(s.blobs, hit)
		c.resolveBlob(hit, s)
	}
	for _, miss := range pk.MissHashes {
		blob, ok := s.blobs[miss]
		if !ok {

			continue
		}
		resp.Blobs = append(resp.Blobs, protocol.CacheBlob{Hash: miss, Payload: blob})
		c.resolveBlob(miss, s)
	}
	s.blobMu.Unlock()

	if len(resp.Blobs) > 0 {
		s.writePacket(resp)
	}
	return nil
}

func (c *ClientCacheBlobStatusHandler) resolveBlob(hash uint64, s *Session) {
	leftover := make([]map[uint64]struct{}, 0, len(s.openChunkTransactions))
	for _, m := range s.openChunkTransactions {
		delete(m, hash)
		if len(m) != 0 {
			leftover = append(leftover, m)
		}
	}
	s.openChunkTransactions = leftover
	delete(s.blobs, hash)
}
