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

package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
)

type Magma struct {
	solid
	bassDrum
}

func (Magma) LightEmissionLevel() uint8 {
	return 3
}

func (Magma) EntityStepOn(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if fireProof, ok := e.(interface{ FireProof() bool }); ok && fireProof.FireProof() {
		return
	}

	if sneaking, ok := e.(interface{ Sneaking() bool }); ok && sneaking.Sneaking() {
		return
	}
	if l, ok := e.(livingEntity); ok {
		l.Hurt(1, MagmaDamageSource{})
	}
}

func (m Magma) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, pickaxeHarvestable, pickaxeEffective, oneOf(m))
}

func (Magma) EncodeItem() (name string, meta int16) {
	return "minecraft:magma", 0
}

func (Magma) EncodeBlock() (string, map[string]any) {
	return "minecraft:magma", nil
}

type MagmaDamageSource struct{}

func (MagmaDamageSource) ReducedByResistance() bool { return true }
func (MagmaDamageSource) ReducedByArmour() bool     { return true }
func (MagmaDamageSource) Fire() bool                { return true }
func (MagmaDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool {
	return e == enchantment.FireProtection
}
func (MagmaDamageSource) IgnoreTotem() bool { return false }
