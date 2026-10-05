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
	"github.com/df-mc/dragonfly/server/player/dialogue"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type NPCRequestHandler struct {
	dialogue        dialogue.Dialogue
	entityRuntimeID uint64
}

func (h *NPCRequestHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.NPCRequest)
	if h.entityRuntimeID == 0 {

		return nil
	}
	switch pk.RequestType {
	case packet.NPCRequestActionExecuteAction:
		if err := h.dialogue.Submit(uint(pk.ActionType), c, tx); err != nil {
			return fmt.Errorf("error submitting dialogue: %w", err)
		}
	case packet.NPCRequestActionExecuteClosingCommands:
		h.dialogue.Close(c, tx)
	}
	return nil
}
