package test

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/config"
	"github.com/yemingfeng/sdb/internal/store"
	"testing"
)

func TestStore(t *testing.T) {
	conf := config.NewTestConfig()
	s := store.NewStore(conf)

	iter := s.NewIter(&pebble.IterOptions{})
	defer iter.Close()

	lowerBound := []byte("as:aDE=:0000000000000000000000000000000000000000000000000000000000000004:0000000000000000000000000000000000000000000000000000002147483651:")
	upperBound := []byte("as:aDE=:0000000000000000000000000000000000000000000000000000000000000004:0000000000000000000000000000000000000000000000000000002147483658:")
	iter.SetBounds(lowerBound, upperBound)

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		t.Logf("key=%s", key)
	}
}
