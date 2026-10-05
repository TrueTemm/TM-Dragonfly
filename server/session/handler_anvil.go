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
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	anvilInputSlot = 0x1

	anvilMaterialSlot = 0x2
)

func (h *ItemStackRequestHandler) handleCraftRecipeOptional(a *protocol.CraftRecipeOptionalStackRequestAction, s *Session, filterStrings []string, co Controllable, tx *world.Tx) (err error) {

	if !s.containerOpened.Load() {
		return fmt.Errorf("no anvil container opened")
	}

	pos := *s.openedPos.Load()
	anvil, ok := tx.Block(pos).(block.Anvil)
	if !ok {
		return fmt.Errorf("no anvil container opened")
	}
	if index := int(a.FilterStringIndex); len(filterStrings) > 0 && (index < 0 || index >= len(filterStrings)) {
		return fmt.Errorf("filter string index %v is out of bounds", a.FilterStringIndex)
	}

	input, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerAnvilInput},
		Slot:      anvilInputSlot,
	}, s, tx)
	if input.Empty() {
		return fmt.Errorf("no item in input input slot")
	}
	material, _ := h.itemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerAnvilMaterial},
		Slot:      anvilMaterialSlot,
	}, s, tx)
	result := input

	anvilCost := input.AnvilCost()
	if !material.Empty() {
		anvilCost += material.AnvilCost()
	}

	var actionCost, renameCost, repairCount int
	if !material.Empty() {

		if repairable, ok := input.Item().(item.Repairable); ok && repairable.RepairableBy(material) {
			result, actionCost, repairCount, err = repairItemWithMaterial(input, material, result)
			if err != nil {
				return err
			}
		} else {
			_, book := material.Item().(item.EnchantedBook)
			_, durable := input.Item().(item.Durable)

			enchantedBook := book && len(material.Enchantments()) > 0
			if !enchantedBook && (input.Item() != material.Item() || !durable) {
				return fmt.Errorf("input item is not repairable/same type or material item is not an enchanted book")
			}

			if durable && !enchantedBook {
				result, actionCost = repairItemWithDurable(input, material, result)
			}

			var hasCompatible, hasIncompatible bool
			result, hasCompatible, hasIncompatible, actionCost = mergeEnchantments(input, material, result, actionCost, enchantedBook)

			if !durable && hasIncompatible && !hasCompatible {
				return fmt.Errorf("no compatible enchantments but have incompatible ones")
			}
		}
	}

	if len(filterStrings) > 0 {
		renameCost = 1
		actionCost += renameCost
		result = result.WithCustomName(filterStrings[int(a.FilterStringIndex)])
	}

	cost := actionCost + anvilCost
	if cost <= 0 {
		return fmt.Errorf("no action was taken")
	}

	if renameCost == actionCost && renameCost > 0 && cost >= 40 {
		cost = 39
	}

	c := co.GameMode().CreativeInventory()
	if cost >= 40 && !c {
		return fmt.Errorf("impossible cost")
	}

	level := co.ExperienceLevel()
	if level < cost && !c {
		return fmt.Errorf("not enough experience")
	} else if !c {
		co.SetExperienceLevel(level - cost)
	}

	if !result.Empty() {
		updatedAnvilCost := result.AnvilCost()
		if !material.Empty() && updatedAnvilCost < material.AnvilCost() {
			updatedAnvilCost = material.AnvilCost()
		}
		if renameCost != actionCost || renameCost == 0 {
			updatedAnvilCost = updatedAnvilCost*2 + 1
		}
		result = result.WithAnvilCost(updatedAnvilCost)
	}

	if !c && rand.Float64() < 0.12 {
		damaged := anvil.Break()
		if _, ok := damaged.(block.Air); ok {
			tx.PlaySound(pos.Vec3Centre(), sound.AnvilBreak{})
		} else {
			tx.PlaySound(pos.Vec3Centre(), sound.AnvilUse{})
		}
		defer tx.SetBlock(pos, damaged, nil)
	} else {
		tx.PlaySound(pos.Vec3Centre(), sound.AnvilUse{})
	}

	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerAnvilInput},
		Slot:      anvilInputSlot,
	}, item.Stack{}, s, tx)
	if repairCount > 0 {
		h.setItemInSlot(protocol.StackRequestSlotInfo{
			Container: protocol.FullContainerName{ContainerID: protocol.ContainerAnvilMaterial},
			Slot:      anvilMaterialSlot,
		}, material.Grow(-repairCount), s, tx)
	} else {
		h.setItemInSlot(protocol.StackRequestSlotInfo{
			Container: protocol.FullContainerName{ContainerID: protocol.ContainerAnvilMaterial},
			Slot:      anvilMaterialSlot,
		}, item.Stack{}, s, tx)
	}
	return h.createResults(s, tx, result)
}

func repairItemWithMaterial(input item.Stack, material item.Stack, result item.Stack) (item.Stack, int, int, error) {

	delta := min(input.MaxDurability()-input.Durability(), input.MaxDurability()/4)
	if delta <= 0 {
		return item.Stack{}, 0, 0, fmt.Errorf("input item is already fully repaired")
	}

	var cost, count int
	for ; delta > 0 && count < material.Count(); count, delta = count+1, min(result.MaxDurability()-result.Durability(), result.MaxDurability()/4) {
		result = result.WithDurability(result.Durability() + delta)
		cost++
	}
	return result, cost, count, nil
}

func repairItemWithDurable(input item.Stack, durable item.Stack, result item.Stack) (item.Stack, int) {
	durability := input.Durability() + durable.Durability() + input.MaxDurability()*12/100
	if durability > input.MaxDurability() {
		durability = input.MaxDurability()
	}

	var cost int
	if durability > input.Durability() {
		result = result.WithDurability(durability)
		cost += 2
	}
	return result, cost
}

func mergeEnchantments(input item.Stack, material item.Stack, result item.Stack, cost int, enchantedBook bool) (item.Stack, bool, bool, int) {
	var hasCompatible, hasIncompatible bool
	for _, enchant := range material.Enchantments() {

		enchantType := enchant.Type()
		compatible := enchantType.CompatibleWithItem(input.Item())
		if _, ok := input.Item().(item.EnchantedBook); ok {
			compatible = true
		}

		for _, otherEnchant := range input.Enchantments() {
			if otherType := otherEnchant.Type(); enchantType != otherType && !enchantType.CompatibleWithEnchantment(otherType) {
				compatible = false
				cost++
			}
		}

		if !compatible {
			hasIncompatible = true
			continue
		}
		hasCompatible = true

		resultLevel := enchant.Level()
		levelCost := resultLevel

		if existingEnchant, ok := input.Enchantment(enchantType); ok {
			if existingEnchant.Level() > resultLevel || (existingEnchant.Level() == resultLevel && resultLevel == enchantType.MaxLevel()) {

				hasIncompatible = true
				continue
			} else if existingEnchant.Level() == resultLevel {

				resultLevel++
			}

			levelCost = resultLevel - existingEnchant.Level()
		}

		rarityCost := enchantType.Rarity().Cost()
		if enchantedBook {
			rarityCost = max(1, rarityCost/2)
		}

		result = result.WithEnchantments(item.NewEnchantment(enchantType, resultLevel))

		cost += rarityCost * levelCost
		if input.Count() > 1 {
			cost = 40
		}
	}
	return result, hasCompatible, hasIncompatible, cost
}
