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

type InfestedDeepslate struct {
	solid
	flute

	Axis cube.Axis
}

func (i InfestedDeepslate) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, silkTouchOnlyDrop(i)).withBlastResistance(0.75)
}

func (i InfestedDeepslate) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, i)
	if !used {
		return
	}
	i.Axis = face.Axis()

	place(tx, pos, i, user, ctx)
	return placed(ctx)
}

func (InfestedDeepslate) EncodeItem() (name string, meta int16) {
	return "minecraft:infested_deepslate", 0
}

func (i InfestedDeepslate) EncodeBlock() (string, map[string]any) {
	return "minecraft:infested_deepslate", map[string]any{"pillar_axis": i.Axis.String()}
}

func allInfestedDeepslate() (s []world.Block) {
	for _, axis := range cube.Axes() {
		s = append(s, InfestedDeepslate{Axis: axis})
	}
	return
}
