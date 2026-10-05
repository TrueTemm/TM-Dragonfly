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
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type Wall struct {
	transparent
	sourceWaterDisplacer

	Block world.Block

	NorthConnection WallConnectionType

	EastConnection WallConnectionType

	SouthConnection WallConnectionType

	WestConnection WallConnectionType

	Post bool
}

func (w Wall) EncodeItem() (string, int16) {
	name := encodeWallBlock(w.Block)
	return "minecraft:" + name + "_wall", 0
}

func (w Wall) EncodeBlock() (string, map[string]any) {
	properties := map[string]any{
		"wall_connection_type_north": w.NorthConnection.String(),
		"wall_connection_type_east":  w.EastConnection.String(),
		"wall_connection_type_south": w.SouthConnection.String(),
		"wall_connection_type_west":  w.WestConnection.String(),
		"wall_post_bit":              boolByte(w.Post),
	}
	name := encodeWallBlock(w.Block)
	return "minecraft:" + name + "_wall", properties
}

func (w Wall) Model() world.BlockModel {
	return model.Wall{
		NorthConnection: w.NorthConnection.Height(),
		EastConnection:  w.EastConnection.Height(),
		SouthConnection: w.SouthConnection.Height(),
		WestConnection:  w.WestConnection.Height(),
		Post:            w.Post,
	}
}

func (w Wall) BreakInfo() BreakInfo {
	breakable, ok := w.Block.(Breakable)
	if !ok {
		panic("wall block is not breakable")
	}
	return newBreakInfo(breakable.BreakInfo().Hardness, pickaxeHarvestable, pickaxeEffective, oneOf(w)).withBlastResistance(breakable.BreakInfo().BlastResistance)
}

func (w Wall) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	w, connectionsUpdated := w.calculateConnections(tx, pos)
	w, postUpdated := w.calculatePost(tx, pos)
	if connectionsUpdated || postUpdated {
		tx.SetBlock(pos, w, nil)
	}
}

func (w Wall) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, w)
	if !used {
		return
	}
	w, _ = w.calculateConnections(tx, pos)
	w, _ = w.calculatePost(tx, pos)
	place(tx, pos, w, user, ctx)
	return placed(ctx)
}

func (Wall) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (w Wall) ConnectionType(direction cube.Direction) WallConnectionType {
	switch direction {
	case cube.North:
		return w.NorthConnection
	case cube.East:
		return w.EastConnection
	case cube.South:
		return w.SouthConnection
	case cube.West:
		return w.WestConnection
	}
	panic("unknown direction")
}

func (w Wall) WithConnectionType(direction cube.Direction, connection WallConnectionType) Wall {
	switch direction {
	case cube.North:
		w.NorthConnection = connection
	case cube.East:
		w.EastConnection = connection
	case cube.South:
		w.SouthConnection = connection
	case cube.West:
		w.WestConnection = connection
	}
	return w
}

func (w Wall) calculateConnections(tx *world.Tx, pos cube.Pos) (Wall, bool) {
	var updated bool
	abovePos := pos.Add(cube.Pos{0, 1, 0})
	above := tx.Block(abovePos)
	for _, face := range cube.HorizontalFaces() {
		sidePos := pos.Side(face)
		side := tx.Block(sidePos)

		connected := side.Model().FaceSolid(sidePos, face.Opposite(), tx)
		if !connected {
			if _, ok := tx.Block(sidePos).(Wall); ok {
				connected = true
			} else if gate, ok := tx.Block(sidePos).(WoodFenceGate); ok {
				connected = gate.Facing.Face().Axis() != face.Axis()
			} else if _, ok := tx.Block(sidePos).Model().(model.Thin); ok {
				connected = true
			}
		}
		var connectionType WallConnectionType
		if connected {

			connectionType = ShortWallConnection()
			boxes := above.Model().BBox(abovePos, tx)
			for _, bb := range boxes {
				if bb.Min().Y() == 0 {
					xOverlap := bb.Min().X() < 0.75 && bb.Max().X() > 0.25
					zOverlap := bb.Min().Z() < 0.75 && bb.Max().Z() > 0.25
					var tall bool
					switch face {
					case cube.FaceNorth:
						tall = xOverlap && bb.Max().Z() > 0.75
					case cube.FaceEast:
						tall = bb.Min().X() < 0.25 && zOverlap
					case cube.FaceSouth:
						tall = xOverlap && bb.Min().Z() < 0.25
					case cube.FaceWest:
						tall = bb.Max().X() > 0.75 && zOverlap
					}
					if tall {
						connectionType = TallWallConnection()
						break
					}
				}
			}

		}
		if w.ConnectionType(face.Direction()) != connectionType {
			updated = true
			w = w.WithConnectionType(face.Direction(), connectionType)
		}
	}
	return w, updated
}

func (w Wall) calculatePost(tx *world.Tx, pos cube.Pos) (Wall, bool) {
	var updated bool
	abovePos := pos.Add(cube.Pos{0, 1, 0})
	above := tx.Block(abovePos)
	connections := len(sliceutil.Filter(cube.HorizontalFaces(), func(face cube.Face) bool {
		return w.ConnectionType(face.Direction()) != NoWallConnection()
	}))
	var post bool
	switch above := above.(type) {
	case Lantern:

		post = !above.Hanging
	case Sign:

		post = !above.Attach.hanging
	case Torch:

		post = above.Facing == cube.FaceDown
	case WoodTrapdoor:

		if above.Open {
			switch above.Facing {
			case cube.North:
				post = w.NorthConnection != NoWallConnection()
			case cube.East:
				post = w.EastConnection != NoWallConnection()
			case cube.South:
				post = w.SouthConnection != NoWallConnection()
			case cube.West:
				post = w.WestConnection != NoWallConnection()
			}
		}
	case Wall:

		post = above.Post
	}
	if !post {

		if connections < 2 {
			post = true
		} else {
			switch {
			case w.NorthConnection != NoWallConnection() && w.SouthConnection != NoWallConnection():
				post = w.EastConnection != NoWallConnection() || w.WestConnection != NoWallConnection()
			case w.EastConnection != NoWallConnection() && w.WestConnection != NoWallConnection():
				post = w.NorthConnection != NoWallConnection() || w.SouthConnection != NoWallConnection()
			default:
				post = true
			}
		}
	}
	if w.Post != post {
		updated = true
		w.Post = post
	}
	return w, updated
}

func allWalls() (walls []world.Block) {
	for _, block := range WallBlocks() {
		if _, hash := block.Hash(); hash > math.MaxUint16 {
			name, _ := block.EncodeBlock()
			panic(fmt.Errorf("hash of block %s exceeds 16 bytes", name))
		}
		for _, north := range WallConnectionTypes() {
			for _, east := range WallConnectionTypes() {
				for _, south := range WallConnectionTypes() {
					for _, west := range WallConnectionTypes() {
						walls = append(walls, Wall{
							Block:           block,
							NorthConnection: north,
							EastConnection:  east,
							SouthConnection: south,
							WestConnection:  west,
							Post:            false,
						})
						walls = append(walls, Wall{
							Block:           block,
							NorthConnection: north,
							EastConnection:  east,
							SouthConnection: south,
							WestConnection:  west,
							Post:            true,
						})
					}
				}
			}
		}
	}
	return
}
