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

package cube

import (
	"fmt"
	"iter"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Pos [3]int

func (p Pos) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p[0], p[1], p[2])
}

func (p Pos) X() int {
	return p[0]
}

func (p Pos) Y() int {
	return p[1]
}

func (p Pos) Z() int {
	return p[2]
}

func (p Pos) OutOfBounds(r Range) bool {
	y := p[1]
	return y > r[1] || y < r[0]
}

func (p Pos) Within(min, max Pos) bool {
	return p[0] >= min[0] && p[0] <= max[0] &&
		p[1] >= min[1] && p[1] <= max[1] &&
		p[2] >= min[2] && p[2] <= max[2]
}

func (p Pos) Add(pos Pos) Pos {
	return Pos{p[0] + pos[0], p[1] + pos[1], p[2] + pos[2]}
}

func (p Pos) Sub(pos Pos) Pos {
	return Pos{p[0] - pos[0], p[1] - pos[1], p[2] - pos[2]}
}

func (p Pos) Vec3() mgl64.Vec3 {
	return mgl64.Vec3{float64(p[0]), float64(p[1]), float64(p[2])}
}

func (p Pos) Vec3Middle() mgl64.Vec3 {
	return mgl64.Vec3{float64(p[0]) + 0.5, float64(p[1]), float64(p[2]) + 0.5}
}

func (p Pos) Vec3Centre() mgl64.Vec3 {
	return mgl64.Vec3{float64(p[0]) + 0.5, float64(p[1]) + 0.5, float64(p[2]) + 0.5}
}

func (p Pos) Side(face Face) Pos {
	switch face {
	case FaceUp:
		p[1]++
	case FaceDown:
		p[1]--
	case FaceNorth:
		p[2]--
	case FaceSouth:
		p[2]++
	case FaceWest:
		p[0]--
	case FaceEast:
		p[0]++
	}
	return p
}

func (p Pos) Face(other Pos) Face {
	face, _ := p.NeighbourFace(other)
	return face
}

func (p Pos) NeighbourFace(other Pos) (Face, bool) {
	switch other.Sub(p) {
	case Pos{0, 1, 0}:
		return FaceUp, true
	case Pos{0, -1, 0}:
		return FaceDown, true
	case Pos{0, 0, -1}:
		return FaceNorth, true
	case Pos{0, 0, 1}:
		return FaceSouth, true
	case Pos{-1, 0, 0}:
		return FaceWest, true
	case Pos{1, 0, 0}:
		return FaceEast, true
	}
	return FaceUp, false
}

func (p Pos) Neighbours(f func(neighbour Pos), r Range) {
	if p.OutOfBounds(r) {
		return
	}
	p[0]++
	f(p)
	p[0] -= 2
	f(p)
	p[0]++
	p[1]++
	if p[1] <= r[1] {
		f(p)
	}
	p[1] -= 2
	if p[1] >= r[0] {
		f(p)
	}
	p[1]++
	p[2]++
	f(p)
	p[2] -= 2
	f(p)
}

func PosFromVec3(vec3 mgl64.Vec3) Pos {
	return Pos{int(math.Floor(vec3[0])), int(math.Floor(vec3[1])), int(math.Floor(vec3[2]))}
}

func Min(p1, p2 Pos) Pos {
	return Pos{min(p1[0], p2[0]), min(p1[1], p2[1]), min(p1[2], p2[2])}
}

func Max(p1, p2 Pos) Pos {
	return Pos{max(p1[0], p2[0]), max(p1[1], p2[1]), max(p1[2], p2[2])}
}

func Range3D(p1, p2 Pos) iter.Seq[Pos] {
	max := Max(p1, p2)
	min := Min(p1, p2)
	return func(yield func(Pos) bool) {
		for x := min[0]; x <= max[0]; x++ {
			for y := min[1]; y <= max[1]; y++ {
				for z := min[2]; z <= max[2]; z++ {
					if !yield(Pos{x, y, z}) {
						return
					}
				}
			}
		}
	}
}
