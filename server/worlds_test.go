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
	"log/slog"
	"testing"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

func TestWorldManagerPristine(t *testing.T) {
	conf := Config{Log: slog.Default(), WorldsFolder: t.TempDir()}
	srv := conf.New()
	defer srv.world.Close()

	pos := cube.Pos{3, -61, 3}
	blockAt := func(w *world.World) world.Block {
		b, _ := world.Call(context.Background(), w, func(tx *world.Tx) (world.Block, error) { return tx.Block(pos), nil })
		return b
	}
	setStone := func(w *world.World) {
		_, _ = world.Call(context.Background(), w, func(tx *world.Tx) (struct{}, error) {
			tx.SetBlock(pos, block.Stone{}, nil)
			return struct{}{}, nil
		})
	}

	if _, err := srv.LoadWorld("arena"); !errors.Is(err, ErrWorldNotFound) {
		t.Fatalf("LoadWorld of a missing world: got %v, want ErrWorldNotFound", err)
	}
	w, err := srv.CreateWorld("arena")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := blockAt(w).(block.Grass); !ok {
		t.Fatalf("CreateWorld must generate flat terrain, got %T", blockAt(w))
	}
	setStone(w)
	if err := srv.UnloadWorld("arena"); err != nil {
		t.Fatal(err)
	}

	w, err = srv.LoadWorld("arena")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := blockAt(w).(block.Stone); !ok {
		t.Fatalf("the writable edit was not saved, got %T", blockAt(w))
	}
	far := cube.Pos{5000, -61, 5000}
	if b, _ := world.Call(context.Background(), w, func(tx *world.Tx) (world.Block, error) { return tx.Block(far), nil }); b != (block.Air{}) {
		t.Fatalf("a loaded world must not generate terrain, got %T outside the map", b)
	}
	_, _ = world.Call(context.Background(), w, func(tx *world.Tx) (struct{}, error) {
		tx.SetBlock(pos, block.Dirt{}, nil)
		return struct{}{}, nil
	})
	if err := srv.UnloadWorld("arena"); err != nil {
		t.Fatal(err)
	}
	w, err = srv.LoadWorld("arena")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := blockAt(w).(block.Stone); !ok {
		t.Fatalf("a read-only world must come back pristine after unload, got %T", blockAt(w))
	}
	if names := srv.WorldNames(); len(names) != 2 || names[1] != "arena" {
		t.Fatalf("WorldNames = %v", names)
	}
	if err := srv.UnloadWorld("arena"); err != nil {
		t.Fatal(err)
	}
	if err := srv.UnloadWorld(MainWorldName); err == nil {
		t.Fatal("the main world must not be unloadable")
	}
}
