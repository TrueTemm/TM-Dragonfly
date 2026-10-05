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

package block

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Leaves struct {
	leaves
	sourceWaterDisplacer

	Type LeavesType

	Persistent bool

	ShouldUpdate bool
}

func (l Leaves) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, l)
	if !used {
		return
	}
	l.Persistent = true

	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

func findLog(pos cube.Pos, tx *world.Tx, visited *[]cube.Pos, distance int) bool {
	for _, v := range *visited {
		if v == pos {
			return false
		}
	}
	*visited = append(*visited, pos)

	if log, ok := tx.Block(pos).(Log); ok && !log.Stripped {
		return true
	}
	if _, ok := tx.Block(pos).(Leaves); !ok || distance > 6 {
		return false
	}
	logFound := false
	pos.Neighbours(func(neighbour cube.Pos) {
		if !logFound && findLog(neighbour, tx, visited, distance+1) {
			logFound = true
		}
	}, tx.Range())
	return logFound
}

func (l Leaves) RandomTick(pos cube.Pos, tx *world.Tx, _ *rand.Rand) {
	if !l.Persistent && l.ShouldUpdate {
		if findLog(pos, tx, &[]cube.Pos{}, 0) {
			l.ShouldUpdate = false
			tx.SetBlock(pos, l, nil)
			return
		}
		ctx := tx.Event()
		if tx.World().Handler().HandleLeavesDecay(ctx, pos); ctx.Cancelled() {

			l.ShouldUpdate = false
			tx.SetBlock(pos, l, nil)
			return
		}
		tx.SetBlock(pos, nil, nil)
		for _, drop := range l.BreakInfo().Drops(item.ToolNone{}, nil) {
			dropItem(tx, drop, pos.Vec3Centre())
		}
	}
}

func (l Leaves) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if !l.Persistent && !l.ShouldUpdate {
		l.ShouldUpdate = true
		tx.SetBlock(pos, l, nil)
	}
}

func (l Leaves) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 60, true)
}

func (l Leaves) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, func(t item.Tool) bool {
		return t.ToolType() == item.TypeShears || t.ToolType() == item.TypeHoe
	}, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if t.ToolType() == item.TypeShears || hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(l, 1)}
		}
		fortune := fortuneLevel(enchantments)
		var drops []item.Stack

		stickChances := []float64{0.02, 0.022222222, 0.025, 0.033333333}
		if rand.Float64() < stickChances[min(fortune, 3)] {
			drops = append(drops, item.NewStack(item.Stick{}, rand.IntN(2)+1))
		}
		if wood, ok := l.Type.Wood(); ok && (wood == OakWood() || wood == DarkOakWood()) {
			appleChances := []float64{0.005, 0.005555556, 0.00625, 0.008333333}
			if rand.Float64() < appleChances[min(fortune, 3)] {
				drops = append(drops, item.NewStack(item.Apple{}, 1))
			}
		}
		return drops
	})
}

func (Leaves) CompostChance() float64 {
	return 0.3
}

func (l Leaves) EncodeItem() (name string, meta int16) {
	return "minecraft:" + l.Type.String(), 0
}

func (Leaves) LightDiffusionLevel() uint8 {
	return 1
}

func (Leaves) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (l Leaves) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + l.Type.String(), map[string]any{"persistent_bit": l.Persistent, "update_bit": l.ShouldUpdate}
}

func allLeaves() (leaves []world.Block) {
	f := func(persistent, update bool) {
		for _, t := range LeavesTypes() {
			leaves = append(leaves, Leaves{Type: t, Persistent: persistent, ShouldUpdate: update})
		}
	}
	f(true, true)
	f(true, false)
	f(false, true)
	f(false, false)
	return
}
