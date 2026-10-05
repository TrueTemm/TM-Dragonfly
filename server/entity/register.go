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
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
)

var DefaultRegistry = conf.New([]world.EntityType{
	AreaEffectCloudType,
	ArrowType,
	BottleOfEnchantingType,
	EggType,
	EndCrystalType,
	EnderPearlType,
	FishingHookType,
	ExperienceOrbType,
	FallingBlockType,
	FireworkType,
	ItemType,
	LightningType,
	LingeringPotionType,
	SnowballType,
	SplashPotionType,
	TNTType,
	TextType,
})

var conf = world.EntityRegistryConfig{
	TNT:                NewTNT,
	Egg:                NewEgg,
	EndCrystal:         NewEndCrystal,
	Snowball:           NewSnowball,
	BottleOfEnchanting: NewBottleOfEnchanting,
	EnderPearl:         NewEnderPearl,
	FishingHook:        NewFishingHook,
	FallingBlock:       NewFallingBlock,
	Lightning:          NewLightning,
	Firework: func(opts world.EntitySpawnOpts, firework world.Item, owner world.Entity, sidewaysVelocityMultiplier, upwardsAcceleration float64, attached bool) *world.EntityHandle {
		return newFirework(opts, firework.(item.Firework), owner, sidewaysVelocityMultiplier, upwardsAcceleration, attached)
	},
	Item: func(opts world.EntitySpawnOpts, it any) *world.EntityHandle {
		return NewItem(opts, it.(item.Stack))
	},
	LingeringPotion: func(opts world.EntitySpawnOpts, t any, owner world.Entity) *world.EntityHandle {
		return NewLingeringPotion(opts, t.(potion.Potion), owner)
	},
	SplashPotion: func(opts world.EntitySpawnOpts, t any, owner world.Entity) *world.EntityHandle {
		return NewSplashPotion(opts, t.(potion.Potion), owner)
	},
	Arrow: func(opts world.EntitySpawnOpts, arrow world.ArrowSpawnConfig) *world.EntityHandle {
		tip := arrow.Tip.(potion.Potion)
		conf := arrowConf
		conf.Damage, conf.Potion, conf.Owner = arrow.Damage, tip, arrow.Owner.H()
		conf.KnockBackForceAddend = float64(arrow.PunchLevel) * enchantment.Punch.KnockBackMultiplier()
		conf.DisablePickup = arrow.DisablePickup
		if arrow.ObtainArrowOnPickup {
			conf.PickupItem = item.NewStack(item.Arrow{Tip: tip}, 1)
		}
		conf.Critical = arrow.Critical
		conf.PiercingLevel = arrow.PiercingLevel
		return opts.New(ArrowType, conf)
	},
}
