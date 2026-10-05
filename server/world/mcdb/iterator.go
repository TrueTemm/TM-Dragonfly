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
	"encoding/binary"
	"fmt"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/goleveldb/leveldb/iterator"
)

type ColumnIterator struct {
	dbIter iterator.Iterator
	db     *DB
	r      *IteratorRange

	err error

	current *chunk.Column
	pos     world.ChunkPos
	dim     world.Dimension
	seen    map[dbKey]struct{}
}

func newColumnIterator(db *DB, r *IteratorRange) *ColumnIterator {
	return &ColumnIterator{
		db:     db,
		dbIter: db.ldb.NewIterator(nil, nil),
		seen:   make(map[dbKey]struct{}),
		r:      r,
	}
}

func (iter *ColumnIterator) Next() bool {
	if iter.err != nil || !iter.dbIter.Next() {
		iter.current = nil
		iter.dim = nil
		return false
	}
	k := iter.dbIter.Key()
	kLen := len(k)
	if (kLen != 9 && kLen != 13) || (k[kLen-1] != keyVersion && k[kLen-1] != keyVersionOld) {
		return iter.Next()
	}
	iter.dim = world.Dimension(world.Overworld)
	if kLen > 9 {
		var ok bool
		id := int(binary.LittleEndian.Uint32(k[8:12]))
		if iter.dim, ok = world.DimensionByID(id); !ok {
			iter.err = fmt.Errorf("unknown dimension id %v", id)
			return false
		}
	}
	iter.pos = world.ChunkPos{
		int32(binary.LittleEndian.Uint32(k[:4])),
		int32(binary.LittleEndian.Uint32(k[4:8])),
	}
	if !iter.r.within(iter.pos, iter.dim) {
		return iter.Next()
	}
	key := dbKey{dim: iter.dim, pos: iter.pos}
	if _, ok := iter.seen[key]; ok {

		return iter.Next()
	}
	iter.current, iter.err = iter.db.LoadColumn(iter.pos, iter.dim)
	if iter.err != nil {
		iter.err = fmt.Errorf("load chunk %v: %w", iter.pos, iter.err)
		return false
	}
	iter.seen[key] = struct{}{}
	return true
}

func (iter *ColumnIterator) Column() *chunk.Column {
	return iter.current
}

func (iter *ColumnIterator) Position() world.ChunkPos {
	return iter.pos
}

func (iter *ColumnIterator) Dimension() world.Dimension {
	return iter.dim
}

func (iter *ColumnIterator) Release() {
	iter.dbIter.Release()
}

func (iter *ColumnIterator) Error() error {
	return iter.err
}

type IteratorRange struct {
	Min, Max world.ChunkPos

	Dimension world.Dimension
}

func (r *IteratorRange) within(pos world.ChunkPos, dim world.Dimension) bool {
	if dim != r.Dimension && r.Dimension != nil {
		return false
	}
	return ((r.Min == world.ChunkPos{}) && (r.Max == world.ChunkPos{})) ||
		pos[0] >= r.Min[0] && pos[0] < r.Max[0] && pos[1] >= r.Min[1] && pos[1] < r.Max[1]
}
