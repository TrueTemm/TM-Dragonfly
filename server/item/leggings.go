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
	"github.com/df-mc/dragonfly/server/world"
	"image/color"
)

type Leggings struct {
	Tier ArmourTier

	Trim ArmourTrim
}

func (l Leggings) Use(_ *world.Tx, _ User, ctx *UseContext) bool {
	ctx.SwapHeldWithArmour(2)
	return false
}

func (l Leggings) MaxCount() int {
	return 1
}

func (l Leggings) DefencePoints() float64 {
	switch l.Tier.Name() {
	case "leather":
		return 2
	case "copper", "golden":
		return 3
	case "chainmail":
		return 4
	case "iron":
		return 5
	case "diamond", "netherite":
		return 6
	}
	panic("invalid leggings tier")
}

func (l Leggings) Toughness() float64 {
	return l.Tier.Toughness()
}

func (l Leggings) KnockBackResistance() float64 {
	return l.Tier.KnockBackResistance()
}

func (l Leggings) EnchantmentValue() int {
	return l.Tier.EnchantmentValue()
}

func (l Leggings) Leggings() bool {
	return true
}

func (l Leggings) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: int(l.Tier.BaseDurability() + l.Tier.BaseDurability()/2.75),
		BrokenItem:    simpleItem(Stack{}),
	}
}

func (l Leggings) RepairableBy(i Stack) bool {
	return armourTierRepairable(l.Tier)(i)
}

func (l Leggings) SmeltInfo() SmeltInfo {
	switch l.Tier.(type) {
	case ArmourTierIron, ArmourTierChain:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ArmourTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ArmourTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (l Leggings) WithTrim(trim ArmourTrim) world.Item {
	l.Trim = trim
	return l
}

func (l Leggings) EncodeItem() (name string, meta int16) {
	return "minecraft:" + l.Tier.Name() + "_leggings", 0
}

func (l Leggings) DecodeNBT(data map[string]any) any {
	if t, ok := l.Tier.(ArmourTierLeather); ok {
		if v, ok := data["customColor"].(int32); ok {
			t.Colour = rgbaFromInt32(v)
			l.Tier = t
		}
	}
	l.Trim = readTrim(data)
	return l
}

func (l Leggings) EncodeNBT() map[string]any {
	m := map[string]any{}
	if t, ok := l.Tier.(ArmourTierLeather); ok && t.Colour != (color.RGBA{}) {
		m["customColor"] = int32FromRGBA(t.Colour)
	}
	writeTrim(m, l.Trim)
	return m
}
