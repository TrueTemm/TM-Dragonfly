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

type Hoe struct {
	Tier ToolTier
}

func (h Hoe) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	if b, ok := tx.Block(pos).(tillable); ok {
		if res, ok := b.Till(); ok {
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

type tillable interface {
	Till() (world.Block, bool)
}

func (h Hoe) MaxCount() int {
	return 1
}

func (h Hoe) AttackDamage() float64 {
	return h.Tier.BaseAttackDamage + 1
}

func (h Hoe) ToolType() ToolType {
	return TypeHoe
}

func (h Hoe) HarvestLevel() int {
	return h.Tier.HarvestLevel
}

func (h Hoe) BaseMiningEfficiency(world.Block) float64 {
	return h.Tier.BaseMiningEfficiency
}

func (h Hoe) EnchantmentValue() int {
	return h.Tier.EnchantmentValue
}

func (h Hoe) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability:    h.Tier.Durability,
		BrokenItem:       simpleItem(Stack{}),
		AttackDurability: 2,
		BreakDurability:  1,
	}
}

func (h Hoe) SmeltInfo() SmeltInfo {
	switch h.Tier {
	case ToolTierIron:
		return newOreSmeltInfo(NewStack(IronNugget{}, 1), 0.1)
	case ToolTierGold:
		return newOreSmeltInfo(NewStack(GoldNugget{}, 1), 0.1)
	case ToolTierCopper:
		return newOreSmeltInfo(NewStack(CopperNugget{}, 1), 0.1)
	}
	return SmeltInfo{}
}

func (h Hoe) FuelInfo() FuelInfo {
	if h.Tier == ToolTierWood {
		return newFuelInfo(time.Second * 10)
	}
	return FuelInfo{}
}

func (h Hoe) RepairableBy(i Stack) bool {
	return toolTierRepairable(h.Tier)(i)
}

func (h Hoe) EncodeItem() (name string, meta int16) {
	return "minecraft:" + h.Tier.Name + "_hoe", 0
}
