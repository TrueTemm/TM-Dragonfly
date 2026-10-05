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

package v844

import (
	"github.com/df-mc/dragonfly/multiversion/v827"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func applyDeltas844(p packet.Pool) {

	p[IDStartGame] = func() packet.Packet { return &StartGame{} }

	p[IDServerBoundPackSettingChange] = func() packet.Packet { return &packet.ServerBoundPackSettingChange{} }
}

func NewServerPool() packet.Pool {
	p := v827.NewServerPool()
	applyDeltas844(p)
	return p
}

func NewClientPool() packet.Pool {
	p := v827.NewClientPool()
	applyDeltas844(p)
	return p
}
