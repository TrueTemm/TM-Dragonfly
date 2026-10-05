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

package cmd

import (
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"slices"
)

type Target interface {
	Position() mgl64.Vec3
}

type NamedTarget interface {
	Target

	Name() string
}

func targets(tx *world.Tx) (entities []Target, players []NamedTarget) {
	if tx == nil {
		return nil, nil
	}
	ent := sliceutil.Convert[Target](slices.Collect(tx.Entities()))
	pl := sliceutil.Convert[NamedTarget](slices.Collect(tx.Players()))
	return ent, pl
}
