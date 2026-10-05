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

package block_test

import (
	"context"
	"testing"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

func spreadOf(liquid world.Liquid, ticks int) int {
	w := world.Config{Synchronous: true}.New()
	defer func() { _ = w.Close() }()

	source := cube.Pos{0, 40, 0}
	w.Do(func(tx *world.Tx) {
		for x := -6; x <= 6; x++ {
			for z := -6; z <= 6; z++ {
				tx.SetBlock(cube.Pos{x, 39, z}, block.Stone{}, nil)
			}
		}
		tx.SetLiquid(source, liquid)
	}).Wait(context.Background())

	for range ticks {
		w.AdvanceTick()
	}

	covered := 0
	w.Do(func(tx *world.Tx) {
		for x := -6; x <= 6; x++ {
			for z := -6; z <= 6; z++ {
				pos := cube.Pos{x, 40, z}
				if pos == source {
					continue
				}
				if _, ok := tx.Block(pos).(world.Liquid); ok {
					covered++
				}
			}
		}
	}).Wait(context.Background())
	return covered
}

func TestLiquidsSpreadEveryTime(t *testing.T) {
	for _, c := range []struct {
		name   string
		liquid world.Liquid
		ticks  int
		want   int
	}{
		{name: "water", liquid: block.Water{Depth: 8, Still: true}, ticks: 60, want: 24},
		{name: "lava", liquid: block.Lava{Depth: 8, Still: true}, ticks: 300, want: 8},
	} {
		t.Run(c.name, func(t *testing.T) {
			for run := range 20 {
				if got := spreadOf(c.liquid, c.ticks); got < c.want {
					t.Fatalf("run %d: %v covered %d blocks after %d ticks, want at least %d",
						run, c.name, got, c.ticks, c.want)
				}
			}
		})
	}
}
