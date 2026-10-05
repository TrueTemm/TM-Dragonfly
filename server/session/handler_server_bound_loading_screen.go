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
	"sync/atomic"
)

type ServerBoundLoadingScreenHandler struct {
	currentID  atomic.Uint32
	expectedID atomic.Uint32
}

func (h *ServerBoundLoadingScreenHandler) Handle(p packet.Packet, s *Session, _ *world.Tx, _ Controllable) error {
	pk := p.(*packet.ServerBoundLoadingScreen)
	v, ok := pk.LoadingScreenID.Value()
	expected := h.expectedID.Load()

	switch {
	case !ok || expected == 0:
		return nil
	case v != expected:
		return fmt.Errorf("expected loading screen ID %d, got %d", expected, v)
	case pk.Type == packet.LoadingScreenTypeEnd:
		s.changingDimension.Store(false)
		h.expectedID.Store(0)
	}

	return nil
}
