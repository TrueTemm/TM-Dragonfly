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

package v898

import (
	"github.com/df-mc/dragonfly/multiversion/v859"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func applyDeltas898(p packet.Pool) {
	p[IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[IDText] = func() packet.Packet { return &Text{} }

	p[IDAnimate] = func() packet.Packet { return &packet.Animate{} }
	p[IDCommandRequest] = func() packet.Packet { return &packet.CommandRequest{} }
	p[IDCommandOutput] = func() packet.Packet { return &packet.CommandOutput{} }

	p[IDInteract] = func() packet.Packet { return &packet.Interact{} }
}

func NewServerPool() packet.Pool {
	p := v859.NewServerPool()
	applyDeltas898(p)
	return p
}

func NewClientPool() packet.Pool {
	p := v859.NewClientPool()
	applyDeltas898(p)
	return p
}
