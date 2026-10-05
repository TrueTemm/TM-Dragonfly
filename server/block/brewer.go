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
	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"sync"
	"time"
)

type brewer struct {
	mu sync.Mutex

	viewers   map[ContainerViewer]struct{}
	inventory *inventory.Inventory

	duration   time.Duration
	fuelAmount int32
	fuelTotal  int32
}

func newBrewer() *brewer {
	b := &brewer{viewers: make(map[ContainerViewer]struct{})}
	b.inventory = inventory.New(5, func(slot int, _, item item.Stack) {
		b.mu.Lock()
		defer b.mu.Unlock()
		for viewer := range b.viewers {
			viewer.ViewSlotChange(slot, item)
		}
	})
	return b
}

func (b *brewer) InsertItem(h Hopper, pos cube.Pos, tx *world.Tx) bool {
	for sourceSlot, sourceStack := range h.inventory.Slots() {
		var slot int

		if sourceStack.Empty() {
			continue
		}

		if h.Facing == cube.FaceDown {
			if !recipe.ValidBrewingReagent(sourceStack.Item()) {

				continue
			}
			slot = 0
		} else if _, ok := sourceStack.Item().(item.BlazePowder); ok {
			slot = 4
		} else {
			_, okPotion := sourceStack.Item().(item.Potion)
			_, okSplash := sourceStack.Item().(item.SplashPotion)
			_, okLingering := sourceStack.Item().(item.LingeringPotion)
			_, okBottle := sourceStack.Item().(item.GlassBottle)
			if !okPotion && !okSplash && !okLingering && !okBottle {
				continue
			}
			for brewingSlot, brewingStack := range b.inventory.Slots() {
				if brewingSlot == 0 || brewingSlot == 4 {
					continue
				}
				if brewingStack.Count() == brewingStack.MaxCount() || !brewingStack.Comparable(sourceStack) {
					continue
				}

				slot = brewingSlot
				break
			}

			if slot == 0 {
				continue
			}
		}

		stack := sourceStack.Grow(-sourceStack.Count() + 1)
		it, _ := b.Inventory(tx, pos).Item(slot)

		if !sourceStack.Comparable(it) {

			continue
		}
		if it.Count() == it.MaxCount() {

			continue
		}
		if !it.Empty() {
			stack = it.Grow(1)
		}

		_ = b.Inventory(tx, pos).SetItem(slot, stack)
		_ = h.inventory.SetItem(sourceSlot, sourceStack.Grow(-1))
		return true

	}
	return false
}

func (b *brewer) ExtractItem(h Hopper, pos cube.Pos, tx *world.Tx) bool {
	for sourceSlot, sourceStack := range b.inventory.Slots() {
		if sourceStack.Empty() || sourceSlot == 0 || sourceSlot == 4 {
			continue
		}
		_, err := h.inventory.AddItem(sourceStack.Grow(-sourceStack.Count() + 1))
		if err != nil {

			continue
		}
		_ = b.Inventory(tx, pos).SetItem(sourceSlot, sourceStack.Grow(-1))
		return true
	}
	return false
}

func (b *brewer) Duration() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.duration
}

func (b *brewer) Fuel() (fuel, maxFuel int32) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fuelAmount, b.fuelTotal
}

func (b *brewer) Inventory(*world.Tx, cube.Pos) *inventory.Inventory {
	return b.inventory
}

func (b *brewer) AddViewer(v ContainerViewer, _ *world.Tx, _ cube.Pos) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.viewers[v] = struct{}{}
}

func (b *brewer) RemoveViewer(v ContainerViewer, _ *world.Tx, _ cube.Pos) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.viewers, v)
}

func (b *brewer) setDuration(duration time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.duration = duration
}

func (b *brewer) setFuel(fuel, maxFuel int32) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fuelAmount, b.fuelTotal = fuel, maxFuel
}

func (b *brewer) tickBrewing(block string, pos cube.Pos, tx *world.Tx) {
	b.mu.Lock()

	left, _ := b.inventory.Item(1)
	middle, _ := b.inventory.Item(2)
	right, _ := b.inventory.Item(3)

	prevDuration := b.duration
	prevFuelAmount := b.fuelAmount
	prevFuelTotal := b.fuelTotal

	fuel, _ := b.inventory.Item(4)

	if _, ok := fuel.Item().(item.BlazePowder); ok && b.fuelAmount <= 0 {
		defer b.inventory.SetItem(4, fuel.Grow(-1))
		b.fuelAmount, b.fuelTotal = 20, 20
	}

	ingredient, _ := b.inventory.Item(0)

	leftOutput, leftAffected := recipe.Perform(block, left.Item(), ingredient.Item())
	middleOutput, middleAffected := recipe.Perform(block, middle.Item(), ingredient.Item())
	rightOutput, rightAffected := recipe.Perform(block, right.Item(), ingredient.Item())

	if b.fuelAmount > 0 {

		if leftAffected || middleAffected || rightAffected {

			if b.duration == 0 {
				b.duration = time.Second * 20
			}
			b.duration -= time.Millisecond * 50

			if b.duration <= 0 {

				if leftAffected {
					defer b.inventory.SetItem(1, leftOutput[0])
				}
				if middleAffected {
					defer b.inventory.SetItem(2, middleOutput[0])
				}
				if rightAffected {
					defer b.inventory.SetItem(3, rightOutput[0])
				}

				defer b.inventory.SetItem(0, ingredient.Grow(-1))
				tx.PlaySound(pos.Vec3Centre(), sound.PotionBrewed{})

				b.fuelAmount--
				b.duration = 0
			}
		} else {

			b.duration = 0
		}
	} else {

		b.duration, b.fuelAmount, b.fuelTotal = 0, 0, 0
	}

	for v := range b.viewers {
		v.ViewBrewingUpdate(prevDuration, b.duration, prevFuelAmount, b.fuelAmount, prevFuelTotal, b.fuelTotal)
	}

	b.mu.Unlock()
}
