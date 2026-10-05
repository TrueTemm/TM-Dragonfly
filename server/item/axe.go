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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type Axe struct {
	Tier ToolTier
}

func (a Axe) UseOnBlock(pos cube.Pos, _ cube.Face, _ mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	if s, ok := tx.Block(pos).(Strippable); ok {
		if res, so, ok := s.Strip(); ok {
			tx.SetBlock(pos, res, nil)
			tx.PlaySound(pos.Vec3(), sound.ItemUseOn{Block: res})
			if so != nil {
				tx.PlaySound(pos.Vec3(), so)
			}

			ctx.DamageItem(1)
			return true
		}
	}
	return false
}

type Strippable interface {
	Strip() (world.Block, world.Sound, bool)
}

func (a Axe) MaxCount() int {
	return 1
}

func (a Axe) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    a.Tier.Durability,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 2,
		BreakDurability:  1,
	}
}

func (a Axe) SmeltInfo() SmeltInfo {
	switch a.Tier {
	case ToolTierIron:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ToolTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ToolTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (a Axe) FuelInfo() FuelInfo {
	if a.Tier == ToolTierWood {
		return newFuelInfo(time.Second * 10)
	}
	return FuelInfo{}
}

func (a Axe) AttackDamage() float64 {
	return a.Tier.BaseAttackDamage + 2
}

func (a Axe) ToolType() ToolType {
	return TypeAxe
}

func (a Axe) HarvestLevel() int {
	return a.Tier.HarvestLevel
}

func (a Axe) BaseMiningEfficiency(world.Block) float64 {
	return a.Tier.BaseMiningEfficiency
}

func (a Axe) RepairableBy(i Stack) bool {
	return toolTierRepairable(a.Tier)(i)
}

func (a Axe) EnchantmentValue() int {
	return a.Tier.EnchantmentValue
}

func (a Axe) EncodeItem() (name string, meta int16) {
	return "minecraft:" + a.Tier.Name + "_axe", 0
}
