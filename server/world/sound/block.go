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

package sound

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type BlockPlace struct {
	Block world.Block

	sound
}

type BlockBreaking struct {
	Block world.Block

	sound
}

type GlassBreak struct{ sound }

type Fizz struct{ sound }

type AnvilLand struct{ sound }

type AnvilUse struct{ sound }

type AnvilBreak struct{ sound }

type ChestOpen struct{ sound }

type ChestClose struct{ sound }

type EnderChestOpen struct{ sound }

type EnderChestClose struct{ sound }

type BarrelOpen struct{ sound }

type BarrelClose struct{ sound }

type Deny struct{ sound }

type ShulkerBoxOpen struct{ sound }

type ShulkerBoxClose struct{ sound }

type DoorOpen struct {
	Block world.Block

	sound
}

type DoorClose struct {
	Block world.Block

	sound
}

type TrapdoorOpen struct {
	Block world.Block

	sound
}

type TrapdoorClose struct {
	Block world.Block

	sound
}

type FenceGateOpen struct {
	Block world.Block

	sound
}

type FenceGateClose struct {
	Block world.Block

	sound
}

type DoorCrash struct{ sound }

type Click struct{ sound }

type Ignite struct{ sound }

type TNT struct{ sound }

type FireExtinguish struct{ sound }

type Note struct {
	sound

	Instrument Instrument

	Pitch int
}

type MusicDiscPlay struct {
	sound

	DiscType DiscType
}

type MusicDiscEnd struct{ sound }

type ItemAdd struct{ sound }

type ItemFrameRemove struct{ sound }

type ItemFrameRotate struct{ sound }

type FurnaceCrackle struct{ sound }

type CampfireCrackle struct{ sound }

type BlastFurnaceCrackle struct{ sound }

type SmokerCrackle struct{ sound }

type ComposterEmpty struct{ sound }

type ComposterFill struct{ sound }

type ComposterFillLayer struct{ sound }

type ComposterReady struct{ sound }

type PotionBrewed struct{ sound }

type PowerOn struct{ sound }

type PowerOff struct{ sound }

type LecternBookPlace struct{ sound }

type SignWaxed struct{ sound }

type WaxedSignFailedInteraction struct{ sound }

type WaxRemoved struct{ sound }

type CopperScraped struct{ sound }

type DecoratedPotInserted struct {
	sound

	Progress float64
}

type DecoratedPotInsertFailed struct{ sound }

type EnderEyePlaced struct{ sound }

type EndPortalCreated struct{ sound }

type sound struct{}

func (sound) Play(*world.World, mgl64.Vec3) {}
