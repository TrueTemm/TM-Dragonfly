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
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	grindstoneFirstInputSlot = 0x10

	grindstoneSecondInputSlot = 0x11
)

func (h *ItemStackRequestHandler) handleGrindstoneCraft(s *Session, tx *world.Tx, c Controllable) error {

	if !s.containerOpened.Load() {
		return fmt.Errorf("no grindstone container opened")
	}
	if _, ok := tx.Block(*s.openedPos.Load()).(block.Grindstone); !ok {
		return fmt.Errorf("no grindstone container opened")
	}

	firstInput, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerGrindstoneInput},
		Slot:      grindstoneFirstInputSlot,
	}, s, tx)
	secondInput, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerGrindstoneAdditional},
		Slot:      grindstoneSecondInputSlot,
	}, s, tx)
	if firstInput.Empty() && secondInput.Empty() {
		return fmt.Errorf("input item(s) are empty")
	}
	if firstInput.Count() > 1 || secondInput.Count() > 1 {
		return fmt.Errorf("input item(s) are not single items")
	}

	resultStack := nonZeroItem(firstInput, secondInput)
	if !firstInput.Empty() && !secondInput.Empty() {
		name, meta := firstInput.Item().EncodeItem()
		name2, meta2 := secondInput.Item().EncodeItem()
		if name != name2 || meta != meta2 {
			return fmt.Errorf("input items must be the same type")
		}
		if _, ok := firstInput.Item().(item.Durable); !ok {
			return fmt.Errorf("input items must be durable")
		}

		resultStack = firstInput.WithEnchantments(secondInput.Enchantments()...)

		maxDurability := firstInput.MaxDurability()
		firstDurability, secondDurability := firstInput.Durability(), secondInput.Durability()

		resultStack = resultStack.WithDurability(firstDurability + secondDurability + maxDurability*5/100)
	}

	for _, o := range entity.NewExperienceOrbs(entity.EyePosition(c), experienceFromEnchantments(resultStack)) {
		tx.AddEntity(o)
	}

	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerGrindstoneInput},
		Slot:      grindstoneFirstInputSlot,
	}, item.Stack{}, s, tx)
	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerGrindstoneAdditional},
		Slot:      grindstoneSecondInputSlot,
	}, item.Stack{}, s, tx)
	return h.createResults(s, tx, stripPossibleEnchantments(resultStack))
}

type curseEnchantment interface {
	Curse() bool
}

func experienceFromEnchantments(stack item.Stack) int {
	var totalCost int
	for _, enchant := range stack.Enchantments() {
		if _, ok := enchant.Type().(curseEnchantment); ok {
			continue
		}
		cost, _ := enchant.Type().Cost(enchant.Level())
		totalCost += cost
	}
	if totalCost == 0 {

		return 0
	}

	minExperience := int(math.Ceil(float64(totalCost) / 2))
	return minExperience + rand.IntN(minExperience)
}

func stripPossibleEnchantments(stack item.Stack) item.Stack {
	for _, enchant := range stack.Enchantments() {
		if _, ok := enchant.Type().(curseEnchantment); ok {
			continue
		}
		stack = stack.WithoutEnchantments(enchant.Type())
	}
	return stack.WithAnvilCost(0)
}

func nonZeroItem(first, second item.Stack) item.Stack {
	if first.Empty() {
		return second
	}
	return first
}
