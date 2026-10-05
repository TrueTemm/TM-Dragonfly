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

package world

import (
	"sync"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/goleveldb/leveldb"
	"github.com/google/uuid"
	"io"
)

type Provider interface {
	io.Closer

	Settings() *Settings

	SaveSettings(*Settings)

	LoadPlayerSpawnPosition(uuid uuid.UUID) (pos cube.Pos, exists bool, err error)

	SavePlayerSpawnPosition(uuid uuid.UUID, pos cube.Pos) error

	LoadColumn(pos ChunkPos, dim Dimension) (*chunk.Column, error)

	StoreColumn(pos ChunkPos, dim Dimension, col *chunk.Column) error
}

var _ Provider = (*NopProvider)(nil)

type NopProvider struct {
	Set *Settings
}

func (n NopProvider) Settings() *Settings {
	if n.Set == nil {
		return defaultSettings()
	}
	return n.Set
}
func (NopProvider) SaveSettings(*Settings) {}
func (NopProvider) LoadColumn(ChunkPos, Dimension) (*chunk.Column, error) {
	return nil, leveldb.ErrNotFound
}
func (NopProvider) StoreColumn(ChunkPos, Dimension, *chunk.Column) error { return nil }
func (NopProvider) LoadPlayerSpawnPosition(uuid.UUID) (cube.Pos, bool, error) {
	return cube.Pos{}, false, nil
}
func (NopProvider) SavePlayerSpawnPosition(uuid.UUID, cube.Pos) error { return nil }
func (NopProvider) Close() error                                      { return nil }

type lockedProvider struct {
	mu sync.Mutex
	p  Provider
}

func (l *lockedProvider) Settings() *Settings {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.p.Settings()
}

func (l *lockedProvider) SaveSettings(s *Settings) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.p.SaveSettings(s)
}

func (l *lockedProvider) LoadPlayerSpawnPosition(id uuid.UUID) (cube.Pos, bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.p.LoadPlayerSpawnPosition(id)
}

func (l *lockedProvider) SavePlayerSpawnPosition(id uuid.UUID, pos cube.Pos) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.p.SavePlayerSpawnPosition(id, pos)
}

func (l *lockedProvider) LoadColumn(pos ChunkPos, dim Dimension) (*chunk.Column, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.p.LoadColumn(pos, dim)
}

func (l *lockedProvider) StoreColumn(pos ChunkPos, dim Dimension, col *chunk.Column) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.p.StoreColumn(pos, dim, col)
}

func (l *lockedProvider) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.p.Close()
}
