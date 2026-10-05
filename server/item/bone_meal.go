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

package item

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

type BoneMealResult int

const (
	BoneMealResultNone BoneMealResult = iota

	BoneMealResultSmall

	BoneMealResultArea
)

type BoneMeal struct{}

type BoneMealAffected interface {
	BoneMeal(pos cube.Pos, tx *world.Tx) BoneMealResult
}

func (b BoneMeal) UseOnBlock(pos cube.Pos, _ cube.Face, _ mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	if bm, ok := tx.Block(pos).(BoneMealAffected); ok {
		result := bm.BoneMeal(pos, tx)
		if result == BoneMealResultNone {
			return false
		}

		ctx.SubtractFromCount(1)
		tx.AddParticle(pos.Vec3(), particle.BoneMeal{
			Area: result == BoneMealResultArea,
		})
		return true
	}
	return false
}

func (b BoneMeal) EncodeItem() (name string, meta int16) {
	return "minecraft:bone_meal", 0
}
