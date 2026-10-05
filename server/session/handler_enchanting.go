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
	"math"
	"math/rand/v2"
	"slices"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	enchantingInputSlot = 0x0e

	enchantingLapisSlot = 0x0f
)

func (h *ItemStackRequestHandler) handleEnchant(a *protocol.CraftRecipeStackRequestAction, s *Session, tx *world.Tx, c Controllable) error {

	if a.RecipeNetworkID > 2 {
		return fmt.Errorf("invalid recipe network id: %d", a.RecipeNetworkID)
	}

	input, err := s.ui.Item(enchantingInputSlot)
	if err != nil {
		return err
	}
	if input.Count() > 1 {
		return fmt.Errorf("enchanting tables only accept one item at a time")
	}

	allCosts, allEnchants := s.determineAvailableEnchantments(tx, c, *s.openedPos.Load(), input)
	if len(allEnchants) == 0 {
		return fmt.Errorf("can't enchant non-enchantable item")
	}

	cost := int(a.RecipeNetworkID + 1)
	requirement := allCosts[a.RecipeNetworkID]
	enchants := allEnchants[a.RecipeNetworkID]

	if !c.GameMode().CreativeInventory() {

		if c.ExperienceLevel() < requirement {
			return fmt.Errorf("not enough levels to meet requirement")
		}
		if c.ExperienceLevel() < cost {
			return fmt.Errorf("not enough levels to meet cost")
		}

		lapis, err := s.ui.Item(enchantingLapisSlot)
		if err != nil {
			return err
		}
		if _, ok := lapis.Item().(item.LapisLazuli); !ok {
			return fmt.Errorf("lapis lazuli was not input")
		}
		if lapis.Count() < cost {
			return fmt.Errorf("not enough lapis lazuli to meet cost")
		}

		c.SetExperienceLevel(c.ExperienceLevel() - cost)
		h.setItemInSlot(protocol.StackRequestSlotInfo{
			Container: protocol.FullContainerName{ContainerID: protocol.ContainerEnchantingMaterial},
			Slot:      enchantingLapisSlot,
		}, lapis.Grow(-cost), s, tx)
	}

	c.ResetEnchantmentSeed()

	h.setItemInSlot(protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerEnchantingInput},
		Slot:      enchantingInputSlot,
	}, item.Stack{}, s, tx)

	return h.createResults(s, tx, input.WithEnchantments(enchants...))
}

func (s *Session) sendEnchantmentOptions(tx *world.Tx, c Controllable, pos cube.Pos, stack item.Stack) {

	selectedCosts, selectedEnchants := s.determineAvailableEnchantments(tx, c, pos, stack)
	if len(selectedEnchants) == 0 {

		return
	}

	options := make([]protocol.EnchantmentOption, 0, 3)
	for i := 0; i < 3; i++ {

		enchants := make([]protocol.EnchantmentInstance, 0, len(selectedEnchants[i]))
		for _, enchant := range selectedEnchants[i] {
			id, _ := item.EnchantmentID(enchant.Type())
			enchants = append(enchants, protocol.EnchantmentInstance{
				Type:  byte(id),
				Level: byte(enchant.Level()),
			})
		}

		options = append(options, protocol.EnchantmentOption{
			Name:            enchantNames[rand.IntN(len(enchantNames))],
			Cost:            uint8(selectedCosts[i]),
			RecipeNetworkID: uint32(i),
			Enchantments: protocol.ItemEnchantments{
				Slot:         int32(i),
				Enchantments: [3][]protocol.EnchantmentInstance{1: enchants},
			},
		})
	}

	s.writePacket(&packet.PlayerEnchantOptions{Options: options})
}

func (s *Session) determineAvailableEnchantments(tx *world.Tx, c Controllable, pos cube.Pos, stack item.Stack) ([]int, [][]item.Enchantment) {

	enchantable, ok := stack.Item().(item.Enchantable)
	if !ok {

		return nil, nil
	}
	if len(stack.Enchantments()) > 0 {

		return nil, nil
	}

	seed := uint64(c.EnchantmentSeed())
	random := rand.New(rand.NewPCG(seed, seed))
	bookshelves := searchBookshelves(tx, pos)
	value := enchantable.EnchantmentValue()

	baseCost := random.IntN(8) + 1 + (bookshelves >> 1) + random.IntN(bookshelves+1)

	upperLevelCost := max(baseCost/3, 1)
	middleLevelCost := baseCost*2/3 + 1
	lowerLevelCost := max(baseCost, bookshelves*2)

	return []int{
		upperLevelCost,
		middleLevelCost,
		lowerLevelCost,
	}, [][]item.Enchantment{
		createEnchantments(random, stack, value, upperLevelCost),
		createEnchantments(random, stack, value, middleLevelCost),
		createEnchantments(random, stack, value, lowerLevelCost),
	}
}

type treasureEnchantment interface {
	item.EnchantmentType
	Treasure() bool
}

func createEnchantments(random *rand.Rand, stack item.Stack, value, level int) []item.Enchantment {

	randomBonus := (random.Float64() + random.Float64() - 1.0) * 0.15

	cost := level + 1 + random.IntN(value/4+1) + random.IntN(value/4+1)
	cost = clamp(int(math.Round(float64(cost)+float64(cost)*randomBonus)), 1, math.MaxInt32)

	it := stack.Item()
	_, book := it.(item.Book)

	availableEnchants := make([]item.Enchantment, 0, len(item.Enchantments()))
	for _, enchant := range item.Enchantments() {
		if t, ok := enchant.(treasureEnchantment); ok && t.Treasure() {

			continue
		}
		if !book && !enchant.CompatibleWithItem(it) {

			continue
		}

		for i := enchant.MaxLevel(); i > 0; i-- {

			if minCost, maxCost := enchant.Cost(i); cost >= minCost && cost <= maxCost {

				availableEnchants = append(availableEnchants, item.NewEnchantment(enchant, i))
				break
			}
		}
	}
	if len(availableEnchants) == 0 {

		return nil
	}

	selectedEnchants := make([]item.Enchantment, 0, len(availableEnchants))

	enchant := weightedRandomEnchantment(random, availableEnchants)
	selectedEnchants = append(selectedEnchants, enchant)

	ind := slices.Index(availableEnchants, enchant)
	availableEnchants = slices.Delete(availableEnchants, ind, ind+1)

	for random.IntN(50) <= cost {

		lastEnchant := selectedEnchants[len(selectedEnchants)-1]
		if availableEnchants = sliceutil.Filter(availableEnchants, func(enchant item.Enchantment) bool {
			return lastEnchant.Type().CompatibleWithEnchantment(enchant.Type())
		}); len(availableEnchants) == 0 {

			break
		}

		enchant = weightedRandomEnchantment(random, availableEnchants)
		selectedEnchants = append(selectedEnchants, enchant)

		ind = slices.Index(availableEnchants, enchant)
		availableEnchants = slices.Delete(availableEnchants, ind, ind+1)

		cost /= 2
	}
	return selectedEnchants
}

func searchBookshelves(tx *world.Tx, pos cube.Pos) (shelves int) {
	for x := -1; x <= 1; x++ {
		for z := -1; z <= 1; z++ {
			for y := 0; y <= 1; y++ {
				if x == 0 && z == 0 {

					continue
				}
				if _, ok := tx.Block(pos.Add(cube.Pos{x, y, z})).(block.Air); !ok {

					continue
				}

				if _, ok := tx.Block(pos.Add(cube.Pos{x * 2, y, z * 2})).(block.Bookshelf); ok {
					shelves++
				}
				if x != 0 && z != 0 {

					if _, ok := tx.Block(pos.Add(cube.Pos{x * 2, y, z})).(block.Bookshelf); ok {
						shelves++
					}

					if _, ok := tx.Block(pos.Add(cube.Pos{x, y, z * 2})).(block.Bookshelf); ok {
						shelves++
					}
				}

				if shelves >= 15 {

					return 15
				}
			}
		}
	}
	return shelves
}

func weightedRandomEnchantment(rs *rand.Rand, enchants []item.Enchantment) item.Enchantment {
	var totalWeight int
	for _, e := range enchants {
		totalWeight += e.Type().Rarity().Weight()
	}
	r := rs.IntN(totalWeight)
	for _, e := range enchants {
		r -= e.Type().Rarity().Weight()
		if r < 0 {
			return e
		}
	}
	panic("should never happen")
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
