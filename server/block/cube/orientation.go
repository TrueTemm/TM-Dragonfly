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

import "math"

type Orientation int

func OrientationFromYaw(yaw float64) Orientation {
	yaw = math.Mod(yaw, 360)
	return Orientation(math.Round(yaw / 360 * 16))
}

func (o Orientation) Yaw() float64 {
	return float64(o) / 16 * 360
}

func (o Orientation) Opposite() Orientation {
	return OrientationFromYaw(o.Yaw() + 180)
}

func (o Orientation) RotateLeft() Orientation {
	return OrientationFromYaw(o.Yaw() - 90)
}

func (o Orientation) RotateRight() Orientation {
	return OrientationFromYaw(o.Yaw() + 90)
}
