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

package v800

import (
	"image/color"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

func ARGB800(io protocol.IO, x *color.RGBA) {
	val := int32(x.A) | int32(x.R)<<8 | int32(x.G)<<16 | int32(x.B)<<24
	io.Int32(&val)
	*x = color.RGBA{A: byte(val), R: byte(val >> 8), G: byte(val >> 16), B: byte(val >> 24)}
}

type PlayerListEntry800 struct {
	v786.PlayerListEntry786

	PlayerColour color.RGBA
}

func (x *PlayerListEntry800) Marshal(io protocol.IO) {
	x.PlayerListEntry786.Marshal(io)
	ARGB800(io, &x.PlayerColour)
}

func PlayerListRemoveEntry800(io protocol.IO, x *PlayerListEntry800) {
	io.UUID(&x.UUID)
}

type PlayerList struct {
	ActionType byte

	Entries []PlayerListEntry800
}

func (*PlayerList) ID() uint32 { return v786.IDPlayerList }

func (pk *PlayerList) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	switch pk.ActionType {
	case v786.PlayerListActionAdd:
		protocol.Slice(io, &pk.Entries)
	case v786.PlayerListActionRemove:
		protocol.FuncIOSlice(io, &pk.Entries, PlayerListRemoveEntry800)
	default:
		io.UnknownEnumOption(pk.ActionType, "player list action type")
	}
	if pk.ActionType == v786.PlayerListActionAdd {
		for i := 0; i < len(pk.Entries); i++ {
			io.Bool(&pk.Entries[i].Skin.Trusted)
		}
	}
}

func FromLatestPlayerList800(pk *packet.PlayerList) *PlayerList {
	old := v786.FromLatestPlayerList(pk)
	entries := make([]PlayerListEntry800, len(old.Entries))
	for i := range old.Entries {
		entries[i] = PlayerListEntry800{PlayerListEntry786: old.Entries[i]}
		if i < len(pk.Entries) {
			entries[i].PlayerColour = pk.Entries[i].PlayerColour
		}
	}
	return &PlayerList{ActionType: old.ActionType, Entries: entries}
}

func ToLatestPlayerList800(pk *PlayerList) *packet.PlayerList {
	old := &v786.PlayerList{ActionType: pk.ActionType, Entries: make([]v786.PlayerListEntry786, len(pk.Entries))}
	for i := range pk.Entries {
		old.Entries[i] = pk.Entries[i].PlayerListEntry786
	}
	out := v786.ToLatestPlayerList(old)
	for i := range out.Entries {
		out.Entries[i].PlayerColour = pk.Entries[i].PlayerColour
	}
	return out
}
