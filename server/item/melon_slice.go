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
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type MelonSlice struct{}

func (m MelonSlice) AlwaysConsumable() bool {
	return false
}

func (m MelonSlice) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}

func (m MelonSlice) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(2, 1.2)
	return Stack{}
}

func (MelonSlice) CompostChance() float64 {
	return 0.5
}

func (m MelonSlice) EncodeItem() (name string, meta int16) {
	return "minecraft:melon_slice", 0
}
