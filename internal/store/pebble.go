package store

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
)

var pebbleLogger = util.GetLogger("pebble")

type PebbleStore struct {
	db *pebble.DB
}

func NewPebbleStore() *PebbleStore {
	dbPath := conf.Conf.Store.Path + "/pebble"
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		pebbleLogger.Fatalf("failed to open file: %+v", err)
	}
	pebbleLogger.Printf("db init %s complete", dbPath)

	return &PebbleStore{db: db}
}

func (store *PebbleStore) NewBatch() Batch {
	return &PebbleBatch{batch: store.db.NewIndexedBatch(), log: &pb.Log{LogEntries: make([]*pb.LogEntry, 0)}}
}

func (store *PebbleStore) Close() error {
	return store.db.Close()
}

type PebbleBatch struct {
	batch *pebble.Batch
	log   *pb.Log
}

func (batch *PebbleBatch) Get(key []byte) ([]byte, error) {
	value, closer, err := batch.batch.Get(key)
	if err == pebble.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err = closer.Close(); err != nil {
		return nil, err
	}
	return value, err
}

func (batch *PebbleBatch) Set(key []byte, value []byte) error {
	batch.log.LogEntries = append(batch.log.LogEntries, &pb.LogEntry{Op: pb.Op_OP_SET, Key: key, Value: value})
	return batch.batch.Set(key, value, nil)
}

func (batch *PebbleBatch) Del(key []byte) error {
	batch.log.LogEntries = append(batch.log.LogEntries, &pb.LogEntry{Op: pb.Op_OP_DEL, Key: key})
	return batch.batch.Delete(key, nil)
}

func (batch *PebbleBatch) Iterate(opt *PrefixIteratorOption, handle func([]byte, []byte) error) error {
	keyUpperBound := func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil
	}

	var it = batch.batch.NewIter(&pebble.IterOptions{
		LowerBound: opt.Prefix,
		UpperBound: keyUpperBound(opt.Prefix),
	})
	defer func() {
		_ = it.Close()
	}()

	if opt.Offset >= 0 {
		i := 0
		for it.First(); i < int(opt.Offset) && it.Valid(); it.Next() {
			i++
		}

		i = 0
		for ; it.Valid(); it.Next() {
			err := handle(util.Copy2(it.Key()), util.Copy2(it.Value()))
			if err != nil {
				return err
			}
			i++
			if opt.Limit > 0 && i == int(opt.Limit) {
				break
			}
		}
	} else {
		i := 0
		for it.Last(); i < int(-opt.Offset-1) && it.Valid(); it.Prev() {
			i++
		}

		i = 0
		for ; it.Valid(); it.Prev() {
			err := handle(util.Copy2(it.Key()), util.Copy2(it.Value()))
			if err != nil {
				return err
			}
			i++
			if opt.Limit > 0 && i == int(opt.Limit) {
				break
			}
		}
	}
	return nil
}

func (batch *PebbleBatch) Commit() error {
	return Apply(batch.log)
}

func (batch *PebbleBatch) ApplyCommit() error {
	return batch.batch.Commit(pebble.Sync)
}

func (batch *PebbleBatch) Close() {
	_ = batch.batch.Close()
}
