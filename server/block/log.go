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
	"time"
)

type Log struct {
	solid
	bass

	Wood WoodType

	Stripped bool

	Axis cube.Axis
}

func (l Log) FlammabilityInfo() FlammabilityInfo {
	if !l.Wood.Flammable() {
		return newFlammabilityInfo(0, 0, false)
	}
	return newFlammabilityInfo(5, 5, true)
}

func (l Log) BreakInfo() BreakInfo {
	breakInfo := newBreakInfo(2, alwaysHarvestable, axeEffective, oneOf(l))
	if l.Wood == MangroveWood() && !l.Stripped {

		return breakInfo.withBlastResistance(0.4)
	}
	return breakInfo
}

func (Log) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(item.Charcoal{}, 1), 0.15)
}

func (l Log) FuelInfo() item.FuelInfo {
	if !l.Wood.Flammable() {
		return item.FuelInfo{}
	}
	return newFuelInfo(time.Second * 15)
}

func (l Log) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, l)
	if !used {
		return
	}
	l.Axis = face.Axis()

	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

func (l Log) Strip() (world.Block, world.Sound, bool) {
	return Log{Axis: l.Axis, Wood: l.Wood, Stripped: true}, nil, !l.Stripped
}

func (l Log) EncodeItem() (name string, meta int16) {
	if !l.Stripped {
		switch l.Wood {
		case CrimsonWood(), WarpedWood():
			return "minecraft:" + l.Wood.String() + "_stem", 0
		default:
			return "minecraft:" + l.Wood.String() + "_log", 0
		}
	}
	switch l.Wood {
	case CrimsonWood(), WarpedWood():
		return "minecraft:stripped_" + l.Wood.String() + "_stem", 0
	default:
		return "minecraft:stripped_" + l.Wood.String() + "_log", 0
	}
}

func (l Log) EncodeBlock() (name string, properties map[string]any) {
	if !l.Stripped {
		switch l.Wood {
		case CrimsonWood(), WarpedWood():
			return "minecraft:" + l.Wood.String() + "_stem", map[string]any{"pillar_axis": l.Axis.String()}
		default:
			return "minecraft:" + l.Wood.String() + "_log", map[string]any{"pillar_axis": l.Axis.String()}
		}
	}
	switch l.Wood {
	case CrimsonWood(), WarpedWood():
		return "minecraft:stripped_" + l.Wood.String() + "_stem", map[string]any{"pillar_axis": l.Axis.String()}
	default:
		return "minecraft:stripped_" + l.Wood.String() + "_log", map[string]any{"pillar_axis": l.Axis.String()}
	}
}

func allLogs() (logs []world.Block) {
	for _, w := range WoodTypes() {
		if w == BambooWood() {
			continue
		}
		for axis := cube.Axis(0); axis < 3; axis++ {
			logs = append(logs, Log{Axis: axis, Stripped: true, Wood: w})
			logs = append(logs, Log{Axis: axis, Stripped: false, Wood: w})
		}
	}
	return
}
