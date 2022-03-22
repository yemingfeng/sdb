package service

import (
	"github.com/yemingfeng/sdb/internal/store"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"math"
)

var setCollection = store.NewCollection(pb.DataType_SET)

func SPush(key []byte, values [][]byte) error {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		if err := setCollection.UpsertRow(&store.Row{
			Key:   key,
			Id:    value,
			Value: value,
		}, batch); err != nil {
			return err
		}
		if err := PAdd(pb.DataType_SET, key, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func SPop(key []byte, values [][]byte) error {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		if err := setCollection.DelRowById(key, value, batch); err != nil {
			return err
		}
	}
	// delete if not element at key
	rows, err := setCollection.PageWithBatch(key, 0, 1, batch)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		if err := PDel(pb.DataType_SET, key, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func SExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i, value := range values {
		exist, err := setCollection.ExistRowById(key, value)
		if err != nil {
			return nil, err
		}
		res[i] = exist
	}
	return res, nil
}

func SDel(key []byte) error {
	lock(LSet, key)
	defer unlock(LSet, key)

	batch := store.NewBatch()
	defer batch.Close()

	if err := setCollection.DelAll(key, batch); err != nil {
		return err
	}
	if err := PDel(pb.DataType_SET, key, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func SCount(key []byte) (uint32, error) {
	return setCollection.Count(key)
}

func SMembers(key []byte) ([][]byte, error) {
	rows, err := setCollection.Page(key, 0, math.MaxUint32)
	if err != nil {
		return nil, err
	}
	res := make([][]byte, len(rows))
	for i := range rows {
		res[i] = rows[i].Value
	}
	return res, nil
}
