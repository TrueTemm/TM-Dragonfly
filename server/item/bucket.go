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

package item

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type BucketContent struct {
	liquid world.Liquid
	milk   bool
}

func LiquidBucketContent(l world.Liquid) BucketContent {
	return BucketContent{liquid: l}
}

func MilkBucketContent() BucketContent {
	return BucketContent{milk: true}
}

func (b BucketContent) Liquid() (world.Liquid, bool) {
	return b.liquid, b.liquid != nil
}

func (b BucketContent) String() string {
	if b.milk {
		return "milk"
	} else if b.liquid != nil {
		return b.liquid.LiquidType()
	}
	return ""
}

func (b BucketContent) LiquidType() string {
	if b.liquid != nil {
		return b.liquid.LiquidType()
	}
	return "milk"
}

type Bucket struct {
	Content BucketContent
}

func (b Bucket) MaxCount() int {
	if b.Empty() {
		return 16
	}
	return 1
}

func (b Bucket) AlwaysConsumable() bool {
	return b.Content.milk
}

func (b Bucket) CanConsume() bool {
	return b.Content.milk
}

func (b Bucket) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}

func (b Bucket) Consume(_ *world.Tx, c Consumer) Stack {
	for _, effect := range c.Effects() {
		c.RemoveEffect(effect.Type())
	}
	return NewStack(Bucket{}, 1)
}

func (b Bucket) Empty() bool {
	return b.Content.liquid == nil && !b.Content.milk
}

func (b Bucket) FuelInfo() FuelInfo {
	if liq := b.Content.liquid; liq != nil && liq.LiquidType() == "lava" {
		return newFuelInfo(time.Second * 1000).WithResidue(NewStack(Bucket{}, 1))
	}
	return FuelInfo{}
}

func (b Bucket) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	if b.Content.milk {
		return false
	}
	if b.Empty() {
		return b.fillFrom(pos, tx, ctx)
	}
	liq := b.Content.liquid.WithDepth(8, false)
	if bl := tx.Block(pos); canDisplace(bl, liq) || replaceableWith(bl, liq) {
		tx.SetLiquid(pos, liq)
	} else if bl := tx.Block(pos.Side(face)); canDisplace(bl, liq) || replaceableWith(bl, liq) {
		tx.SetLiquid(pos.Side(face), liq)
	} else {
		return false
	}

	tx.PlaySound(pos.Vec3Centre(), sound.BucketEmpty{Liquid: b.Content.liquid})
	ctx.NewItem = NewStack(Bucket{}, 1)
	ctx.NewItemSurvivalOnly = true
	ctx.SubtractFromCount(1)
	return true
}

func (b Bucket) fillFrom(pos cube.Pos, tx *world.Tx, ctx *UseContext) bool {
	liquid, ok := tx.Liquid(pos)
	if !ok {
		return false
	}
	if liquid.LiquidDepth() != 8 || liquid.LiquidFalling() {

		return false
	}
	tx.SetLiquid(pos, nil)
	tx.PlaySound(pos.Vec3Centre(), sound.BucketFill{Liquid: liquid})

	ctx.NewItem = NewStack(Bucket{Content: LiquidBucketContent(liquid)}, 1)
	ctx.NewItemSurvivalOnly = true
	ctx.SubtractFromCount(1)
	return true
}

func (b Bucket) EncodeItem() (name string, meta int16) {
	if !b.Empty() {
		return "minecraft:" + b.Content.String() + "_bucket", 0
	}
	return "minecraft:bucket", 0
}

type replaceable interface {
	ReplaceableBy(b world.Block) bool
}

func replaceableWith(b world.Block, with world.Block) bool {
	if r, ok := b.(replaceable); ok {
		return r.ReplaceableBy(with)
	}
	return false
}

func canDisplace(b world.Block, liq world.Liquid) bool {
	if d, ok := b.(world.LiquidDisplacer); ok {
		return d.CanDisplace(liq)
	}
	return false
}
