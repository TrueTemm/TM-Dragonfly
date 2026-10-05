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
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	ScoreboardActionModify = iota
	ScoreboardActionRemove
)

type ScoreboardEntry786 struct {
	EntryID        int64
	ObjectiveName  string
	Score          int32
	IdentityType   byte
	EntityUniqueID int64
	DisplayName    string
}

func (x *ScoreboardEntry786) Marshal(r protocol.IO) {
	r.Varint64(&x.EntryID)
	r.String(&x.ObjectiveName)
	r.Int32(&x.Score)
	r.Uint8(&x.IdentityType)
	switch x.IdentityType {
	case protocol.ScoreboardIdentityEntity, protocol.ScoreboardIdentityPlayer:
		r.Varint64(&x.EntityUniqueID)
	case protocol.ScoreboardIdentityFakePlayer:
		r.String(&x.DisplayName)
	default:
		r.UnknownEnumOption(x.IdentityType, "scoreboard entry identity type")
	}
}

func scoreRemoveEntry786(r protocol.IO, x *ScoreboardEntry786) {
	r.Varint64(&x.EntryID)
	r.String(&x.ObjectiveName)
	r.Int32(&x.Score)
}

type SetScore struct {
	ActionType byte

	Entries []ScoreboardEntry786
}

func (*SetScore) ID() uint32 {
	return IDSetScore
}

func (pk *SetScore) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	switch pk.ActionType {
	case ScoreboardActionRemove:
		protocol.FuncIOSlice(io, &pk.Entries, scoreRemoveEntry786)
	case ScoreboardActionModify:
		protocol.Slice(io, &pk.Entries)
	default:
		io.UnknownEnumOption(pk.ActionType, "set score action type")
	}
}

func fromLatestSetScore(pk *packet.SetScore) []packet.Packet {
	var removed, changed []ScoreboardEntry786
	for _, e := range pk.Entries {
		entry := ScoreboardEntry786{
			EntryID:        e.EntryID,
			ObjectiveName:  e.ObjectiveName,
			Score:          e.Score,
			IdentityType:   e.IdentityType,
			EntityUniqueID: e.EntityUniqueID,
			DisplayName:    e.DisplayName,
		}
		if e.IdentityType == protocol.ScoreboardIdentityRemove {
			removed = append(removed, entry)
			continue
		}
		changed = append(changed, entry)
	}
	var out []packet.Packet
	if len(removed) > 0 {
		out = append(out, &SetScore{ActionType: ScoreboardActionRemove, Entries: removed})
	}
	if len(changed) > 0 {
		out = append(out, &SetScore{ActionType: ScoreboardActionModify, Entries: changed})
	}
	return out
}

func toLatestSetScore(pk *SetScore) *packet.SetScore {
	out := &packet.SetScore{Entries: make([]protocol.ScoreboardEntry, 0, len(pk.Entries))}
	for _, e := range pk.Entries {
		entry := protocol.ScoreboardEntry{
			EntryID:        e.EntryID,
			ObjectiveName:  e.ObjectiveName,
			Score:          e.Score,
			IdentityType:   e.IdentityType,
			EntityUniqueID: e.EntityUniqueID,
			DisplayName:    e.DisplayName,
		}
		if pk.ActionType == ScoreboardActionRemove {
			entry.IdentityType = protocol.ScoreboardIdentityRemove
		}
		out.Entries = append(out.Entries, entry)
	}
	return out
}
