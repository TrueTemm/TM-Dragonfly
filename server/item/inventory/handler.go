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

package inventory

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
)

type Holder interface{}

type Context = event.Context[Holder]

type Handler interface {
	HandleTake(ctx *Context, slot int, it item.Stack)

	HandlePlace(ctx *Context, slot int, it item.Stack)

	HandleDrop(ctx *Context, slot int, it item.Stack)
}

var _ Handler = NopHandler{}

type NopHandler struct{}

func (NopHandler) HandleTake(*Context, int, item.Stack)  {}
func (NopHandler) HandlePlace(*Context, int, item.Stack) {}
func (NopHandler) HandleDrop(*Context, int, item.Stack)  {}
