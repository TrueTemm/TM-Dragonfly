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

package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Shulker struct {
	Facing cube.Face

	Progress int32
}

func (s Shulker) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	peak := ShulkerPhysicalPeak(s.Progress)
	return []cube.BBox{full.ExtendTowards(s.Facing, peak)}
}

func ShulkerPhysicalPeak(progress int32) float64 {
	t := float64(progress) / 10.0
	return (1.0 - (1.0-t)*(1.0-t)*(1.0-t)) * 0.5
}

func (Shulker) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
