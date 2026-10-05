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

package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"math"
	"math/rand/v2"
	"sync"
	"time"
)

type smelter struct {
	mu sync.Mutex

	viewers   map[ContainerViewer]struct{}
	inventory *inventory.Inventory

	remainingDuration time.Duration
	cookDuration      time.Duration
	maxDuration       time.Duration
	experience        int
}

func newSmelter() *smelter {
	s := &smelter{viewers: make(map[ContainerViewer]struct{})}
	s.inventory = inventory.New(3, func(slot int, _, item item.Stack) {
		s.mu.Lock()
		defer s.mu.Unlock()
		for viewer := range s.viewers {
			viewer.ViewSlotChange(slot, item)
		}
	})
	return s
}

func (s *smelter) InsertItem(h Hopper, pos cube.Pos, tx *world.Tx) bool {
	for sourceSlot, sourceStack := range h.inventory.Slots() {
		var slot int

		if sourceStack.Empty() {
			continue
		}

		if h.Facing != cube.FaceDown {
			slot = 1
		} else {
			slot = 0
		}

		stack := sourceStack.Grow(-sourceStack.Count() + 1)
		it, _ := s.Inventory(tx, pos).Item(slot)
		if slot == 1 {
			if fuel, ok := sourceStack.Item().(item.Fuel); !ok || fuel.FuelInfo().Duration == 0 {

				continue
			}
		}
		if !sourceStack.Comparable(it) {

			continue
		}
		if it.Count() == it.MaxCount() {

			continue
		}
		if !it.Empty() {
			stack = it.Grow(1)
		}

		_ = s.Inventory(tx, pos).SetItem(slot, stack)
		_ = h.inventory.SetItem(sourceSlot, sourceStack.Grow(-1))
		return true
	}

	return false
}

func (s *smelter) ExtractItem(h Hopper, pos cube.Pos, tx *world.Tx) bool {
	for sourceSlot, sourceStack := range s.inventory.Slots() {
		if sourceStack.Empty() {
			continue
		}

		if sourceSlot == 0 {
			continue
		}

		if sourceSlot == 1 {
			fuel, ok := sourceStack.Item().(item.Fuel)
			if ok && fuel.FuelInfo().Duration.Seconds() != 0 {
				continue
			}
		}

		_, err := h.inventory.AddItem(sourceStack.Grow(-sourceStack.Count() + 1))
		if err != nil {

			continue
		}

		_ = s.Inventory(tx, pos).SetItem(sourceSlot, sourceStack.Grow(-1))
		return true
	}

	return false
}

func (s *smelter) Durations() (remaining time.Duration, max time.Duration, cook time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remainingDuration, s.maxDuration, s.cookDuration
}

func (s *smelter) Experience() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.experience
}

func (s *smelter) ResetExperience() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	xp := s.experience
	s.experience = 0
	return xp
}

func (s *smelter) Inventory(*world.Tx, cube.Pos) *inventory.Inventory {
	return s.inventory
}

func (s *smelter) AddViewer(v ContainerViewer, _ *world.Tx, _ cube.Pos) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.viewers[v] = struct{}{}
}

func (s *smelter) RemoveViewer(v ContainerViewer, _ *world.Tx, _ cube.Pos) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.viewers) == 0 {

		return
	}
	delete(s.viewers, v)
}

func (s *smelter) setExperience(xp int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.experience = xp
}

func (s *smelter) setDurations(remaining, max, cook time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.remainingDuration, s.maxDuration, s.cookDuration = remaining, max, cook
}

func (s *smelter) tickSmelting(requirement, decrement time.Duration, lit bool, supported func(item.SmeltInfo) bool) bool {
	s.mu.Lock()

	prevCookDuration := s.cookDuration
	prevRemainingDuration := s.remainingDuration
	prevMaxDuration := s.maxDuration

	input, _ := s.inventory.Item(0)
	fuel, _ := s.inventory.Item(1)
	product, _ := s.inventory.Item(2)

	var inputInfo item.SmeltInfo
	if i, ok := input.Item().(item.Smeltable); ok && supported(i.SmeltInfo()) {
		inputInfo = i.SmeltInfo()
	}

	var fuelInfo item.FuelInfo
	if f, ok := fuel.Item().(item.Fuel); ok {
		fuelInfo = f.FuelInfo()
		if fuelInfo.Residue.Empty() {

			fuelInfo.Residue = fuel.Grow(-1)
		}
	}

	canSmelt := input.Count() > 0 && (inputInfo.Product.Comparable(product)) && !inputInfo.Product.Empty() && product.Count() < product.MaxCount()
	if s.remainingDuration <= 0 && canSmelt && fuelInfo.Duration > 0 && fuel.Count() > 0 {
		s.remainingDuration, s.maxDuration, lit = fuelInfo.Duration, fuelInfo.Duration, true
		defer s.inventory.SetItem(1, fuelInfo.Residue)
	}

	if s.remainingDuration > 0 {

		s.remainingDuration -= time.Millisecond * 50

		switch {
		case canSmelt:

			s.cookDuration += time.Millisecond * 50

			if s.cookDuration >= requirement {

				defer s.inventory.SetItem(0, input.Grow(-1))
				defer s.inventory.SetItem(2, item.NewStack(inputInfo.Product.Item(), product.Count()+inputInfo.Product.Count()))

				xp := inputInfo.Experience * float64(inputInfo.Product.Count())
				earned := math.Floor(inputInfo.Experience)
				if chance := xp - earned; chance > 0 && rand.Float64() < chance {
					earned++
				}

				s.cookDuration -= requirement
				s.experience += int(earned)
			}
		case s.remainingDuration == 0:

			s.maxDuration = 0
		default:

			s.cookDuration = 0
		}
	} else {

		s.maxDuration, lit = 0, false
	}

	if s.cookDuration > 0 && !lit {
		s.cookDuration -= decrement
	}

	for v := range s.viewers {
		v.ViewFurnaceUpdate(prevCookDuration, s.cookDuration, prevRemainingDuration, s.remainingDuration, prevMaxDuration, s.maxDuration)
	}

	s.mu.Unlock()
	return lit
}
