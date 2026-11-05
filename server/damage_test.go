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

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"os"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type damageTestHandler struct{ player.NopHandler }

func (damageTestHandler) HandleHurt(_ *player.Context, _ *float64, _ bool, immunity *time.Duration, _ world.DamageSource) {
	*immunity = 0
}

func newDamageTestServer(t *testing.T) (*Server, *minecraft.Listener, <-chan uuid.UUID) {
	t.Helper()
	cfg := (UserConfig{}).minecraftListenConfig(Config{Compression: packet.SnappyCompression, AuthDisabled: true, Allower: allower{}})
	l, err := cfg.ListenNetwork(legacyRakNet{Network: rakNetNetwork(nil)}, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := (Config{
		Log:                     log,
		AuthDisabled:            true,
		DisableResourceBuilding: true,
		FallDamage:              true,
		DefaultGameMode:         world.GameModeSurvival,
		WorldsFolder:            t.TempDir(),
		MaxChunkRadius:          2,
		Listeners: []func(Config) (Listener, error){
			func(Config) (Listener, error) { return listener{Listener: l, log: log}, nil },
		},
	}).New()
	srv.Listen()
	t.Cleanup(func() { srv.Close() })
	joined := make(chan uuid.UUID, 1)
	go func() {
		for p := range srv.Accept() {
			p.Handle(damageTestHandler{})
			joined <- p.UUID()
		}
	}()
	return srv, l, joined
}

func damageTestPlayer(t *testing.T, srv *Server, id uuid.UUID, apply func(*world.Tx, *player.Player)) {
	t.Helper()
	handle, ok := srv.Player(id)
	if !ok {
		t.Fatal("damage test player missing")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := player.Call(ctx, handle, func(tx *world.Tx, p *player.Player) (struct{}, error) {
		apply(tx, p)
		return struct{}{}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDamageSourcesOnRakNet(t *testing.T) {
	data, err := os.ReadFile("../multiversion/conformance/testdata/damage.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Packets map[string][]byte `json:"packets"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	sources := map[string]world.DamageSource{
		"attack": entity.AttackDamageSource{}, "projectile": entity.ProjectileDamageSource{},
		"explosion": entity.ExplosionDamageSource{}, "fall": entity.FallDamageSource{},
		"glide": entity.GlideDamageSource{}, "lightning": entity.LightningDamageSource{},
		"void": entity.VoidDamageSource{}, "suffocation": entity.SuffocationDamageSource{},
		"drowning": entity.DrowningDamageSource{}, "fire": block.FireDamageSource{},
		"lava": block.LavaDamageSource{}, "magma": block.MagmaDamageSource{},
		"cactus": block.DamageSource{Block: block.Cactus{}}, "poison": effect.PoisonDamageSource{},
		"fatal_poison": effect.PoisonDamageSource{Fatal: true}, "wither": effect.WitherDamageSource{},
		"instant_damage": effect.InstantDamageSource{}, "starvation": player.StarvationDamageSource{},
		"thorns": enchantment.ThornsDamageSource{},
	}
	srv, l, joined := newDamageTestServer(t)
	cfg := (UserConfig{}).minecraftListenConfig(Config{Allower: allower{}})
	for _, p := range cfg.AcceptedProtocols {
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			type frame struct {
				id   uint32
				data []byte
			}
			frames := make(chan frame, 256)
			client, err := dialTestClient(minecraft.Dialer{
				Protocol: spawnProbeProtocol{Protocol: p}, EnableLegacyAuth: true, DisablePortFollowing: true,
				PacketFunc: func(h packet.Header, payload []byte, src, dst net.Addr) {
					if src.String() == l.Addr().String() && (h.PacketID == packet.IDUpdateAttributes || h.PacketID == packet.IDActorEvent) {
						frames <- frame{id: h.PacketID, data: append([]byte(nil), payload...)}
					}
				},
			}, l.Addr().String(), 8*time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			if err := client.DoSpawnTimeout(8 * time.Second); err != nil {
				t.Fatal(err)
			}
			var id uuid.UUID
			select {
			case id = <-joined:
			case <-time.After(8 * time.Second):
				t.Fatal("damage session not started")
			}
			go func() {
				for {
					if _, err := client.ReadPacket(); err != nil {
						return
					}
				}
			}()
			expected := fixtures[fmt.Sprint(p.ID())].Packets
			if len(expected) == 0 {
				t.Fatal("historical damage packets missing")
			}
			checkFrames := func(t *testing.T, healthName, eventName string) {
				t.Helper()
				health, hurt := false, false
				timeout := time.NewTimer(4 * time.Second)
				defer timeout.Stop()
				for !health || !hurt {
					select {
					case frame := <-frames:
						if frame.id == packet.IDUpdateAttributes && bytes.Equal(frame.data, expected[healthName]) {
							health = true
						}
						if frame.id == packet.IDActorEvent && bytes.Equal(frame.data, expected[eventName]) {
							hurt = true
						}
					case <-timeout.C:
						t.Fatalf("missing historical damage packets health=%v hurt=%v", health, hurt)
					}
				}
			}
			for name, source := range sources {
				t.Run(name, func(t *testing.T) {
					damageTestPlayer(t, srv, id, func(_ *world.Tx, p *player.Player) {
						p.Hurt(4, source)
						if p.Health() != 16 {
							t.Errorf("server health=%v want=16", p.Health())
						}
					})
					checkFrames(t, "health_hurt", "hurt")
					damageTestPlayer(t, srv, id, func(_ *world.Tx, p *player.Player) { p.Heal(4, entity.FoodHealingSource{}) })
				})
			}
			t.Run("death", func(t *testing.T) {
				damageTestPlayer(t, srv, id, func(_ *world.Tx, p *player.Player) {
					p.Hurt(25, entity.VoidDamageSource{})
					if !p.Dead() {
						t.Error("lethal damage did not kill player")
					}
				})
				checkFrames(t, "health_dead", "death")
			})
		})
	}
}

func TestFallDamageFromLegacyMovement(t *testing.T) {
	srv, l, joined := newDamageTestServer(t)
	cfg := (UserConfig{}).minecraftListenConfig(Config{Allower: allower{}})
	for _, p := range cfg.AcceptedProtocols {
		for _, pause := range []string{"none", "stationary", "rotation"} {
			t.Run(fmt.Sprintf("%d/%s", p.ID(), pause), func(t *testing.T) {
				client, err := dialTestClient(minecraft.Dialer{Protocol: spawnProbeProtocol{Protocol: p}, EnableLegacyAuth: true, DisablePortFollowing: true}, l.Addr().String(), 8*time.Second)
				if err != nil {
					t.Fatal(err)
				}
				defer client.Close()
				if err := client.DoSpawnTimeout(8 * time.Second); err != nil {
					t.Fatal(err)
				}
				var id uuid.UUID
				select {
				case id = <-joined:
				case <-time.After(8 * time.Second):
					t.Fatal("fall session not started")
				}
				go func() {
					for {
						if _, err := client.ReadPacket(); err != nil {
							return
						}
					}
				}()
				damageTestPlayer(t, srv, id, func(tx *world.Tx, p *player.Player) {
					tx.SetBlock(cube.Pos{0, 63, 0}, block.Stone{}, nil)
					p.Teleport(mgl64.Vec3{0.5, 70, 0.5})
				})
				tick := uint64(1)
				yaw := float32(0)
				send := func(y float32) {
					pk := &packet.PlayerAuthInput{Position: mgl32.Vec3{0.5, y + 1.62, 0.5}, Tick: tick, Yaw: yaw}
					tick++
					if err := client.WritePacket(pk); err != nil {
						t.Fatal(err)
					}
				}
				send(70)
				for y := float32(69); y >= 64; y-- {
					send(y)
					if pause != "none" && y == 67 {
						if pause == "rotation" {
							yaw = 30
						}
						send(y)
					}
				}
				if err := client.Flush(); err != nil {
					t.Fatal(err)
				}
				deadline := time.Now().Add(4 * time.Second)
				var health, y float64
				for {
					damageTestPlayer(t, srv, id, func(_ *world.Tx, p *player.Player) { health = p.Health(); y = p.Position().Y() })
					if math.Abs(y-64) < 0.001 || time.Now().After(deadline) {
						break
					}
					time.Sleep(time.Millisecond)
				}
				if math.Abs(y-64) > 0.001 || health != 17 {
					t.Fatalf("landing y=%v health=%v want y=64 health=17", y, health)
				}
			})
		}
	}
}
