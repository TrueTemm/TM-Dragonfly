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

type Deepslate struct {
	solid
	bassDrum

	Type DeepslateType

	Axis cube.Axis
}

func (d Deepslate) BreakInfo() BreakInfo {
	if d.Type == NormalDeepslate() {
		return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, silkTouchOneOf(Deepslate{Type: CobbledDeepslate()}, d)).withBlastResistance(6)
	}
	return newBreakInfo(3.5, pickaxeHarvestable, pickaxeEffective, oneOf(d)).withBlastResistance(6)
}

func (d Deepslate) SmeltInfo() item.SmeltInfo {
	if d.Type == CobbledDeepslate() {
		return newSmeltInfo(item.NewStack(Deepslate{}, 1), 0.1)
	}
	return item.SmeltInfo{}
}

func (d Deepslate) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, d)
	if !used {
		return
	}
	if d.Type == NormalDeepslate() {
		d.Axis = face.Axis()
	}

	place(tx, pos, d, user, ctx)
	return placed(ctx)
}

func (d Deepslate) EncodeItem() (name string, meta int16) {
	return "minecraft:" + d.Type.String(), 0
}

func (d Deepslate) EncodeBlock() (string, map[string]any) {
	if d.Type == NormalDeepslate() {
		return "minecraft:deepslate", map[string]any{"pillar_axis": d.Axis.String()}
	}
	return "minecraft:" + d.Type.String(), nil
}

func allDeepslate() (s []world.Block) {
	for _, t := range DeepslateTypes() {
		axes := []cube.Axis{0}
		if t == NormalDeepslate() {
			axes = cube.Axes()
		}
		for _, axis := range axes {
			s = append(s, Deepslate{Type: t, Axis: axis})
		}
	}
	return
}
