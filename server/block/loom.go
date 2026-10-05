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

type Loom struct {
	solid
	bass

	Facing cube.Direction
}

func (Loom) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (l Loom) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, oneOf(l))
}

func (Loom) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (l Loom) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, l)
	if !used {
		return
	}
	l.Facing = user.Rotation().Direction().Opposite()
	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

func (Loom) EncodeItem() (name string, meta int16) {
	return "minecraft:loom", 0
}

func (l Loom) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:loom", map[string]interface{}{"direction": int32(horizontalDirection(l.Facing))}
}

func allLooms() (looms []world.Block) {
	for _, d := range cube.Directions() {
		looms = append(looms, Loom{Facing: d})
	}
	return
}
