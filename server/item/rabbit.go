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

import "github.com/df-mc/dragonfly/server/world"

type Rabbit struct {
	defaultFood

	Cooked bool
}

func (r Rabbit) Consume(_ *world.Tx, c Consumer) Stack {
	if r.Cooked {
		c.Saturate(5, 6)
	} else {
		c.Saturate(3, 1.8)
	}
	return Stack{}
}

func (r Rabbit) SmeltInfo() SmeltInfo {
	if r.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Rabbit{Cooked: true}, 1), 0.35)
}

func (r Rabbit) EncodeItem() (name string, meta int16) {
	if r.Cooked {
		return "minecraft:cooked_rabbit", 0
	}
	return "minecraft:rabbit", 0
}
