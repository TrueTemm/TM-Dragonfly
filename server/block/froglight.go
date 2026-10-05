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

type Froglight struct {
	solid

	Type FroglightType

	Axis cube.Axis
}

func (f Froglight) LightEmissionLevel() uint8 {
	return 15
}

func (f Froglight) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, f)
	if !used {
		return
	}
	f.Axis = face.Axis()

	place(tx, pos, f, user, ctx)
	return placed(ctx)
}

func (f Froglight) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, oneOf(f))
}

func (f Froglight) EncodeItem() (name string, meta int16) {
	return "minecraft:" + f.Type.String() + "_froglight", 0
}

func (f Froglight) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + f.Type.String() + "_froglight", map[string]any{"pillar_axis": f.Axis.String()}
}

func allFroglight() (froglight []world.Block) {
	for _, axis := range cube.Axes() {
		for _, t := range FroglightTypes() {
			froglight = append(froglight, Froglight{Type: t, Axis: axis})
		}
	}
	return
}
