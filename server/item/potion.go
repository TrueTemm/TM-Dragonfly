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
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type Potion struct {
	Type potion.Potion
}

func (p Potion) MaxCount() int {
	return 1
}

func (p Potion) AlwaysConsumable() bool {
	return true
}

func (p Potion) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}

func (p Potion) Consume(_ *world.Tx, c Consumer) Stack {
	for _, effect := range p.Type.Effects() {
		c.AddEffect(effect)
	}
	return NewStack(GlassBottle{}, 1)
}

func (p Potion) EncodeItem() (name string, meta int16) {
	return "minecraft:potion", int16(p.Type.Uint8())
}
