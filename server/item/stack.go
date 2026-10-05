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

package item

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/df-mc/dragonfly/server/world"
)

type Stack struct {
	id int32

	item  world.Item
	count int

	customName string
	lore       []string

	damage      int
	unbreakable bool

	anvilCost int

	data map[string]any

	enchantments map[EnchantmentType]Enchantment
}

func NewStack(t world.Item, count int) Stack {
	if count < 0 {
		panic("cannot use negative count for item stack")
	}
	if t == nil {
		panic("cannot have a stack with item type nil")
	}
	return Stack{item: t, count: count, id: newID()}
}

func (s Stack) Count() int {
	return s.count
}

func (s Stack) MaxCount() int {
	if counter, ok := s.item.(MaxCounter); ok {
		return counter.MaxCount()
	}
	return 64
}

func (s Stack) Grow(n int) Stack {
	s.count += n
	if s.count < 0 {
		s.count = 0
	}
	s.id = newID()
	return s
}

func (s Stack) Durability() int {
	if durable, ok := s.Item().(Durable); ok {
		return durable.DurabilityInfo().MaxDurability - s.damage
	}
	return -1
}

func (s Stack) MaxDurability() int {
	if durable, ok := s.Item().(Durable); ok {
		return durable.DurabilityInfo().MaxDurability
	}
	return -1
}

func (s Stack) Damage(d int) Stack {
	durable, ok := s.Item().(Durable)
	if !ok || s.unbreakable {
		return s
	}
	durability := s.Durability()
	info := durable.DurabilityInfo()
	if durability-d <= 0 {
		if info.Persistent {

			return s
		}

		return info.BrokenItem()
	}
	if durability-d > info.MaxDurability {

		s.damage, d = 0, 0
	}
	s.damage += d
	return s
}

func (s Stack) WithDurability(d int) Stack {
	durable, ok := s.Item().(Durable)
	if !ok {
		return s
	}
	maxDurability := durable.DurabilityInfo().MaxDurability
	if d > maxDurability {

		s.damage = 0
		return s
	}
	if d == 0 {

		return durable.DurabilityInfo().BrokenItem()
	}
	s.damage = maxDurability - d
	return s
}

func (s Stack) Unbreakable() bool {
	return s.unbreakable
}

func (s Stack) AsUnbreakable() Stack {
	if _, ok := s.Item().(Durable); !ok {
		return s
	}
	s.unbreakable = true
	return s
}

func (s Stack) AsBreakable() Stack {
	if _, ok := s.Item().(Durable); !ok {
		return s
	}
	s.unbreakable = false
	return s
}

func (s Stack) Empty() bool {
	if s.Count() == 0 || s.item == nil {
		return true
	}
	name, _ := s.item.EncodeItem()
	return name == "minecraft:air"
}

func (s Stack) Item() world.Item {
	if s.Empty() || s.item == nil {
		return nil
	}
	return s.item
}

func (s Stack) AttackDamage() float64 {
	if weapon, ok := s.Item().(Weapon); ok {

		return weapon.AttackDamage()
	}
	return 1.0
}

func (s Stack) WithCustomName(a ...any) Stack {
	s.customName = format(a)
	if nameable, ok := s.Item().(nameable); ok {
		s.item = nameable.WithName(a...)
	}
	return s
}

func (s Stack) CustomName() string {
	return s.customName
}

func (s Stack) WithLore(lines ...string) Stack {
	s.lore = lines
	return s
}

func (s Stack) Lore() []string {
	if s.Empty() {
		return nil
	}
	return s.lore
}

func (s Stack) WithValue(key string, val any) Stack {
	s.data = cloneMap(s.data)
	if val != nil {
		s.data[key] = val
	} else {
		delete(s.data, key)
		if len(s.data) == 0 {
			s.data = nil
		}
	}
	return s
}

func (s Stack) Value(key string) (val any, ok bool) {
	if s.Empty() {
		return nil, false
	}
	val, ok = s.data[key]
	return val, ok
}

func (s Stack) WithEnchantments(enchants ...Enchantment) Stack {
	if _, ok := s.item.(Book); ok {
		s.item = EnchantedBook{}
	}
	s.enchantments = cloneMap(s.enchantments)
	for _, enchant := range enchants {
		if _, ok := s.Item().(EnchantedBook); !ok && !enchant.t.CompatibleWithItem(s.item) {

			continue
		}
		compatible := true
		for _, otherEnchant := range s.enchantments {
			addingType := enchant.t
			existingType := otherEnchant.Type()
			addingAcceptsExisting := addingType.CompatibleWithEnchantment(existingType)
			existingAcceptsAdding := existingType.CompatibleWithEnchantment(addingType)
			if addingType != existingType && (!addingAcceptsExisting || !existingAcceptsAdding) {
				compatible = false
				break
			}
		}
		if !compatible {

			continue
		}
		s.enchantments[enchant.t] = enchant
	}
	return s
}

func (s Stack) WithForcedEnchantments(enchants ...Enchantment) Stack {
	s.enchantments = cloneMap(s.enchantments)
	for _, enchant := range enchants {
		s.enchantments[enchant.t] = enchant
	}
	return s
}

func (s Stack) WithoutEnchantments(enchants ...EnchantmentType) Stack {
	s.enchantments = cloneMap(s.enchantments)
	for _, enchant := range enchants {
		delete(s.enchantments, enchant)
	}
	if _, ok := s.item.(EnchantedBook); ok && len(s.enchantments) == 0 {
		s.item = Book{}
	}
	return s
}

func (s Stack) Enchantment(enchant EnchantmentType) (Enchantment, bool) {
	if s.Empty() {
		return Enchantment{}, false
	}
	ench, ok := s.enchantments[enchant]
	return ench, ok
}

func (s Stack) Enchantments() []Enchantment {
	if s.Empty() {
		return nil
	}
	e := slices.Collect(maps.Values(s.enchantments))
	sort.Slice(e, func(i, j int) bool {
		id1, _ := EnchantmentID(e[i].t)
		id2, _ := EnchantmentID(e[j].t)
		return id1 < id2
	})
	return e
}

func (s Stack) AnvilCost() int {
	return s.anvilCost
}

func (s Stack) WithAnvilCost(anvilCost int) Stack {
	i := s.Item()
	_, repairable := i.(Repairable)
	_, enchantedBook := i.(EnchantedBook)
	if !repairable && !enchantedBook {

		return s
	}
	s.anvilCost = anvilCost
	return s
}

func (s Stack) WithItem(t world.Item) Stack {
	cp := NewStack(t, s.count).
		Damage(s.damage).
		WithCustomName(s.customName).
		WithLore(s.lore...).
		WithEnchantments(s.Enchantments()...).
		WithAnvilCost(s.anvilCost)
	cp.unbreakable = s.unbreakable && s.MaxDurability() != -1
	cp.data = s.data
	return cp
}

func (s Stack) AddStack(s2 Stack) (a, b Stack) {
	if s.Count() >= s.MaxCount() {

		return s, s2
	}
	if !s.Comparable(s2) {

		return s, s2
	}
	diff := s.MaxCount() - s.Count()
	if s2.Count() < diff {
		diff = s2.Count()
	}

	s.count, s2.count = s.count+diff, s2.count-diff
	s.id, s2.id = newID(), newID()
	return s, s2
}

func (s Stack) Equal(s2 Stack) bool {
	return s.Comparable(s2) && s.count == s2.count && s.damage == s2.damage
}

func (s Stack) Comparable(s2 Stack) bool {
	if s.Empty() || s2.Empty() {
		return true
	}

	name, meta := s.Item().EncodeItem()
	name2, meta2 := s2.Item().EncodeItem()
	if name != name2 || meta != meta2 || s.anvilCost != s2.anvilCost || s.customName != s2.customName {
		return false
	}
	for !slices.Equal(s.lore, s2.lore) {
		return false
	}
	if len(s.enchantments) != len(s2.enchantments) {
		return false
	}
	for i := range s.enchantments {
		if s.enchantments[i] != s2.enchantments[i] {
			return false
		}
	}
	if !reflect.DeepEqual(s.data, s2.data) {
		return false
	}
	if nbt, ok := s.Item().(world.NBTer); ok {
		nbt2, ok := s2.Item().(world.NBTer)
		return ok && reflect.DeepEqual(nbt.EncodeNBT(), nbt2.EncodeNBT())
	}
	return true
}

func (s Stack) String() string {
	if s.item == nil {
		return fmt.Sprintf("Stack<nil> x%v", s.count)
	}
	return fmt.Sprintf("Stack<%T%+v>(custom name='%v', lore='%v', damage=%v, anvilCost=%v) x%v", s.item, s.item, s.customName, s.lore, s.damage, s.anvilCost, s.count)
}

func (s Stack) Values() map[string]any {
	if s.Empty() {
		return nil
	}
	return maps.Clone(s.data)
}

func cloneMap[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		m = make(M)
	}
	return maps.Clone(m)
}

var stackID = new(int32)

func newID() int32 {
	return atomic.AddInt32(stackID, 1)
}

func id(s Stack) int32 {
	if s.Empty() {
		return 0
	}
	return s.id
}

func format(a []any) string {
	return strings.TrimSuffix(fmt.Sprintln(a...), "\n")
}
