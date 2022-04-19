package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/yemingfeng/sdb/internal/conf"
	util2 "github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
)

var levelLogger = util2.GetLogger("level")

type LevelStore struct {
	db *leveldb.DB
}

func NewLevelStore() *LevelStore {
	dbPath := conf.Conf.Store.Path + "/level"
	db, err := leveldb.OpenFile(dbPath, &opt.Options{Filter: filter.NewBloomFilter(10)})
	if err != nil {
		levelLogger.Fatalf("failed to open file: %+v", err)
	}
	levelLogger.Printf("db init %s complete", dbPath)

	return &LevelStore{db: db}
}

func (store *LevelStore) NewBatch() Batch {
	transaction, _ := store.db.OpenTransaction()
	return &LevelBatch{transaction: transaction, log: &pb.Log{LogEntries: make([]*pb.LogEntry, 0)}}
}

func (store *LevelStore) Close() error {
	return store.db.Close()
}

type LevelBatch struct {
	transaction *leveldb.Transaction
	log         *pb.Log
}

func (batch *LevelBatch) Get(key []byte) ([]byte, error) {
	value, err := batch.transaction.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return value, err
}

func (batch *LevelBatch) Set(key []byte, value []byte) error {
	batch.log.LogEntries = append(batch.log.LogEntries, &pb.LogEntry{Op: pb.Op_OP_SET, Key: key, Value: value})
	return batch.transaction.Put(key, value, &opt.WriteOptions{Sync: true})
}

func (batch *LevelBatch) Del(key []byte) error {
	batch.log.LogEntries = append(batch.log.LogEntries, &pb.LogEntry{Op: pb.Op_OP_DEL, Key: key})
	return batch.transaction.Delete(key, &opt.WriteOptions{Sync: true})
}

func (batch *LevelBatch) Iterate(opt *PrefixIteratorOption, handle func([]byte, []byte) error) error {
	it := batch.transaction.NewIterator(util.BytesPrefix(opt.Prefix), nil)
	defer func() {
		it.Release()
	}()

	if opt.Offset >= 0 {
		i := 0
		for it.First(); i < int(opt.Offset) && it.Valid(); it.Next() {
			i++
		}

		i = 0
		for ; it.Valid(); it.Next() {
			err := handle(util2.Copy2(it.Key()), util2.Copy2(it.Value()))
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
			err := handle(util2.Copy2(it.Key()), util2.Copy2(it.Value()))
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

func (batch *LevelBatch) Commit() error {
	return Apply(batch.log)
}

func (batch *LevelBatch) ApplyCommit() error {
	return batch.transaction.Commit()
}

func (batch *LevelBatch) Close() {
	batch.transaction.Discard()
}
