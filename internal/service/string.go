package service

import (
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/pb"
	"strconv"
)

var stringCollection = collection.NewCollection(pb.DataType_STRING)

func Set(key []byte, value []byte) error {
	batch := collection.NewBatch()
	defer batch.Close()
	if err := stringCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value}, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func MSet(keys [][]byte, values [][]byte) error {
	batch := collection.NewBatch()
	defer batch.Close()
	for i := range keys {
		if err := stringCollection.UpsertRow(&collection.Row{
			Key:   keys[i],
			Id:    keys[i],
			Value: values[i],
		}, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func SetNX(key []byte, value []byte) error {
	batch := collection.NewBatch()
	defer batch.Close()

	exist, err := stringCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if exist {
		return err
	}
	if err := stringCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value}, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func Get(key []byte) ([]byte, error) {
	row, err := stringCollection.GetRowById(key, key)
	if err != nil || row == nil {
		return nil, err
	}
	return row.Value, nil
}

func MGet(keys [][]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	for i := range keys {
		row, err := stringCollection.GetRowById(keys[i], keys[i])
		if err != nil {
			return nil, err
		}
		if row != nil {
			values[i] = row.Value
		}
	}
	return values, nil
}

func Del(key []byte) error {
	batch := collection.NewBatch()
	defer batch.Close()
	if err := stringCollection.DelRowById(key, key, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func Incr(key []byte, delta int32) error {
	lock(LString, key)
	defer unlock(LString, key)

	batch := collection.NewBatch()
	defer batch.Close()

	row, err := stringCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	var valueInt = 0
	if row != nil {
		valueInt, err = strconv.Atoi(string(row.Value))
		if err != nil {
			return err
		}
	}
	valueInt = valueInt + int(delta)

	if err := stringCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: []byte(strconv.Itoa(valueInt))}, batch); err != nil {
		return err
	}
	return batch.Commit()
}
