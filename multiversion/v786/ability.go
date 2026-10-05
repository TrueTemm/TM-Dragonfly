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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func AbilityData786(r protocol.IO, x *protocol.AbilityData) {
	r.Int64(&x.EntityUniqueID)
	r.Uint8(&x.PlayerPermissions)
	r.Uint8(&x.CommandPermissions)
	n := uint8(len(x.Layers))
	r.Uint8(&n)
	if len(x.Layers) != int(n) {
		x.Layers = make([]protocol.AbilityLayer, n)
	}
	for i := range x.Layers {
		abilityLayer786(r, &x.Layers[i])
	}
}

func abilityLayer786(r protocol.IO, x *protocol.AbilityLayer) {
	r.Uint16(&x.Type)
	r.Uint32(&x.Abilities)
	r.Uint32(&x.Values)
	r.Float32(&x.FlySpeed)
	if p := ProtoOf(r); p == 0 || p >= 776 {
		r.Float32(&x.VerticalFlySpeed)
	}
	r.Float32(&x.WalkSpeed)
}

func EntityLinks786(r protocol.IO, x *[]protocol.EntityLink) {
	protocol.FuncIOSlice(r, x, entityLink786)
}

func entityLink786(r protocol.IO, x *protocol.EntityLink) {
	r.Varint64(&x.RiddenEntityUniqueID)
	r.Varint64(&x.RiderEntityUniqueID)
	r.Uint8(&x.Type)
	r.Bool(&x.Immediate)
	r.Bool(&x.RiderInitiated)
	if p := ProtoOf(r); p == 0 || p >= 712 {
		r.Float32(&x.VehicleAngularVelocity)
	}
}
