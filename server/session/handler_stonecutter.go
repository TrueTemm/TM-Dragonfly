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
	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const stonecutterInputSlot = 0x03

func (h *ItemStackRequestHandler) handleStonecutting(a *protocol.CraftRecipeStackRequestAction, s *Session, tx *world.Tx) error {
	craft, ok := s.recipes[a.RecipeNetworkID]
	if !ok {
		return fmt.Errorf("recipe with network id %v does not exist", a.RecipeNetworkID)
	}
	if _, shapeless := craft.(recipe.Shapeless); !shapeless {
		return fmt.Errorf("recipe with network id %v is not a shapeless recipe", a.RecipeNetworkID)
	}
	if craft.Block() != "stonecutter" {
		return fmt.Errorf("recipe with network id %v is not a stonecutter recipe", a.RecipeNetworkID)
	}

	timesCrafted := int(a.NumberOfCrafts)
	if timesCrafted < 1 {
		return fmt.Errorf("times crafted must be at least 1")
	}

	expectedInputs := craft.Input()
	input, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerStonecutterInput},
		Slot:      stonecutterInputSlot,
	}, s, tx)
	if input.Count() < timesCrafted {
		return fmt.Errorf("input item count is less than number of crafts")
	}
	if !matchingStacks(input, expectedInputs[0]) {
		return fmt.Errorf("input item is not the same as expected input")
	}

	output := craft.Output()
	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerStonecutterInput},
		Slot:      stonecutterInputSlot,
	}, input.Grow(-timesCrafted), s, tx)
	return h.createResults(s, tx, repeatStacks(output, timesCrafted)...)
}
