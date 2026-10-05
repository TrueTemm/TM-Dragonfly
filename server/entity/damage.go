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

package entity

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
)

type (
	AttackDamageSource struct {
		Attacker world.Entity
	}

	VoidDamageSource struct{}

	SuffocationDamageSource struct{}

	DrowningDamageSource struct{}

	FallDamageSource struct{}

	GlideDamageSource struct{}

	LightningDamageSource struct{}

	ProjectileDamageSource struct {
		Projectile, Owner world.Entity
	}

	ExplosionDamageSource struct {
		Source world.ExplosionSource
	}
)

func (FallDamageSource) ReducedByArmour() bool     { return false }
func (FallDamageSource) ReducedByResistance() bool { return true }
func (FallDamageSource) Fire() bool                { return false }
func (FallDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool {
	return e == enchantment.FeatherFalling
}
func (FallDamageSource) IgnoreTotem() bool                { return false }
func (GlideDamageSource) ReducedByArmour() bool           { return false }
func (GlideDamageSource) ReducedByResistance() bool       { return true }
func (GlideDamageSource) Fire() bool                      { return false }
func (GlideDamageSource) IgnoreTotem() bool               { return false }
func (LightningDamageSource) ReducedByArmour() bool       { return true }
func (LightningDamageSource) ReducedByResistance() bool   { return true }
func (LightningDamageSource) Fire() bool                  { return false }
func (LightningDamageSource) IgnoreTotem() bool           { return false }
func (AttackDamageSource) ReducedByArmour() bool          { return true }
func (AttackDamageSource) ReducedByResistance() bool      { return true }
func (AttackDamageSource) Fire() bool                     { return false }
func (AttackDamageSource) IgnoreTotem() bool              { return false }
func (VoidDamageSource) ReducedByResistance() bool        { return false }
func (VoidDamageSource) ReducedByArmour() bool            { return false }
func (VoidDamageSource) Fire() bool                       { return false }
func (VoidDamageSource) IgnoreTotem() bool                { return true }
func (SuffocationDamageSource) ReducedByResistance() bool { return false }
func (SuffocationDamageSource) ReducedByArmour() bool     { return false }
func (SuffocationDamageSource) Fire() bool                { return false }
func (SuffocationDamageSource) IgnoreTotem() bool         { return false }
func (DrowningDamageSource) ReducedByResistance() bool    { return false }
func (DrowningDamageSource) ReducedByArmour() bool        { return false }
func (DrowningDamageSource) Fire() bool                   { return false }
func (DrowningDamageSource) IgnoreTotem() bool            { return false }
func (ProjectileDamageSource) ReducedByResistance() bool  { return true }
func (ProjectileDamageSource) ReducedByArmour() bool      { return true }
func (ProjectileDamageSource) Fire() bool                 { return false }
func (ProjectileDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool {
	return e == enchantment.ProjectileProtection
}
func (ProjectileDamageSource) IgnoreTotem() bool        { return false }
func (ExplosionDamageSource) ReducedByResistance() bool { return true }
func (ExplosionDamageSource) ReducedByArmour() bool     { return true }
func (ExplosionDamageSource) Fire() bool                { return false }
func (ExplosionDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool {
	return e == enchantment.BlastProtection
}
func (ExplosionDamageSource) IgnoreTotem() bool { return false }
