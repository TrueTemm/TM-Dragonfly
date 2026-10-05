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
	"github.com/df-mc/dragonfly/multiversion/v844"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func applyDeltas859(p packet.Pool) {
	p[IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[IDAnimate] = func() packet.Packet { return &Animate{} }

	p[IDRequestChunkRadius] = func() packet.Packet { return &packet.RequestChunkRadius{} }
	p[IDRequestPermissions] = func() packet.Packet { return &packet.RequestPermissions{} }
}

func NewServerPool() packet.Pool {
	p := v844.NewServerPool()
	applyDeltas859(p)
	return p
}

func NewClientPool() packet.Pool {
	p := v844.NewClientPool()
	applyDeltas859(p)
	return p
}
