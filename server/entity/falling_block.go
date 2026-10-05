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

func NewFallingBlock(opts world.EntitySpawnOpts, block world.Block) *world.EntityHandle {
	conf := fallingBlockConf
	conf.Block = block
	return opts.New(FallingBlockType, conf)
}

var fallingBlockConf = FallingBlockBehaviourConfig{
	Gravity: 0.04,
	Drag:    0.02,
}

var FallingBlockType fallingBlockType

type fallingBlockType struct{}

func (t fallingBlockType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}
func (fallingBlockType) EncodeEntity() string   { return "minecraft:falling_block" }
func (fallingBlockType) NetworkOffset() float64 { return 0.49 }
func (fallingBlockType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.49, 0, -0.49, 0.49, 0.98, 0.49)
}

func (fallingBlockType) DecodeNBT(m map[string]any, data *world.EntityData) {
	conf := fallingBlockConf
	conf.Block = nbtconv.Block(m, "FallingBlock")
	conf.DistanceFallen = nbtconv.Float64(m, "FallDistance")
	data.Data = conf.New()
}

func (fallingBlockType) EncodeNBT(data *world.EntityData) map[string]any {
	b := data.Data.(*FallingBlockBehaviour)
	return map[string]any{"FallDistance": b.passive.fallDistance, "FallingBlock": nbtconv.WriteBlock(b.block)}
}
