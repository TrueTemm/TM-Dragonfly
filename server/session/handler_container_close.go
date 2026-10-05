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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ContainerCloseHandler struct{}

func (h *ContainerCloseHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.ContainerClose)

	c.MoveItemsToInventory()

	var containerType byte
	switch pk.WindowID {
	case 0:

		s.invOpened = false
	case byte(s.openedWindowID.Load()):
		containerType = byte(s.openedContainerID.Load())
		s.closeCurrentContainer(tx, true)
	case 0xff:

		s.invOpened = false
		if s.containerOpened.Load() {
			s.closeCurrentContainer(tx, false)
		}
		return nil
	default:
		containerType = pk.ContainerType
	}
	s.writePacket(&packet.ContainerClose{
		WindowID:      pk.WindowID,
		ContainerType: containerType,
	})
	return nil
}
