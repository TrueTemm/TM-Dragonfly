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
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type GoldenApple struct{}

func (e GoldenApple) AlwaysConsumable() bool {
	return true
}

func (e GoldenApple) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}

func (e GoldenApple) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(4, 9.6)
	prev := c.Absorption()
	c.AddEffect(effect.New(effect.Absorption, 1, 2*time.Minute))
	c.SetAbsorption(max(prev, min(prev+4, 16)))
	c.AddEffect(effect.New(effect.Regeneration, 2, 5*time.Second))
	return Stack{}
}

func (e GoldenApple) EncodeItem() (name string, meta int16) {
	return "minecraft:golden_apple", 0
}
