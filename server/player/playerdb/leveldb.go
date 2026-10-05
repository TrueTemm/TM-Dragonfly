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

package playerdb

import (
	"encoding/json"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/goleveldb/leveldb"
	"github.com/df-mc/goleveldb/leveldb/opt"
	"github.com/google/uuid"
	"os"
)

type Provider struct {
	db *leveldb.DB
}

func NewProvider(path string) (*Provider, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0777)
	}
	db, err := leveldb.OpenFile(path, &opt.Options{Compression: opt.SnappyCompression})
	if err != nil {
		return nil, err
	}
	return &Provider{db: db}, nil
}

func (p *Provider) Save(id uuid.UUID, d player.Config, w *world.World) error {
	b, err := json.Marshal(p.toJson(d, w))
	if err != nil {
		return err
	}
	return p.db.Put(id[:], b, nil)
}

func (p *Provider) Load(id uuid.UUID, world func(world.Dimension) *world.World) (player.Config, *world.World, error) {
	b, err := p.db.Get(id[:], nil)
	if err != nil {
		return player.Config{}, nil, err
	}
	var d jsonData
	err = json.Unmarshal(b, &d)
	if err != nil {
		return player.Config{}, nil, err
	}
	conf, w := p.fromJson(d, world)

	return conf, w, nil
}

func (p *Provider) Close() error {
	return p.db.Close()
}
