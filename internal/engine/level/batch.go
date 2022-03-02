package level

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/yemingfeng/sdb/internal/engine"
	util2 "github.com/yemingfeng/sdb/internal/util"
)

type LevelBatch struct {
	db          *leveldb.DB
	transaction *leveldb.Transaction
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
	return batch.transaction.Put(key, value, &opt.WriteOptions{Sync: true})
}

func (batch *LevelBatch) Del(key []byte) error {
	return batch.transaction.Delete(key, &opt.WriteOptions{Sync: true})
}

func (batch *LevelBatch) Commit() error {
	return batch.transaction.Commit()
}

func (batch *LevelBatch) Iterate(opt *engine.PrefixIteratorOption, handle func([]byte, []byte) error) error {
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

func (batch *LevelBatch) Reset() {
	batch.transaction.Discard()
	transaction, _ := batch.db.OpenTransaction()
	batch.transaction = transaction
}

func (batch *LevelBatch) Close() {
	batch.transaction.Discard()
}
