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

type Porkchop struct {
	defaultFood

	Cooked bool
}

func (p Porkchop) Consume(_ *world.Tx, c Consumer) Stack {
	if p.Cooked {
		c.Saturate(8, 12.8)
	} else {
		c.Saturate(3, 1.8)
	}
	return Stack{}
}

func (p Porkchop) SmeltInfo() SmeltInfo {
	if p.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Porkchop{Cooked: true}, 1), 0.35)
}

func (p Porkchop) EncodeItem() (name string, meta int16) {
	if p.Cooked {
		return "minecraft:cooked_porkchop", 0
	}
	return "minecraft:porkchop", 0
}
