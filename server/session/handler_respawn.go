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
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type RespawnHandler struct{}

func (*RespawnHandler) Handle(p packet.Packet, _ *Session, _ *world.Tx, c Controllable) error {
	pk := p.(*packet.Respawn)
	if pk.EntityRuntimeID != selfEntityRuntimeID {
		return errSelfRuntimeID
	}
	if pk.State != packet.RespawnStateClientReadyToSpawn {
		return fmt.Errorf("respawn state must always be %v, but got %v", packet.RespawnStateClientReadyToSpawn, pk.State)
	}
	c.Respawn()
	return nil
}
