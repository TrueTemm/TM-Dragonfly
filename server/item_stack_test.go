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
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestCreativeItemRequestsOnRakNet(t *testing.T) {
	data, err := os.ReadFile("../multiversion/conformance/testdata/item_stack.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Fixtures map[string][]byte `json:"fixtures"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	cfg := (UserConfig{}).minecraftListenConfig(Config{Compression: packet.SnappyCompression, AuthDisabled: true, Allower: allower{}})
	for _, p := range cfg.AcceptedProtocols {
		t.Run(fmt.Sprint(p.ID()), func(t *testing.T) {
			l, err := cfg.ListenNetwork(legacyRakNet{Network: rakNetNetwork(nil)}, "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer l.Close()
			done, release := make(chan error, 1), make(chan struct{})
			defer close(release)
			go func() {
				raw, err := l.Accept()
				if err != nil {
					done <- err
					return
				}
				conn := raw.(*minecraft.Conn)
				defer conn.Close()
				if err := conn.StartGameTimeout(minecraft.GameData{WorldName: "Creative request regression", EntityUniqueID: 1, EntityRuntimeID: 1}, 5*time.Second); err != nil {
					done <- err
					return
				}
				_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				for {
					pk, err := conn.ReadPacket()
					if err != nil {
						done <- err
						return
					}
					req, ok := pk.(*packet.ItemStackRequest)
					if !ok {
						continue
					}
					if len(req.Requests) != 1 || len(req.Requests[0].Actions) != 3 {
						done <- fmt.Errorf("creative request changed: %+v", req)
						return
					}
					actions := req.Requests[0].Actions
					creative, ok := actions[0].(*protocol.CraftCreativeStackRequestAction)
					if !ok || creative.CreativeItemNetworkID != 391 {
						done <- fmt.Errorf("creative action changed: %+v", actions[0])
						return
					}
					result, ok := actions[1].(*protocol.CraftResultsDeprecatedStackRequestAction)
					if !ok || len(result.ResultItems) != 1 || result.ResultItems[0].Identifier != "minecraft:stone_sword" {
						done <- fmt.Errorf("stone sword result changed: %+v", actions[1])
						return
					}
					take, ok := actions[2].(*protocol.TakeStackRequestAction)
					if !ok || take.Count != 1 || take.Source.StackNetworkID != -3 {
						done <- fmt.Errorf("transfer changed: %+v", actions[2])
						return
					}
					err = conn.WritePacket(&packet.ItemStackResponse{Responses: []protocol.ItemStackResponse{{RequestID: -3, Status: protocol.ItemStackResponseStatusOK}}})
					if err == nil {
						err = conn.Flush()
					}
					done <- err
					<-release
					return
				}
			}()
			client, err := dialTestClient(minecraft.Dialer{Protocol: spawnProbeProtocol{Protocol: p}, EnableLegacyAuth: true, DisablePortFollowing: true}, l.Addr().String(), 5*time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			if err := client.DoSpawnTimeout(5 * time.Second); err != nil {
				t.Fatal(err)
			}
			payload := fixtures[fmt.Sprint(p.ID())].Fixtures["creative_take"]
			if len(payload) == 0 {
				t.Fatal("historical creative fixture missing")
			}
			buf := new(bytes.Buffer)
			if err := (&packet.Header{PacketID: packet.IDItemStackRequest}).Write(buf); err != nil {
				t.Fatal(err)
			}
			buf.Write(payload)
			if _, err := client.Write(buf.Bytes()); err != nil {
				t.Fatal(err)
			}
			if err := client.Flush(); err != nil {
				t.Fatal(err)
			}
			_ = client.SetReadDeadline(time.Now().Add(5 * time.Second))
			for {
				pk, err := client.ReadPacket()
				if err != nil {
					t.Fatal(err)
				}
				if pk.ID() == packet.IDItemStackResponse {
					response := &packet.ItemStackResponse{Responses: reflect.ValueOf(pk).Elem().FieldByName("Responses").Interface().([]protocol.ItemStackResponse)}
					if len(response.Responses) != 1 || response.Responses[0].RequestID != -3 || response.Responses[0].Status != protocol.ItemStackResponseStatusOK {
						t.Fatalf("response changed: %+v", response)
					}
					break
				}
			}
			if err := <-done; err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreativeItemsInSession(t *testing.T) {
	data, err := os.ReadFile("../multiversion/conformance/testdata/item_stack.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Fixtures map[string][]byte `json:"fixtures"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	cfg := (UserConfig{}).minecraftListenConfig(Config{Compression: packet.SnappyCompression, AuthDisabled: true, Allower: allower{}})
	l, err := cfg.ListenNetwork(legacyRakNet{Network: rakNetNetwork(nil)}, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	log := slog.New(slog.DiscardHandler)
	srv := (Config{
		Log:                     log,
		AuthDisabled:            true,
		DisableResourceBuilding: true,
		DefaultGameMode:         world.GameModeCreative,
		WorldsFolder:            t.TempDir(),
		MaxChunkRadius:          2,
		Listeners: []func(Config) (Listener, error){
			func(Config) (Listener, error) { return listener{Listener: l, log: log}, nil },
		},
	}).New()
	srv.Listen()
	defer srv.Close()
	joined := make(chan uuid.UUID, 1)
	go func() {
		for p := range srv.Accept() {
			joined <- p.H().UUID()
		}
	}()
	for _, p := range cfg.AcceptedProtocols {
		for _, want := range []struct {
			name  string
			count byte
			place bool
		}{{"minecraft:stone_sword", 1, false}, {"minecraft:white_wool", 64, true}, {"minecraft:snowball", 16, true}} {
			t.Run(fmt.Sprintf("%d/%s", p.ID(), want.name), func(t *testing.T) {
				var creativeID uint32
				for i, it := range creative.Items() {
					name, _ := it.Stack.Item().EncodeItem()
					if name == want.name {
						creativeID = uint32(i + 1)
						break
					}
				}
				if creativeID == 0 {
					t.Fatalf("%s absent from creative registry", want.name)
				}

				client, err := dialTestClient(minecraft.Dialer{Protocol: spawnProbeProtocol{Protocol: p}, EnableLegacyAuth: true, DisablePortFollowing: true}, l.Addr().String(), 8*time.Second)
				if err != nil {
					t.Fatal(err)
				}
				defer client.Close()
				timer := time.AfterFunc(10*time.Second, func() { client.Close() })
				defer timer.Stop()
				if err := client.DoSpawnTimeout(8 * time.Second); err != nil {
					t.Fatal(err)
				}
				var playerID uuid.UUID
				select {
				case playerID = <-joined:
				case <-time.After(8 * time.Second):
					t.Fatal("session not started")
				}
				pk := p.Packets(true)[packet.IDItemStackRequest]()
				pk.Marshal(p.NewReader(bytes.NewBuffer(fixtures[fmt.Sprint(p.ID())].Fixtures["creative_take"]), 0, true))
				req := p.ConvertToLatest(pk, nil)[0].(*packet.ItemStackRequest)
				req.Requests[0].Actions[0].(*protocol.CraftCreativeStackRequestAction).CreativeItemNetworkID = creativeID
				result := req.Requests[0].Actions[1].(*protocol.CraftResultsDeprecatedStackRequestAction)
				result.ResultItems[0].Identifier = want.name
				result.ResultItems[0].Count = uint16(want.count)
				req.Requests[0].Actions[2].(*protocol.TakeStackRequestAction).Count = want.count
				if want.place {
					take := req.Requests[0].Actions[2].(*protocol.TakeStackRequestAction)
					place := &protocol.PlaceStackRequestAction{}
					place.Count, place.Source, place.Destination = take.Count, take.Source, take.Destination
					req.Requests[0].Actions[2] = place
				}
				buf := new(bytes.Buffer)
				if err := (&packet.Header{PacketID: packet.IDItemStackRequest}).Write(buf); err != nil {
					t.Fatal(err)
				}
				pk.Marshal(p.NewWriter(buf, 0))
				if _, err := client.Write(buf.Bytes()); err != nil {
					t.Fatal(err)
				}
				if err := client.Flush(); err != nil {
					t.Fatal(err)
				}
				for {
					pk, err := client.ReadPacket()
					if err != nil {
						t.Fatal(err)
					}
					if pk.ID() != packet.IDItemStackResponse {
						continue
					}
					responses := reflect.ValueOf(pk).Elem().FieldByName("Responses").Interface().([]protocol.ItemStackResponse)
					if len(responses) != 1 || responses[0].Status != protocol.ItemStackResponseStatusOK {
						t.Fatalf("creative request rejected: %+v", responses)
					}
					break
				}
				var heldCount int
				name, err := world.Call(context.Background(), srv.World(), func(tx *world.Tx) (string, error) {
					for player := range srv.Players(tx) {
						if player.H().UUID() != playerID {
							continue
						}
						hand, _ := player.HeldItems()
						if hand.Empty() {
							return "", nil
						}
						heldCount = hand.Count()
						name, _ := hand.Item().EncodeItem()
						return name, nil
					}
					return "", nil
				})
				if err != nil || name != want.name || heldCount != int(want.count) {
					t.Fatalf("server held item = %s x%d want %s x%d: %v", name, heldCount, want.name, want.count, err)
				}
			})
		}
	}
}
