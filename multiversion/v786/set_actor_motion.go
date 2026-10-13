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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SetActorMotion struct {
	EntityRuntimeID uint64
	Velocity        mgl32.Vec3
	Tick            uint64
}

func (*SetActorMotion) ID() uint32 { return IDSetActorMotion }
func (pk *SetActorMotion) Marshal(io protocol.IO) {
	io.ActorRuntimeID(&pk.EntityRuntimeID)
	io.Vec3(&pk.Velocity)
	if p := ProtoOf(io); p == 0 || p >= 662 { // added in 662
		io.Varuint64(&pk.Tick)
	}
}
func toLatestSetActorMotion(pk *SetActorMotion) *packet.SetActorMotion {
	return &packet.SetActorMotion{EntityRuntimeID: pk.EntityRuntimeID, Velocity: pk.Velocity, Tick: pk.Tick}
}
func fromLatestSetActorMotion(pk *packet.SetActorMotion) *SetActorMotion {
	return &SetActorMotion{EntityRuntimeID: pk.EntityRuntimeID, Velocity: pk.Velocity, Tick: pk.Tick}
}
