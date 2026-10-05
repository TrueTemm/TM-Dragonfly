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

package mcdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/dragonfly/server/world/mcdb/leveldat"
	"github.com/df-mc/goleveldb/leveldb"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type DB struct {
	conf Config
	ldb  *leveldb.DB
	dir  string
	ldat *leveldat.Data
	set  *world.Settings
}

func Open(dir string) (*DB, error) {
	var conf Config
	return conf.Open(dir)
}

func (db *DB) SetBlockRegistry(br world.BlockRegistry) {
	if br == nil {
		br = world.DefaultBlockRegistry
	}
	br.Finalize()
	db.conf.Blocks = br
}

func (db *DB) Settings() *world.Settings {
	return db.set
}

func (db *DB) SaveSettings(s *world.Settings) {
	db.ldat.PutSettings(s)
}

type playerData struct {
	UUID         string `nbt:"MsaId"`
	ServerID     string `nbt:"ServerId"`
	SelfSignedID string `nbt:"SelfSignedId"`
}

func (db *DB) LoadPlayerSpawnPosition(id uuid.UUID) (pos cube.Pos, exists bool, err error) {
	serverData, _, exists, err := db.loadPlayerData(id)
	if !exists || err != nil {
		return cube.Pos{}, exists, err
	}
	x, y, z := serverData["SpawnX"], serverData["SpawnY"], serverData["SpawnZ"]
	if x == nil || y == nil || z == nil {
		return cube.Pos{}, true, fmt.Errorf("error reading spawn fields from server data for player %v", id)
	}
	return cube.Pos{int(x.(int32)), int(y.(int32)), int(z.(int32))}, true, nil
}

func (db *DB) loadPlayerData(id uuid.UUID) (serverData map[string]interface{}, key string, exists bool, err error) {
	data, err := db.ldb.Get([]byte("player_"+id.String()), nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, "", false, nil
	} else if err != nil {
		return nil, "", true, fmt.Errorf("error reading player data for uuid %v: %w", id, err)
	}

	var d playerData
	if err := nbt.UnmarshalEncoding(data, &d, nbt.LittleEndian); err != nil {
		return nil, "", true, fmt.Errorf("error decoding player data for uuid %v: %w", id, err)
	}
	if d.UUID != id.String() || d.ServerID == "" {
		return nil, d.ServerID, true, fmt.Errorf("invalid player data for uuid %v: %v", id, d)
	}
	serverDB, err := db.ldb.Get([]byte(d.ServerID), nil)
	if err != nil {
		return nil, d.ServerID, true, fmt.Errorf("error reading server data for player %v (%v): %w", id, d.ServerID, err)
	}

	if err := nbt.UnmarshalEncoding(serverDB, &serverData, nbt.LittleEndian); err != nil {
		return nil, d.ServerID, true, fmt.Errorf("error decoding server data for player %v", id)
	}
	return serverData, d.ServerID, true, nil
}

func (db *DB) SavePlayerSpawnPosition(id uuid.UUID, pos cube.Pos) error {
	_, err := db.ldb.Get([]byte("player_"+id.String()), nil)
	d := make(map[string]interface{})
	k := "player_server_" + id.String()

	if errors.Is(err, leveldb.ErrNotFound) {
		data, err := nbt.MarshalEncoding(playerData{UUID: id.String(), ServerID: k}, nbt.LittleEndian)
		if err != nil {
			panic(err)
		}
		if err := db.ldb.Put([]byte("player_"+id.String()), data, nil); err != nil {
			return fmt.Errorf("write player data (uuid=%v): %w", id, err)
		}
	} else if d, k, _, err = db.loadPlayerData(id); err != nil {
		return err
	}
	d["SpawnX"], d["SpawnY"], d["SpawnZ"] = int32(pos.X()), int32(pos.Y()), int32(pos.Z())

	data, err := nbt.MarshalEncoding(d, nbt.LittleEndian)
	if err != nil {
		panic(err)
	}
	if err = db.ldb.Put([]byte(k), data, nil); err != nil {
		return fmt.Errorf("write server data for player %v: %w", id, err)
	}
	return nil
}

func (db *DB) LoadColumn(pos world.ChunkPos, dim world.Dimension) (*chunk.Column, error) {
	k := dbKey{pos: pos, dim: dim}
	col, err := db.column(k)
	if err != nil {
		return nil, fmt.Errorf("load column %v (%v): %w", pos, dim, err)
	}
	return col, nil
}

const chunkVersion = 42

func (db *DB) column(k dbKey) (*chunk.Column, error) {
	var cdata chunk.SerialisedData
	col := new(chunk.Column)

	ver, err := db.version(k)
	if err != nil {
		return nil, fmt.Errorf("read version: %w", err)
	}
	if ver != chunkVersion {
		db.conf.Log.Debug("column: unsupported chunk version, trying to load anyway", "X", k.pos[0], "Z", k.pos[1], "dimension", fmt.Sprint(k.dim), "ver", ver)
	}
	cdata.Biomes, err = db.biomes(k)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {

		return nil, fmt.Errorf("read biomes: %w", err)
	}
	cdata.SubChunks, err = db.subChunks(k)
	if err != nil {
		return nil, fmt.Errorf("read sub chunks: %w", err)
	}
	col.Chunk, err = chunk.DiskDecode(db.conf.Blocks, cdata, k.dim.Range())
	if err != nil {
		return nil, fmt.Errorf("decode chunk data: %w", err)
	}
	col.Entities, err = db.entities(k)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {

		return nil, fmt.Errorf("read entities: %w", err)
	}
	col.BlockEntities, err = db.blockEntities(k)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {

		return nil, fmt.Errorf("read block entities: %w", err)
	}
	col.ScheduledBlocks, col.Tick, err = db.scheduledUpdates(k)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, fmt.Errorf("read scheduled updates: %w", err)
	}
	return col, nil
}

func (db *DB) version(k dbKey) (byte, error) {
	p, err := db.ldb.Get(k.Sum(keyVersion), nil)
	if errors.Is(err, leveldb.ErrNotFound) {

		if p, err = db.ldb.Get(k.Sum(keyVersionOld), nil); err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}
	if n := len(p); n != 1 {
		return 0, fmt.Errorf("expected 1 version byte, got %v", n)
	}
	return p[0], nil
}

func (db *DB) biomes(k dbKey) ([]byte, error) {
	biomes, err := db.ldb.Get(k.Sum(key3DData), nil)
	if err != nil {
		return nil, err
	}

	if n := len(biomes); n <= 512 {
		return nil, fmt.Errorf("expected at least 513 bytes for 3D data, got %v", n)
	}
	return biomes[512:], nil
}

func (db *DB) subChunks(k dbKey) ([][]byte, error) {
	r := k.dim.Range()
	sub := make([][]byte, (r.Height()>>4)+1)

	var err error
	for i := range sub {
		y := uint8(i + (r[0] >> 4))
		sub[i], err = db.ldb.Get(k.Sum(keySubChunkData, y), nil)
		if errors.Is(err, leveldb.ErrNotFound) {

			continue
		} else if err != nil {
			return nil, fmt.Errorf("sub chunk %v: %w", int8(i), err)
		}
	}
	return sub, nil
}

func (db *DB) entities(k dbKey) ([]chunk.Entity, error) {

	ids, err := db.ldb.Get(append([]byte(keyEntityIdentifiers), index(k.pos, k.dim)...), nil)
	if err != nil {

		return db.entitiesOld(k)
	}
	entities := make([]chunk.Entity, 0, len(ids)/8)
	for i := 0; i < len(ids); i += 8 {
		id := int64(binary.LittleEndian.Uint64(ids[i : i+8]))
		data, err := db.ldb.Get(entityIndex(id), nil)
		if err != nil {

			db.conf.Log.Debug("read entity: "+err.Error(), "ID", id)
			continue
		}
		ent := chunk.Entity{ID: id, Data: make(map[string]any)}
		if err = nbt.UnmarshalEncoding(data, &ent.Data, nbt.LittleEndian); err != nil {
			db.conf.Log.Error("decode entity nbt: "+err.Error(), "ID", id)
		}
		entities = append(entities, ent)
	}
	return entities, nil
}

func (db *DB) entitiesOld(k dbKey) ([]chunk.Entity, error) {
	data, err := db.ldb.Get(k.Sum(keyEntitiesOld), nil)
	if err != nil {
		return nil, err
	}
	var entities []chunk.Entity

	buf := bytes.NewBuffer(data)
	dec, ok := nbt.NewDecoderWithEncoding(buf, nbt.LittleEndian), false

	for buf.Len() != 0 {
		ent := chunk.Entity{Data: make(map[string]any)}
		if err := dec.Decode(&ent.Data); err != nil {
			return nil, fmt.Errorf("decode entity nbt: %w", err)
		}
		ent.ID, ok = ent.Data["UniqueID"].(int64)
		if !ok {
			db.conf.Log.Error("missing unique ID field, generating random", "data", fmt.Sprint(ent.Data))
			ent.ID = rand.Int64()
		}
		entities = append(entities, ent)
	}
	return entities, nil
}

func (db *DB) blockEntities(k dbKey) ([]chunk.BlockEntity, error) {
	var blockEntities []chunk.BlockEntity

	data, err := db.ldb.Get(k.Sum(keyBlockEntities), nil)
	if err != nil {
		return blockEntities, err
	}

	buf := bytes.NewBuffer(data)
	dec := nbt.NewDecoderWithEncoding(buf, nbt.LittleEndian)

	for buf.Len() != 0 {
		be := chunk.BlockEntity{Data: make(map[string]any)}
		if err := dec.Decode(&be.Data); err != nil {
			return blockEntities, fmt.Errorf("decode nbt: %w", err)
		}
		be.Pos = blockPosFromNBT(be.Data)
		blockEntities = append(blockEntities, be)
	}
	return blockEntities, nil
}

func (db *DB) scheduledUpdates(k dbKey) ([]chunk.ScheduledBlockUpdate, int64, error) {
	data, err := db.ldb.Get(k.Sum(keyPendingScheduledTicks), nil)
	if err != nil {
		return nil, 0, err
	}
	var m scheduledUpdates
	if err := nbt.UnmarshalEncoding(data, &m, nbt.LittleEndian); err != nil {
		return nil, 0, fmt.Errorf("read nbt: %s", err.Error())
	}
	updates := make([]chunk.ScheduledBlockUpdate, len(m.TickList))
	for i, tick := range m.TickList {
		t, _ := tick["time"].(int64)
		bl, _ := tick["blockState"].(map[string]any)
		bpe := chunk.BlockPaletteEncoding{Blocks: db.conf.Blocks}
		block, err := bpe.DecodeBlockState(bl)
		if err != nil {
			db.conf.Log.Error("read scheduled updates: decode block state: " + err.Error())
			continue
		}
		updates[i] = chunk.ScheduledBlockUpdate{Pos: blockPosFromNBT(tick), Block: block, Tick: t}
	}
	return updates, int64(m.CurrentTick), nil
}

func (db *DB) StoreColumn(pos world.ChunkPos, dim world.Dimension, col *chunk.Column) error {
	k := dbKey{pos: pos, dim: dim}
	if err := db.storeColumn(k, col); err != nil {
		return fmt.Errorf("store column %v (%v): %w", pos, dim, err)
	}
	return nil
}

func (db *DB) storeColumn(k dbKey, col *chunk.Column) error {
	data := chunk.Encode(col.Chunk, chunk.DiskEncoding)
	n := 7 + len(data.SubChunks) + len(col.Entities)
	batch := leveldb.MakeBatch(n)

	db.storeVersion(batch, k, chunkVersion)
	db.storeBiomes(batch, k, data.Biomes)
	db.storeSubChunks(batch, k, data.SubChunks, col.Chunk.Range())
	db.storeFinalisation(batch, k, finalisationPopulated)
	db.storeEntities(batch, k, col.Entities)
	db.storeBlockEntities(batch, k, col.BlockEntities)
	db.storeScheduledUpdates(batch, k, col.Tick, col.ScheduledBlocks)

	return db.ldb.Write(batch, nil)
}

func (db *DB) storeVersion(batch *leveldb.Batch, k dbKey, ver uint8) {
	batch.Put(k.Sum(keyVersion), []byte{ver})
}

var emptyHeightmap = make([]byte, 512)

func (db *DB) storeBiomes(batch *leveldb.Batch, k dbKey, biomes []byte) {
	batch.Put(k.Sum(key3DData), append(emptyHeightmap, biomes...))
}

func (db *DB) storeSubChunks(batch *leveldb.Batch, k dbKey, subChunks [][]byte, r cube.Range) {
	for i, sub := range subChunks {
		batch.Put(k.Sum(keySubChunkData, byte(i+(r[0]>>4))), sub)
	}
}

func (db *DB) storeFinalisation(batch *leveldb.Batch, k dbKey, finalisation uint32) {
	p := make([]byte, 4)
	binary.LittleEndian.PutUint32(p, finalisation)
	batch.Put(k.Sum(keyFinalisation), p)
}

func (db *DB) storeEntities(batch *leveldb.Batch, k dbKey, entities []chunk.Entity) {
	idsKey := append([]byte(keyEntityIdentifiers), index(k.pos, k.dim)...)

	var previousIDs []int64
	digpPrev, err := db.ldb.Get(idsKey, nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		db.conf.Log.Error("store entities: read chunk entity IDs: " + err.Error())
	}
	if err == nil {
		for i := 0; i < len(digpPrev); i += 8 {
			previousIDs = append(previousIDs, int64(binary.LittleEndian.Uint64(digpPrev[i:])))
		}
	}

	newIDs := make([]int64, 0, len(entities))
	for _, e := range entities {
		e.Data["UniqueID"] = e.ID
		b, err := nbt.MarshalEncoding(e.Data, nbt.LittleEndian)
		if err != nil {
			db.conf.Log.Error("store entities: encode NBT: " + err.Error())
			continue
		}
		batch.Put(entityIndex(e.ID), b)
		newIDs = append(newIDs, e.ID)
	}

	for _, uniqueID := range previousIDs {
		if !slices.Contains(newIDs, uniqueID) {
			batch.Delete(entityIndex(uniqueID))
		}
	}
	if len(entities) == 0 {
		batch.Delete(idsKey)
	} else {

		ids := make([]byte, 0, 8*len(newIDs))
		for _, uniqueID := range newIDs {
			ids = binary.LittleEndian.AppendUint64(ids, uint64(uniqueID))
		}
		batch.Put(idsKey, ids)
	}

	batch.Delete(k.Sum(keyEntitiesOld))
}

func entityIndex(id int64) []byte {
	return binary.LittleEndian.AppendUint64([]byte(keyEntity), uint64(id))
}

func (db *DB) storeBlockEntities(batch *leveldb.Batch, k dbKey, blockEntities []chunk.BlockEntity) {
	if len(blockEntities) == 0 {
		batch.Delete(k.Sum(keyBlockEntities))
		return
	}

	buf := bytes.NewBuffer(nil)
	enc := nbt.NewEncoderWithEncoding(buf, nbt.LittleEndian)
	for _, b := range blockEntities {
		b.Data["x"], b.Data["y"], b.Data["z"] = int32(b.Pos[0]), int32(b.Pos[1]), int32(b.Pos[2])
		if err := enc.Encode(b.Data); err != nil {
			db.conf.Log.Error("store block entities: encode nbt: " + err.Error())
		}
	}
	batch.Put(k.Sum(keyBlockEntities), buf.Bytes())
}

func (db *DB) storeScheduledUpdates(batch *leveldb.Batch, k dbKey, tick int64, updates []chunk.ScheduledBlockUpdate) {
	if len(updates) == 0 {
		batch.Delete(k.Sum(keyPendingScheduledTicks))
		return
	}
	list := make([]map[string]any, len(updates))
	bpe := chunk.BlockPaletteEncoding{Blocks: db.conf.Blocks}
	for i, update := range updates {
		list[i] = map[string]any{
			"x": int32(update.Pos[0]), "y": int32(update.Pos[1]), "z": int32(update.Pos[2]),
			"time": update.Tick, "blockState": bpe.EncodeBlockState(update.Block),
		}
	}
	b, err := nbt.MarshalEncoding(scheduledUpdates{CurrentTick: int32(tick), TickList: list}, nbt.LittleEndian)
	if err != nil {
		db.conf.Log.Error("store scheduled updates: encode nbt: " + err.Error())
		return
	}
	batch.Put(k.Sum(keyPendingScheduledTicks), b)
}

type scheduledUpdates struct {
	CurrentTick int32            `nbt:"currentTick"`
	TickList    []map[string]any `nbt:"tickList"`
}

func (db *DB) NewColumnIterator(r *IteratorRange) *ColumnIterator {
	if r == nil {
		r = &IteratorRange{}
	}
	return newColumnIterator(db, r)
}

func (db *DB) Close() error {
	db.ldat.LastPlayed = time.Now().Unix()

	var ldat leveldat.LevelDat
	if err := ldat.Marshal(*db.ldat); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := ldat.WriteFile(filepath.Join(db.dir, "level.dat")); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := os.WriteFile(filepath.Join(db.dir, "levelname.txt"), []byte(db.ldat.LevelName), 0644); err != nil {
		return fmt.Errorf("close: write levelname.txt: %w", err)
	}
	return db.ldb.Close()
}

type dbKey struct {
	pos world.ChunkPos
	dim world.Dimension
}

func (k dbKey) Sum(p ...byte) []byte {
	return append(index(k.pos, k.dim), p...)
}

func index(position world.ChunkPos, d world.Dimension) []byte {
	dim, _ := world.DimensionID(d)
	x, z := uint32(position[0]), uint32(position[1])
	b := make([]byte, 12)

	binary.LittleEndian.PutUint32(b, x)
	binary.LittleEndian.PutUint32(b[4:], z)
	if dim == 0 {
		return b[:8]
	}
	binary.LittleEndian.PutUint32(b[8:], uint32(dim))
	return b
}

func blockPosFromNBT(data map[string]any) cube.Pos {
	x, _ := data["x"].(int32)
	y, _ := data["y"].(int32)
	z, _ := data["z"].(int32)
	return cube.Pos{int(x), int(y), int(z)}
}
