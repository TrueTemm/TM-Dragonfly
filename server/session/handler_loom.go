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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	loomInputSlot = 0x09

	loomDyeSlot = 0x0a

	loomPatternSlot = 0x0b
)

func (h *ItemStackRequestHandler) handleLoomCraft(a *protocol.CraftLoomRecipeStackRequestAction, s *Session, tx *world.Tx) error {

	if _, ok := tx.Block(*s.openedPos.Load()).(block.Loom); !ok || !s.containerOpened.Load() {
		return fmt.Errorf("no loom container opened")
	}
	timesCrafted := int(a.TimesCrafted)
	if timesCrafted < 1 {
		return fmt.Errorf("times crafted must be least 1")
	}

	input, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerLoomInput},
		Slot:      loomInputSlot,
	}, s, tx)
	if input.Count() < timesCrafted {
		return fmt.Errorf("input item count is less than times crafted")
	}
	b, ok := input.Item().(block.Banner)
	if !ok {
		return fmt.Errorf("input item is not a banner")
	}
	if b.Illager {
		return fmt.Errorf("input item is an illager banner")
	}

	dye, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerLoomDye},
		Slot:      loomDyeSlot,
	}, s, tx)
	if dye.Count() < timesCrafted {
		return fmt.Errorf("dye item count is less than times crafted")
	}
	d, ok := dye.Item().(item.Dye)
	if !ok {
		return fmt.Errorf("dye item is not a dye")
	}

	expectedPattern, exists := block.BannerPatternByID(a.Pattern)
	if !exists {
		return fmt.Errorf("unknown banner pattern id %q", a.Pattern)
	}

	pattern, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerLoomMaterial},
		Slot:      loomPatternSlot,
	}, s, tx)
	if expectedPatternItem, hasPatternItem := expectedPattern.Item(); hasPatternItem {
		if pattern.Empty() {
			return fmt.Errorf("pattern item is empty but the pattern is required")
		}
		p, ok := pattern.Item().(item.BannerPattern)
		if !ok {
			return fmt.Errorf("pattern item is not a banner pattern")
		}
		if expectedPatternItem != p.Type {
			return fmt.Errorf("pattern item does not match the expected pattern")
		}
	}

	b.Patterns = append(b.Patterns, block.BannerPatternLayer{
		Type:   expectedPattern,
		Colour: d.Colour,
	})
	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerLoomInput},
		Slot:      loomInputSlot,
	}, input.Grow(-timesCrafted), s, tx)
	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerLoomDye},
		Slot:      loomDyeSlot,
	}, dye.Grow(-timesCrafted), s, tx)

	return h.createResults(s, tx, input.Grow(timesCrafted-input.Count()).WithItem(b))
}
