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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/internal/packbuilder"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/playerdb"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/biome"
	"github.com/df-mc/dragonfly/server/world/generator"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

type Config struct {
	Log *slog.Logger

	Listeners []func(conf Config) (Listener, error)

	Name string

	SubName string

	Resources []*resource.Pack

	ResourcesRequired bool

	DisableResourceBuilding bool

	Allower Allower

	AuthDisabled bool

	MuteEmoteChat bool

	MaxPlayers int

	MaxChunkRadius int

	JoinMessage, QuitMessage, ShutdownMessage chat.Translation

	StatusProvider minecraft.ServerStatusProvider

	Compression packet.Compression

	PlayerProvider player.Provider

	WorldProvider world.Provider

	ReadOnlyWorld bool

	WorldsFolder string

	AlwaysDay bool

	FallDamage bool

	DefaultGameMode world.GameMode

	Generator func(dim world.Dimension) world.Generator

	RandomTickSpeed int

	SaveInterval time.Duration

	ChunkUnloadInterval time.Duration

	ChunkLoadWorkers int

	Entities world.EntityRegistry

	Blocks world.BlockRegistry
}

func (conf Config) New() *Server {
	if conf.Log == nil {
		conf.Log = slog.Default()
	}
	if len(conf.Listeners) == 0 {
		conf.Log.Warn("config: no listeners set, no connections will be accepted")
	}
	conf.Name, conf.SubName = srvID(), srvSub()
	conf.StatusProvider = statusProvider{}
	if conf.PlayerProvider == nil {
		conf.PlayerProvider = player.NopProvider{}
	}
	if conf.Allower == nil {
		conf.Allower = allower{}
	}
	if conf.WorldProvider == nil {
		conf.WorldProvider = world.NopProvider{}
	}
	if conf.Generator == nil {
		conf.Generator = loadGenerator
	}
	if conf.MaxChunkRadius == 0 {
		conf.MaxChunkRadius = 12
	}
	if conf.ShutdownMessage.Zero() {
		conf.ShutdownMessage = chat.MessageServerDisconnect
	}
	if len(conf.Entities.Types()) == 0 {
		conf.Entities = entity.DefaultRegistry
	}
	if conf.Blocks == nil {
		conf.Blocks = world.DefaultBlockRegistry
	}
	if conf.Compression == nil {
		conf.Compression = packet.SnappyCompression
	}

	conf.Blocks.Finalize()
	world.DefaultBlockRegistry.Finalize()

	if !conf.DisableResourceBuilding {
		if pack, ok := packbuilder.BuildResourcePack(conf.Blocks); ok {
			conf.Resources = append(conf.Resources, pack)
		}
	}

	conf.Resources = slices.Clone(conf.Resources)

	srv := &Server{
		conf:     conf,
		incoming: make(chan incoming),
		p:        make(map[uuid.UUID]*onlinePlayer),
		world:    &world.World{},
	}
	for _, lf := range conf.Listeners {
		l, err := lf(conf)
		if err != nil {
			conf.Log.Error("create listener: " + err.Error())
			continue
		}
		srv.listeners = append(srv.listeners, l)
	}

	creative_registerCreativeItems()
	recipe_registerVanilla()

	srv.world = srv.newWorld(conf.WorldProvider, MainWorldName, conf.Generator(world.Overworld), false)
	srv.wm.worlds = map[string]*world.World{MainWorldName: srv.world}
	srv.conf.Log.Info("Opened world.", "name", srv.world.Name())

	return srv
}

type UserConfig struct {
	Network struct {
		Address string
	}
	Server struct {
		Name string

		SubName string

		AuthEnabled bool

		DisableJoinQuitMessages bool

		MuteEmoteChat bool
	}
	World struct {
		SaveData bool

		Folder string

		WorldsFolder string

		AlwaysDay bool

		FallDamage bool

		DefaultGameMode string

		RandomTickSpeed int
	}
	Players struct {
		MaxCount int

		MaximumChunkRadius int

		SaveData bool

		Folder string
	}
	Debug struct {
		PprofAddress string
	}
	Resources struct {
		AutoBuildPack bool

		Folder string

		Required bool
	}
}

func (uc UserConfig) Config(log *slog.Logger) (Config, error) {
	var err error
	conf := Config{
		Log:                     log,
		Name:                    uc.Server.Name,
		SubName:                 uc.Server.SubName,
		ResourcesRequired:       uc.Resources.Required,
		AuthDisabled:            !uc.Server.AuthEnabled,
		MuteEmoteChat:           uc.Server.MuteEmoteChat,
		MaxPlayers:              uc.Players.MaxCount,
		MaxChunkRadius:          uc.Players.MaximumChunkRadius,
		DisableResourceBuilding: !uc.Resources.AutoBuildPack,
		RandomTickSpeed:         uc.World.RandomTickSpeed,
		AlwaysDay:               uc.World.AlwaysDay,
		FallDamage:              uc.World.FallDamage,
	}
	if uc.World.DefaultGameMode != "" {
		mode, ok := gameModeByName(uc.World.DefaultGameMode)
		if !ok {
			return conf, fmt.Errorf("config: unknown DefaultGameMode %q (survival, creative, adventure or spectator)", uc.World.DefaultGameMode)
		}
		conf.DefaultGameMode = mode
	}
	if uc.World.SaveData {
		conf.WorldsFolder = uc.World.WorldsFolder
		if conf.WorldsFolder != "" {

			if err = os.MkdirAll(conf.WorldsFolder, 0777); err != nil {
				return conf, fmt.Errorf("create worlds folder: %w", err)
			}
		}
	}
	if conf.RandomTickSpeed <= 0 {

		conf.RandomTickSpeed = -1
	}
	if !uc.Server.DisableJoinQuitMessages {
		conf.JoinMessage, conf.QuitMessage = chat.MessageJoin, chat.MessageQuit
	}
	if uc.World.SaveData {
		conf.WorldProvider, err = mcdb.Config{Log: log}.Open(uc.World.Folder)
		if err != nil {
			return conf, fmt.Errorf("create world provider: %w", err)
		}
	}
	conf.Resources, err = loadResources(uc.Resources.Folder)
	if err != nil {
		return conf, fmt.Errorf("load resources: %w", err)
	}

	folder, _ := filepath.Abs(uc.Resources.Folder)
	for i, pack := range conf.Resources {
		log.Info("Loaded resource pack.", "order", i+1, "name", pack.Name(), "uuid", pack.UUID().String(), "version", pack.Version(), "size", pack.Len(), "folder", folder)
	}
	if len(conf.Resources) == 0 {
		log.Info("No resource packs found.", "folder", uc.Resources.Folder)
	}
	if uc.Players.SaveData {
		conf.PlayerProvider, err = playerdb.NewProvider(uc.Players.Folder)
		if err != nil {
			return conf, fmt.Errorf("create player provider: %w", err)
		}
	}
	conf.Listeners = append(conf.Listeners, uc.listenerFunc)
	return conf, nil
}

func loadResources(dir string) ([]*resource.Pack, error) {
	_ = os.MkdirAll(dir, 0777)

	resources, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	packs := make([]*resource.Pack, 0, len(resources))
	for _, entry := range resources {
		name := entry.Name()
		if !entry.IsDir() {
			switch strings.ToLower(filepath.Ext(name)) {
			case ".mcpack", ".zip":
			default:
				continue
			}
		}
		pack, err := resource.ReadPath(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("compile resource (%v): %w", name, err)
		}
		packs = append(packs, pack)
	}
	return newestOfEach(packs), nil
}

func newestOfEach(packs []*resource.Pack) []*resource.Pack {
	newest := make(map[string]*resource.Pack)
	var order []string
	for _, p := range packs {
		id := p.UUID().String()
		current, seen := newest[id]
		if !seen {
			order = append(order, id)
		}
		if !seen || newerVersion(p.Version(), current.Version()) {
			newest[id] = p
		}
	}
	kept := make([]*resource.Pack, 0, len(order))
	for _, id := range order {
		kept = append(kept, newest[id])
	}
	if len(kept) != len(packs) {
		for _, p := range packs {
			if newest[p.UUID().String()] != p {
				slog.Warn("An older version of a resource pack is ignored: delete it.", "name", p.Name(), "uuid", p.UUID().String(), "version", p.Version())
			}
		}
	}
	return kept
}

func newerVersion(a, b string) bool {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		x, _ := strconv.Atoi(as[i])
		y, _ := strconv.Atoi(bs[i])
		if x != y {
			return x > y
		}
	}
	return len(as) > len(bs)
}

func loadGenerator(world.Dimension) world.Generator {
	return generator.NewFlat(biome.Plains{}, []world.Block{block.Grass{}, block.Dirt{}, block.Dirt{}, block.Bedrock{}})
}

func gameModeByName(name string) (world.GameMode, bool) {
	switch strings.ToLower(name) {
	case "survival", "s", "0":
		return world.GameModeSurvival, true
	case "creative", "c", "1":
		return world.GameModeCreative, true
	case "adventure", "a", "2":
		return world.GameModeAdventure, true
	case "spectator", "sp", "3":
		return world.GameModeSpectator, true
	}
	return nil, false
}

func DefaultConfig() UserConfig {
	c := UserConfig{}
	c.Network.Address = ":19132"
	c.Server.Name = "A TM-Dragonfly server"
	c.Server.SubName = "Powered by TM-Dragonfly"
	c.Server.AuthEnabled = true
	c.World.SaveData = true
	c.World.Folder = "worlds/world"
	c.World.WorldsFolder = "worlds"
	c.World.AlwaysDay = true
	c.World.FallDamage = false
	c.World.DefaultGameMode = "survival"
	c.World.RandomTickSpeed = 0

	c.Players.MaximumChunkRadius = 10
	c.Players.SaveData = true
	c.Players.Folder = "players"
	c.Resources.AutoBuildPack = true
	c.Resources.Folder = "resources"
	c.Resources.Required = false
	return c
}

//go:linkname creative_registerCreativeItems github.com/df-mc/dragonfly/server/item/creative.registerCreativeItems
func creative_registerCreativeItems()

//go:linkname recipe_registerVanilla github.com/df-mc/dragonfly/server/item/recipe.registerVanilla
func recipe_registerVanilla()
