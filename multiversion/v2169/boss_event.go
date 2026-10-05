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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type BossEvent struct {
	packet.BossEvent

	PlayerUniqueID int64
}

func (pk *BossEvent) Marshal(io protocol.IO) {
	io.ActorUniqueID(&pk.BossEntityUniqueID)
	io.ActorUniqueID(&pk.PlayerUniqueID)
	io.Uint8(&pk.EventType)
	io.String(&pk.BossBarTitle)
	io.String(&pk.FilteredBossBarTitle)
	io.Float32(&pk.HealthPercentage)
	io.Uint8(&pk.Colour)
	io.Uint8(&pk.Overlay)
}

func fromLatestBossEvent(pk *packet.BossEvent) *BossEvent {
	return &BossEvent{BossEvent: *pk}
}

func toLatestBossEvent(pk *BossEvent) *packet.BossEvent {
	out := pk.BossEvent
	return &out
}
