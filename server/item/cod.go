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

type Cod struct {
	defaultFood

	Cooked bool
}

func (c Cod) Consume(_ *world.Tx, co Consumer) Stack {
	if c.Cooked {
		co.Saturate(5, 6)
	} else {
		co.Saturate(2, 0.4)
	}
	return Stack{}
}

func (c Cod) SmeltInfo() SmeltInfo {
	if c.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Cod{Cooked: true}, 1), 0.35)
}

func (c Cod) EncodeItem() (name string, meta int16) {
	if c.Cooked {
		return "minecraft:cooked_cod", 0
	}
	return "minecraft:cod", 0
}
