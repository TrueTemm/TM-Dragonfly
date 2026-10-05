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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type ClientCacheBlobStatus struct {
	MissHashes []uint64

	HitHashes []uint64
}

func (*ClientCacheBlobStatus) ID() uint32 {
	return IDClientCacheBlobStatus
}

func (pk *ClientCacheBlobStatus) Marshal(io protocol.IO) {
	missLen, hitLen := uint32(len(pk.MissHashes)), uint32(len(pk.HitHashes))
	io.Varuint32(&missLen)
	io.Varuint32(&hitLen)
	protocol.FuncSliceOfLen(io, missLen, &pk.MissHashes, io.Uint64)
	protocol.FuncSliceOfLen(io, hitLen, &pk.HitHashes, io.Uint64)
}
