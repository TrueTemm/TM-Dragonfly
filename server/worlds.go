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
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
)

const worldNoon = 6000

var worldNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

var ErrWorldNotLoaded = errors.New("world not loaded")

var ErrWorldNotFound = errors.New("world not found on disk")

type WorldOption func(*worldOptions)

type worldOptions struct {
	writable  bool
	generator world.Generator
}

func Writable() WorldOption { return func(o *worldOptions) { o.writable = true } }

func WithGenerator(g world.Generator) WorldOption {
	return func(o *worldOptions) { o.generator = g }
}

type worldManager struct {
	mu     sync.RWMutex
	worlds map[string]*world.World
}

const MainWorldName = "world"

func (srv *Server) Worlds() []*world.World {
	srv.wm.mu.RLock()
	defer srv.wm.mu.RUnlock()
	names := make([]string, 0, len(srv.wm.worlds))
	for name := range srv.wm.worlds {
		if name != MainWorldName {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	worlds := []*world.World{srv.world}
	for _, name := range names {
		worlds = append(worlds, srv.wm.worlds[name])
	}
	return worlds
}

func (srv *Server) WorldNames() []string {
	names := make([]string, 0, 8)
	for _, w := range srv.Worlds() {
		names = append(names, srv.WorldName(w))
	}
	return names
}

func (srv *Server) WorldByName(name string) (*world.World, bool) {
	srv.wm.mu.RLock()
	defer srv.wm.mu.RUnlock()
	w, ok := srv.wm.worlds[name]
	return w, ok
}

func (srv *Server) WorldName(w *world.World) string {
	srv.wm.mu.RLock()
	defer srv.wm.mu.RUnlock()
	for name, other := range srv.wm.worlds {
		if other == w {
			return name
		}
	}
	return ""
}

func (srv *Server) WorldsFolder() string {
	return srv.conf.WorldsFolder
}

func (srv *Server) LoadWorld(name string, opts ...WorldOption) (*world.World, error) {
	if !worldNamePattern.MatchString(name) {
		return nil, fmt.Errorf("load world: invalid world name %q", name)
	}
	o := worldOptions{generator: world.NopGenerator{}}
	for _, opt := range opts {
		opt(&o)
	}
	srv.wm.mu.Lock()
	defer srv.wm.mu.Unlock()
	if w, ok := srv.wm.worlds[name]; ok {
		return w, nil
	}
	if srv.conf.WorldsFolder == "" {
		return nil, fmt.Errorf("load world %s: %w (world saving is disabled)", name, ErrWorldNotFound)
	}
	dir := filepath.Join(srv.conf.WorldsFolder, name)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("load world %s: %w", name, ErrWorldNotFound)
	}
	return srv.openWorld(name, dir, o)
}

func (srv *Server) CreateWorld(name string, opts ...WorldOption) (*world.World, error) {
	if !worldNamePattern.MatchString(name) {
		return nil, fmt.Errorf("create world: invalid world name %q", name)
	}
	o := worldOptions{writable: true, generator: srv.conf.Generator(world.Overworld)}
	for _, opt := range opts {
		opt(&o)
	}
	srv.wm.mu.Lock()
	defer srv.wm.mu.Unlock()
	if _, ok := srv.wm.worlds[name]; ok {
		return nil, fmt.Errorf("create world %s: already loaded", name)
	}
	if srv.conf.WorldsFolder == "" {
		return nil, fmt.Errorf("create world %s: world saving is disabled", name)
	}
	dir := filepath.Join(srv.conf.WorldsFolder, name)
	if _, err := os.Stat(dir); err == nil {
		return nil, fmt.Errorf("create world %s: already exists on disk", name)
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, fmt.Errorf("create world %s: %w", name, err)
	}
	return srv.openWorld(name, dir, o)
}

func (srv *Server) openWorld(name, dir string, o worldOptions) (*world.World, error) {
	db, err := mcdb.Config{Log: srv.conf.Log, Blocks: srv.conf.Blocks}.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open world %s: %w", name, err)
	}
	w := srv.newWorld(db, name, o.generator, !o.writable)
	if srv.wm.worlds == nil {
		srv.wm.worlds = map[string]*world.World{MainWorldName: srv.world}
	}
	srv.wm.worlds[name] = w
	srv.conf.Log.Info("Loaded world.", "name", name, "level", w.Name(), "read-only", !o.writable)
	return w, nil
}

func (srv *Server) UnloadWorld(name string) error {
	if name == MainWorldName {
		return fmt.Errorf("unload world: the main world cannot be unloaded")
	}
	srv.wm.mu.Lock()
	w, ok := srv.wm.worlds[name]
	if ok {
		delete(srv.wm.worlds, name)
	}
	srv.wm.mu.Unlock()
	if !ok {
		return fmt.Errorf("unload world %s: %w", name, ErrWorldNotLoaded)
	}

	mainWorld := srv.world
	spawn := mainWorld.Spawn().Vec3Middle()
	_, _ = world.Call(context.Background(), w, func(tx *world.Tx) (struct{}, error) {
		for e := range tx.Entities() {
			if p, ok := e.(*player.Player); ok {
				p.TransferTo(mainWorld, spawn)
			}
		}
		return struct{}{}, nil
	})

	_, _ = world.Call(context.Background(), w, func(*world.Tx) (struct{}, error) { return struct{}{}, nil })
	if err := w.Close(); err != nil {
		return fmt.Errorf("unload world %s: %w", name, err)
	}
	srv.conf.Log.Info("Unloaded world.", "name", name)
	return nil
}

func (srv *Server) ListWorldFolders() []string {
	if srv.conf.WorldsFolder == "" {
		return nil
	}
	entries, err := os.ReadDir(srv.conf.WorldsFolder)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && worldNamePattern.MatchString(e.Name()) {
			names = append(names, e.Name())
		}
	}
	return names
}

func (srv *Server) closeWorlds() {
	srv.wm.mu.Lock()
	worlds := srv.wm.worlds
	srv.wm.worlds = nil
	srv.wm.mu.Unlock()
	for name, w := range worlds {
		if name == MainWorldName {
			continue
		}

		done := make(chan error, 1)
		go func() { done <- w.Close() }()
		select {
		case err := <-done:
			if err != nil {
				srv.conf.Log.Error("Close world " + name + ": " + err.Error())
			}
		case <-time.After(closeWorldTimeout):
			srv.conf.Log.Error("Timed out closing a world; closing the rest anyway.", "world", name, "after", closeWorldTimeout.String())
		}
	}
}

const closeWorldTimeout = 5 * time.Second

func (srv *Server) newWorld(provider world.Provider, name string, generator world.Generator, readOnly bool) *world.World {
	logger := srv.conf.Log.With("world", name)
	conf := world.Config{
		Log:                 logger,
		Dim:                 world.Overworld,
		Provider:            provider,
		Generator:           generator,
		RandomTickSpeed:     srv.conf.RandomTickSpeed,
		ReadOnly:            readOnly || srv.conf.ReadOnlyWorld,
		SaveInterval:        srv.conf.SaveInterval,
		ChunkUnloadInterval: srv.conf.ChunkUnloadInterval,
		ChunkLoadWorkers:    srv.conf.ChunkLoadWorkers,
		Entities:            srv.conf.Entities,
		Blocks:              srv.conf.Blocks,
	}
	w := conf.New()
	w.SetFallDamage(srv.conf.FallDamage)
	if srv.conf.DefaultGameMode != nil {
		w.SetDefaultGameMode(srv.conf.DefaultGameMode)
	}

	w.StopWeatherCycle()
	w.StopThundering()
	w.StopRaining()
	if srv.conf.AlwaysDay {
		w.StopTime()
		w.SetTime(worldNoon)
	} else {

		w.StartTime()
	}
	return w
}
