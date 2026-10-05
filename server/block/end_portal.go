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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/portal"
)

type EndPortal struct {
	transparent
}

func (EndPortal) Model() world.BlockModel {
	return model.Empty{}
}

func (EndPortal) LightEmissionLevel() uint8 {
	return 15
}

func (EndPortal) HasLiquidDrops() bool {
	return false
}

func (EndPortal) Portal() world.Dimension {
	return world.End
}

func (EndPortal) EncodeNBT() map[string]any {
	return map[string]any{"id": "EndPortal"}
}

func (e EndPortal) DecodeNBT(map[string]any) any {
	return e
}

func (EndPortal) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if portal.EndPortalRingIntact(tx, pos) {
		return
	}
	portal.DeactivateEndPortal(tx, pos)
}

func (EndPortal) EntityInside(_ cube.Pos, tx *world.Tx, e world.Entity) {
	if t, ok := e.(portalTraveller); ok {
		t.TravelThroughPortal(tx, world.End)
	}
}

func (EndPortal) EncodeBlock() (string, map[string]any) {
	return "minecraft:end_portal", nil
}
