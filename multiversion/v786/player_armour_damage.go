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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	PlayerArmourDamageFlagHelmet = 1 << iota
	PlayerArmourDamageFlagChestplate
	PlayerArmourDamageFlagLeggings
	PlayerArmourDamageFlagBoots
	PlayerArmourDamageFlagBody
)

type PlayerArmourDamage struct {
	Bitset uint8

	HelmetDamage int32

	ChestplateDamage int32

	LeggingsDamage int32

	BootsDamage int32

	BodyDamage int32
}

func (pk *PlayerArmourDamage) ID() uint32 {
	return IDPlayerArmourDamage
}

func (pk *PlayerArmourDamage) Marshal(io protocol.IO) {
	io.Uint8(&pk.Bitset)
	if pk.Bitset&PlayerArmourDamageFlagHelmet != 0 {
		io.Varint32(&pk.HelmetDamage)
	} else {
		pk.HelmetDamage = 0
	}
	if pk.Bitset&PlayerArmourDamageFlagChestplate != 0 {
		io.Varint32(&pk.ChestplateDamage)
	} else {
		pk.ChestplateDamage = 0
	}
	if pk.Bitset&PlayerArmourDamageFlagLeggings != 0 {
		io.Varint32(&pk.LeggingsDamage)
	} else {
		pk.LeggingsDamage = 0
	}
	if pk.Bitset&PlayerArmourDamageFlagBoots != 0 {
		io.Varint32(&pk.BootsDamage)
	} else {
		pk.BootsDamage = 0
	}
	if pk.Bitset&PlayerArmourDamageFlagBody != 0 {
		io.Varint32(&pk.BodyDamage)
	} else {
		pk.BodyDamage = 0
	}
}
