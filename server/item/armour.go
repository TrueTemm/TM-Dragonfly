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

import "image/color"

type (
	Armour interface {
		DefencePoints() float64

		Toughness() float64

		KnockBackResistance() float64
	}

	ArmourTier interface {
		BaseDurability() float64

		Toughness() float64

		KnockBackResistance() float64

		EnchantmentValue() int

		Name() string
	}

	HelmetType interface {
		Armour
		Helmet() bool
	}

	ChestplateType interface {
		Armour
		Chestplate() bool
	}

	LeggingsType interface {
		Armour
		Leggings() bool
	}

	BootsType interface {
		Armour
		Boots() bool
	}
)

type ArmourTierLeather struct {
	Colour color.RGBA
}

func (ArmourTierLeather) BaseDurability() float64      { return 55 }
func (ArmourTierLeather) Toughness() float64           { return 0 }
func (ArmourTierLeather) KnockBackResistance() float64 { return 0 }
func (ArmourTierLeather) EnchantmentValue() int        { return 15 }
func (ArmourTierLeather) Name() string                 { return "leather" }

type ArmourTierCopper struct{}

func (ArmourTierCopper) BaseDurability() float64      { return 121 }
func (ArmourTierCopper) Toughness() float64           { return 0 }
func (ArmourTierCopper) KnockBackResistance() float64 { return 0 }
func (ArmourTierCopper) EnchantmentValue() int        { return 8 }
func (ArmourTierCopper) Name() string                 { return "copper" }

type ArmourTierGold struct{}

func (ArmourTierGold) BaseDurability() float64      { return 77 }
func (ArmourTierGold) Toughness() float64           { return 0 }
func (ArmourTierGold) KnockBackResistance() float64 { return 0 }
func (ArmourTierGold) EnchantmentValue() int        { return 25 }
func (ArmourTierGold) Name() string                 { return "golden" }

type ArmourTierChain struct{}

func (ArmourTierChain) BaseDurability() float64      { return 165 }
func (ArmourTierChain) Toughness() float64           { return 0 }
func (ArmourTierChain) KnockBackResistance() float64 { return 0 }
func (ArmourTierChain) EnchantmentValue() int        { return 12 }
func (ArmourTierChain) Name() string                 { return "chainmail" }

type ArmourTierIron struct{}

func (ArmourTierIron) BaseDurability() float64      { return 165 }
func (ArmourTierIron) Toughness() float64           { return 0 }
func (ArmourTierIron) KnockBackResistance() float64 { return 0 }
func (ArmourTierIron) EnchantmentValue() int        { return 9 }
func (ArmourTierIron) Name() string                 { return "iron" }

type ArmourTierDiamond struct{}

func (ArmourTierDiamond) BaseDurability() float64      { return 363 }
func (ArmourTierDiamond) Toughness() float64           { return 2 }
func (ArmourTierDiamond) KnockBackResistance() float64 { return 0 }
func (ArmourTierDiamond) EnchantmentValue() int        { return 10 }
func (ArmourTierDiamond) Name() string                 { return "diamond" }

type ArmourTierNetherite struct{}

func (ArmourTierNetherite) BaseDurability() float64      { return 407 }
func (ArmourTierNetherite) Toughness() float64           { return 3 }
func (ArmourTierNetherite) KnockBackResistance() float64 { return 0.1 }
func (ArmourTierNetherite) EnchantmentValue() int        { return 15 }
func (ArmourTierNetherite) Name() string                 { return "netherite" }

func ArmourTiers() []ArmourTier {
	return []ArmourTier{ArmourTierLeather{}, ArmourTierCopper{}, ArmourTierGold{}, ArmourTierChain{}, ArmourTierIron{}, ArmourTierDiamond{}, ArmourTierNetherite{}}
}

func armourTierRepairable(tier ArmourTier) func(Stack) bool {
	return func(stack Stack) bool {
		var ok bool
		switch tier.(type) {
		case ArmourTierLeather:
			_, ok = stack.Item().(Leather)
		case ArmourTierCopper:
			_, ok = stack.Item().(CopperIngot)
		case ArmourTierGold:
			_, ok = stack.Item().(GoldIngot)
		case ArmourTierChain, ArmourTierIron:
			_, ok = stack.Item().(IronIngot)
		case ArmourTierDiamond:
			_, ok = stack.Item().(Diamond)
		case ArmourTierNetherite:
			_, ok = stack.Item().(NetheriteIngot)
		}
		return ok
	}
}
