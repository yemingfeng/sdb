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

func NewBatch() engine.Batch {
	return store.NewBatch()
}

//// Close todo call
//func Close() error {
//	return store.Close()
//}
