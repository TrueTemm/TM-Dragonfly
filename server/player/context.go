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

package player

import (
	"github.com/df-mc/dragonfly/server/world"
)

type Context struct {
	*world.Context
	p *Player
}

func NewEventContext(tx *world.Tx, p *Player) *Context {
	if tx == nil || p == nil || p.tx != tx {
		panic("player: transaction and player do not belong to the same callback")
	}
	_ = tx.World()
	return &Context{Context: tx.Event(), p: p}
}

func (ctx *Context) Player() *Player { return ctx.p }

func (ctx *Context) Defer(f func(ctx *Context)) *world.Task {
	return ctx.DeferErr(func(ctx *Context) error {
		f(ctx)
		return nil
	})
}

func (ctx *Context) DeferErr(f func(ctx *Context) error) *world.Task {
	h := ctx.p.H()
	return ctx.Context.DeferErr(func(tx *world.Tx) error {
		if e, ok := h.Entity(tx); ok {
			return f(NewEventContext(tx, e.(*Player)))
		}
		if h.Closed() {
			return world.ErrEntityClosed
		}
		return world.ErrEntityNotInWorld
	})
}
