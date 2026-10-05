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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlaySound struct {
	packet.PlaySound
}

func (pk *PlaySound) Marshal(io protocol.IO) {
	io.String(&pk.SoundName)
	io.SoundPos(&pk.Position)
	io.Float32(&pk.Volume)
	io.Float32(&pk.Pitch)
	io.Varint32(&pk.LoopCount)
	protocol.OptionalFunc(io, &pk.Handle, io.Uint64)
}

func fromLatestPlaySound(pk *packet.PlaySound) *PlaySound {
	return &PlaySound{PlaySound: *pk}
}

func toLatestPlaySound(pk *PlaySound) *packet.PlaySound {
	out := pk.PlaySound
	return &out
}
