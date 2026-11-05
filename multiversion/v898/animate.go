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

package v898

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Animate struct{ packet.Animate }

func (pk *Animate) Marshal(io protocol.IO) {
	var swingSource protocol.Optional[string]
	if pk.SwingSource != 0 {
		swingSource = protocol.Option(swingSourceToString(pk.SwingSource))
	}
	io.Uint8(&pk.ActionType)
	io.ActorRuntimeID(&pk.EntityRuntimeID)
	io.Float32(&pk.Data)
	protocol.OptionalFunc(io, &swingSource, io.String)
	if val, ok := swingSource.Value(); ok {
		swingSourceFromString(io, &pk.SwingSource, val)
	}
}

func swingSourceFromString(io protocol.IO, x *uint8, s string) {
	switch s {
	case "none":
		*x = packet.AnimateSwingSourceNone
	case "build":
		*x = packet.AnimateSwingSourceBuild
	case "mine":
		*x = packet.AnimateSwingSourceMine
	case "interact":
		*x = packet.AnimateSwingSourceInteract
	case "attack":
		*x = packet.AnimateSwingSourceAttack
	case "useitem":
		*x = packet.AnimateSwingSourceUseItem
	case "throwitem":
		*x = packet.AnimateSwingSourceThrowItem
	case "dropitem":
		*x = packet.AnimateSwingSourceDropItem
	case "event":
		*x = packet.AnimateSwingSourceEvent
	default:
		io.InvalidValue(s, "swingSource", "unknown source")
	}
}

func swingSourceToString(x uint8) string {
	switch x {
	case packet.AnimateSwingSourceNone:
		return "none"
	case packet.AnimateSwingSourceBuild:
		return "build"
	case packet.AnimateSwingSourceMine:
		return "mine"
	case packet.AnimateSwingSourceInteract:
		return "interact"
	case packet.AnimateSwingSourceAttack:
		return "attack"
	case packet.AnimateSwingSourceUseItem:
		return "useitem"
	case packet.AnimateSwingSourceThrowItem:
		return "throwitem"
	case packet.AnimateSwingSourceDropItem:
		return "dropitem"
	case packet.AnimateSwingSourceEvent:
		return "event"
	default:
		return "unknown"
	}
}
