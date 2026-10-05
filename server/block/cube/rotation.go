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
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type Rotation [2]float64

func (r Rotation) Yaw() float64 {
	return r[0]
}

func (r Rotation) Pitch() float64 {
	return r[1]
}

func (r Rotation) Elem() (yaw, pitch float64) {
	return r[0], r[1]
}

func (r Rotation) Add(r2 Rotation) Rotation {
	return Rotation{r[0] + r2[0], r[1] + r2[1]}.fix()
}

func (r Rotation) Opposite() Rotation {
	fixed := r.fix()
	return Rotation{fixed[0] + 180, -fixed[1]}.fix()
}

func (r Rotation) Neg() Rotation {
	fixed := r.fix()
	return Rotation{-fixed[0], -fixed[1]}
}

func (r Rotation) Direction() Direction {
	yaw := r.fix().Yaw()
	switch {
	case yaw > 45 && yaw <= 135:
		return West
	case yaw > -45 && yaw <= 45:
		return South
	case yaw > -135 && yaw <= -45:
		return East
	case yaw <= -135 || yaw > 135:
		return North
	}
	return 0
}

func (r Rotation) Orientation() Orientation {
	const step = 360 / 16.0

	yaw := r.fix().Yaw()
	if yaw < -step/2 {
		yaw += 360
	}
	return Orientation(math.Round(yaw / step))
}

func (r Rotation) Vec3() mgl64.Vec3 {
	yaw, pitch := r.fix().Elem()
	yawRad, pitchRad := mgl64.DegToRad(yaw), mgl64.DegToRad(pitch)

	m := math.Cos(pitchRad)
	return mgl64.Vec3{
		-m * math.Sin(yawRad),
		-math.Sin(pitchRad),
		m * math.Cos(yawRad),
	}
}

func (r Rotation) fix() Rotation {
	signYaw := math.Copysign(180, r[0])

	return Rotation{
		math.Mod(r[0]+signYaw, 360) - signYaw,
		max(min(r[1], 90), -90),
	}
}
