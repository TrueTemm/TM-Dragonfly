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
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Anvil struct {
	gravityAffected
	transparent

	Type AnvilType

	Facing cube.Direction
}

func (a Anvil) Model() world.BlockModel {
	return model.Anvil{Facing: a.Facing}
}

func (a Anvil) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(a)).withBlastResistance(1200)
}

func (Anvil) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (a Anvil) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, a)
	if !used {
		return
	}
	a.Facing = user.Rotation().Direction().RotateLeft()
	place(tx, pos, a, user, ctx)
	return placed(ctx)
}

func (a Anvil) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	a.fall(a, pos, tx)
}

func (Anvil) Damage() (damagePerBlock, maxDamage float64) {
	return 2, 40
}

func (a Anvil) Break() world.Block {
	switch a.Type {
	case UndamagedAnvil():
		a.Type = SlightlyDamagedAnvil()
	case SlightlyDamagedAnvil():
		a.Type = VeryDamagedAnvil()
	case VeryDamagedAnvil():
		return Air{}
	}
	return a
}

func (Anvil) Landed(tx *world.Tx, pos cube.Pos) {
	tx.PlaySound(pos.Vec3Centre(), sound.AnvilLand{})
}

func (a Anvil) EncodeItem() (name string, meta int16) {
	return "minecraft:" + a.Type.String(), 0
}

func (a Anvil) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + a.Type.String(), map[string]any{
		"minecraft:cardinal_direction": a.Facing.String(),
	}
}

func allAnvils() (anvils []world.Block) {
	for _, t := range AnvilTypes() {
		for _, d := range cube.Directions() {
			anvils = append(anvils, Anvil{Type: t, Facing: d})
		}
	}
	return
}
