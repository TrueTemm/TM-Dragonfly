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
)

type ConcretePowder struct {
	gravityAffected
	solid
	snare

	Colour item.Colour
}

func (c ConcretePowder) Solidifies(pos cube.Pos, tx *world.Tx) bool {
	_, water := tx.Block(pos).(Water)
	return water
}

func (c ConcretePowder) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	for i := cube.Face(0); i < 6; i++ {
		if _, ok := tx.Block(pos.Side(i)).(Water); ok {
			tx.SetBlock(pos, Concrete{Colour: c.Colour}, nil)
			return
		}
	}
	c.fall(c, pos, tx)
}

func (c ConcretePowder) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(c))
}

func (c ConcretePowder) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.Colour.String() + "_concrete_powder", 0
}

func (c ConcretePowder) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + c.Colour.String() + "_concrete_powder", nil
}

func allConcretePowder() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, ConcretePowder{Colour: c})
	}
	return b
}
