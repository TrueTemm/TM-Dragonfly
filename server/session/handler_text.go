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

type TextHandler struct{}

func (TextHandler) Handle(p packet.Packet, s *Session, _ *world.Tx, c Controllable) error {
	pk := p.(*packet.Text)

	if pk.TextType != packet.TextTypeChat {
		return fmt.Errorf("TextType should always be Chat (%v), but got %v", packet.TextTypeChat, pk.TextType)
	}
	if pk.SourceName != s.conn.IdentityData().DisplayName {
		return fmt.Errorf("SourceName must be equal to DisplayName")
	}
	if pk.XUID != s.conn.IdentityData().XUID {
		return fmt.Errorf("XUID must be equal to player's XUID")
	}
	c.Chat(pk.Message)
	return nil
}
