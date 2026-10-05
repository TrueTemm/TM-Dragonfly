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
	"time"

	"github.com/df-mc/dragonfly/server/world"
)

type Sword struct {
	Tier ToolTier
}

func (s Sword) AttackDamage() float64 {
	return s.Tier.BaseAttackDamage + 3
}

func (s Sword) MaxCount() int {
	return 1
}

func (s Sword) ToolType() ToolType {
	return TypeSword
}

func (s Sword) HarvestLevel() int {
	return s.Tier.HarvestLevel
}

func (s Sword) EnchantmentValue() int {
	return s.Tier.EnchantmentValue
}

func (s Sword) BaseMiningEfficiency(b world.Block) float64 {
	if _, ok := b.(interface{ Cobweb() }); ok {
		return 15
	}
	return 1.5
}

func (s Sword) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    s.Tier.Durability,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 1,
		BreakDurability:  2,
	}
}

func (s Sword) SmeltInfo() SmeltInfo {
	switch s.Tier {
	case ToolTierIron:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ToolTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ToolTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (s Sword) FuelInfo() FuelInfo {
	if s.Tier == ToolTierWood {
		return newFuelInfo(time.Second * 10)
	}
	return FuelInfo{}
}

func (s Sword) RepairableBy(i Stack) bool {
	return toolTierRepairable(s.Tier)(i)
}

func (s Sword) EncodeItem() (name string, meta int16) {
	return "minecraft:" + s.Tier.Name + "_sword", 0
}
