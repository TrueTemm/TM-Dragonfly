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

package portal

import (
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

type Nether struct {
	w, h      int
	framed    bool
	axis      cube.Axis
	tx        *world.Tx
	spawnPos  cube.Pos
	positions []cube.Pos
}

const (
	minimumNetherPortalWidth, maximumNetherPortalWidth   = 2, 21
	minimumNetherPortalHeight, maximumNetherPortalHeight = 3, 21

	minimumNetherPortalArea = minimumNetherPortalWidth * minimumNetherPortalHeight
)

func NetherPortalFromPos(tx *world.Tx, pos cube.Pos) (Nether, bool) {
	if tx.World().Dimension() == world.End {
		return Nether{}, false
	}

	axis, positions, width, height, completed, ok := multiAxisScan(pos, tx, matchesNetherPortalInterior)
	if !ok {
		axis, positions, width, height, completed, ok = multiAxisScan(pos, tx, matchesNetherPortal)
	}
	if !ok {
		return Nether{}, false
	}
	return Nether{
		w: width, h: height,
		spawnPos:  pos,
		positions: positions,
		framed:    completed,
		axis:      axis,
		tx:        tx,
	}, ok
}

var PortalsEnabled = false

func ActivateNetherPortal(tx *world.Tx, pos cube.Pos) bool {

	if !PortalsEnabled {
		return false
	}
	p, ok := NetherPortalFromPos(tx, pos)
	if !ok || !p.Framed() || p.Activated() {
		return false
	}
	ctx := tx.Event()
	positions := append([]cube.Pos(nil), p.Positions()...)
	if tx.World().Handler().HandlePortalActivate(ctx, world.Nether, positions); ctx.Cancelled() {
		return false
	}
	p.Activate()
	return true
}

func DeactivateNetherPortal(tx *world.Tx, pos cube.Pos) bool {
	_, positions, ok := connectedNetherPortal(tx, pos)
	if !ok {
		return false
	}
	deactivate(tx, positions)
	return true
}

func FindOrCreateNetherPortal(tx *world.Tx, pos cube.Pos, radius int) (Nether, bool) {
	n, ok := FindNetherPortal(tx, pos, radius)
	if ok {
		return n, true
	}
	return CreateNetherPortal(tx, pos)
}

type portalBlock interface {
	Portal() world.Dimension
}

type frameBlock interface {
	Frame(dimension world.Dimension) bool
}

func FindNetherPortal(tx *world.Tx, pos cube.Pos, radius int) (Nether, bool) {
	if tx.World().Dimension() == world.End {
		return Nether{}, false
	}

	closest, closestDist, found := Nether{}, math.MaxFloat64, false
	seen := make(map[cube.Pos]struct{})
	for selectedPos := range tx.BlocksWithin(pos, radius, portal(cube.X), portal(cube.Z)) {
		if _, ok := seen[selectedPos]; ok {

			continue
		}
		if n, ok := NetherPortalFromPos(tx, selectedPos); ok && n.Framed() && n.Activated() {
			for _, p := range n.Positions() {
				seen[p] = struct{}{}
				if dist := p.Vec3().Sub(pos.Vec3()).Len(); dist < closestDist {
					closestDist, closest, found = dist, n, true
				}
			}
		}
	}
	return closest, found
}

func CreateNetherPortal(tx *world.Tx, pos cube.Pos) (Nether, bool) {
	if tx.World().Dimension() == world.End {
		return Nether{}, false
	}

	resultPos, random, distance, a, r := pos, rand.IntN(4), -1.0, 0, tx.Range()
	searchValidArea := func(directions int, valid func(pos cube.Pos, riv int, coEff1, coEff2 int) bool) {
		for tempX := pos.X() - 16; tempX <= pos.X()+16; tempX++ {
			offsetX := float64(tempX-pos.X()) + 0.5
			for tempZ := pos.Z() - 16; tempZ <= pos.Z()+16; tempZ++ {
				offsetZ := float64(tempZ-pos.Z()) + 0.5
				for tempY := r.Max() - 1; tempY >= r.Min(); tempY-- {
					entryPos := cube.Pos{tempX, tempY, tempZ}
					if tx.Block(entryPos) != air() {
						continue
					}

					for tempY > r.Min() && tx.Block(entryPos.Side(cube.FaceDown)) == air() {
						tempY--
						entryPos[1]--
					}

					for riv := random; riv < random+directions; riv++ {
						coEff1 := riv % 2
						coEff2 := 1 - coEff1

						if !valid(entryPos, riv, coEff1, coEff2) {
							break
						}

						offsetY := float64(tempY-pos.Y()) + 0.5
						newDist := offsetX*offsetX + offsetY*offsetY + offsetZ*offsetZ
						if distance < 0.0 || newDist < distance {
							distance = newDist
							a = riv % directions
							resultPos = cube.Pos{tempX, tempY, tempZ}
						}
					}
				}
			}
		}
	}

	searchValidArea(4, func(pos cube.Pos, riv int, coEff1, coEff2 int) bool {
		if riv%4 >= 2 {
			coEff1 = -coEff1
			coEff2 = -coEff2
		}

		for safeSpace1 := range 3 {
			for safeSpace2 := -1; safeSpace2 < 3; safeSpace2++ {
				for height := -1; height < 4; height++ {
					b := tx.Block(cube.Pos{
						pos.X() + safeSpace2*coEff1 + safeSpace1*coEff2,
						pos.Y() + height,
						pos.Z() + safeSpace2*coEff2 - safeSpace1*coEff1,
					})
					_, solid := b.Model().(model.Solid)
					if height < 0 && !solid || height >= 0 && b != air() {
						return false
					}
				}
			}
		}
		return true
	})

	if distance < 0.0 {

		searchValidArea(2, func(pos cube.Pos, riv int, coEff1, coEff2 int) bool {
			for safeSpace := range 3 {
				for height := -1; height < 4; height++ {
					b := tx.Block(cube.Pos{
						pos.X() + safeSpace*coEff1,
						pos.Y() + height,
						pos.Z() + safeSpace*coEff2,
					})
					_, solid := b.Model().(model.Solid)
					if height < 0 && !solid || height >= 0 && b != air() {
						return false
					}
				}
			}
			return true
		})
	}

	coEff1 := a % 2
	coEff2 := 1 - coEff1
	if a%4 >= 2 {
		coEff1 = -coEff1
		coEff2 = -coEff2
	}

	axis := cube.X
	if coEff1 == 0 {
		axis = cube.Z
	}

	ob, pb := obsidian(), portal(axis)
	blocks := make(map[cube.Pos]world.Block)
	var affected []cube.Pos
	setBlock := func(pos cube.Pos, b world.Block) {
		if _, ok := blocks[pos]; !ok {
			affected = append(affected, pos)
		}
		blocks[pos] = b
	}

	if distance < 0.0 {

		resultPos[1] = min(max(resultPos[1], 70), r.Max()-10)
		for safeBeforeAfter := -1; safeBeforeAfter <= 1; safeBeforeAfter++ {
			for safeWidth := range 2 {
				for height := -1; height < 3; height++ {
					entryPos := cube.Pos{
						resultPos.X() + safeWidth*coEff1 + safeBeforeAfter*coEff2,
						resultPos.Y() + height,
						resultPos.Z() + safeWidth*coEff2 - safeBeforeAfter*coEff1,
					}

					if height < 0 {
						setBlock(entryPos, ob)
					} else {
						setBlock(entryPos, nil)
					}
				}
			}
		}
	}

	var positions []cube.Pos
	for width := -1; width < 3; width++ {
		for height := -1; height < 4; height++ {
			entryPos := cube.Pos{
				resultPos.X() + width*coEff1,
				resultPos.Y() + height,
				resultPos.Z() + width*coEff2,
			}

			if width == -1 || width == 2 || height == -1 || height == 3 {
				setBlock(entryPos, ob)
				continue
			}
			positions = append(positions, entryPos)
			setBlock(entryPos, pb)
		}
	}

	ctx := tx.Event()
	handlerPositions := append([]cube.Pos(nil), affected...)
	if tx.World().Handler().HandlePortalCreate(ctx, world.Nether, handlerPositions); ctx.Cancelled() {
		return Nether{}, false
	}
	for _, pos := range affected {
		tx.SetBlock(pos, blocks[pos], nil)
	}

	return Nether{
		w:         minimumNetherPortalWidth,
		h:         minimumNetherPortalHeight,
		framed:    true,
		spawnPos:  resultPos,
		positions: positions,
		axis:      axis,
		tx:        tx,
	}, true
}

func (n Nether) Activate() {
	for _, pos := range n.Positions() {
		n.tx.SetBlock(pos, portal(n.axis), nil)
	}
}

func deactivate(tx *world.Tx, positions []cube.Pos) {
	for _, pos := range positions {
		tx.SetBlock(pos, nil, nil)
	}
}

func (n Nether) Framed() bool {
	return n.framed
}

func (n Nether) Activated() bool {
	for _, pos := range n.Positions() {
		if n.tx.Block(pos) != portal(n.axis) {
			return false
		}
	}
	return true
}

func (n Nether) Spawn() cube.Pos {
	return n.spawnPos
}

func (n Nether) Positions() []cube.Pos {
	return n.positions
}
