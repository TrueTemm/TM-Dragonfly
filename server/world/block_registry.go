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
	"fmt"
	"maps"
	"math"
	"math/bits"
	"slices"
	"sort"
	"sync"

	"github.com/brentp/intintmap"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/segmentio/fasthash/fnv1"
)

var DefaultBlockRegistry = &BasicBlockRegistry{
	blockProperties: make(map[string]map[string]any),
	stateRuntimeIDs: make(map[stateHash]uint32),
	customBlocks:    make(map[string]CustomBlock),
}

type BlockRegistry interface {
	chunk.BlockRegistry

	BlockByRuntimeID(rid uint32) (Block, bool)

	BlockByRuntimeIDOrAir(rid uint32) Block

	BlockRuntimeID(block Block) (rid uint32)

	RegisterBlock(block Block)

	RegisterBlockState(blockState BlockState)

	CustomBlocks() map[string]CustomBlock

	BlockByName(name string, properties map[string]any) (Block, bool)

	Blocks() []Block

	Air() Block

	Finalize()

	BitSize() int

	BlockHash(b Block) uint64

	RuntimeIDToHash(runtimeID uint32) (hash uint32, ok bool)
}

const (
	blockFlagNBT uint16 = 1 << iota
	blockFlagRandomTick
	blockFlagLiquid
	blockFlagLiquidDisplacing
)

type blockInfo uint16

func (b *blockInfo) set(flag uint16) {
	*b |= blockInfo(flag)
}

func (b blockInfo) get(flag uint16) bool {
	return uint16(b)&flag != 0
}

func (b *blockInfo) setLight(light uint8) {

	*b &^= blockInfo(0xF) << 8
	*b |= blockInfo(light&0xF) << 8
}

func (b *blockInfo) setLightFilter(light uint8) {

	*b &^= blockInfo(0xF) << 12
	*b |= blockInfo(light&0xF) << 12
}

func (b blockInfo) getLight() uint8 {
	return uint8((b >> 8) & 0xF)
}

func (b blockInfo) getLightFilter() uint8 {
	return uint8((b >> 12) & 0xF)
}

type BasicBlockRegistry struct {
	mu sync.Mutex

	finalized         bool
	bitSize           int
	hashes            *intintmap.Map
	networkhashToRids map[uint32]uint32
	ridsToNetworkhash []uint32

	stateRuntimeIDs map[stateHash]uint32

	blockProperties map[string]map[string]any

	blocks []Block

	customBlocks map[string]CustomBlock

	blockInfos []blockInfo

	airRID uint32
}

func (br *BasicBlockRegistry) BitSize() int {
	if !br.finalized {
		panic("BlockRegistry.BitSize called on non finalized BlockRegistry")
	}
	return br.bitSize
}

func (br *BasicBlockRegistry) BlockCount() int {
	if !br.finalized {
		panic("BlockRegistry.BlockCount called on non finalized BlockRegistry")
	}
	return len(br.blockInfos)
}

func (br *BasicBlockRegistry) RandomTickBlock(rid uint32) bool {
	if !br.finalized {
		panic("BlockRegistry.RandomTickBlock called on non finalized BlockRegistry")
	}
	return br.blockInfos[rid].get(blockFlagRandomTick)
}

func (br *BasicBlockRegistry) FilteringBlock(rid uint32) uint8 {
	if !br.finalized {
		panic("BlockRegistry.FilteringBlock called on non finalized BlockRegistry")
	}
	return br.blockInfos[rid].getLightFilter()
}

func (br *BasicBlockRegistry) LightBlock(rid uint32) uint8 {
	if !br.finalized {
		panic("BlockRegistry.LightBlock called on non finalized BlockRegistry")
	}
	return br.blockInfos[rid].getLight()
}

func (br *BasicBlockRegistry) NBTBlock(rid uint32) bool {
	if !br.finalized {
		panic("BlockRegistry.NBTBlock called on non finalized BlockRegistry")
	}
	return br.blockInfos[rid].get(blockFlagNBT)
}

func (br *BasicBlockRegistry) LiquidDisplacingBlock(rid uint32) bool {
	if !br.finalized {
		panic("BlockRegistry.LiquidDisplacingBlock called on non finalized BlockRegistry")
	}
	return br.blockInfos[rid].get(blockFlagLiquidDisplacing)
}

func (br *BasicBlockRegistry) LiquidBlock(rid uint32) bool {
	if !br.finalized {
		panic("BlockRegistry.LiquidBlock called on non finalized BlockRegistry")
	}
	return br.blockInfos[rid].get(blockFlagLiquid)
}

func (br *BasicBlockRegistry) Blocks() []Block {
	if !br.finalized {
		panic("BlockRegistry.Blocks called on non finalized BlockRegistry")
	}
	return slices.Clone(br.blocks)
}

func (br *BasicBlockRegistry) HashToRuntimeID(hash uint32) (rid uint32, ok bool) {
	if !br.finalized {
		panic("BlockRegistry.HashToRuntimeID called on non finalized BlockRegistry")
	}
	rid, ok = br.networkhashToRids[hash]
	return
}

func (br *BasicBlockRegistry) RuntimeIDToHash(runtimeID uint32) (hash uint32, ok bool) {
	if !br.finalized {
		panic("BlockRegistry.RuntimeIDToHash called on non finalized BlockRegistry")
	}
	if runtimeID >= uint32(len(br.ridsToNetworkhash)) {
		return 0, false
	}
	return br.ridsToNetworkhash[runtimeID], true
}

func (br *BasicBlockRegistry) Clone() *BasicBlockRegistry {
	br.mu.Lock()
	defer br.mu.Unlock()

	br2 := &BasicBlockRegistry{
		blockProperties: make(map[string]map[string]any, len(br.blockProperties)),
		stateRuntimeIDs: make(map[stateHash]uint32, len(br.stateRuntimeIDs)),
		customBlocks:    make(map[string]CustomBlock, len(br.customBlocks)),
		finalized:       br.finalized,
		bitSize:         br.bitSize,
		airRID:          br.airRID,
	}

	for k, v := range br.blockProperties {
		br2.blockProperties[k] = maps.Clone(v)
	}
	maps.Copy(br2.stateRuntimeIDs, br.stateRuntimeIDs)
	maps.Copy(br2.customBlocks, br.customBlocks)

	br2.blocks = make([]Block, len(br.blocks))
	copy(br2.blocks, br.blocks)
	br2.blockInfos = append([]blockInfo(nil), br.blockInfos...)

	if br.finalized {
		br2.hashes = intintmap.New(len(br.blocks), 0.999)
		for rid, b := range br2.blocks {
			if _, hash := b.Hash(); hash == math.MaxUint64 {
				continue
			}
			br2.hashes.Put(int64(br2.BlockHash(b)), int64(rid))
		}
		br2.networkhashToRids = make(map[uint32]uint32, len(br.networkhashToRids))
		maps.Copy(br2.networkhashToRids, br.networkhashToRids)
		br2.ridsToNetworkhash = append([]uint32(nil), br.ridsToNetworkhash...)
	}

	return br2
}

func NewBlockRegistry() BlockRegistry {

	br := DefaultBlockRegistry.Clone()
	br.finalized = false
	br.bitSize = 0
	br.hashes = nil
	br.networkhashToRids = nil
	br.ridsToNetworkhash = nil
	br.blockInfos = nil
	return br
}

func (br *BasicBlockRegistry) RegisterBlock(b Block) {
	br.mu.Lock()
	defer br.mu.Unlock()

	if br.finalized {
		panic("BlockRegistry.RegisterBlock called on finalized BlockRegistry")
	}
	name, properties := b.EncodeBlock()
	if _, ok := b.(CustomBlock); ok {
		br.registerBlockStateLocked(BlockState{Name: name, Properties: properties})
	}
	rid, ok := br.stateRuntimeIDs[stateHash{name: name, properties: hashProperties(properties)}]
	if !ok {

		panic(fmt.Sprintf("block state returned is not registered (%v {%#v})", name, properties))
	}
	if _, ok := br.blocks[rid].(unknownBlock); !ok {
		panic(fmt.Sprintf("block with name and properties %v {%#v} already registered", name, properties))
	}
	br.blocks[rid] = b
	if c, ok := b.(CustomBlock); ok {
		if _, ok := br.customBlocks[name]; !ok {
			br.customBlocks[name] = c
		}
	}
}

func (br *BasicBlockRegistry) RegisterBlockState(s BlockState) {
	br.mu.Lock()
	defer br.mu.Unlock()
	br.registerBlockStateLocked(s)
}

func (br *BasicBlockRegistry) registerBlockStateLocked(s BlockState) {
	if br.finalized {
		panic("BlockRegistry.RegisterBlockState called on finalized BlockRegistry")
	}
	h := stateHash{name: s.Name, properties: hashProperties(s.Properties)}
	if _, ok := br.stateRuntimeIDs[h]; ok {
		panic(fmt.Sprintf("cannot register the same state twice (%+v)", s))
	}
	if _, ok := br.blockProperties[s.Name]; !ok {
		br.blockProperties[s.Name] = s.Properties
	}
	rid := uint32(len(br.blocks))
	br.blocks = append(br.blocks, unknownBlock{BlockState: s})
	br.stateRuntimeIDs[h] = rid
}

func (br *BasicBlockRegistry) Finalize() {
	br.mu.Lock()
	defer br.mu.Unlock()

	if br.finalized {

		return
	}

	br.bitSize = bits.Len64(uint64(len(br.blocks)))
	sort.SliceStable(br.blocks, func(i, j int) bool {
		var nameOne string
		if b1, ok := br.blocks[i].(unknownBlock); ok {
			nameOne = b1.Name
		} else {
			nameOne, _ = br.blocks[i].EncodeBlock()
		}
		var nameTwo string
		if b2, ok := br.blocks[j].(unknownBlock); ok {
			nameTwo = b2.Name
		} else {
			nameTwo, _ = br.blocks[j].EncodeBlock()
		}
		return fnv1.HashString64(nameOne) < fnv1.HashString64(nameTwo)
	})

	br.blockInfos = make([]blockInfo, len(br.blocks))
	br.hashes = intintmap.New(len(br.blocks), 0.999)
	br.networkhashToRids = make(map[uint32]uint32, len(br.blocks))
	br.ridsToNetworkhash = make([]uint32, len(br.blocks))
	br.stateRuntimeIDs = make(map[stateHash]uint32, len(br.blocks))
	networkHashScratch := make([]byte, 0, 0xff)
	foundAir := false

	for idx, b := range br.blocks {
		rid := uint32(idx)
		name, properties := b.EncodeBlock()
		h := stateHash{name: name, properties: hashProperties(properties)}
		if name == "minecraft:air" {
			br.airRID = rid
			foundAir = true
		}
		if _, ok := br.stateRuntimeIDs[h]; ok {
			panic(fmt.Sprintf("cannot register the same state twice (%s %+v)", name, properties))
		}
		br.stateRuntimeIDs[h] = rid

		var info blockInfo

		info.setLightFilter(15)
		if diffuser, ok := b.(lightDiffuser); ok {
			info.setLightFilter(diffuser.LightDiffusionLevel())
		}
		if emitter, ok := b.(lightEmitter); ok {
			info.setLight(emitter.LightEmissionLevel())
		}
		if _, ok := b.(NBTer); ok {
			info.set(blockFlagNBT)
		}
		if _, ok := b.(RandomTicker); ok {
			info.set(blockFlagRandomTick)
		}
		if _, ok := b.(Liquid); ok {
			info.set(blockFlagLiquid)
		}
		if _, ok := b.(LiquidDisplacer); ok {
			info.set(blockFlagLiquidDisplacing)
		}
		br.blockInfos[rid] = info

		if _, hash := b.Hash(); hash != math.MaxUint64 {
			h := int64(br.BlockHash(b))
			if other, ok := br.hashes.Get(h); ok {
				panic(fmt.Sprintf("block %#v with hash %v already registered by %#v", b, h, br.blocks[other]))
			}
			br.hashes.Put(h, int64(rid))
		}
		var netHash uint32
		netHash, networkHashScratch = networkBlockHash(name, properties, networkHashScratch)
		if other, ok := br.networkhashToRids[netHash]; ok {
			otherName, otherProperties := br.blocks[other].EncodeBlock()
			panic(fmt.Sprintf("network block hash collision for (%s %+v) and (%s %+v)", name, properties, otherName, otherProperties))
		}
		br.networkhashToRids[netHash] = rid
		br.ridsToNetworkhash[rid] = netHash
	}
	if !foundAir {
		panic("BlockRegistry.Finalize: no minecraft:air block state registered")
	}
	br.finalized = true
}

func (br *BasicBlockRegistry) AirRuntimeID() uint32 {
	if !br.finalized {
		panic("BlockRegistry.AirRuntimeID called on non finalized BlockRegistry")
	}
	return br.airRID
}

func (br *BasicBlockRegistry) RuntimeIDToState(runtimeID uint32) (name string, properties map[string]any, found bool) {
	if !br.finalized {
		panic("BlockRegistry.RuntimeIDToState called on non finalized BlockRegistry")
	}
	if runtimeID >= uint32(len(br.blocks)) {
		return "", nil, false
	}
	name, properties = br.blocks[runtimeID].EncodeBlock()
	return name, properties, true
}

func (br *BasicBlockRegistry) StateToRuntimeID(name string, properties map[string]any) (runtimeID uint32, found bool) {
	if !br.finalized {
		panic("BlockRegistry.StateToRuntimeID called on non finalized BlockRegistry")
	}
	if rid, ok := br.stateRuntimeIDs[stateHash{name: name, properties: hashProperties(properties)}]; ok {
		return rid, true
	}
	rid, ok := br.stateRuntimeIDs[stateHash{name: name, properties: hashProperties(br.blockProperties[name])}]
	return rid, ok
}

func (br *BasicBlockRegistry) BlockHash(b Block) uint64 {
	base, hash := b.Hash()
	return base | (hash << uint64(br.bitSize))
}

func (br *BasicBlockRegistry) BlockRuntimeID(b Block) uint32 {
	if !br.finalized {
		panic("BlockRegistry.BlockRuntimeID called on non finalized BlockRegistry")
	}
	if b == nil {
		return br.airRID
	}
	if _, hash := b.Hash(); hash != math.MaxUint64 {

		if rid, ok := br.hashes.Get(int64(br.BlockHash(b))); ok {
			return uint32(rid)
		}
		panic(fmt.Sprintf("cannot find block by non-0 hash of block %#v", b))
	}
	return br.slowBlockRuntimeID(b)
}

func (br *BasicBlockRegistry) BlockByRuntimeIDOrAir(rid uint32) Block {
	bl, _ := br.BlockByRuntimeID(rid)
	return bl
}

func (br *BasicBlockRegistry) slowBlockRuntimeID(b Block) uint32 {
	name, properties := b.EncodeBlock()
	rid, ok := br.stateRuntimeIDs[stateHash{name: name, properties: hashProperties(properties)}]
	if !ok {
		panic(fmt.Sprintf("cannot find block by (name + properties): %#v", b))
	}
	return rid
}

func (br *BasicBlockRegistry) BlockByRuntimeID(rid uint32) (Block, bool) {
	if !br.finalized {
		panic("BlockRegistry.BlockByRuntimeID called on non finalized BlockRegistry")
	}
	if rid >= uint32(len(br.blocks)) {
		return br.Air(), false
	}
	return br.blocks[rid], true
}

func (br *BasicBlockRegistry) BlockByName(name string, properties map[string]any) (Block, bool) {
	if !br.finalized {
		panic("BlockRegistry.BlockByName called on non finalized BlockRegistry")
	}
	rid, ok := br.stateRuntimeIDs[stateHash{name: name, properties: hashProperties(properties)}]
	if !ok {
		return nil, false
	}
	return br.blocks[rid], true
}

func (br *BasicBlockRegistry) CustomBlocks() map[string]CustomBlock {
	return maps.Clone(br.customBlocks)
}

func (br *BasicBlockRegistry) Air() Block {
	if !br.finalized {
		panic("BlockRegistry.Air called on non finalized BlockRegistry")
	}
	if br.airRID >= uint32(len(br.blocks)) {

		panic("BlockRegistry.Air: air runtime ID out of range")
	}
	return br.blocks[br.airRID]
}
