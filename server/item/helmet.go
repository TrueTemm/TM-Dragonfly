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

type Helmet struct {
	Tier ArmourTier

	Trim ArmourTrim
}

func (h Helmet) Use(_ *world.Tx, _ User, ctx *UseContext) bool {
	ctx.SwapHeldWithArmour(0)
	return false
}

func (h Helmet) MaxCount() int {
	return 1
}

func (h Helmet) DefencePoints() float64 {
	switch h.Tier.Name() {
	case "leather":
		return 1
	case "copper", "golden", "chainmail", "iron":
		return 2
	case "diamond", "netherite":
		return 3
	}
	panic("invalid helmet tier")
}

func (h Helmet) KnockBackResistance() float64 {
	return h.Tier.KnockBackResistance()
}

func (h Helmet) Toughness() float64 {
	return h.Tier.Toughness()
}

func (h Helmet) EnchantmentValue() int {
	return h.Tier.EnchantmentValue()
}

func (h Helmet) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: int(h.Tier.BaseDurability()),
		BrokenItem:    simpleItem(Stack{}),
	}
}

func (h Helmet) SmeltInfo() SmeltInfo {
	switch h.Tier.(type) {
	case ArmourTierIron, ArmourTierChain:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ArmourTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ArmourTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (h Helmet) RepairableBy(i Stack) bool {
	return armourTierRepairable(h.Tier)(i)
}

func (h Helmet) Helmet() bool {
	return true
}

func (h Helmet) WithTrim(trim ArmourTrim) world.Item {
	h.Trim = trim
	return h
}

func (h Helmet) EncodeItem() (name string, meta int16) {
	return "minecraft:" + h.Tier.Name() + "_helmet", 0
}

func (h Helmet) DecodeNBT(data map[string]any) any {
	if t, ok := h.Tier.(ArmourTierLeather); ok {
		if v, ok := data["customColor"].(int32); ok {
			t.Colour = rgbaFromInt32(v)
			h.Tier = t
		}
	}
	h.Trim = readTrim(data)
	return h
}

func (h Helmet) EncodeNBT() map[string]any {
	m := map[string]any{}
	if t, ok := h.Tier.(ArmourTierLeather); ok && t.Colour != (color.RGBA{}) {
		m["customColor"] = int32FromRGBA(t.Colour)
	}
	writeTrim(m, h.Trim)
	return m
}

func readTrim(m map[string]any) ArmourTrim {
	if trim, ok := m["Trim"].(map[string]any); ok {
		material, _ := trim["Material"].(string)
		pattern, _ := trim["Pattern"].(string)
		template, ok := smithingTemplateFromString(pattern)
		trimMaterial, ok2 := trimMaterialFromString(material)
		if ok && ok2 {
			return ArmourTrim{Template: template, Material: trimMaterial}
		}
	}
	return ArmourTrim{}
}

func writeTrim(m map[string]any, t ArmourTrim) {
	if !t.Zero() {
		m["Trim"] = map[string]any{
			"Material": t.Material.TrimMaterial(),
			"Pattern":  t.Template.String(),
		}
	}
}
