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
	"image"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/customblock"
)

type Block interface {
	EncodeBlock() (string, map[string]any)

	Hash() (uint64, uint64)

	Model() BlockModel
}

type CustomBlock interface {
	Block
	Properties() customblock.Properties
}

type CustomBlockBuildable interface {
	CustomBlock

	Name() string

	Geometry() []byte

	Textures() map[string]image.Image
}

type Liquid interface {
	Block

	LiquidDepth() int

	SpreadDecay() int

	WithDepth(depth int, falling bool) Liquid

	LiquidFalling() bool

	BlastResistance() float64

	LiquidType() string

	Harden(pos cube.Pos, tx *Tx, flownIntoBy *cube.Pos) bool

	LiquidRemoveBlock(pos cube.Pos, tx *Tx, removed Block)
}

func RegisterBlock(b Block) {
	DefaultBlockRegistry.RegisterBlock(b)
}

func BlockHash(b Block) uint64 {
	return DefaultBlockRegistry.BlockHash(b)
}

func BlockRuntimeID(b Block) uint32 {
	return DefaultBlockRegistry.BlockRuntimeID(b)
}

func BlockByRuntimeID(rid uint32) (Block, bool) {
	return DefaultBlockRegistry.BlockByRuntimeID(rid)
}

func BlockByName(name string, properties map[string]any) (Block, bool) {
	return DefaultBlockRegistry.BlockByName(name, properties)
}

func Blocks() []Block {
	return DefaultBlockRegistry.Blocks()
}

func CustomBlocks() map[string]CustomBlock {
	return DefaultBlockRegistry.CustomBlocks()
}

type RandomTicker interface {
	RandomTick(pos cube.Pos, tx *Tx, r *rand.Rand)
}

type ScheduledTicker interface {
	ScheduledTick(pos cube.Pos, tx *Tx, r *rand.Rand)
}

type TickerBlock interface {
	NBTer
	Tick(currentTick int64, pos cube.Pos, tx *Tx)
}

type NeighbourUpdateTicker interface {
	NeighbourUpdateTick(pos, changedNeighbour cube.Pos, tx *Tx)
}

type NBTer interface {
	DecodeNBT(data map[string]any) any

	EncodeNBT() map[string]any
}

type LiquidDisplacer interface {
	CanDisplace(b Liquid) bool

	SideClosed(pos, side cube.Pos, tx *Tx) bool
}

type lightEmitter interface {
	LightEmissionLevel() uint8
}

type lightDiffuser interface {
	LightDiffusionLevel() uint8
}

type replaceableBlock interface {
	ReplaceableBy(b Block) bool
}

func replaceable(w *World, c *Column, pos cube.Pos, with Block) bool {
	if r, ok := w.blockInChunk(c, pos).(replaceableBlock); ok {
		return r.ReplaceableBy(with)
	}
	return false
}

type BlockAction interface {
	BlockAction()
}
