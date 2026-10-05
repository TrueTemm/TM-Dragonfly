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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)

const defaultExplosionSize = 4

type ExplosionSource interface {
	Position() mgl64.Vec3

	Size() float64
}

type EntityExplosionSource struct {
	Entity Entity

	ExplosionSize float64
}

func (e EntityExplosionSource) Position() mgl64.Vec3 {
	return e.Entity.Position()
}

func (e EntityExplosionSource) Size() float64 {
	if e.ExplosionSize == 0 {
		return defaultExplosionSize
	}
	return e.ExplosionSize
}

type BlockExplosionSource struct {
	Block Block

	Pos cube.Pos

	ExplosionSize float64
}

func (b BlockExplosionSource) Position() mgl64.Vec3 {
	return b.Pos.Vec3Centre()
}

func (b BlockExplosionSource) Size() float64 {
	if b.ExplosionSize == 0 {
		return defaultExplosionSize
	}
	return b.ExplosionSize
}
