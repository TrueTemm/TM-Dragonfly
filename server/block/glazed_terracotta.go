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

type GlazedTerracotta struct {
	solid
	bassDrum

	Colour item.Colour

	Facing cube.Direction
}

func (t GlazedTerracotta) BreakInfo() BreakInfo {
	return newBreakInfo(1.4, pickaxeHarvestable, pickaxeEffective, oneOf(t))
}

func (t GlazedTerracotta) EncodeItem() (name string, meta int16) {
	return "minecraft:" + t.Colour.SilverString() + "_glazed_terracotta", 0
}

func (t GlazedTerracotta) EncodeBlock() (name string, properties map[string]any) {
	if t.Facing == unknownDirection {
		return "minecraft:" + t.Colour.SilverString() + "_glazed_terracotta", map[string]any{"facing_direction": int32(0)}
	}
	return "minecraft:" + t.Colour.SilverString() + "_glazed_terracotta", map[string]any{"facing_direction": int32(2 + t.Facing)}
}

func (t GlazedTerracotta) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, t)
	if !used {
		return
	}
	t.Facing = user.Rotation().Direction().Opposite()

	place(tx, pos, t, user, ctx)
	return placed(ctx)
}

func allGlazedTerracotta() (b []world.Block) {
	for _, dir := range append(cube.Directions(), unknownDirection) {
		for _, c := range item.Colours() {
			b = append(b, GlazedTerracotta{Colour: c, Facing: dir})
		}
	}
	return b
}
