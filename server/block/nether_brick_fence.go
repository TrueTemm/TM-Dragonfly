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
	"github.com/df-mc/dragonfly/server/world"
)

type NetherBrickFence struct {
	transparent
	sourceWaterDisplacer
}

func (n NetherBrickFence) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(n)).withBlastResistance(6)
}

func (NetherBrickFence) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (n NetherBrickFence) Model() world.BlockModel {
	return model.Fence{}
}

func (NetherBrickFence) EncodeItem() (name string, meta int16) {
	return "minecraft:nether_brick_fence", 0
}

func (NetherBrickFence) EncodeBlock() (string, map[string]any) {
	return "minecraft:nether_brick_fence", nil
}
