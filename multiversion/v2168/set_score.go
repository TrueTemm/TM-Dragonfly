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

package v2168

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v2169"
)

type ScoreboardEntry struct {
	EntryID        int64
	ObjectiveName  string
	Score          int32
	IdentityType   byte
	EntityUniqueID int64
	DisplayName    string
}

func (x *ScoreboardEntry) Marshal(r protocol.IO) {
	variant := uint32(x.IdentityType)
	r.Varuint32(&variant)
	x.IdentityType = byte(variant)

	typeNames := [...]string{"remove", "changeplayer", "changeentity", "changefakeplayer"}
	if variant >= uint32(len(typeNames)) {
		r.UnknownEnumOption(variant, "scoreboard entry variant")
		return
	}
	typeName := typeNames[variant]
	r.String(&typeName)
	r.Varint64(&x.EntryID)
	switch x.IdentityType {
	case protocol.ScoreboardIdentityRemove:
		objective := protocol.Optional[string]{}
		if x.ObjectiveName != "" {
			objective = protocol.Option(x.ObjectiveName)
		}

		v2169.DoubleOptionalFunc(r, &objective, r.String)
		x.ObjectiveName, _ = objective.Value()
	case protocol.ScoreboardIdentityEntity, protocol.ScoreboardIdentityPlayer:
		r.String(&x.ObjectiveName)
		r.Int32(&x.Score)
		r.ActorUniqueID(&x.EntityUniqueID)
	case protocol.ScoreboardIdentityFakePlayer:
		r.String(&x.ObjectiveName)
		r.Int32(&x.Score)
		r.String(&x.DisplayName)
	}
}

type SetScore struct {
	Entries []ScoreboardEntry
}

func (*SetScore) ID() uint32 { return packet.IDSetScore }

func (pk *SetScore) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Entries)
}

func fromLatestSetScore(pk *packet.SetScore) *SetScore {
	out := &SetScore{Entries: make([]ScoreboardEntry, len(pk.Entries))}
	for i, e := range pk.Entries {
		out.Entries[i] = ScoreboardEntry{
			EntryID: e.EntryID, ObjectiveName: e.ObjectiveName, Score: e.Score,
			IdentityType: e.IdentityType, EntityUniqueID: e.EntityUniqueID, DisplayName: e.DisplayName,
		}
	}
	return out
}

func toLatestSetScore(pk *SetScore) *packet.SetScore {
	out := &packet.SetScore{Entries: make([]protocol.ScoreboardEntry, len(pk.Entries))}
	for i, e := range pk.Entries {
		out.Entries[i] = protocol.ScoreboardEntry{
			EntryID: e.EntryID, ObjectiveName: e.ObjectiveName, Score: e.Score,
			IdentityType: e.IdentityType, EntityUniqueID: e.EntityUniqueID, DisplayName: e.DisplayName,
		}
	}
	return out
}
