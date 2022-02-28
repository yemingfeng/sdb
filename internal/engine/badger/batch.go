package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type BadgerBatch struct {
	db          *badger.DB
	transaction *badger.Txn
}

func (batch *BadgerBatch) Get(key []byte) ([]byte, error) {
	item, err := batch.transaction.Get(key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (batch *BadgerBatch) Set(key []byte, value []byte) error {
	return batch.transaction.Set(key, value)
}

func (batch *BadgerBatch) Del(key []byte) error {
	return batch.transaction.Delete(key)
}

func (batch *BadgerBatch) Commit() error {
	return batch.transaction.Commit()
}

func (batch *BadgerBatch) Reset() {
	batch.transaction.Discard()
	batch.transaction = batch.db.NewTransaction(true)
}

func (batch *BadgerBatch) Close() {
	batch.transaction.Discard()
}
