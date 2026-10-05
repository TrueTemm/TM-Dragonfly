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
	"github.com/df-mc/dragonfly/server/world"
)

type DirtPath struct {
	tilledGrass
	transparent
}

func (p DirtPath) Till() (world.Block, bool) {
	return Farmland{}, true
}

func (p DirtPath) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	up := pos.Side(cube.FaceUp)
	if tx.Block(up).Model().FaceSolid(up, cube.FaceDown, tx) {

		tx.SetBlock(pos, Dirt{}, nil)
	}
}

func (p DirtPath) BreakInfo() BreakInfo {
	return newBreakInfo(0.65, alwaysHarvestable, shovelEffective, silkTouchOneOf(Dirt{}, p))
}

func (DirtPath) EncodeItem() (name string, meta int16) {
	return "minecraft:grass_path", 0
}

func (DirtPath) EncodeBlock() (string, map[string]any) {
	return "minecraft:grass_path", nil
}
