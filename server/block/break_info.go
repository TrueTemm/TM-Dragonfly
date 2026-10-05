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

package block

import (
	"math"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
)

type Breakable interface {
	BreakInfo() BreakInfo
}

type BreakContext struct {
	HasteLevel int

	ConduitPowerLevel int

	MiningFatigueLevel int

	Underwater bool

	AquaAffinity bool

	Flying bool

	Airborne bool
}

func BreakDuration(b world.Block, i item.Stack, ctx BreakContext) time.Duration {
	breakable, ok := b.(Breakable)
	if !ok {
		return math.MaxInt64
	}
	info := breakable.BreakInfo()
	if info.Hardness <= 0 {
		return 0
	}
	t, ok := i.Item().(item.Tool)
	if !ok {
		t = item.ToolNone{}
	}

	canHarvest := info.Harvestable(t)
	speed := 1.0
	if info.Effective(t) {
		speed = t.BaseMiningEfficiency(b)
		if !canHarvest {

			speed = 1
		} else if e, ok := i.Enchantment(enchantment.Efficiency); ok {
			speed += enchantment.Efficiency.Addend(e.Level())
		}
	}

	positive := max(ctx.HasteLevel, ctx.ConduitPowerLevel)
	if positive > 0 {
		speed *= 0.2*float64(positive) + 1
	}
	if ctx.MiningFatigueLevel > 0 {
		speed *= math.Pow(0.3, float64(ctx.MiningFatigueLevel))
	}
	if ctx.Underwater && !ctx.AquaAffinity {
		speed /= 5
	}
	if ctx.Airborne && !ctx.Flying {
		speed /= 5
	}

	damage := speed / info.Hardness
	if canHarvest {
		damage /= 30
	} else {
		damage /= 100
	}
	if positive > 0 {
		damage *= math.Pow(1.2, float64(positive))
	}
	if ctx.MiningFatigueLevel > 0 {
		damage *= math.Pow(0.7, float64(ctx.MiningFatigueLevel))
	}
	if damage >= 1 {

		return 0
	}
	return time.Duration(math.Ceil(1/damage)) * time.Second / 20
}

func BreaksInstantly(b world.Block) bool {
	breakable, ok := b.(Breakable)
	return ok && breakable.BreakInfo().Hardness <= 0
}

type BreakInfo struct {
	Hardness float64

	Harvestable func(t item.Tool) bool

	Effective func(t item.Tool) bool

	Drops func(t item.Tool, enchantments []item.Enchantment) []item.Stack

	BreakHandler func(pos cube.Pos, w *world.Tx, u item.User)

	XPDrops XPDropRange

	BlastResistance float64
}

func newBreakInfo(hardness float64, harvestable func(item.Tool) bool, effective func(item.Tool) bool, drops func(item.Tool, []item.Enchantment) []item.Stack) BreakInfo {
	return BreakInfo{
		Hardness:        hardness,
		BlastResistance: hardness,
		Harvestable:     harvestable,
		Effective:       effective,
		Drops:           drops,
	}
}

func (b BreakInfo) withXPDropRange(min, max int) BreakInfo {
	b.XPDrops = XPDropRange{min, max}
	return b
}

func (b BreakInfo) withBlastResistance(res float64) BreakInfo {
	b.BlastResistance = res
	return b
}

func (b BreakInfo) withBreakHandler(handler func(pos cube.Pos, w *world.Tx, u item.User)) BreakInfo {
	b.BreakHandler = handler
	return b
}

type XPDropRange [2]int

func (r XPDropRange) RandomValue() int {
	diff := r[1] - r[0]

	return rand.IntN(diff+1) + r[0]
}

var pickaxeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypePickaxe
}

var axeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeAxe
}

var shearsEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeShears
}

var swordEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeSword
}

var shovelEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeShovel
}

var hoeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeHoe
}

var nothingEffective = func(item.Tool) bool {
	return false
}

var alwaysHarvestable = func(t item.Tool) bool {
	return true
}

var neverHarvestable = func(t item.Tool) bool {
	return false
}

var pickaxeHarvestable = pickaxeEffective

func simpleDrops(s ...item.Stack) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(item.Tool, []item.Enchantment) []item.Stack {
		return s
	}
}

func oneOf(i ...world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(item.Tool, []item.Enchantment) []item.Stack {
		var s []item.Stack
		for _, it := range i {
			s = append(s, item.NewStack(it, 1))
		}
		return s
	}
}

func hasSilkTouch(enchantments []item.Enchantment) bool {
	return slices.IndexFunc(enchantments, func(i item.Enchantment) bool {
		return i.Type() == enchantment.SilkTouch
	}) != -1
}

func silkTouchOneOf(normal, silkTouch world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(silkTouch, 1)}
		}
		return []item.Stack{item.NewStack(normal, 1)}
	}
}

func silkTouchDrop(normal, silkTouch item.Stack) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{silkTouch}
		}
		return []item.Stack{normal}
	}
}

func silkTouchOnlyDrop(it world.Item) func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(it, 1)}
		}
		return nil
	}
}

func fortuneLevel(enchantments []item.Enchantment) int {
	index := slices.IndexFunc(enchantments, func(i item.Enchantment) bool {
		return i.Type() == enchantment.Fortune
	})
	if index == -1 {
		return 0
	}
	return enchantments[index].Level()
}

func fortuneOreCount(base int, enchantments []item.Enchantment) int {
	fortune := fortuneLevel(enchantments)
	if fortune == 0 || rand.IntN(fortune+2) < 2 {
		return base
	}
	multiplier := rand.IntN(fortune) + 2
	return base * multiplier
}

func fortuneDiscreteCount(minCount, maxCount, capCount int, enchantments []item.Enchantment) int {
	fortune := fortuneLevel(enchantments)
	maxWithFortune := maxCount + fortune
	return min(capCount, rand.IntN(maxWithFortune-minCount+1)+minCount)
}

func fortuneBinomial(attempts int) int {
	count := 0
	for range attempts {
		if rand.IntN(15) < 8 {
			count++
		}
	}
	return count
}

func oreDrops(drop, block world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(block, 1)}
		}
		return []item.Stack{item.NewStack(drop, fortuneOreCount(1, enchantments))}
	}
}

func multiOreDrops(drop, block world.Item, minCount, maxCount int) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(block, 1)}
		}
		baseCount := rand.IntN(maxCount-minCount+1) + minCount
		return []item.Stack{item.NewStack(drop, fortuneOreCount(baseCount, enchantments))}
	}
}

func discreteDrops(drop, block world.Item, minCount, maxCount, capCount int) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(block, 1)}
		}
		return []item.Stack{item.NewStack(drop, fortuneDiscreteCount(minCount, maxCount, capCount, enchantments))}
	}
}

func grassDrops(grass world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if t.ToolType() == item.TypeShears || hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(grass, 1)}
		}
		if rand.Float32() < 0.125 {
			count := 1
			if fortune := fortuneLevel(enchantments); fortune > 0 {
				count += rand.IntN(fortune*2 + 1)
			}
			return []item.Stack{item.NewStack(WheatSeeds{}, count)}
		}
		return nil
	}
}

func cropSeedDrops(seed, crop world.Item, growth int) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if growth < 7 {
			return []item.Stack{item.NewStack(seed, 1)}
		}
		seedCount := fortuneBinomial(3 + fortuneLevel(enchantments))
		if seedCount == 0 {
			return []item.Stack{item.NewStack(crop, 1)}
		}
		return []item.Stack{item.NewStack(crop, 1), item.NewStack(seed, seedCount)}
	}
}

func breakBlock(b world.Block, pos cube.Pos, tx *world.Tx) {
	breakBlockNoDrops(b, pos, tx)
	if breakable, ok := b.(Breakable); ok {
		for _, drop := range breakable.BreakInfo().Drops(item.ToolNone{}, nil) {
			dropItem(tx, drop, pos.Vec3Centre())
		}
	}
}

func breakBlockNoDrops(b world.Block, pos cube.Pos, tx *world.Tx) {

	tx.SetBlock(pos, nil, nil)
	if breakable, ok := b.(Breakable); ok {
		breakHandler := breakable.BreakInfo().BreakHandler
		if breakHandler != nil {
			breakHandler(pos, tx, nil)
		}
	}
	tx.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: b})
}
