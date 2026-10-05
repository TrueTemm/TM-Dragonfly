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

package v776

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

func applyDeltas776(p packet.Pool) {

	delete(p, v786.IDUpdateClientOptions)
	delete(p, v786.IDPlayerVideoCapture)
	delete(p, v786.IDPlayerUpdateEntityOverrides)

	p[v786.IDLevelSoundEvent] = func() packet.Packet { return &LevelSoundEvent{} }
	p[v786.IDClientMovementPredictionSync] = func() packet.Packet { return &ClientMovementPredictionSync{} }
	p[v786.IDSetHud] = func() packet.Packet { return &SetHud{} }
}

func NewClientPool() packet.Pool {
	p := v786.NewClientPool()
	applyDeltas776(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v786.NewServerPool()
	applyDeltas776(p)
	return p
}
