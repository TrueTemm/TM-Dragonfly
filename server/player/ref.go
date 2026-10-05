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
	"context"
	"time"

	"github.com/df-mc/dragonfly/server/world"
)

type Ref = world.EntityRef[*Player]

func NewRef(h *world.EntityHandle) Ref { return world.NewEntityRef[*Player](h) }

func Do(h *world.EntityHandle, f func(tx *world.Tx, p *Player)) *world.Task {
	return NewRef(h).Do(f)
}

func DoAfter(h *world.EntityHandle, delay time.Duration, f func(tx *world.Tx, p *Player)) *world.Task {
	return NewRef(h).DoAfter(delay, f)
}

func Call[T any](ctx context.Context, h *world.EntityHandle, f func(tx *world.Tx, p *Player) (T, error)) (T, error) {
	return world.CallRef(ctx, NewRef(h), f)
}

func (p *Player) Do(f func(tx *world.Tx, p *Player)) *world.Task {
	if p == nil {
		return world.NewFinishedTask(world.ErrEntityClosed)
	}
	return Do(p.handle, f)
}

func (p *Player) DoAfter(delay time.Duration, f func(tx *world.Tx, p *Player)) *world.Task {
	if p == nil {
		return world.NewFinishedTask(world.ErrEntityClosed)
	}
	return DoAfter(p.handle, delay, f)
}
