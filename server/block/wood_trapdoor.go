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
	"math"
	"time"
)

type WoodTrapdoor struct {
	transparent
	bass
	sourceWaterDisplacer

	Wood WoodType

	Facing cube.Direction

	Open bool

	Top bool
}

func (t WoodTrapdoor) FlammabilityInfo() FlammabilityInfo {
	if !t.Wood.Flammable() {
		return newFlammabilityInfo(0, 0, false)
	}
	return newFlammabilityInfo(0, 0, true)
}

func (t WoodTrapdoor) Model() world.BlockModel {
	return model.Trapdoor{Facing: t.Facing, Top: t.Top, Open: t.Open}
}

func (t WoodTrapdoor) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, t)
	if !used {
		return false
	}
	t.Facing = user.Rotation().Direction().Opposite()
	t.Top = (clickPos.Y() > 0.5 && face != cube.FaceUp) || face == cube.FaceDown

	place(tx, pos, t, user, ctx)
	return placed(ctx)
}

func (t WoodTrapdoor) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, _ item.User, _ *item.UseContext) bool {
	t.Open = !t.Open
	tx.SetBlock(pos, t, nil)
	if t.Open {
		tx.PlaySound(pos.Vec3Centre(), sound.TrapdoorOpen{Block: t})
		return true
	}
	tx.PlaySound(pos.Vec3Centre(), sound.TrapdoorClose{Block: t})
	return true
}

func (t WoodTrapdoor) BreakInfo() BreakInfo {
	return newBreakInfo(3, alwaysHarvestable, axeEffective, oneOf(t))
}

func (t WoodTrapdoor) FuelInfo() item.FuelInfo {
	if !t.Wood.Flammable() {
		return item.FuelInfo{}
	}
	return newFuelInfo(time.Second * 15)
}

func (t WoodTrapdoor) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (t WoodTrapdoor) EncodeItem() (name string, meta int16) {
	if t.Wood == OakWood() {
		return "minecraft:trapdoor", 0
	}
	return "minecraft:" + t.Wood.String() + "_trapdoor", 0
}

func (t WoodTrapdoor) EncodeBlock() (name string, properties map[string]any) {
	if t.Wood == OakWood() {
		return "minecraft:trapdoor", map[string]any{"direction": int32(math.Abs(float64(t.Facing) - 3)), "open_bit": t.Open, "upside_down_bit": t.Top}
	}
	return "minecraft:" + t.Wood.String() + "_trapdoor", map[string]any{"direction": int32(math.Abs(float64(t.Facing) - 3)), "open_bit": t.Open, "upside_down_bit": t.Top}
}

func allTrapdoors() (trapdoors []world.Block) {
	for _, w := range WoodTypes() {
		for i := cube.Direction(0); i <= 3; i++ {
			trapdoors = append(trapdoors, WoodTrapdoor{Wood: w, Facing: i, Open: false, Top: false})
			trapdoors = append(trapdoors, WoodTrapdoor{Wood: w, Facing: i, Open: false, Top: true})
			trapdoors = append(trapdoors, WoodTrapdoor{Wood: w, Facing: i, Open: true, Top: true})
			trapdoors = append(trapdoors, WoodTrapdoor{Wood: w, Facing: i, Open: true, Top: false})
		}
	}
	return
}
