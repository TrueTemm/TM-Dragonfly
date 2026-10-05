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

package model

import (
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/mcrandom"
	"github.com/go-gl/mathgl/mgl64"
)

func randomOffset(pos cube.Pos, mn, mx float32, steps int) mgl64.Vec3 {
	seed := offsetSeed(int32(pos.X()), int32(pos.Z()))
	s0 := mcrandom.MixStafford13(seed)
	s1 := mcrandom.MixStafford13(seed + 0x9E3779B97F4A7C15)
	prng := mcrandom.NewXoroshiro128PlusPlus(s0, s1)
	x := offsetValue(mn, mx, steps, randomFloat32(prng.Next()))
	prng.Next()
	z := offsetValue(mn, mx, steps, randomFloat32(prng.Next()))
	return mgl64.Vec3{float64(x), 0, float64(z)}
}

func offsetSeed(x, z int32) uint64 {
	v := int64(z)*116129781 ^ int64(x)*0x2fc20f
	v = int64(int32(uint64(v*(v*42317861+11)) >> 16))
	return uint64(v) ^ 0x6A09E667F3BCC909
}

func randomFloat32(random uint64) float32 {
	return float32(random>>40) * (1.0 / 16777216.0)
}

func offsetValue(mn, mx float32, steps int, random float32) float32 {
	if mn >= mx {
		return mn
	}
	if steps == 1 {
		return (mn + mx) * 0.5
	}
	if steps > 1 {
		index := float32(math.Floor(float64(float32(steps) * random)))
		return mn + index*(mx-mn)/float32(steps-1)
	}
	return mn + (mx-mn)*random
}
