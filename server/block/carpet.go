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

type Carpet struct {
	carpet
	transparent
	sourceWaterDisplacer

	Colour item.Colour
}

func (c Carpet) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 20, true)
}

func (Carpet) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c Carpet) BreakInfo() BreakInfo {
	return newBreakInfo(0.1, alwaysHarvestable, nothingEffective, oneOf(c))
}

func (c Carpet) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.Colour.String() + "_carpet", 0
}

func (c Carpet) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + c.Colour.String() + "_carpet", nil
}

func (Carpet) HasLiquidDrops() bool {
	return true
}

func (c Carpet) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Air); ok {
		breakBlock(c, pos, tx)
	}
}

func (c Carpet) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, c)
	if !used {
		return
	}

	if _, ok := tx.Block(pos.Side(cube.FaceDown)).(Air); ok {
		return
	}

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func allCarpet() (carpets []world.Block) {
	for _, c := range item.Colours() {
		carpets = append(carpets, Carpet{Colour: c})
	}
	return
}
