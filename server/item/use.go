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

type UseContext struct {
	Damage int

	CountSub int

	IgnoreBBox bool

	NewItem Stack

	ConsumedItems []Stack

	NewItemSurvivalOnly bool

	FirstFunc func(comparable func(Stack) bool) (Stack, bool)

	SwapHeldWithArmour func(i int)
}

func (ctx *UseContext) Consume(s Stack) {
	ctx.ConsumedItems = append(ctx.ConsumedItems, s)
}

func (ctx *UseContext) DamageItem(d int) { ctx.Damage += d }

func (ctx *UseContext) SubtractFromCount(d int) { ctx.CountSub += d }
