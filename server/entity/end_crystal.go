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
)

func NewEndCrystal(opts world.EntitySpawnOpts) *world.EntityHandle {
	return EndCrystalConfig{}.New(opts)
}

type EndCrystalConfig struct {
	ExplosionSize float64

	ShowBase bool
}

func (c EndCrystalConfig) New(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(EndCrystalType, c)
}

func (c EndCrystalConfig) Apply(data *world.EntityData) {
	if c.ExplosionSize == 0 {
		c.ExplosionSize = 6
	}
	data.Data = endCrystalBehaviour{
		showBase:      c.ShowBase,
		explosionSize: c.ExplosionSize,
	}
}

var EndCrystalType endCrystalType

type endCrystalType struct{}

func (endCrystalType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return Open(tx, handle, data)
}

func (endCrystalType) EncodeEntity() string {
	return "minecraft:ender_crystal"
}

func (endCrystalType) BBox(world.Entity) cube.BBox {
	return cube.Box(-1, 0, -1, 1, 2, 1)
}

func (endCrystalType) DecodeNBT(m map[string]any, data *world.EntityData) {
	b := endCrystalBehaviour{
		showBase:      nbtconv.Bool(m, "ShowBottom"),
		explosionSize: nbtconv.Float64(m, "ExplosionSize"),
	}
	if b.explosionSize == 0 {
		b.explosionSize = 6
	}
	x, hasX := m["BlockTargetX"].(int32)
	y, hasY := m["BlockTargetY"].(int32)
	z, hasZ := m["BlockTargetZ"].(int32)
	if hasX && hasY && hasZ {
		b.beamTarget = cube.Pos{int(x), int(y), int(z)}
		b.hasBeamTarget = true
	}
	b.Apply(data)
}

func (endCrystalType) EncodeNBT(data *world.EntityData) map[string]any {
	b := data.Data.(endCrystalBehaviour)
	m := map[string]any{
		"ShowBottom":    boolByte(b.showBase),
		"ExplosionSize": b.explosionSize,
	}
	if b.hasBeamTarget {
		m["BlockTargetX"] = int32(b.beamTarget[0])
		m["BlockTargetY"] = int32(b.beamTarget[1])
		m["BlockTargetZ"] = int32(b.beamTarget[2])
	}
	return m
}
