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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type NetherBricks struct {
	solid
	bassDrum

	Type NetherBricksType
}

func (n NetherBricks) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(n)).withBlastResistance(6)
}

func (n NetherBricks) SmeltInfo() item.SmeltInfo {
	if n.Type == NormalNetherBricks() {
		return newSmeltInfo(item.NewStack(NetherBricks{Type: CrackedNetherBricks()}, 1), 0.1)
	}
	return item.SmeltInfo{}
}

func (n NetherBricks) EncodeItem() (id string, meta int16) {
	return "minecraft:" + n.Type.String(), 0
}

func (n NetherBricks) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + n.Type.String(), nil
}

func allNetherBricks() (netherBricks []world.Block) {
	for _, t := range NetherBricksTypes() {
		netherBricks = append(netherBricks, NetherBricks{Type: t})
	}
	return
}
