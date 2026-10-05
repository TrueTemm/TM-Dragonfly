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
	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	smithingInputSlot = 0x33

	smithingMaterialSlot = 0x34

	smithingTemplateSlot = 0x35
)

func (h *ItemStackRequestHandler) handleSmithing(a *protocol.CraftRecipeStackRequestAction, s *Session, tx *world.Tx) error {

	craft, ok := s.recipes[a.RecipeNetworkID]
	if !ok {
		return fmt.Errorf("recipe with network id %v does not exist", a.RecipeNetworkID)
	}
	if craft.Block() != "smithing_table" {
		return fmt.Errorf("recipe with network id %v is not a smithing table recipe", a.RecipeNetworkID)
	}
	switch craft.(type) {
	case recipe.SmithingTransform, recipe.SmithingTrim:
	default:
		return fmt.Errorf("recipe with network id %v is not a smithing recipe", a.RecipeNetworkID)
	}

	expectedInputs := craft.Input()
	input, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerSmithingTableInput},
		Slot:      smithingInputSlot,
	}, s, tx)
	if !matchingStacks(input, expectedInputs[0]) {
		return fmt.Errorf("input item is not the same as expected input")
	}
	material, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerSmithingTableMaterial},
		Slot:      smithingMaterialSlot,
	}, s, tx)
	if !matchingStacks(material, expectedInputs[1]) {
		return fmt.Errorf("material item is not the same as expected material")
	}
	template, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerSmithingTableTemplate},
		Slot:      smithingTemplateSlot,
	}, s, tx)
	if !matchingStacks(template, expectedInputs[2]) {
		return fmt.Errorf("template item is not the same as expected template")
	}

	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerSmithingTableInput},
		Slot:      smithingInputSlot,
	}, input.Grow(-1), s, tx)
	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerSmithingTableMaterial},
		Slot:      smithingMaterialSlot,
	}, material.Grow(-1), s, tx)
	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerSmithingTableTemplate},
		Slot:      smithingTemplateSlot,
	}, template.Grow(-1), s, tx)

	result := input.Grow(1 - input.Count())
	if _, ok = craft.(recipe.SmithingTrim); ok {
		var trim item.ArmourTrim
		if t, ok := template.Item().(item.SmithingTemplate); ok {
			trim.Template = t.Template
		} else {
			return fmt.Errorf("template item is not a smithing template")
		}
		if trim.Material, ok = material.Item().(item.ArmourTrimMaterial); !ok {
			return fmt.Errorf("material item is not an armour trim material")
		}
		trimmable, ok := input.Item().(item.Trimmable)
		if !ok {
			return fmt.Errorf("input item is not trimmable")
		}
		return h.createResults(s, tx, result.WithItem(trimmable.WithTrim(trim)))
	}
	return h.createResults(s, tx, result.WithItem(craft.Output()[0].Item()))
}
