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

type Wood struct {
	solid
	bass

	Wood WoodType

	Stripped bool

	Axis cube.Axis
}

func (w Wood) FlammabilityInfo() FlammabilityInfo {
	if !w.Wood.Flammable() {
		return newFlammabilityInfo(0, 0, false)
	}
	return newFlammabilityInfo(5, 5, true)
}

func (w Wood) BreakInfo() BreakInfo {
	return newBreakInfo(2.0, alwaysHarvestable, axeEffective, oneOf(w))
}

func (Wood) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(item.Charcoal{}, 1), 0.15)
}

func (w Wood) FuelInfo() item.FuelInfo {
	if !w.Wood.Flammable() {
		return item.FuelInfo{}
	}
	return newFuelInfo(time.Second * 15)
}

func (w Wood) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, w)
	if !used {
		return
	}
	w.Axis = face.Axis()

	place(tx, pos, w, user, ctx)
	return placed(ctx)
}

func (w Wood) Strip() (world.Block, world.Sound, bool) {
	return Wood{Axis: w.Axis, Wood: w.Wood, Stripped: true}, nil, !w.Stripped
}

func (w Wood) EncodeItem() (name string, meta int16) {
	switch w.Wood {
	case CrimsonWood(), WarpedWood():
		if w.Stripped {
			return "minecraft:stripped_" + w.Wood.String() + "_hyphae", 0
		}
		return "minecraft:" + w.Wood.String() + "_hyphae", 0
	default:
		if w.Stripped {
			return "minecraft:stripped_" + w.Wood.String() + "_wood", 0
		}
		return "minecraft:" + w.Wood.String() + "_wood", 0
	}
}

func (w Wood) EncodeBlock() (name string, properties map[string]any) {
	switch w.Wood {
	case CrimsonWood(), WarpedWood():
		if w.Stripped {
			return "minecraft:stripped_" + w.Wood.String() + "_hyphae", map[string]any{"pillar_axis": w.Axis.String()}
		}
		return "minecraft:" + w.Wood.String() + "_hyphae", map[string]any{"pillar_axis": w.Axis.String()}
	default:
		if w.Stripped {
			return "minecraft:stripped_" + w.Wood.String() + "_wood", map[string]any{"pillar_axis": w.Axis.String()}
		}
		return "minecraft:" + w.Wood.String() + "_wood", map[string]any{"pillar_axis": w.Axis.String()}
	}
}

func allWood() (wood []world.Block) {
	for _, w := range WoodTypes() {
		if w == BambooWood() {
			continue
		}
		for axis := cube.Axis(0); axis < 3; axis++ {
			wood = append(wood, Wood{Axis: axis, Stripped: true, Wood: w})
			wood = append(wood, Wood{Axis: axis, Stripped: false, Wood: w})
		}
	}
	return
}
