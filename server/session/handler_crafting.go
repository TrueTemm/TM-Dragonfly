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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"math"
	"slices"
)

var CraftingEnabled = false

var errCraftingDisabled = fmt.Errorf("crafting is disabled")

func (h *ItemStackRequestHandler) handleCraft(a *protocol.CraftRecipeStackRequestAction, s *Session, tx *world.Tx) error {
	craft, ok := s.recipes[a.RecipeNetworkID]
	if !ok {

		return h.tryDynamicCraft(s, tx, int(a.NumberOfCrafts))
	}
	_, shaped := craft.(recipe.Shaped)
	_, shapeless := craft.(recipe.Shapeless)
	if !shaped && !shapeless {
		return fmt.Errorf("recipe with network id %v is not a shaped or shapeless recipe", a.RecipeNetworkID)
	}
	if craft.Block() != "crafting_table" {
		return fmt.Errorf("recipe with network id %v is not a crafting table recipe", a.RecipeNetworkID)
	}

	timesCrafted := int(a.NumberOfCrafts)
	if timesCrafted < 1 {
		return fmt.Errorf("times crafted must be at least 1")
	}

	size := s.craftingSize()
	offset := s.craftingOffset()
	consumed := make([]bool, size)
	for _, expected := range craft.Input() {
		var processed bool
		for slot := offset; slot < offset+size; slot++ {
			if consumed[slot-offset] {

				continue
			}
			has, _ := s.ui.Item(int(slot))
			if has.Empty() != expected.Empty() || has.Count() < expected.Count()*timesCrafted {

				continue
			}
			if !matchingStacks(has, expected) {

				continue
			}
			processed, consumed[slot-offset] = true, true
			st := has.Grow(-expected.Count() * timesCrafted)
			h.setItemInSlot(protocol.StackRequestSlotInfo{
				Container: protocol.FullContainerName{ContainerID: protocol.ContainerCraftingInput},
				Slot:      byte(slot),
			}, st, s, tx)
			break
		}
		if !processed {
			return fmt.Errorf("recipe %v: could not consume expected item: %v", a.RecipeNetworkID, expected)
		}
	}
	return h.createResults(s, tx, repeatStacks(craft.Output(), timesCrafted)...)
}

func (h *ItemStackRequestHandler) handleAutoCraft(a *protocol.AutoCraftRecipeStackRequestAction, s *Session, tx *world.Tx) error {
	craft, ok := s.recipes[a.RecipeNetworkID]
	if !ok {

		return h.tryDynamicCraft(s, tx, int(a.NumberOfCrafts))
	}
	_, shaped := craft.(recipe.Shaped)
	_, shapeless := craft.(recipe.Shapeless)
	if !shaped && !shapeless {
		return fmt.Errorf("recipe with network id %v is not a shaped or shapeless recipe", a.RecipeNetworkID)
	}
	if craft.Block() != "crafting_table" {
		return fmt.Errorf("recipe with network id %v is not a crafting table recipe", a.RecipeNetworkID)
	}

	timesCrafted := int(a.NumberOfCrafts)
	if timesCrafted < 1 {
		return fmt.Errorf("times crafted must be at least 1")
	}

	flattenedInputs := make([]recipe.Item, 0, len(craft.Input()))
	for _, i := range craft.Input() {
		if i.Empty() {

			continue
		}

		if ind := slices.IndexFunc(flattenedInputs, func(it recipe.Item) bool {
			return matchingStacks(it, i)
		}); ind >= 0 {
			flattenedInputs[ind] = grow(i, flattenedInputs[ind].Count())
			continue
		}
		flattenedInputs = append(flattenedInputs, i)
	}

	for _, expected := range flattenedInputs {
		remaining := expected.Count() * timesCrafted

		for id, inv := range map[byte]*inventory.Inventory{
			protocol.ContainerCraftingInput:              s.ui,
			protocol.ContainerCombinedHotBarAndInventory: s.inv,
		} {
			for slot, has := range inv.Slots() {
				if has.Empty() {

					continue
				}
				if !matchingStacks(has, expected) {

					continue
				}

				removal := has.Count()
				if remaining < removal {
					removal = remaining
				}
				remaining -= removal

				has = has.Grow(-removal)
				h.setItemInSlot(protocol.StackRequestSlotInfo{
					Container: protocol.FullContainerName{ContainerID: id},
					Slot:      byte(slot),
				}, has, s, tx)
				if remaining == 0 {

					break
				}
			}
			if remaining == 0 {

				break
			}
		}
		if remaining != 0 {
			return fmt.Errorf("recipe %v: could not consume expected item: %v", a.RecipeNetworkID, expected)
		}
	}

	return h.createResults(s, tx, repeatStacks(craft.Output(), timesCrafted)...)
}

func (h *ItemStackRequestHandler) handleCreativeCraft(a *protocol.CraftCreativeStackRequestAction, s *Session, tx *world.Tx, c Controllable) error {
	if !c.GameMode().CreativeInventory() {
		return fmt.Errorf("can only craft creative items in gamemode creative/spectator")
	}
	index := a.CreativeItemNetworkID - 1
	if int(index) >= len(creative.Items()) {
		return fmt.Errorf("creative item with network ID %v does not exist", index)
	}
	it := creative.Items()[index].Stack
	it = it.Grow(it.MaxCount() - 1)
	return h.createResults(s, tx, it)
}

func (s *Session) craftingSize() uint32 {
	if s.openedContainerID.Load() == 1 {
		return craftingGridSizeLarge
	}
	return craftingGridSizeSmall
}

func (s *Session) craftingOffset() uint32 {
	if s.openedContainerID.Load() == 1 {
		return craftingGridLargeOffset
	}
	return craftingGridSmallOffset
}

func matchingStacks(has, expected recipe.Item) bool {
	switch expected := expected.(type) {
	case item.Stack:
		switch has := has.(type) {
		case recipe.ItemTag:
			name, _ := expected.Item().EncodeItem()
			return has.Contains(name)
		case item.Stack:
			_, variants := expected.Value("variants")
			if !variants {
				return has.Comparable(expected)
			}
			nameOne, _ := has.Item().EncodeItem()
			nameTwo, _ := expected.Item().EncodeItem()
			return nameOne == nameTwo
		}
		panic(fmt.Errorf("client has unexpected recipe item %T", has))
	case recipe.ItemTag:
		switch has := has.(type) {
		case item.Stack:
			name, _ := has.Item().EncodeItem()
			return expected.Contains(name)
		case recipe.ItemTag:
			return has.Tag() == expected.Tag()
		}
		panic(fmt.Errorf("client has unexpected recipe item %T", has))
	}
	panic(fmt.Errorf("tried to match with unexpected recipe item %T", expected))
}

func repeatStacks(items []item.Stack, repetitions int) []item.Stack {
	output := make([]item.Stack, 0, len(items))
	for _, o := range items {
		count, maxCount := o.Count(), o.MaxCount()
		total := count * repetitions

		stacks := int(math.Ceil(float64(total) / float64(maxCount)))
		for i := 0; i < stacks; i++ {
			inc := min(total, maxCount)
			total -= inc

			output = append(output, o.Grow(inc-count))
		}
	}
	return output
}

func grow(i recipe.Item, count int) recipe.Item {
	switch i := i.(type) {
	case item.Stack:
		return i.Grow(count)
	case recipe.ItemTag:
		return recipe.NewItemTag(i.Tag(), i.Count()+count)
	}
	panic(fmt.Errorf("unexpected recipe item %T", i))
}

func (h *ItemStackRequestHandler) tryDynamicCraft(s *Session, tx *world.Tx, timesCrafted int) error {
	if timesCrafted < 1 {
		return fmt.Errorf("times crafted must be at least 1")
	}

	size := s.craftingSize()
	offset := s.craftingOffset()

	input := make([]recipe.Item, size)
	for i := uint32(0); i < size; i++ {
		slot := offset + i
		it, _ := s.ui.Item(int(slot))
		if it.Empty() {
			input[i] = item.Stack{}
		} else {
			input[i] = it
		}
	}

	for _, dynamicRecipe := range recipe.DynamicRecipes() {
		if dynamicRecipe.Block() != "crafting_table" {
			continue
		}

		output, ok := dynamicRecipe.Match(input)
		if !ok {
			continue
		}

		minStackCount := math.MaxInt
		for i := uint32(0); i < size; i++ {
			slot := offset + i
			it, _ := s.ui.Item(int(slot))
			if !it.Empty() {
				if it.Count() < minStackCount {
					minStackCount = it.Count()
				}
			}
		}

		if minStackCount < timesCrafted {
			timesCrafted = minStackCount
		}

		for i := uint32(0); i < size; i++ {
			slot := offset + i
			it, _ := s.ui.Item(int(slot))
			if !it.Empty() {

				st := it.Grow(-1 * timesCrafted)
				h.setItemInSlot(protocol.StackRequestSlotInfo{
					Container: protocol.FullContainerName{ContainerID: protocol.ContainerCraftingInput},
					Slot:      byte(slot),
				}, st, s, tx)
			}
		}

		return h.createResults(s, tx, repeatStacks(output, timesCrafted)...)
	}

	return fmt.Errorf("no matching recipe found for crafting grid")
}
