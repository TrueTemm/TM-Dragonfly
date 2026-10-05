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

package conformance

import (
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v800"
)

func TestPlayerListNeverPassedThrough(t *testing.T) {
	for _, proto := range protocols() {
		if !legacy(proto) {

			continue
		}
		out := proto.ConvertFromLatest(&packet.PlayerList{}, nil)
		if len(out) != 1 {
			t.Fatalf("protocol %v: PlayerList converted into %v packets, want 1", proto.ID(), len(out))
		}
		switch got := out[0].(type) {
		case *packet.PlayerList:
			t.Errorf("protocol %v: PlayerList passed through in the VENDORED wire — the packet-level "+
				"ActionType byte every protocol 786-898 expects is missing, so the client reads the entry "+
				"count as ActionType and the player never joins the list (see v800/player_list.go)", proto.ID())
		case *v786.PlayerList:

			if proto.ID() > 786 {
				t.Errorf("protocol %v: PlayerList converted to 786's type, but every protocol from 800 up "+
					"needs the trailing PlayerColour per entry", proto.ID())
			}
		case *v800.PlayerList:
			if proto.ID() == 786 {
				t.Errorf("protocol 786: PlayerList converted to the 800-era type, which writes a " +
					"PlayerColour 786 has no field for")
			}
		default:
			t.Errorf("protocol %v: PlayerList converted to unexpected type %T", proto.ID(), got)
		}
	}
}

func TestInteractPoolType(t *testing.T) {
	for _, proto := range protocols() {
		pk := proto.Packets(true)[v786.IDInteract]()
		_, latest := pk.(*packet.Interact)
		if want := proto.ID() >= 898; latest != want {
			t.Errorf("protocol %v: client pool decodes Interact as %T (latest=%v), want latest=%v — 898 is "+
				"where the conditional Vec3 became an Optional", proto.ID(), pk, latest, want)
		}
	}
}

func TestCommandDataFamilyAt898(t *testing.T) {
	for _, proto := range protocols() {
		for _, pk := range []packet.Packet{
			&packet.AvailableCommands{}, &packet.CommandOutput{}, &packet.CommandRequest{},
		} {
			out := proto.ConvertFromLatest(pk, nil)
			if len(out) != 1 {
				t.Fatalf("protocol %v: %T converted into %v packets, want 1", proto.ID(), pk, len(out))
			}
			latest := isLatest(out[0])
			if want := proto.ID() >= 898; latest != want {
				t.Errorf("protocol %v: %T is emitted as %T (latest=%v), want latest=%v — protocol/command.go's "+
					"nested types changed at 898", proto.ID(), pk, out[0], latest, want)
			}
		}
	}
}
