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
)

type SuspiciousStew struct {
	defaultFood

	Type StewType
}

func (SuspiciousStew) MaxCount() int {
	return 1
}

func (SuspiciousStew) AlwaysConsumable() bool {
	return true
}

func (s SuspiciousStew) EncodeItem() (name string, meta int16) {
	return "minecraft:suspicious_stew", int16(s.Type.Uint8())
}

func (s SuspiciousStew) Consume(_ *world.Tx, c Consumer) Stack {
	for _, effect := range s.Type.Effects() {
		c.AddEffect(effect)
	}
	c.Saturate(6, 7.2)

	return NewStack(Bowl{}, 1)
}
