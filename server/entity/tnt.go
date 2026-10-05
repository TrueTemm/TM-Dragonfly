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
	"math"
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func NewTNT(opts world.EntitySpawnOpts, fuse time.Duration) *world.EntityHandle {
	conf := tntConf
	conf.ExistenceDuration = fuse
	if opts.Velocity.Len() == 0 {
		angle := rand.Float64() * math.Pi * 2
		opts.Velocity = mgl64.Vec3{-math.Sin(angle) * 0.02, 0.1, -math.Cos(angle) * 0.02}
	}
	return opts.New(TNTType, conf)
}

var tntConf = PassiveBehaviourConfig{
	Gravity: 0.04,
	Drag:    0.02,
	Expire:  explodeTNT,
}

func explodeTNT(e *Ent, tx *world.Tx) {
	block.ExplosionConfig{ItemDropChance: 1}.Explode(tx, world.EntityExplosionSource{
		Entity:        e,
		ExplosionSize: 4,
	})
}

var TNTType tntType

type tntType struct{}

func (t tntType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (tntType) EncodeEntity() string   { return "minecraft:tnt" }
func (tntType) NetworkOffset() float64 { return 0.49 }
func (tntType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.49, 0, -0.49, 0.49, 0.98, 0.49)
}

func (t tntType) DecodeNBT(m map[string]any, data *world.EntityData) {
	conf := tntConf
	conf.ExistenceDuration = nbtconv.TickDuration[int16](m, "Fuse")
	data.Data = conf.New()
}

func (tntType) EncodeNBT(data *world.EntityData) map[string]any {
	return map[string]any{"Fuse": int16(min(data.Data.(*PassiveBehaviour).Fuse()/(time.Second/20), math.MaxInt16))}
}
