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

package v2168

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v2169"
)

func applyDeltas2168(p packet.Pool) {
	p[packet.IDItemStackResponse] = func() packet.Packet { return &ItemStackResponse{} }
	p[packet.IDSetScore] = func() packet.Packet { return &SetScore{} }
}

func NewClientPool() packet.Pool {
	p := v2169.NewClientPool()
	applyDeltas2168(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v2169.NewServerPool()
	applyDeltas2168(p)
	return p
}
