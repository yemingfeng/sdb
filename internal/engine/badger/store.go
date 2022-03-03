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

func (store *BadgerStore) NewBatch() engine.Batch {
	return &BadgerBatch{db: store.db, transaction: store.db.NewTransaction(true)}
}

func (store *BadgerStore) Close() error {
	return store.db.Close()
}
