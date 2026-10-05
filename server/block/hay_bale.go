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

type HayBale struct {
	solid

	Axis cube.Axis
}

func (HayBale) Instrument() sound.Instrument {
	return sound.Banjo()
}

func (h HayBale) EntityLand(_ cube.Pos, _ *world.Tx, e world.Entity, distance *float64) {
	if _, ok := e.(fallDistanceEntity); ok {
		*distance *= 0.2
	}
}

func (h HayBale) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, h)
	if !used {
		return
	}
	h.Axis = face.Axis()

	place(tx, pos, h, user, ctx)
	return placed(ctx)
}

func (HayBale) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 20, false)
}

func (h HayBale) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, hoeEffective, oneOf(h))
}

func (HayBale) CompostChance() float64 {
	return 0.85
}

func (HayBale) EncodeItem() (name string, meta int16) {
	return "minecraft:hay_block", 0
}

func (h HayBale) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:hay_block", map[string]interface{}{"pillar_axis": h.Axis.String(), "deprecated": int32(0)}
}

func allHayBales() (haybale []world.Block) {
	for _, a := range cube.Axes() {
		haybale = append(haybale, HayBale{Axis: a})
	}
	return
}
