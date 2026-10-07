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
	"reflect"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/df-mc/dragonfly/multiversion/native"
	"github.com/df-mc/dragonfly/multiversion/v1001"
	"github.com/df-mc/dragonfly/multiversion/v2168"
	"github.com/df-mc/dragonfly/multiversion/v2169"
	"github.com/df-mc/dragonfly/multiversion/v671"
	"github.com/df-mc/dragonfly/multiversion/v685"
	"github.com/df-mc/dragonfly/multiversion/v686"
	"github.com/df-mc/dragonfly/multiversion/v712"
	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v748"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v776"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v800"
	"github.com/df-mc/dragonfly/multiversion/v818"
	"github.com/df-mc/dragonfly/multiversion/v819"
	"github.com/df-mc/dragonfly/multiversion/v827"
	"github.com/df-mc/dragonfly/multiversion/v844"
	"github.com/df-mc/dragonfly/multiversion/v859"
	"github.com/df-mc/dragonfly/multiversion/v898"
	"github.com/df-mc/dragonfly/multiversion/v924"
	"github.com/df-mc/dragonfly/multiversion/v944"
	"github.com/df-mc/dragonfly/multiversion/v975"
)

func protocols() []minecraft.Protocol {
	return []minecraft.Protocol{
		v671.Protocol{}, v685.Protocol{}, v686.Protocol{}, v712.Protocol{}, v729.Protocol{}, v748.Protocol{}, v766.Protocol{}, v776.Protocol{}, v786.Protocol{}, v800.Protocol{}, v818.Protocol{}, v819.Protocol{},
		v827.Protocol{}, v844.Protocol{}, v859.Protocol{}, v898.Protocol{}, v924.Protocol{}, v944.Protocol{}, v975.Protocol{}, v1001.Protocol{}, v2168.Protocol{}, v2169.Protocol{}, native.Protocol{},
	}
}

func legacy(proto minecraft.Protocol) bool { return proto.ID() < 900 }

func TestStartGameUsesBlockNetworkIDHashes(t *testing.T) {
	for _, proto := range protocols() {
		out := proto.ConvertFromLatest(&packet.StartGame{}, nil)
		if len(out) != 1 {
			t.Fatalf("protocol %v: StartGame converted into %v packets, want exactly 1", proto.ID(), len(out))
		}
		remapped := blockpalette.Supported(uint32(proto.ID()))
		field := reflect.ValueOf(out[0]).Elem().FieldByName("UseBlockNetworkIDHashes")
		if !field.IsValid() {
			t.Fatalf("protocol %v: converted StartGame (%T) has no UseBlockNetworkIDHashes field", proto.ID(), out[0])
		}
		if field.Bool() != remapped {
			t.Errorf("protocol %v: StartGame.UseBlockNetworkIDHashes = %v but blockpalette.Supported = %v — the "+
				"flag and the chunk palette must agree; a mismatch renders the world invisible and uncollidable "+
				"(see this test's doc comment)", proto.ID(), field.Bool(), remapped)
		}
	}
}

func TestStartGameReportsOwnVersion(t *testing.T) {
	seen := make(map[string]int32)
	for _, proto := range protocols() {
		out := proto.ConvertFromLatest(&packet.StartGame{}, nil)
		if len(out) != 1 {
			t.Fatalf("protocol %v: StartGame converted into %v packets, want exactly 1", proto.ID(), len(out))
		}
		if _, native := out[0].(*packet.StartGame); native {
			if legacy(proto) {
				t.Fatalf("protocol %v: a 1.21.x package must convert StartGame, but it passed through", proto.ID())
			}

			continue
		}
		v := reflect.ValueOf(out[0]).Elem()
		base, game := v.FieldByName("BaseGameVersion"), v.FieldByName("GameVersion")
		if !base.IsValid() || !game.IsValid() {
			t.Fatalf("protocol %v: converted StartGame (%T) is missing a version field", proto.ID(), out[0])
		}
		if base.String() == "" || base.String() != game.String() {
			t.Errorf("protocol %v: BaseGameVersion=%q GameVersion=%q, want both set and equal",
				proto.ID(), base.String(), game.String())
			continue
		}

		if other, dup := seen[base.String()]; dup && !(other == 818 && proto.ID() == 819) {
			t.Errorf("protocol %v reports version %q, already claimed by protocol %v — a transcribed sibling "+
				"package almost certainly kept its ancestor's version string", proto.ID(), base.String(), other)
			continue
		}
		seen[base.String()] = proto.ID()
	}
}
