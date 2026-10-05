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

type MuddyMangroveRoots struct {
	solid

	Axis cube.Axis
}

func (m MuddyMangroveRoots) BreakInfo() BreakInfo {
	return newBreakInfo(0.7, alwaysHarvestable, shovelEffective, oneOf(m))
}

func (MuddyMangroveRoots) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, NetherSprouts, PinkPetals, BambooSapling, Bamboo:
		return true
	}
	return false
}

func (m MuddyMangroveRoots) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, m)
	if !used {
		return
	}
	m.Axis = face.Axis()

	place(tx, pos, m, user, ctx)
	return placed(ctx)
}

func (MuddyMangroveRoots) EncodeItem() (name string, meta int16) {
	return "minecraft:muddy_mangrove_roots", 0
}

func (m MuddyMangroveRoots) EncodeBlock() (string, map[string]any) {
	return "minecraft:muddy_mangrove_roots", map[string]any{"pillar_axis": m.Axis.String()}
}

func allMuddyMangroveRoots() (roots []world.Block) {
	for _, axis := range cube.Axes() {
		roots = append(roots, MuddyMangroveRoots{Axis: axis})
	}
	return
}
