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
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"slices"
)

var orbSplitSizes = []int{2477, 1237, 617, 307, 149, 73, 37, 17, 7, 3, 1}

func NewExperienceOrbs(pos mgl64.Vec3, amount int) (orbs []*world.EntityHandle) {
	for amount > 0 {
		size := orbSplitSizes[slices.IndexFunc(orbSplitSizes, func(value int) bool {
			return amount >= value
		})]

		orbs = append(orbs, NewExperienceOrb(world.EntitySpawnOpts{Position: pos}, size))
		amount -= size
	}
	return
}

func NewExperienceOrb(opts world.EntitySpawnOpts, xp int) *world.EntityHandle {
	conf := experienceOrbConf
	conf.Experience = xp
	if opts.Velocity.Len() == 0 {
		opts.Velocity = mgl64.Vec3{(rand.Float64()*0.2 - 0.1) * 2, rand.Float64() * 0.4, (rand.Float64()*0.2 - 0.1) * 2}
	}
	return opts.New(ExperienceOrbType, conf)
}

var experienceOrbConf = ExperienceOrbBehaviourConfig{
	Gravity: 0.04,
	Drag:    0.02,
}

var ExperienceOrbType experienceOrbType

type experienceOrbType struct{}

func (t experienceOrbType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (experienceOrbType) EncodeEntity() string { return "minecraft:xp_orb" }
func (experienceOrbType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (experienceOrbType) DecodeNBT(m map[string]any, data *world.EntityData) {
	conf := experienceOrbConf
	conf.Experience = int(nbtconv.Int32(m, "Value"))
	data.Data = conf.New()
}

func (experienceOrbType) EncodeNBT(data *world.EntityData) map[string]any {
	return map[string]any{"Value": int32(data.Data.(*ExperienceOrbBehaviour).Experience())}
}
