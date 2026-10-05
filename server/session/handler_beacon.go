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
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const beaconInputSlot = 0x1b

func (h *ItemStackRequestHandler) handleBeaconPayment(a *protocol.BeaconPaymentStackRequestAction, s *Session, tx *world.Tx) error {
	slot := protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerBeaconPayment},
		Slot:      beaconInputSlot,
	}

	if !s.containerOpened.Load() {
		return fmt.Errorf("no beacon container opened")
	}
	pos := *s.openedPos.Load()
	beacon, ok := tx.Block(pos).(block.Beacon)
	if !ok {
		return fmt.Errorf("no beacon container opened")
	}

	payment, _ := h.itemInSlot(slot, s, tx)
	if payable, ok := payment.Item().(item.BeaconPayment); !ok || !payable.PayableForBeacon() {
		return fmt.Errorf("item %#v in beacon slot cannot be used as payment", payment)
	}

	if !h.validBeaconEffect(a.PrimaryEffect, beacon) {
		return fmt.Errorf("primary effect selected is not allowed: %v for level %v", a.PrimaryEffect, beacon.Level())
	} else if !h.validBeaconEffect(a.SecondaryEffect, beacon) || (beacon.Level() < 4 && a.SecondaryEffect != 0) {
		return fmt.Errorf("secondary effect selected is not allowed: %v for level %v", a.SecondaryEffect, beacon.Level())
	}

	primary, pOk := effect.ByID(int(a.PrimaryEffect))
	secondary, sOk := effect.ByID(int(a.SecondaryEffect))
	if pOk {
		beacon.Primary = primary.(effect.LastingType)
	}
	if sOk {
		beacon.Secondary = secondary.(effect.LastingType)
	}
	tx.SetBlock(pos, beacon, nil)

	h.setItemInSlot(slot, item.Stack{}, s, tx)
	h.ignoreDestroy = true
	return nil
}

func (h *ItemStackRequestHandler) validBeaconEffect(id int32, beacon block.Beacon) bool {
	switch id {
	case 1, 3:
		return beacon.Level() >= 1
	case 8, 11:
		return beacon.Level() >= 2
	case 5:
		return beacon.Level() >= 3
	case 10:
		return beacon.Level() >= 4
	case 0:
		return true
	}
	return false
}
