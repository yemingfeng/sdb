package collection

import (
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/engine/badger"
	"github.com/yemingfeng/sdb/internal/engine/level"
	"github.com/yemingfeng/sdb/internal/engine/pebble"
	"log"
)

const (
	PEBBLE string = "pebble"
	BADGER string = "badger"
	LEVEL  string = "level"
)

var store engine.Store

func init() {
	if conf.Conf.Store.Engine == PEBBLE {
		store = pebble.NewPebbleStore()
	} else if conf.Conf.Store.Engine == BADGER {
		store = badger.NewBadgerStore()
	} else if conf.Conf.Store.Engine == LEVEL {
		store = level.NewLevelStore()
	} else {
		log.Fatalf("not match store engine: %s", conf.Conf.Store.Engine)
	}
}

func Get(key []byte) ([]byte, error) {
	return store.Get(key)
}

func NewBatch() engine.Batch {
	return store.NewBatch()
}

func Iterate(prefix []byte, offset int32, limit uint32, handle func([]byte, []byte) error) error {
	return store.Iterate(&engine.PrefixIteratorOption{Prefix: prefix, Offset: offset, Limit: limit}, handle)
}

//// Close todo call
//func Close() error {
//	return store.Close()
//}
