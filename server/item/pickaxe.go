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
	"time"
)

type Pickaxe struct {
	Tier ToolTier
}

func (p Pickaxe) ToolType() ToolType {
	return TypePickaxe
}

func (p Pickaxe) HarvestLevel() int {
	return p.Tier.HarvestLevel
}

func (p Pickaxe) BaseMiningEfficiency(world.Block) float64 {
	return p.Tier.BaseMiningEfficiency
}

func (p Pickaxe) MaxCount() int {
	return 1
}

func (p Pickaxe) AttackDamage() float64 {
	return p.Tier.BaseAttackDamage + 1
}

func (p Pickaxe) EnchantmentValue() int {
	return p.Tier.EnchantmentValue
}

func (p Pickaxe) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    p.Tier.Durability,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 2,
		BreakDurability:  1,
	}
}

func (p Pickaxe) RepairableBy(i Stack) bool {
	return toolTierRepairable(p.Tier)(i)
}

func (p Pickaxe) SmeltInfo() SmeltInfo {
	switch p.Tier {
	case ToolTierIron:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ToolTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ToolTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (p Pickaxe) FuelInfo() FuelInfo {
	if p.Tier == ToolTierWood {
		return newFuelInfo(time.Second * 10)
	}
	return FuelInfo{}
}

func (p Pickaxe) EncodeItem() (name string, meta int16) {
	return "minecraft:" + p.Tier.Name + "_pickaxe", 0
}
