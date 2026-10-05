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
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type LecternUpdateHandler struct{}

func (LecternUpdateHandler) Handle(p packet.Packet, _ *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.LecternUpdate)
	pos := blockPosFromProtocol(pk.Position)
	if !canReach(c, pos.Vec3Middle()) {
		return fmt.Errorf("block at %v is not within reach", pos)
	}
	if _, ok := tx.Block(pos).(block.Lectern); !ok {
		return fmt.Errorf("block at %v is not a lectern", pos)
	}
	return c.TurnLecternPage(pos, int(pk.Page))
}
