package store

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/config"
	"github.com/yemingfeng/sdb/internal/util"
	"os"
)

var storeLogger = util.GetLogger("store")

func NewStore(config *config.Config) *pebble.DB {
	dbPath := config.Store.Path
	if config.Store.DeleteOld {
		if err := os.RemoveAll(dbPath); err != nil {
			storeLogger.Fatalf("delete old error: %+v", err)
		}
	}
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		storeLogger.Fatalf("failed to open file: %+v", err)
	}
	storeLogger.Printf("db init %s complete", dbPath)

	return db
}

func NewPrefixIterOptions(prefix []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: Next(prefix),
	}
}

func Next(key []byte) []byte {
	upperBound := func(b []byte) []byte {
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
	return upperBound(key)
}
