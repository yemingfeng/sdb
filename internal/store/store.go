package store

import (
	"github.com/yemingfeng/sdb/internal/conf"
	"log"
)

const (
	PEBBLE string = "pebble"
	BADGER string = "badger"
	LEVEL  string = "level"
)

var store Store

func init() {
	if conf.Conf.Store.Engine == PEBBLE {
		store = NewPebbleStore()
	} else if conf.Conf.Store.Engine == BADGER {
		store = NewBadgerStore()
	} else if conf.Conf.Store.Engine == LEVEL {
		store = NewLevelStore()
	} else {
		log.Fatalf("not match store engine: %s", conf.Conf.Store.Engine)
	}
}

func NewBatch() Batch {
	return store.NewBatch()
}

func Stop() {
	if store != nil {
		if err := store.Close(); err != nil {
			log.Printf("shutdown store error: %+v", err)
		}
		log.Println("stop store finished")
	}
}
