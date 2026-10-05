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

type Boots struct {
	Tier ArmourTier

	Trim ArmourTrim
}

func (b Boots) Use(_ *world.Tx, _ User, ctx *UseContext) bool {
	ctx.SwapHeldWithArmour(3)
	return false
}

func (b Boots) MaxCount() int {
	return 1
}

func (b Boots) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: int(b.Tier.BaseDurability() + b.Tier.BaseDurability()/5.5),
		BrokenItem:    simpleItem(Stack{}),
	}
}

func (b Boots) SmeltInfo() SmeltInfo {
	switch b.Tier.(type) {
	case ArmourTierIron, ArmourTierChain:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ArmourTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ArmourTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (b Boots) RepairableBy(i Stack) bool {
	return armourTierRepairable(b.Tier)(i)
}

func (b Boots) DefencePoints() float64 {
	switch b.Tier.Name() {
	case "leather", "copper", "golden", "chainmail":
		return 1
	case "iron":
		return 2
	case "diamond", "netherite":
		return 3
	}
	panic("invalid boots tier")
}

func (b Boots) Toughness() float64 {
	return b.Tier.Toughness()
}

func (b Boots) KnockBackResistance() float64 {
	return b.Tier.KnockBackResistance()
}

func (b Boots) EnchantmentValue() int {
	return b.Tier.EnchantmentValue()
}

func (b Boots) Boots() bool {
	return true
}

func (b Boots) WithTrim(trim ArmourTrim) world.Item {
	b.Trim = trim
	return b
}

func (b Boots) EncodeItem() (name string, meta int16) {
	return "minecraft:" + b.Tier.Name() + "_boots", 0
}

func (b Boots) DecodeNBT(data map[string]any) any {
	if t, ok := b.Tier.(ArmourTierLeather); ok {
		if v, ok := data["customColor"].(int32); ok {
			t.Colour = rgbaFromInt32(v)
			b.Tier = t
		}
	}
	b.Trim = readTrim(data)
	return b
}

func (b Boots) EncodeNBT() map[string]any {
	m := map[string]any{}
	if t, ok := b.Tier.(ArmourTierLeather); ok && t.Colour != (color.RGBA{}) {
		m["customColor"] = int32FromRGBA(t.Colour)
	}
	writeTrim(m, b.Trim)
	return m
}
