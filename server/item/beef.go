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

type Beef struct {
	defaultFood

	Cooked bool
}

func (b Beef) Consume(_ *world.Tx, c Consumer) Stack {
	if b.Cooked {
		c.Saturate(8, 12.8)
	} else {
		c.Saturate(3, 1.8)
	}
	return Stack{}
}

func (b Beef) SmeltInfo() SmeltInfo {
	if b.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Beef{Cooked: true}, 1), 0.35)
}

func (b Beef) EncodeItem() (name string, meta int16) {
	if b.Cooked {
		return "minecraft:cooked_beef", 0
	}
	return "minecraft:beef", 0
}
