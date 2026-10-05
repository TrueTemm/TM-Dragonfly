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
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Bone struct {
	solid

	Axis cube.Axis
}

func (b Bone) Instrument() sound.Instrument {
	return sound.Xylophone()
}

func (b Bone) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, b)
	if !used {
		return
	}
	b.Axis = face.Axis()

	place(tx, pos, b, user, ctx)
	return placed(ctx)
}

func (b Bone) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(b))
}

func (b Bone) EncodeItem() (name string, meta int16) {
	return "minecraft:bone_block", 0
}

func (b Bone) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:bone_block", map[string]any{"pillar_axis": b.Axis.String(), "deprecated": int32(0)}
}

func allBoneBlock() (boneBlocks []world.Block) {
	for _, axis := range cube.Axes() {
		boneBlocks = append(boneBlocks, Bone{Axis: axis})
	}
	return
}
