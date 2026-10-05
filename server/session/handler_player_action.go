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

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerActionHandler struct{}

func (*PlayerActionHandler) Handle(p packet.Packet, s *Session, _ *world.Tx, c Controllable) error {
	pk := p.(*packet.PlayerAction)

	return handlePlayerAction(pk.ActionType, pk.BlockFace, pk.BlockPosition, pk.EntityRuntimeID, s, c)
}

func handlePlayerAction(action int32, face int32, pos protocol.BlockPos, entityRuntimeID uint64, s *Session, c Controllable) error {
	if entityRuntimeID != selfEntityRuntimeID {
		return errSelfRuntimeID
	}
	switch action {
	case protocol.PlayerActionStartSleeping, protocol.PlayerActionRespawn, protocol.PlayerActionDimensionChangeDone:

	case protocol.PlayerActionStopSleeping:
		c.Wake()
	case protocol.PlayerActionStartBreak, protocol.PlayerActionContinueDestroyBlock:
		s.swingingArm.Store(true)
		defer s.swingingArm.Store(false)

		s.breakingPos = cube.Pos{int(pos[0]), int(pos[1]), int(pos[2])}
		c.StartBreaking(s.breakingPos, cube.Face(face))
	case protocol.PlayerActionAbortBreak:
		c.AbortBreaking()
	case protocol.PlayerActionPredictDestroyBlock, protocol.PlayerActionStopBreak:
		s.swingingArm.Store(true)
		defer s.swingingArm.Store(false)
		c.FinishBreaking()
	case protocol.PlayerActionCrackBreak:

	case protocol.PlayerActionStartItemUseOn:

	case protocol.PlayerActionStopItemUseOn:
		c.ReleaseItem()
	case protocol.PlayerActionInteractWithBlock:

	case protocol.PlayerActionCreativePlayerDestroyBlock:

	case protocol.PlayerActionMissedSwing:
		s.swingingArm.Store(true)
		defer s.swingingArm.Store(false)
		c.PunchAir()
	default:
		return fmt.Errorf("unhandled ActionType %v", action)
	}
	return nil
}
