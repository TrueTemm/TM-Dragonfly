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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	endSpawnX = 100
	endSpawnY = 49
	endSpawnZ = 0
)

func EndSpawnPosition(player bool) mgl64.Vec3 {
	y := endSpawnY + 1
	if player {
		y = endSpawnY
	}
	return mgl64.Vec3{float64(endSpawnX) + 0.5, float64(y), float64(endSpawnZ) + 0.5}
}

func GenerateEndSpawnPlatform(tx *world.Tx) {
	ob := obsidian()
	for dx := -2; dx <= 2; dx++ {
		for dz := -2; dz <= 2; dz++ {
			tx.SetBlock(cube.Pos{endSpawnX + dx, endSpawnY - 1, endSpawnZ + dz}, ob, nil)
			for dy := 0; dy < 3; dy++ {
				tx.SetBlock(cube.Pos{endSpawnX + dx, endSpawnY + dy, endSpawnZ + dz}, nil, nil)
			}
		}
	}
}

type endRingFrame struct {
	pos    cube.Pos
	facing cube.Direction
}

type endFrameBlock interface {
	world.Block
	EndPortalFrameState() (eye bool, facing cube.Direction)
}

func ActivateEndPortal(tx *world.Tx, framePos cube.Pos) bool {

	if !PortalsEnabled {
		return false
	}
	f, ok := tx.Block(framePos).(endFrameBlock)
	if !ok {
		return false
	}
	_, facing := f.EndPortalFrameState()

	inward, tangent := facing.Face(), facing.RotateRight().Face()
	base := framePos.Side(inward).Side(inward)
	for _, center := range []cube.Pos{base.Side(tangent.Opposite()), base, base.Side(tangent)} {
		interior, ok := matchEndRing(tx, center)
		if !ok {
			continue
		}
		ep := endPortal()
		active := true
		for _, pos := range interior {
			if tx.Block(pos) != ep {
				active = false
				break
			}
		}
		if active {
			return false
		}
		ctx := tx.Event()
		positions := append([]cube.Pos(nil), interior...)
		if tx.World().Handler().HandlePortalActivate(ctx, world.End, positions); ctx.Cancelled() {
			return false
		}
		for _, pos := range interior {
			if tx.Block(pos) != ep {
				tx.SetBlock(pos, ep, nil)
			}
		}
		return true
	}
	return false
}

func EndPortalRingIntact(tx *world.Tx, portalPos cube.Pos) bool {
	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			if _, ok := matchEndRing(tx, portalPos.Add(cube.Pos{dx, 0, dz})); ok {
				return true
			}
		}
	}
	return false
}

func DeactivateEndPortal(tx *world.Tx, portalPos cube.Pos) {
	ep := endPortal()
	if tx.Block(portalPos) != ep {
		return
	}
	var positions []cube.Pos
	queue := []cube.Pos{portalPos}
	seen := map[cube.Pos]struct{}{portalPos: {}}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if tx.Block(p) != ep {
			continue
		}
		positions = append(positions, p)
		for _, face := range cube.HorizontalFaces() {
			n := p.Side(face)
			if _, ok := seen[n]; ok {
				continue
			}
			seen[n] = struct{}{}
			queue = append(queue, n)
		}
	}
	deactivate(tx, positions)
}

func matchEndRing(tx *world.Tx, center cube.Pos) ([]cube.Pos, bool) {
	for _, want := range expectedEndRingFrames(center) {
		b, ok := tx.Block(want.pos).(endFrameBlock)
		if !ok {
			return nil, false
		}
		eye, facing := b.EndPortalFrameState()
		if !eye || facing != want.facing {
			return nil, false
		}
	}
	return endRingInterior(center), true
}

func expectedEndRingFrames(center cube.Pos) []endRingFrame {
	frames := make([]endRingFrame, 0, 12)
	for _, side := range cube.Directions() {
		base := center.Side(side.Face()).Side(side.Face())
		tangent := side.RotateRight().Face()
		inward := side.Opposite()
		for _, pos := range []cube.Pos{base.Side(tangent.Opposite()), base, base.Side(tangent)} {
			frames = append(frames, endRingFrame{pos: pos, facing: inward})
		}
	}
	return frames
}

func endRingInterior(center cube.Pos) []cube.Pos {
	out := make([]cube.Pos, 0, 9)
	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			out = append(out, center.Add(cube.Pos{dx, 0, dz}))
		}
	}
	return out
}

func endPortal() world.Block {
	p, ok := world.BlockByName("minecraft:end_portal", nil)
	if !ok {
		panic("could not find end_portal block")
	}
	return p
}
