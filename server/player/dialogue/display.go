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

package dialogue

import (
	"encoding/json"
	"github.com/go-gl/mathgl/mgl64"
)

type DisplaySettings struct {
	EntityScale mgl64.Vec3

	EntityOffset mgl64.Vec3

	EntityRotation mgl64.Vec3
}

func (d DisplaySettings) MarshalJSON() ([]byte, error) {

	d.EntityRotation[0], d.EntityRotation[1] = d.EntityRotation[1], d.EntityRotation[0]-32
	m := map[string]any{

		"translate": d.EntityOffset.Mul(-32),

		"rotate": d.EntityRotation,
		"scale":  [3]float64{1, 1, 1},
	}
	if (d.EntityScale != mgl64.Vec3{}) {
		m["scale"] = d.EntityScale
	}
	return json.Marshal(m)
}
