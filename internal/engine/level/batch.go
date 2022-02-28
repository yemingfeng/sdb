package level

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
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

func (batch *LevelBatch) Reset() {
	batch.transaction.Discard()
	transaction, _ := batch.db.OpenTransaction()
	batch.transaction = transaction
}

func (batch *LevelBatch) Close() {
	batch.transaction.Discard()
}
