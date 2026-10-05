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

package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	FishingHookForce = 0.38

	FishingHookHeight = 0.37
)

func NewFishingHook(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	conf := fishingHookConf
	conf.Owner = owner.H()
	return opts.New(FishingHookType, conf)
}

var fishingHookConf = ProjectileBehaviourConfig{
	Gravity: 0.05,
	Drag:    0.01,

	Damage:    0,
	HitForce:  FishingHookForce,
	HitHeight: FishingHookHeight,
}

var FishingHookType fishingHookType

type fishingHookType struct{}

func (t fishingHookType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (fishingHookType) EncodeEntity() string { return "minecraft:fishing_hook" }
func (fishingHookType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}
func (fishingHookType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	data.Data = fishingHookConf.New()
}
func (fishingHookType) EncodeNBT(*world.EntityData) map[string]any { return nil }
