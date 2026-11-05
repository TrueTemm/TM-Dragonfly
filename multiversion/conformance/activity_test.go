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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/df-mc/dragonfly/multiversion/v898"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestHistoricalIdleInput898(t *testing.T) {
	p := v898.Protocol{}
	raw := activityFixtures(t, 898)["auth_idle"]
	buf := bytes.NewBuffer(raw)
	input := p.Packets(true)[packet.IDPlayerAuthInput]()
	input.Marshal(p.NewReader(buf, 0, true))
	if buf.Len() != 0 {
		t.Fatalf("%d movement bytes unread", buf.Len())
	}
	out := p.ConvertToLatest(input, nil)[0].(*packet.PlayerAuthInput)
	if out.Tick != 100 || out.Position != (mgl32.Vec3{-16.5, 65.25, 3}) || !out.InputData.Load(packet.InputFlagJumpDown) {
		t.Fatal("idle input or movement flag changed")
	}
	encoded := new(bytes.Buffer)
	input.Marshal(p.NewWriter(encoded, 0))
	if !bytes.Equal(encoded.Bytes(), raw) {
		t.Fatal("idle input differs from historical wire")
	}
}

func activityFixtures(t *testing.T, version int) map[string][]byte {
	t.Helper()
	data, err := os.ReadFile(fmt.Sprintf("testdata/activity%d.json", version))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string][]byte
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	return fixtures
}

func TestHistoricalSoundActivity(t *testing.T) {
	finaliseOnce.Do(world.DefaultBlockRegistry.Finalize)
	old, modern := activityFixtures(t, 630), activityFixtures(t, 898)
	for _, p := range protocols() {
		if p.ID() > 975 {
			continue
		}
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			for name, kind := range map[string]string{"attack": packet.SoundEventAttackStrong, "attack_no_damage": packet.SoundEventAttackNoDamage, "breathe": packet.SoundEventBreathe, "step": packet.SoundEventStep, "hit": packet.SoundEventHit, "record_relic": packet.SoundEventRecordRelic} {
				t.Run(name, func(t *testing.T) {
					source := &packet.LevelSoundEvent{SoundType: kind, Position: mgl32.Vec3{-16.5, 65.25, 3}, ExtraData: -1, EntityType: "minecraft:player", BabyMob: true, DisableRelativeVolume: true, EntityUniqueID: -12345}
					out := p.ConvertFromLatest(source, nil)
					if len(out) != 1 {
						t.Fatalf("sound converted into %d packets", len(out))
					}
					buf := new(bytes.Buffer)
					out[0].Marshal(p.NewWriter(buf, 0))
					want := old["sound_"+name]
					if p.ID() >= 786 {
						want = modern["sound_"+name]
					}
					if p.ID() >= 975 {
						want = append(append([]byte(nil), want...), 0)
					}
					if !bytes.Equal(buf.Bytes(), want) {
						t.Fatalf("wire differs from historical codec: got %x want %x", buf.Bytes(), want)
					}
					decoded := p.Packets(true)[packet.IDLevelSoundEvent]()
					decoded.Marshal(p.NewReader(buf, 0, true))
					if buf.Len() != 0 {
						t.Fatalf("%d unread bytes", buf.Len())
					}
					back := p.ConvertToLatest(decoded, nil)
					if len(back) != 1 {
						t.Fatal("sound not converted to native")
					}
					got := back[0].(*packet.LevelSoundEvent)
					wantUID := int64(0)
					if p.ID() >= 786 {
						wantUID = source.EntityUniqueID
					}
					if got.SoundType != kind || got.ExtraData != -1 || got.Position != source.Position || got.EntityUniqueID != wantUID || !got.BabyMob || !got.DisableRelativeVolume {
						t.Fatalf("sound changed: %+v", got)
					}
					if source.ExtraData != -1 {
						t.Fatal("source sound changed")
					}
				})
			}
			if out := p.ConvertFromLatest(&packet.LevelSoundEvent{SoundType: "tm:unavailable"}, nil); len(out) != 0 {
				t.Fatal("unknown sound must not become item.use.on")
			}
			unknown := p.Packets(true)[packet.IDLevelSoundEvent]()
			reflect.ValueOf(unknown).Elem().FieldByName("SoundType").SetUint(0xffffffff)
			if out := p.ConvertToLatest(unknown, nil); len(out) != 0 {
				t.Fatal("unknown sound ID must be dropped")
			}
		})
	}
}

func TestHistoricalSwingActivity(t *testing.T) {
	old, modern := activityFixtures(t, 630), activityFixtures(t, 898)
	for _, p := range protocols() {
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			for _, name := range []string{"swing", "attack"} {
				source := &packet.Animate{ActionType: 1, EntityRuntimeID: 27, Data: 0.5, Hand: 1}
				if name == "attack" {
					source.SwingSource = packet.AnimateSwingSourceAttack
				}
				out := p.ConvertFromLatest(source, nil)
				if len(out) != 1 {
					t.Fatalf("animation converted into %d packets", len(out))
				}
				buf := new(bytes.Buffer)
				out[0].Marshal(p.NewWriter(buf, 0))
				want := old["animate_"+name]
				if p.ID() == 859 {
					want = append([]byte{2}, modern["animate_swing"][1:len(modern["animate_swing"])-1]...)
				}
				if p.ID() >= 898 {
					want = modern["animate_"+name]
				}
				if p.ID() >= 2225 {
					want = append(append([]byte(nil), want...), 1)
				}
				if !bytes.Equal(buf.Bytes(), want) {
					t.Fatalf("%s wire differs: got %x want %x", name, buf.Bytes(), want)
				}
				input := p.Packets(true)[packet.IDAnimate]()
				input.Marshal(p.NewReader(buf, 0, true))
				back := p.ConvertToLatest(input, nil)
				if buf.Len() != 0 || len(back) != 1 {
					t.Fatal("animation was not fully decoded")
				}
				got := back[0].(*packet.Animate)
				if got.ActionType != 1 || got.EntityRuntimeID != 27 || source.Hand != 1 {
					t.Fatal("animation identity or source changed")
				}
				if p.ID() >= 898 && got.SwingSource != source.SwingSource {
					t.Fatal("swing source lost")
				}
				if p.ID() < 2225 && got.Hand != 0 {
					t.Fatal("legacy packet has a hand field")
				}
			}
		})
	}
}
