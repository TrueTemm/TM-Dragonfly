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
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type HasFacing interface {
	world.Block
	FacingDirection() cube.Direction
	WithFacing(facing cube.Direction) world.Block
}

type HasAxis interface {
	world.Block
	PillarAxis() cube.Axis
	WithAxis(axis cube.Axis) world.Block
}

type HasColour interface {
	world.Block
	DyeColour() item.Colour
	WithColour(colour item.Colour) world.Block
}

func (b Anvil) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Anvil) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Bed) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Bed) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b BlastFurnace) FacingDirection() cube.Direction {
	return b.Facing
}

func (b BlastFurnace) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Campfire) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Campfire) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Chest) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Chest) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b CocoaBean) FacingDirection() cube.Direction {
	return b.Facing
}

func (b CocoaBean) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b CopperDoor) FacingDirection() cube.Direction {
	return b.Facing
}

func (b CopperDoor) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b CopperGolemStatue) FacingDirection() cube.Direction {
	return b.Facing
}

func (b CopperGolemStatue) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b CopperTrapdoor) FacingDirection() cube.Direction {
	return b.Facing
}

func (b CopperTrapdoor) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b DecoratedPot) FacingDirection() cube.Direction {
	return b.Facing
}

func (b DecoratedPot) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b EndPortalFrame) FacingDirection() cube.Direction {
	return b.Facing
}

func (b EndPortalFrame) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b EnderChest) FacingDirection() cube.Direction {
	return b.Facing
}

func (b EnderChest) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Furnace) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Furnace) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b GlazedTerracotta) FacingDirection() cube.Direction {
	return b.Facing
}

func (b GlazedTerracotta) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Grindstone) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Grindstone) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Ladder) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Ladder) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Lectern) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Lectern) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b LitPumpkin) FacingDirection() cube.Direction {
	return b.Facing
}

func (b LitPumpkin) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Loom) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Loom) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b PinkPetals) FacingDirection() cube.Direction {
	return b.Facing
}

func (b PinkPetals) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Pumpkin) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Pumpkin) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Smoker) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Smoker) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Stairs) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Stairs) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b Stonecutter) FacingDirection() cube.Direction {
	return b.Facing
}

func (b Stonecutter) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b WoodDoor) FacingDirection() cube.Direction {
	return b.Facing
}

func (b WoodDoor) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b WoodFenceGate) FacingDirection() cube.Direction {
	return b.Facing
}

func (b WoodFenceGate) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b WoodTrapdoor) FacingDirection() cube.Direction {
	return b.Facing
}

func (b WoodTrapdoor) WithFacing(facing cube.Direction) world.Block {
	b.Facing = facing
	return b
}

func (b BambooBlock) PillarAxis() cube.Axis {
	return b.Axis
}

func (b BambooBlock) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Basalt) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Basalt) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Bone) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Bone) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b CopperChain) PillarAxis() cube.Axis {
	return b.Axis
}

func (b CopperChain) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Deepslate) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Deepslate) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Froglight) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Froglight) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b HayBale) PillarAxis() cube.Axis {
	return b.Axis
}

func (b HayBale) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b InfestedDeepslate) PillarAxis() cube.Axis {
	return b.Axis
}

func (b InfestedDeepslate) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b IronChain) PillarAxis() cube.Axis {
	return b.Axis
}

func (b IronChain) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Log) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Log) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b MuddyMangroveRoots) PillarAxis() cube.Axis {
	return b.Axis
}

func (b MuddyMangroveRoots) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Portal) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Portal) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b PurpurPillar) PillarAxis() cube.Axis {
	return b.Axis
}

func (b PurpurPillar) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b QuartzPillar) PillarAxis() cube.Axis {
	return b.Axis
}

func (b QuartzPillar) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Wood) PillarAxis() cube.Axis {
	return b.Axis
}

func (b Wood) WithAxis(axis cube.Axis) world.Block {
	b.Axis = axis
	return b
}

func (b Banner) DyeColour() item.Colour {
	return b.Colour
}

func (b Banner) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b Bed) DyeColour() item.Colour {
	return b.Colour
}

func (b Bed) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b Carpet) DyeColour() item.Colour {
	return b.Colour
}

func (b Carpet) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b Concrete) DyeColour() item.Colour {
	return b.Colour
}

func (b Concrete) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b ConcretePowder) DyeColour() item.Colour {
	return b.Colour
}

func (b ConcretePowder) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b GlazedTerracotta) DyeColour() item.Colour {
	return b.Colour
}

func (b GlazedTerracotta) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b StainedGlass) DyeColour() item.Colour {
	return b.Colour
}

func (b StainedGlass) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b StainedGlassPane) DyeColour() item.Colour {
	return b.Colour
}

func (b StainedGlassPane) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b StainedTerracotta) DyeColour() item.Colour {
	return b.Colour
}

func (b StainedTerracotta) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}

func (b Wool) DyeColour() item.Colour {
	return b.Colour
}

func (b Wool) WithColour(colour item.Colour) world.Block {
	b.Colour = colour
	return b
}
