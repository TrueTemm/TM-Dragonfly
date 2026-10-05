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

type EnchantedApple struct{}

func (EnchantedApple) AlwaysConsumable() bool {
	return true
}

func (EnchantedApple) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}

func (EnchantedApple) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(4, 9.6)
	c.AddEffect(effect.New(effect.Absorption, 4, 2*time.Minute))
	c.AddEffect(effect.New(effect.Regeneration, 2, 30*time.Second))
	c.AddEffect(effect.New(effect.FireResistance, 1, 5*time.Minute))
	c.AddEffect(effect.New(effect.Resistance, 1, 5*time.Minute))
	return Stack{}
}

func (EnchantedApple) EncodeItem() (name string, meta int16) {
	return "minecraft:enchanted_golden_apple", 0
}
