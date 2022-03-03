package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/engine"
	"log"
)

type PebbleStore struct {
	db *pebble.DB
}

func NewPebbleStore() *PebbleStore {
	dbPath := conf.Conf.Store.Path + "/pebble"
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		log.Fatalf("failed to open file: %+v", err)
	}
	log.Printf("db init %s complete", dbPath)

	return &PebbleStore{db: db}
}

func (store *PebbleStore) NewBatch() engine.Batch {
	return &PebbleBatch{batch: store.db.NewIndexedBatch()}
}

func (store *PebbleStore) Close() error {
	return store.db.Close()
}
