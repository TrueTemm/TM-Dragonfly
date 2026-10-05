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
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"time"
)

type GoatHorn struct {
	nopReleasable

	Type sound.Horn
}

func (GoatHorn) MaxCount() int {
	return 1
}

func (GoatHorn) Cooldown() time.Duration {
	return time.Second * 7
}

func (g GoatHorn) Use(tx *world.Tx, user User, _ *UseContext) bool {
	tx.PlaySound(user.Position(), sound.GoatHorn{Horn: g.Type})
	user.H().DoAfter(time.Second, g.releaseItem)
	return true
}

func (g GoatHorn) releaseItem(_ *world.Tx, e world.Entity) {
	user := e.(User)
	if !user.UsingItem() {

		return
	}
	held, _ := user.HeldItems()
	if _, ok := held.Item().(GoatHorn); !ok {

		return
	}

	user.ReleaseItem()
}

func (g GoatHorn) EncodeItem() (name string, meta int16) {
	return "minecraft:goat_horn", int16(g.Type.Uint8())
}
