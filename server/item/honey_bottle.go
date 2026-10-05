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
	"github.com/df-mc/dragonfly/server/world"
)

type HoneyBottle struct{}

func (HoneyBottle) MaxCount() int {
	return 16
}

func (HoneyBottle) AlwaysConsumable() bool {
	return true
}

func (HoneyBottle) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration * 5 / 4
}

func (HoneyBottle) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(6, 1.2)
	c.RemoveEffect(effect.Poison)
	return NewStack(GlassBottle{}, 1)
}

func (HoneyBottle) EncodeItem() (name string, meta int16) {
	return "minecraft:honey_bottle", 0
}
