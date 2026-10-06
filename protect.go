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

package main

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type spawnProtect struct {
	player.NopHandler
	o      *ops
	name   string
	spawn  cube.Pos
	radius int
}

func (h spawnProtect) HandleBlockBreak(ctx *player.Context, pos cube.Pos, _ *[]item.Stack, _ *int) {
	if h.inside(pos) {
		ctx.Cancel()
	}
}

func (h spawnProtect) HandleBlockPlace(ctx *player.Context, pos cube.Pos, _ world.Block) {
	if h.inside(pos) {
		ctx.Cancel()
	}
}

func (h spawnProtect) inside(pos cube.Pos) bool {
	if h.radius <= 0 || h.o.is(h.name) {
		return false // off or an operator building
	}
	dx, dz := pos.X()-h.spawn.X(), pos.Z()-h.spawn.Z()
	return dx*dx+dz*dz <= h.radius*h.radius
}
