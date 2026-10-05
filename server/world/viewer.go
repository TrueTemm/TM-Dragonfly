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

package world

import (
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type Viewer interface {
	ViewEntity(e Entity)

	HideEntity(e Entity)

	ViewEntityGameMode(e Entity)

	ViewEntityMovement(e Entity, pos mgl64.Vec3, rot cube.Rotation, onGround bool)

	ViewEntityDisplacement(e Entity, pos mgl64.Vec3, rot cube.Rotation, onGround bool)

	ViewEntityVelocity(e Entity, vel mgl64.Vec3)

	ViewEntityTeleport(e Entity, pos mgl64.Vec3)

	ViewFurnaceUpdate(prevCookTime, cookTime, prevRemainingFuelTime, remainingFuelTime, prevMaxFuelTime, maxFuelTime time.Duration)

	ViewBrewingUpdate(prevBrewTime, brewTime time.Duration, prevFuelAmount, fuelAmount, prevFuelTotal, fuelTotal int32)

	ViewChunk(pos ChunkPos, dim Dimension, c *Column)

	ViewTime(t int)

	ViewTimeCycle(doDayLightCycle bool)

	ViewEntityItems(e Entity)

	ViewEntityArmour(e Entity)

	ViewEntityAction(e Entity, a EntityAction)

	ViewEntityState(e Entity)

	ViewEntityAnimation(e Entity, a EntityAnimation)

	ViewParticle(pos mgl64.Vec3, p Particle)

	ViewSound(pos mgl64.Vec3, s Sound)

	ViewBlockUpdate(pos cube.Pos, b Block, layer int)

	ViewBlockAction(pos cube.Pos, a BlockAction)

	ViewEmote(e Entity, emote uuid.UUID)

	ViewSkin(e Entity)

	ViewWorldSpawn(pos cube.Pos)

	ViewWeather(raining, thunder bool)

	ViewEntityWake(e Entity)
}

type NopViewer struct{}

var _ Viewer = NopViewer{}

func (NopViewer) ViewEntity(Entity)                                                          {}
func (NopViewer) HideEntity(Entity)                                                          {}
func (NopViewer) ViewEntityGameMode(Entity)                                                  {}
func (NopViewer) ViewEntityMovement(Entity, mgl64.Vec3, cube.Rotation, bool)                 {}
func (NopViewer) ViewEntityDisplacement(Entity, mgl64.Vec3, cube.Rotation, bool)             {}
func (NopViewer) ViewEntityVelocity(Entity, mgl64.Vec3)                                      {}
func (NopViewer) ViewEntityTeleport(Entity, mgl64.Vec3)                                      {}
func (NopViewer) ViewChunk(ChunkPos, Dimension, *Column)                                     {}
func (NopViewer) ViewTime(int)                                                               {}
func (NopViewer) ViewTimeCycle(bool)                                                         {}
func (NopViewer) ViewEntityItems(Entity)                                                     {}
func (NopViewer) ViewEntityArmour(Entity)                                                    {}
func (NopViewer) ViewEntityAction(Entity, EntityAction)                                      {}
func (NopViewer) ViewEntityState(Entity)                                                     {}
func (NopViewer) ViewEntityAnimation(Entity, EntityAnimation)                                {}
func (NopViewer) ViewParticle(mgl64.Vec3, Particle)                                          {}
func (NopViewer) ViewSound(mgl64.Vec3, Sound)                                                {}
func (NopViewer) ViewBlockUpdate(cube.Pos, Block, int)                                       {}
func (NopViewer) ViewBlockAction(cube.Pos, BlockAction)                                      {}
func (NopViewer) ViewEmote(Entity, uuid.UUID)                                                {}
func (NopViewer) ViewSkin(Entity)                                                            {}
func (NopViewer) ViewWorldSpawn(cube.Pos)                                                    {}
func (NopViewer) ViewWeather(bool, bool)                                                     {}
func (NopViewer) ViewBrewingUpdate(time.Duration, time.Duration, int32, int32, int32, int32) {}
func (NopViewer) ViewEntityWake(Entity)                                                      {}
func (NopViewer) ViewFurnaceUpdate(time.Duration, time.Duration, time.Duration, time.Duration, time.Duration, time.Duration) {
}
