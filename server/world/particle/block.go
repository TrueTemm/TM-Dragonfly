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

package particle

import (
	"image/color"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Flame struct {
	particle

	Colour color.RGBA
}

type Dust struct {
	particle

	Colour color.RGBA
}

type BlockBreak struct {
	particle

	Block world.Block
}

type PunchBlock struct {
	particle

	Block world.Block

	Face cube.Face
}

type BlockForceField struct{ particle }

type BoneMeal struct {
	particle

	Area bool
}

type Note struct {
	particle

	Instrument sound.Instrument

	Pitch int
}

type DragonEggTeleport struct {
	particle

	Diff cube.Pos
}

type Evaporate struct{ particle }

type WaterDrip struct{ particle }

type LavaDrip struct{ particle }

type Lava struct{ particle }

type DustPlume struct{ particle }

type particle struct{}

func (particle) Spawn(*world.World, mgl64.Vec3) {}
