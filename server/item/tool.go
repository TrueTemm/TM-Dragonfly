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
)

var (
	TypeNone = ToolType{-1}

	TypePickaxe = ToolType{0}

	TypeAxe = ToolType{1}

	TypeHoe = ToolType{2}

	TypeShovel = ToolType{3}

	TypeShears = ToolType{4}

	TypeSword = ToolType{5}

	ToolTierWood = ToolTier{HarvestLevel: 1, Durability: 59, BaseMiningEfficiency: 2, BaseAttackDamage: 1, EnchantmentValue: 15, Name: "wooden"}

	ToolTierGold = ToolTier{HarvestLevel: 1, Durability: 32, BaseMiningEfficiency: 12, BaseAttackDamage: 1, EnchantmentValue: 22, Name: "golden"}

	ToolTierStone = ToolTier{HarvestLevel: 2, Durability: 131, BaseMiningEfficiency: 4, BaseAttackDamage: 2, EnchantmentValue: 5, Name: "stone"}

	ToolTierCopper = ToolTier{HarvestLevel: 2, Durability: 190, BaseMiningEfficiency: 5, BaseAttackDamage: 2, EnchantmentValue: 13, Name: "copper"}

	ToolTierIron = ToolTier{HarvestLevel: 3, Durability: 250, BaseMiningEfficiency: 6, BaseAttackDamage: 3, EnchantmentValue: 14, Name: "iron"}

	ToolTierDiamond = ToolTier{HarvestLevel: 4, Durability: 1561, BaseMiningEfficiency: 8, BaseAttackDamage: 4, EnchantmentValue: 10, Name: "diamond"}

	ToolTierNetherite = ToolTier{HarvestLevel: 4, Durability: 2031, BaseMiningEfficiency: 9, BaseAttackDamage: 5, EnchantmentValue: 15, Name: "netherite"}
)

type (
	Tool interface {
		ToolType() ToolType

		HarvestLevel() int

		BaseMiningEfficiency(b world.Block) float64
	}

	ToolTier struct {
		HarvestLevel int

		BaseMiningEfficiency float64

		BaseAttackDamage float64

		EnchantmentValue int

		Durability int

		Name string
	}

	ToolType struct{ t }
	t        int

	ToolNone struct{}
)

func ToolTiers() []ToolTier {
	return []ToolTier{ToolTierWood, ToolTierGold, ToolTierStone, ToolTierCopper, ToolTierIron, ToolTierDiamond, ToolTierNetherite}
}

func (n ToolNone) ToolType() ToolType { return TypeNone }

func (n ToolNone) HarvestLevel() int { return 0 }

func (n ToolNone) BaseMiningEfficiency(world.Block) float64 { return 1 }

func toolTierRepairable(tier ToolTier) func(Stack) bool {
	return func(stack Stack) bool {
		switch tier {
		case ToolTierWood:
			if planks, ok := stack.Item().(interface{ RepairsWoodTools() bool }); ok {
				return planks.RepairsWoodTools()
			}
		case ToolTierStone:
			if cobblestone, ok := stack.Item().(interface{ RepairsStoneTools() bool }); ok {
				return cobblestone.RepairsStoneTools()
			}
		case ToolTierGold:
			_, ok := stack.Item().(GoldIngot)
			return ok
		case ToolTierCopper:
			_, ok := stack.Item().(CopperIngot)
			return ok
		case ToolTierIron:
			_, ok := stack.Item().(IronIngot)
			return ok
		case ToolTierDiamond:
			_, ok := stack.Item().(Diamond)
			return ok
		case ToolTierNetherite:
			_, ok := stack.Item().(NetheriteIngot)
			return ok
		}
		return false
	}
}
