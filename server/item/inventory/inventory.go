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

package inventory

import (
	"errors"
	"fmt"
	"github.com/df-mc/dragonfly/server/item"
	"math"
	"slices"
	"strings"
	"sync"
)

type Inventory struct {
	mu    sync.RWMutex
	h     Handler
	slots []item.Stack

	f         SlotFunc
	validator SlotValidatorFunc
}

type SlotFunc func(slot int, before, after item.Stack)

type SlotValidatorFunc func(s item.Stack, slot int) bool

var ErrSlotOutOfRange = errors.New("slot is out of range: must be in range 0 <= slot < inventory.Size()")

func New(size int, f SlotFunc) *Inventory {
	if size <= 0 {
		panic("inventory size must be at least 1")
	}
	if f == nil {
		f = func(slot int, before, after item.Stack) {}
	}
	return &Inventory{h: NopHandler{}, slots: make([]item.Stack, size), f: f, validator: func(s item.Stack, slot int) bool { return true }}
}

func (inv *Inventory) Clone(f SlotFunc) *Inventory {
	if f == nil {
		f = func(slot int, before, after item.Stack) {}
	}
	return &Inventory{h: NopHandler{}, slots: inv.Slots(), f: f, validator: func(s item.Stack, slot int) bool { return true }}
}

func (inv *Inventory) SlotFunc(f SlotFunc) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.f = f
}

func (inv *Inventory) SlotValidatorFunc(f SlotValidatorFunc) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.validator = f
}

func (inv *Inventory) Item(slot int) (item.Stack, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	inv.check()
	if !inv.validSlot(slot) {
		return item.Stack{}, ErrSlotOutOfRange
	}
	return inv.slots[slot], nil
}

func (inv *Inventory) SetItem(slot int, item item.Stack) error {
	inv.mu.Lock()

	inv.check()
	if !inv.validSlot(slot) {
		inv.mu.Unlock()
		return ErrSlotOutOfRange
	}
	f := inv.setItem(slot, item)

	inv.mu.Unlock()

	f()
	return nil
}

func (inv *Inventory) Slots() []item.Stack {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	return slices.Clone(inv.slots)
}

func (inv *Inventory) Items() []item.Stack {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	items := make([]item.Stack, 0, len(inv.slots))
	for _, it := range inv.slots {
		if !it.Empty() {
			items = append(items, it)
		}
	}
	return items
}

func (inv *Inventory) First(item item.Stack) (int, bool) {
	return inv.FirstFunc(item.Comparable)
}

func (inv *Inventory) FirstFunc(comparable func(stack item.Stack) bool) (int, bool) {
	for slot, it := range inv.Slots() {
		if !it.Empty() && comparable(it) {
			return slot, true
		}
	}
	return -1, false
}

func (inv *Inventory) FirstEmpty() (int, bool) {
	for slot, it := range inv.Slots() {
		if it.Empty() {
			return slot, true
		}
	}
	return -1, false
}

func (inv *Inventory) Swap(slotA, slotB int) error {
	inv.mu.Lock()

	inv.check()
	if !inv.validSlot(slotA) || !inv.validSlot(slotB) {
		inv.mu.Unlock()
		return ErrSlotOutOfRange
	}
	a, b := inv.slots[slotA], inv.slots[slotB]
	fa, fb := inv.setItem(slotA, b), inv.setItem(slotB, a)

	inv.mu.Unlock()

	fa()
	fb()
	return nil
}

func (inv *Inventory) AddItem(it item.Stack) (n int, err error) {
	if it.Empty() {
		return 0, nil
	}
	first := it.Count()
	emptySlots := make([]int, 0, 16)

	inv.mu.Lock()

	inv.check()
	for slot, invIt := range inv.slots {
		if invIt.Empty() {

			emptySlots = append(emptySlots, slot)
			continue
		}
		a, b := invIt.AddStack(it)
		if it.Count() == b.Count() {

			continue
		}
		f := inv.setItem(slot, a)

		defer f()

		if it = b; it.Empty() {
			inv.mu.Unlock()

			return first, nil
		}
	}
	for _, slot := range emptySlots {
		a, b := it.Grow(-math.MaxInt32).AddStack(it)

		f := inv.setItem(slot, a)

		defer f()

		if it = b; it.Empty() {
			inv.mu.Unlock()

			return first, nil
		}
	}
	inv.mu.Unlock()

	return first - it.Count(), fmt.Errorf("could not add full item stack to inventory")
}

func (inv *Inventory) RemoveItem(it item.Stack) error {
	return inv.RemoveItemFunc(it.Count(), it.Comparable)
}

func (inv *Inventory) RemoveItemFunc(n int, comparable func(stack item.Stack) bool) error {
	inv.mu.Lock()
	inv.check()
	for slot, slotIt := range inv.slots {
		if slotIt.Empty() || !comparable(slotIt) {
			continue
		}
		c := slotIt.Count() - n

		var f func()
		if c <= 0 {
			f = inv.setItem(slot, item.Stack{})
		} else {
			f = inv.setItem(slot, slotIt.Grow(-n))
		}

		defer f()

		if n -= slotIt.Count(); n <= 0 {
			break
		}
	}
	inv.mu.Unlock()

	if n > 0 {
		return fmt.Errorf("could not remove all items from the inventory")
	}
	return nil
}

func (inv *Inventory) ContainsItem(it item.Stack) bool {
	return inv.ContainsItemFunc(it.Count(), it.Comparable)
}

func (inv *Inventory) ContainsItemFunc(n int, comparable func(stack item.Stack) bool) bool {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inv.check()
	for _, slotIt := range inv.slots {
		if !slotIt.Empty() && comparable(slotIt) {
			if n -= slotIt.Count(); n <= 0 {
				break
			}
		}
	}
	return n <= 0
}

func (inv *Inventory) Merge(inv2 *Inventory, f func(int, item.Stack, item.Stack)) *Inventory {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	inv2.mu.RLock()
	defer inv2.mu.RUnlock()

	n := New(len(inv.slots)+len(inv2.slots), f)
	n.slots = make([]item.Stack, 0, len(inv.slots)+len(inv2.slots))
	n.slots = append(n.slots, inv.slots...)
	n.slots = append(n.slots, inv2.slots...)
	return n
}

func (inv *Inventory) Empty() bool {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	inv.check()
	for _, it := range inv.slots {
		if !it.Empty() {
			return false
		}
	}
	return true
}

func (inv *Inventory) Clear() []item.Stack {
	inv.mu.Lock()

	inv.check()

	items := make([]item.Stack, 0, inv.size())
	for slot, i := range inv.slots {
		if !i.Empty() {
			items = append(items, i)
			f := inv.setItem(slot, item.Stack{})

			defer f()
		}
	}
	inv.mu.Unlock()

	return items
}

func (inv *Inventory) Handle(h Handler) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inv.check()
	if h == nil {
		h = NopHandler{}
	}
	inv.h = h
}

func (inv *Inventory) Handler() Handler {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	inv.check()
	return inv.h
}

func (inv *Inventory) setItem(slot int, it item.Stack) func() {
	if !inv.validator(it, slot) {
		return func() {}
	}
	if it.Count() > it.MaxCount() {
		it = it.Grow(it.MaxCount() - it.Count())
	}
	before := inv.slots[slot]
	inv.slots[slot] = it
	return func() {
		inv.f(slot, before, it)
	}
}

func (inv *Inventory) Size() int {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	return inv.size()
}

func (inv *Inventory) size() int {
	return len(inv.slots)
}

func (inv *Inventory) Close() error {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inv.check()
	inv.f = func(int, item.Stack, item.Stack) {}
	return nil
}

func (inv *Inventory) String() string {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	s := make([]string, 0, inv.size())
	for _, it := range inv.slots {
		s = append(s, it.String())
	}
	return "(" + strings.Join(s, ", ") + ")"
}

func (inv *Inventory) validSlot(slot int) bool {
	return slot >= 0 && slot < inv.size()
}

func (inv *Inventory) check() {
	if inv.size() == 0 {
		panic("uninitialised inventory: inventory must be constructed using inventory.New()")
	}
}
