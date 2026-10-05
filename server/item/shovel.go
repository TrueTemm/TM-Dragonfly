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

type Shovel struct {
	Tier ToolTier
}

func (s Shovel) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	if b, ok := tx.Block(pos).(shovellable); ok {
		if res, ok := b.Shovel(); ok {
			if face == cube.FaceDown {

				return false
			}
			if tx.Block(pos.Side(cube.FaceUp)) != air() {

				return false
			}
			tx.SetBlock(pos, res, nil)
			tx.PlaySound(pos.Vec3(), sound.ItemUseOn{Block: res})

			ctx.DamageItem(1)
			return true
		}
	}
	return false
}

type shovellable interface {
	Shovel() (world.Block, bool)
}

func (s Shovel) MaxCount() int {
	return 1
}

func (s Shovel) AttackDamage() float64 {
	return s.Tier.BaseAttackDamage
}

func (s Shovel) ToolType() ToolType {
	return TypeShovel
}

func (s Shovel) HarvestLevel() int {
	return s.Tier.HarvestLevel
}

func (s Shovel) BaseMiningEfficiency(world.Block) float64 {
	return s.Tier.BaseMiningEfficiency
}

func (s Shovel) EnchantmentValue() int {
	return s.Tier.EnchantmentValue
}

func (s Shovel) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    s.Tier.Durability,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 2,
		BreakDurability:  1,
	}
}

func (s Shovel) SmeltInfo() SmeltInfo {
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

func (s Shovel) FuelInfo() FuelInfo {
	if s.Tier == ToolTierWood {
		return newFuelInfo(time.Second * 10)
	}
	return FuelInfo{}
}

func (s Shovel) RepairableBy(i Stack) bool {
	return toolTierRepairable(s.Tier)(i)
}

func (s Shovel) EncodeItem() (name string, meta int16) {
	return "minecraft:" + s.Tier.Name + "_shovel", 0
}
