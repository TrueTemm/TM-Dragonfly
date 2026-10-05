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

type DriedKelp struct{}

func (DriedKelp) AlwaysConsumable() bool {
	return false
}

func (DriedKelp) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration / 2
}

func (DriedKelp) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(1, 0.2)
	return Stack{}
}

func (DriedKelp) CompostChance() float64 {
	return 0.3
}

func (DriedKelp) EncodeItem() (name string, meta int16) {
	return "minecraft:dried_kelp", 0
}
