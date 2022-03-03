package level

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/engine"
	"log"
)

type LevelStore struct {
	db *leveldb.DB
}

func NewLevelStore() *LevelStore {
	dbPath := conf.Conf.Store.Path + "/level"
	db, err := leveldb.OpenFile(dbPath, &opt.Options{Filter: filter.NewBloomFilter(10)})
	if err != nil {
		log.Fatalf("failed to open file: %+v", err)
	}
	log.Printf("db init %s complete", dbPath)

	return &LevelStore{db: db}
}

func (store *LevelStore) NewBatch() engine.Batch {
	transaction, _ := store.db.OpenTransaction()
	return &LevelBatch{db: store.db, transaction: transaction}
}

func (store *LevelStore) Close() error {
	return store.db.Close()
}
