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

type TurtleShell struct{}

func (TurtleShell) Use(_ *world.Tx, _ User, ctx *UseContext) bool {
	ctx.SwapHeldWithArmour(0)
	return false
}

func (TurtleShell) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 276,
		BrokenItem:    simpleItem(Stack{}),
	}
}

func (TurtleShell) RepairableBy(i Stack) bool {
	_, ok := i.Item().(Scute)
	return ok
}

func (TurtleShell) MaxCount() int {
	return 1
}

func (TurtleShell) DefencePoints() float64 {
	return 2
}

func (TurtleShell) Toughness() float64 {
	return 0
}

func (TurtleShell) KnockBackResistance() float64 {
	return 0
}

func (TurtleShell) EnchantmentValue() int {
	return 9
}

func (TurtleShell) Helmet() bool {
	return true
}

func (TurtleShell) EncodeItem() (name string, meta int16) {
	return "minecraft:turtle_helmet", 0
}
