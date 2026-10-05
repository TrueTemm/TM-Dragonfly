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
	"math/rand/v2"
	"time"
)

type Chicken struct {
	defaultFood

	Cooked bool
}

func (c Chicken) Consume(_ *world.Tx, co Consumer) Stack {
	if c.Cooked {
		co.Saturate(6, 7.2)
	} else {
		co.Saturate(2, 1.2)
		if rand.Float64() < 0.3 {
			co.AddEffect(effect.New(effect.Hunger, 1, 30*time.Second))
		}
	}
	return Stack{}
}

func (c Chicken) SmeltInfo() SmeltInfo {
	if c.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Chicken{Cooked: true}, 1), 0.35)
}

func (c Chicken) EncodeItem() (name string, meta int16) {
	if c.Cooked {
		return "minecraft:cooked_chicken", 0
	}
	return "minecraft:chicken", 0
}
