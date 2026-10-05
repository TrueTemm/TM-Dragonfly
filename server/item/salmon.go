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

type Salmon struct {
	defaultFood

	Cooked bool
}

func (s Salmon) Consume(_ *world.Tx, c Consumer) Stack {
	if s.Cooked {
		c.Saturate(6, 9.6)
	} else {
		c.Saturate(2, 0.4)
	}
	return Stack{}
}

func (s Salmon) SmeltInfo() SmeltInfo {
	if s.Cooked {
		return SmeltInfo{}
	}
	return newFoodSmeltInfo(NewStack(Salmon{Cooked: true}, 1), 0.35)
}

func (s Salmon) EncodeItem() (name string, meta int16) {
	if s.Cooked {
		return "minecraft:cooked_salmon", 0
	}
	return "minecraft:salmon", 0
}
