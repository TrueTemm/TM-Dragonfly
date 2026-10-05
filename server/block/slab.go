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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Slab struct {
	Block world.Block

	Top bool

	Double bool
}

func (s Slab) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	id, meta := s.EncodeItem()
	clickedBlock := tx.Block(pos)
	if clickedSlab, ok := clickedBlock.(Slab); ok && !s.Double {
		clickedId, clickedMeta := clickedSlab.EncodeItem()
		if !clickedSlab.Double && id == clickedId && meta == clickedMeta && ((face == cube.FaceUp && !clickedSlab.Top) || (face == cube.FaceDown && clickedSlab.Top)) {

			clickedSlab.Double = true

			place(tx, pos, clickedSlab, user, ctx)
			return placed(ctx)
		}
	}
	if sideSlab, ok := tx.Block(pos.Side(face)).(Slab); ok && !replaceableWith(tx, pos, s) && !s.Double {
		sideId, sideMeta := sideSlab.EncodeItem()

		if !sideSlab.Double && id == sideId && meta == sideMeta {
			sideSlab.Double = true

			place(tx, pos.Side(face), sideSlab, user, ctx)
			return placed(ctx)
		}
	}
	pos, face, used = firstReplaceable(tx, pos, face, s)
	if !used {
		return
	}
	if face == cube.FaceDown || (clickPos[1] > 0.5 && face != cube.FaceUp) {
		s.Top = true
	}

	place(tx, pos, s, user, ctx)
	return placed(ctx)
}

func (s Slab) Instrument() sound.Instrument {
	if _, ok := s.Block.(Planks); ok {
		return sound.Bass()
	}
	if _, ok := s.Block.(BambooMosaic); ok {
		return sound.Bass()
	}
	return sound.BassDrum()
}

func (s Slab) FlammabilityInfo() FlammabilityInfo {
	if flammable, ok := s.Block.(Flammable); ok {
		return flammable.FlammabilityInfo()
	}
	return newFlammabilityInfo(0, 0, false)
}

func (s Slab) FuelInfo() item.FuelInfo {
	if fuel, ok := s.Block.(item.Fuel); ok {
		return fuel.FuelInfo()
	}
	return item.FuelInfo{}
}

func (s Slab) CanDisplace(b world.Liquid) bool {
	water, ok := b.(Water)
	return !s.Double && ok && water.Depth == 8
}

func (s Slab) SideClosed(pos, side cube.Pos, _ *world.Tx) bool {

	return !s.Top && side[1] == pos[1]-1
}

func (s Slab) LightDiffusionLevel() uint8 {
	if s.Double {
		return 15
	}
	return 0
}

func (s Slab) CanRedstoneWireStepDown(cube.Pos, cube.Pos, *world.Tx) bool {
	return s.Double
}

func (s Slab) BreakInfo() BreakInfo {
	hardness, blastResistance, harvestable, effective := 2.0, 6.0, pickaxeHarvestable, pickaxeEffective

	switch block := s.Block.(type) {
	case Stone, Sandstone, Quartz, Purpur, Blackstone, PolishedBlackstoneBrick:

	case EndBricks:
		hardness = 3
	case StoneBricks:
		if block.Type == MossyStoneBricks() {
			hardness = 1.5
		}
	case Breakable:
		breakInfo := block.BreakInfo()
		hardness, blastResistance, harvestable, effective = breakInfo.Hardness, breakInfo.BlastResistance, breakInfo.Harvestable, breakInfo.Effective
	}
	return newBreakInfo(hardness, harvestable, effective, func(tool item.Tool, enchantments []item.Enchantment) []item.Stack {
		single := Slab{Block: s.Block}
		if s.Double {
			return []item.Stack{item.NewStack(single, 2)}
		}
		return []item.Stack{item.NewStack(single, 1)}
	}).withBlastResistance(blastResistance)
}

func (s Slab) Model() world.BlockModel {
	return model.Slab{Double: s.Double, Top: s.Top}
}

func (s Slab) EncodeItem() (string, int16) {
	name, suffix := encodeSlabBlock(s.Block, false)
	return "minecraft:" + name + suffix, 0
}

func (s Slab) EncodeBlock() (string, map[string]any) {
	side := "bottom"
	if s.Top {
		side = "top"
	}
	name, suffix := encodeSlabBlock(s.Block, s.Double)
	return "minecraft:" + name + suffix, map[string]any{"minecraft:vertical_half": side}
}

func allSlabs() (b []world.Block) {
	for _, s := range SlabBlocks() {
		b = append(b, Slab{Block: s, Double: true})
		b = append(b, Slab{Block: s, Top: true, Double: true})
		b = append(b, Slab{Block: s, Top: true})
		b = append(b, Slab{Block: s})
	}
	return
}
