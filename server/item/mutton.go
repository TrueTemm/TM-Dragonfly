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

type Mutton struct {
	defaultFood

	Cooked bool
}

func (m Mutton) Consume(_ *world.Tx, c Consumer) Stack {
	if m.Cooked {
		c.Saturate(6, 9.6)
	} else {
		c.Saturate(2, 1.2)
	}
	return Stack{}
}

func (m Mutton) SmeltInfo() SmeltInfo {
	if m.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Mutton{Cooked: true}, 1), 0.35)
}

func (m Mutton) EncodeItem() (name string, meta int16) {
	if m.Cooked {
		return "minecraft:cooked_mutton", 0
	}
	return "minecraft:mutton", 0
}
