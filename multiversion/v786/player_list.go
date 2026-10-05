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
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	PlayerListActionAdd = iota
	PlayerListActionRemove
)

type PlayerListEntry786 struct {
	UUID           uuid.UUID
	EntityUniqueID int64
	Username       string
	XUID           string
	PlatformChatID string
	BuildPlatform  int32
	Skin           Skin786
	Teacher        bool
	Host           bool
	SubClient      bool
}

func (x *PlayerListEntry786) Marshal(io protocol.IO) {
	io.UUID(&x.UUID)
	io.Varint64(&x.EntityUniqueID)
	io.String(&x.Username)
	io.String(&x.XUID)
	io.String(&x.PlatformChatID)
	io.Int32(&x.BuildPlatform)
	protocol.Single(io, &x.Skin)
	io.Bool(&x.Teacher)
	io.Bool(&x.Host)
	io.Bool(&x.SubClient)
}

type PlayerList struct {
	ActionType byte

	Entries []PlayerListEntry786
}

func (*PlayerList) ID() uint32 {
	return IDPlayerList
}

func (pk *PlayerList) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	switch pk.ActionType {
	case PlayerListActionAdd:
		protocol.Slice(io, &pk.Entries)
	case PlayerListActionRemove:
		protocol.FuncIOSlice(io, &pk.Entries, PlayerListRemoveEntry786)
	default:
		io.UnknownEnumOption(pk.ActionType, "player list action type")
	}
	if pk.ActionType == PlayerListActionAdd {
		for i := 0; i < len(pk.Entries); i++ {
			io.Bool(&pk.Entries[i].Skin.Trusted)
		}
	}
}
