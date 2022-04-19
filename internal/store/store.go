package store

import (
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/util"
)

type Store interface {
	NewBatch() Batch
	Close() error
}

type Batch interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	Del(key []byte) error
	Iterate(opt *PrefixIteratorOption, handle func([]byte, []byte) error) error
	Commit() error
	ApplyCommit() error
	Close()
}

type PrefixIteratorOption struct {
	Prefix []byte

	Offset int32
	Limit  uint32
}

const (
	PEBBLE string = "pebble"
	BADGER string = "badger"
	LEVEL  string = "level"
)

var storeLogger = util.GetLogger("store")
var store Store

func init() {
	if conf.Conf.Store.Engine == PEBBLE {
		store = NewPebbleStore()
	} else if conf.Conf.Store.Engine == BADGER {
		store = NewBadgerStore()
	} else if conf.Conf.Store.Engine == LEVEL {
		store = NewLevelStore()
	} else {
		storeLogger.Fatalf("not match store engine: %s", conf.Conf.Store.Engine)
	}
}

func NewBatch() Batch {
	return store.NewBatch()
}

func Stop() {
	if store != nil {
		if err := store.Close(); err != nil {
			storeLogger.Printf("shutdown store error: %+v", err)
		}
		storeLogger.Println("stop store finished")
	}
}
