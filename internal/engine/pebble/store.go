package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/util"
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

func (store *PebbleStore) Get(key []byte) ([]byte, error) {
	value, closer, err := store.db.Get(key)
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

func (store *PebbleStore) NewBatch() engine.Batch {
	return &PebbleBatch{batch: store.db.NewIndexedBatch()}
}

func (store *PebbleStore) Iterate(opt *engine.PrefixIteratorOption, handle func([]byte, []byte) error) error {
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

	var it = store.db.NewIter(&pebble.IterOptions{
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

func (store *PebbleStore) Close() error {
	return store.db.Close()
}
