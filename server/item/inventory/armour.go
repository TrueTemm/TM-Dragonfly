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
	"fmt"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"math"
	"math/rand/v2"
)

type Armour struct {
	inv *Inventory
}

func NewArmour(f func(slot int, before, after item.Stack)) *Armour {
	inv := New(4, f)
	inv.validator = canAddArmour
	return &Armour{inv: inv}
}

func canAddArmour(s item.Stack, slot int) bool {
	if s.Empty() {
		return true
	}
	switch slot {
	case 0:
		if h, ok := s.Item().(item.HelmetType); ok {
			return h.Helmet()
		}
	case 1:
		if c, ok := s.Item().(item.ChestplateType); ok {
			return c.Chestplate()
		}
	case 2:
		if l, ok := s.Item().(item.LeggingsType); ok {
			return l.Leggings()
		}
	case 3:
		if b, ok := s.Item().(item.BootsType); ok {
			return b.Boots()
		}
	}
	return false
}

func (a *Armour) Set(helmet, chestplate, leggings, boots item.Stack) {
	a.SetHelmet(helmet)
	a.SetChestplate(chestplate)
	a.SetLeggings(leggings)
	a.SetBoots(boots)
}

func (a *Armour) SetHelmet(helmet item.Stack) {
	_ = a.inv.SetItem(0, helmet)
}

func (a *Armour) Helmet() item.Stack {
	i, _ := a.inv.Item(0)
	return i
}

func (a *Armour) SetChestplate(chestplate item.Stack) {
	_ = a.inv.SetItem(1, chestplate)
}

func (a *Armour) Chestplate() item.Stack {
	i, _ := a.inv.Item(1)
	return i
}

func (a *Armour) SetLeggings(leggings item.Stack) {
	_ = a.inv.SetItem(2, leggings)
}

func (a *Armour) Leggings() item.Stack {
	i, _ := a.inv.Item(2)
	return i
}

func (a *Armour) SetBoots(boots item.Stack) {
	_ = a.inv.SetItem(3, boots)
}

func (a *Armour) Boots() item.Stack {
	i, _ := a.inv.Item(3)
	return i
}

func (a *Armour) DamageReduction(dmg float64, src world.DamageSource) float64 {
	var (
		original                 = dmg
		defencePoints, toughness float64
		enchantments             []item.Enchantment
	)

	for _, it := range a.Items() {
		enchantments = append(enchantments, it.Enchantments()...)
		if armour, ok := it.Item().(item.Armour); ok {
			defencePoints += armour.DefencePoints()
			toughness += armour.Toughness()
		}
	}

	dmg -= dmg * enchantment.ProtectionFactor(src, enchantments)
	if src.ReducedByArmour() {

		dmg -= dmg * 0.04 * math.Max(defencePoints*0.2, defencePoints-dmg/(2+toughness/4))
	}
	return original - dmg
}

func (a *Armour) HighestEnchantmentLevel(t item.EnchantmentType) int {
	lvl := 0
	for _, it := range a.Items() {
		if e, ok := it.Enchantment(t); ok && e.Level() > lvl {
			lvl = e.Level()
		}
	}
	return lvl
}

type DamageFunc func(s item.Stack, d int) item.Stack

func (a *Armour) Damage(dmg float64, f DamageFunc) {
	armourDamage := int(math.Max(math.Floor(dmg/4), 1))
	for slot, it := range a.Slots() {
		_ = a.inv.SetItem(slot, f(it, armourDamage))
	}
}

func (a *Armour) ThornsDamage(f DamageFunc) float64 {
	slots := a.Slots()
	dmg := 0.0

	for _, i := range slots {
		thorns, _ := i.Enchantment(enchantment.Thorns)
		if level := float64(thorns.Level()); rand.Float64() < level*0.15 {

			dmg = math.Min(dmg+float64(1+rand.IntN(4)), 4.0)
		}
	}
	if highest := a.HighestEnchantmentLevel(enchantment.Thorns); highest > 10 {

		dmg = float64(highest - 10)
	}
	if dmg > 0 {

		slot := rand.IntN(len(slots))
		_ = a.Inventory().SetItem(slot, f(slots[slot], 2))
	}
	return dmg
}

func (a *Armour) KnockBackResistance() float64 {
	resistance := 0.0
	for _, i := range a.Items() {
		if a, ok := i.Item().(item.Armour); ok {
			resistance += a.KnockBackResistance()
		}
	}
	return resistance
}

func (a *Armour) Slots() []item.Stack {
	return a.inv.Slots()
}

func (a *Armour) Items() []item.Stack {
	return a.inv.Items()
}

func (a *Armour) Clear() []item.Stack {
	return a.inv.Clear()
}

func (a *Armour) String() string {
	return fmt.Sprintf("(helmet: %v, chestplate: %v, leggings: %v, boots: %v)", a.Helmet(), a.Chestplate(), a.Leggings(), a.Boots())
}

func (a *Armour) Inventory() *Inventory {
	return a.inv
}

func (a *Armour) Handle(h Handler) {
	a.inv.Handle(h)
}

func (a *Armour) Close() error {
	return a.inv.Close()
}
