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
	"strconv"
)

type Light struct {
	empty
	replaceable
	transparent
	flowingWaterDisplacer

	Level int
}

func (Light) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (l Light) EncodeItem() (name string, meta int16) {
	return "minecraft:light_block_" + strconv.Itoa(l.Level), 0
}

func (l Light) LightEmissionLevel() uint8 {
	return uint8(l.Level)
}

func (l Light) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:light_block_" + strconv.Itoa(l.Level), nil
}

func allLight() []world.Block {
	m := make([]world.Block, 0, 16)
	for i := 0; i < 16; i++ {
		m = append(m, Light{Level: i})
	}
	return m
}
