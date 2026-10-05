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
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

func NewEnderPearl(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	conf := enderPearlConf
	conf.Owner = owner.H()
	return opts.New(EnderPearlType, conf)
}

var enderPearlConf = ProjectileBehaviourConfig{
	Gravity:  0.03,
	Drag:     0.01,
	Particle: particle.EndermanTeleport{},
	Sound:    sound.Teleport{},
	Hit:      teleport,

	PassesThroughEntities: true,
}

type teleporter interface {
	Teleport(pos mgl64.Vec3)
	Living
}

type pearlTeleporter interface {
	TeleportFromProjectile(pos mgl64.Vec3)
}

const pearlClearance = 0.35

func landing(result trace.Result) mgl64.Vec3 {
	pos := result.Position()
	hit, ok := result.(trace.BlockResult)
	if !ok || hit.Face().Axis() == cube.Y {

		return pos
	}

	at, side := hit.BlockPosition(), hit.BlockPosition().Side(hit.Face())
	out := mgl64.Vec3{float64(side[0] - at[0]), float64(side[1] - at[1]), float64(side[2] - at[2])}
	return pos.Add(out.Mul(pearlClearance))
}

func teleport(e *Ent, tx *world.Tx, target trace.Result) {
	behaviour := e.Behaviour().(*ProjectileBehaviour)
	if behaviour.PortalTravel() {
		return
	}
	owner, _ := behaviour.Owner().Entity(tx)
	if user, ok := owner.(teleporter); ok {
		tx.PlaySound(user.Position(), sound.Teleport{})
		pos := landing(target)
		if pearl, ok := user.(pearlTeleporter); ok {
			pearl.TeleportFromProjectile(pos)
		} else {
			user.Teleport(pos)
		}
		if tx.World().FallDamage() {

			user.Hurt(5, FallDamageSource{})
		}
	}
}

var EnderPearlType enderPearlType

type enderPearlType struct{}

func (t enderPearlType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (enderPearlType) EncodeEntity() string { return "minecraft:ender_pearl" }
func (enderPearlType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}
func (enderPearlType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	data.Data = enderPearlConf.New()
}
func (enderPearlType) EncodeNBT(*world.EntityData) map[string]any { return nil }
