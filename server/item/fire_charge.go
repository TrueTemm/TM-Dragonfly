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
	"github.com/df-mc/dragonfly/server/world/portal"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"time"
)

type FireCharge struct{}

func (f FireCharge) EncodeItem() (name string, meta int16) {
	return "minecraft:fire_charge", 0
}

func (f FireCharge) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user User, ctx *UseContext) bool {
	if l, ok := tx.Block(pos).(ignitable); ok && l.Ignite(pos, tx, user) {
		ctx.SubtractFromCount(1)
		tx.PlaySound(pos.Vec3Centre(), sound.FireCharge{})
		return true
	} else if s := pos.Side(face); tx.Block(s) == air() {
		ctx.SubtractFromCount(1)
		tx.PlaySound(s.Vec3Centre(), sound.FireCharge{})
		if portal.ActivateNetherPortal(tx, s) {
			return true
		}

		flame := fire()
		tx.SetBlock(s, flame, nil)
		tx.ScheduleBlockUpdate(s, flame, time.Duration(30+rand.IntN(10))*time.Second/20)
		return true
	}
	return false
}
