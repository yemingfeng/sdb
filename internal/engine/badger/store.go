package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/engine"
	"log"
)

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore() *BadgerStore {
	dbPath := conf.Conf.Store.Path + "/badger"
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithSyncWrites(true))
	if err != nil {
		log.Fatalf("failed to open file: %+v", err)
	}
	log.Printf("db init %s complete", dbPath)

	return &BadgerStore{db: db}
}

func (store *BadgerStore) Get(key []byte) ([]byte, error) {
	txn := store.db.NewTransaction(true)
	item, err := txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (store *BadgerStore) NewBatch() engine.Batch {
	return &BadgerBatch{db: store.db, transaction: store.db.NewTransaction(true)}
}

func (store *BadgerStore) Iterate(opt *engine.PrefixIteratorOption, handle func([]byte, []byte) error) error {
	return store.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			Reverse:        opt.Offset < 0,
			PrefetchSize:   10,
			PrefetchValues: true})
		defer it.Close()

		prefix := opt.Prefix
		if opt.Offset < 0 {
			opt.Offset = -opt.Offset - 1
			prefix = append(opt.Prefix, 0xFF)
		}

		i := 0
		for it.Seek(prefix); i < int(opt.Offset) && it.ValidForPrefix(opt.Prefix); it.Next() {
			i++
		}

		i = 0
		for ; it.ValidForPrefix(opt.Prefix); it.Next() {
			err := it.Item().Value(func(value []byte) error {
				return handle(it.Item().Key(), value)
			})
			if err != nil {
				return err
			}
			i++
			if opt.Limit > 0 && i == int(opt.Limit) {
				break
			}
		}

		return nil
	})
}

func (store *BadgerStore) Close() error {
	return store.db.Close()
}
