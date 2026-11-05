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
	"fmt"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestHistoricalSoundRegistryCoverage(t *testing.T) {
	counts := map[uint32]int{649: 450, 662: 450, 671: 450, 685: 487, 686: 487, 712: 487, 729: 496, 748: 496, 766: 496, 776: 496, 786: 496, 800: 496, 818: 523, 819: 523, 827: 525, 844: 528, 859: 528, 898: 540, 924: 559, 944: 561, 975: 563}
	if len(soundRegistries) != len(counts) {
		t.Fatal("sound registry version coverage changed")
	}
	for proto, count := range counts {
		t.Run(fmt.Sprint(proto), func(t *testing.T) {
			r := soundRegistries[proto]
			if len(r.IDs) != count || len(r.Revision) != 40 {
				t.Fatal("historical registry or source revision missing")
			}
			if id, ok := SoundIDFor(proto, packet.SoundEventRecordRelic); !ok || id != 469 {
				t.Fatal("reserved sound slots were lost")
			}
			for name, id := range r.IDs {
				source := &packet.LevelSoundEvent{SoundType: name, ExtraData: -1, EntityUniqueID: -12345}
				encoded := FromLatestLevelSoundEvent786(proto, source)
				if encoded == nil || encoded.SoundType != id {
					t.Fatalf("sound %s did not map to %d", name, id)
				}
				decoded := ToLatestLevelSoundEvent786(proto, encoded)
				if decoded == nil || decoded.SoundType != name || decoded.ExtraData != -1 || decoded.EntityUniqueID != source.EntityUniqueID {
					t.Fatalf("sound %s changed", name)
				}
			}
			if out := LevelSoundPacketsFromLatest(proto, &packet.LevelSoundEvent{SoundType: "tm:unknown"}); len(out) != 0 {
				t.Fatal("unknown sound emitted")
			}
			if out := LevelSoundPacketsToLatest(proto, &LevelSoundEvent{SoundType: 0xffffffff}); len(out) != 0 {
				t.Fatal("unknown sound accepted")
			}
			t.Logf("checked %d historical sound IDs", count)
		})
	}
}
