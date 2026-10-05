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
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type SwingArmAction struct{ action }

type HurtAction struct{ action }

type CriticalHitAction struct {
	action

	Count int
}

type EnchantedHitAction struct {
	action

	Count int
}

type DeathAction struct{ action }

type EatAction struct{ action }

type ArrowShakeAction struct {
	Duration time.Duration

	action
}

type PickedUpAction struct {
	Collector world.Entity

	action
}

type FireworkExplosionAction struct{ action }

type TotemUseAction struct{ action }

type action struct{}

func (action) EntityAction() {}
