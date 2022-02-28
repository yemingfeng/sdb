package pebble

import "github.com/cockroachdb/pebble"

type PebbleBatch struct {
	batch *pebble.Batch
}

func (batch *PebbleBatch) Get(key []byte) ([]byte, error) {
	value, closer, err := batch.batch.Get(key)
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

func (batch *PebbleBatch) Set(key []byte, value []byte) error {
	return batch.batch.Set(key, value, nil)
}

func (batch *PebbleBatch) Del(key []byte) error {
	return batch.batch.Delete(key, nil)
}

func (batch *PebbleBatch) Commit() error {
	return batch.batch.Commit(pebble.Sync)
}

func (batch *PebbleBatch) Reset() {
	batch.batch.Reset()
}

func (batch *PebbleBatch) Close() {
	_ = batch.batch.Close()
}
