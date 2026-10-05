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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type (
	Quartz struct {
		solid
		bassDrum

		Smooth bool
	}

	ChiseledQuartz struct {
		solid
		bassDrum
	}

	QuartzPillar struct {
		solid
		bassDrum

		Axis cube.Axis
	}
)

func (q QuartzPillar) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, q)
	if !used {
		return
	}
	q.Axis = face.Axis()

	place(tx, pos, q, user, ctx)
	return placed(ctx)
}

func (q Quartz) BreakInfo() BreakInfo {
	if q.Smooth {
		return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(q)).withBlastResistance(6)
	}
	return newBreakInfo(0.8, pickaxeHarvestable, pickaxeEffective, oneOf(q))
}

func (c ChiseledQuartz) BreakInfo() BreakInfo {
	return newBreakInfo(0.8, pickaxeHarvestable, pickaxeEffective, simpleDrops(item.NewStack(c, 1)))
}

func (q QuartzPillar) BreakInfo() BreakInfo {
	return newBreakInfo(0.8, pickaxeHarvestable, pickaxeEffective, simpleDrops(item.NewStack(q, 1)))
}

func (q Quartz) SmeltInfo() item.SmeltInfo {
	if q.Smooth {
		return item.SmeltInfo{}
	}
	return newSmeltInfo(item.NewStack(Quartz{Smooth: true}, 1), 0.1)
}

func (q Quartz) EncodeItem() (name string, meta int16) {
	if q.Smooth {
		return "minecraft:smooth_quartz", 0
	}
	return "minecraft:quartz_block", 0
}

func (c ChiseledQuartz) EncodeItem() (name string, meta int16) {
	return "minecraft:chiseled_quartz_block", 0
}

func (q QuartzPillar) EncodeItem() (name string, meta int16) {
	return "minecraft:quartz_pillar", 0
}

func (q Quartz) EncodeBlock() (name string, properties map[string]any) {
	if q.Smooth {
		return "minecraft:smooth_quartz", map[string]any{"pillar_axis": "y"}
	}
	return "minecraft:quartz_block", map[string]any{"pillar_axis": "y"}
}

func (ChiseledQuartz) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:chiseled_quartz_block", map[string]any{"pillar_axis": "y"}
}

func (q QuartzPillar) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:quartz_pillar", map[string]any{"pillar_axis": q.Axis.String()}
}

func allQuartz() (quartz []world.Block) {
	quartz = append(quartz, Quartz{})
	quartz = append(quartz, Quartz{Smooth: true})
	quartz = append(quartz, ChiseledQuartz{})
	for _, a := range cube.Axes() {
		quartz = append(quartz, QuartzPillar{Axis: a})
	}
	return
}
