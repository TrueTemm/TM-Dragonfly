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

	"github.com/df-mc/dragonfly/server/entity/effect"
)

type StewType struct {
	stewType
}

func NightVisionPoppyStew() StewType {
	return StewType{0}
}

func JumpBoostStew() StewType {
	return StewType{1}
}

func WeaknessStew() StewType {
	return StewType{2}
}

func BlindnessBluetStew() StewType {
	return StewType{3}
}

func PoisonStew() StewType {
	return StewType{4}
}

func SaturationDandelionStew() StewType {
	return StewType{5}
}

func SaturationOrchidStew() StewType {
	return StewType{6}
}

func FireResistanceStew() StewType {
	return StewType{7}
}

func RegenerationStew() StewType {
	return StewType{8}
}

func WitherStew() StewType {
	return StewType{9}
}

func NightVisionTorchflowerStew() StewType {
	return StewType{10}
}

func BlindnessEyeblossomStew() StewType {
	return StewType{11}
}

func NauseaStew() StewType {
	return StewType{12}
}

func StewTypes() []StewType {
	return []StewType{NightVisionPoppyStew(), JumpBoostStew(), WeaknessStew(), BlindnessBluetStew(), PoisonStew(), SaturationDandelionStew(), SaturationOrchidStew(), FireResistanceStew(), RegenerationStew(), WitherStew(), NightVisionTorchflowerStew(), BlindnessEyeblossomStew(), NauseaStew()}
}

type stewType uint8

func (s stewType) Uint8() uint8 {
	return uint8(s)
}

func (s stewType) Effects() []effect.Effect {
	var effects []effect.Effect
	switch s.Uint8() {
	case 0, 10:
		effects = append(effects, effect.New(effect.NightVision, 1, time.Second*5))
	case 1:
		effects = append(effects, effect.New(effect.JumpBoost, 1, time.Second*5))
	case 2:
		effects = append(effects, effect.New(effect.Weakness, 1, time.Second*7))
	case 3, 11:
		effects = append(effects, effect.New(effect.Blindness, 1, time.Second*6))
	case 4:
		effects = append(effects, effect.New(effect.Poison, 1, time.Second*11))
	case 5, 6:
		effects = append(effects, effect.New(effect.Saturation, 1, time.Second*3/10))
	case 7:
		effects = append(effects, effect.New(effect.FireResistance, 1, time.Second*3))
	case 8:
		effects = append(effects, effect.New(effect.Regeneration, 1, time.Second*7))
	case 9:
		effects = append(effects, effect.New(effect.Wither, 1, time.Second*7))
	case 12:
		effects = append(effects, effect.New(effect.Nausea, 1, time.Second*7))
	default:
		panic("should never happen")
	}

	return effects
}
