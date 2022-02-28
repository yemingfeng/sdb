package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/util"
	"math"
)

var listCollection = collection.NewCollection(pb.DataType_LIST)

func newListIndexes(score []byte, value []byte) []collection.Index {
	return []collection.Index{
		{Name: []byte("score"), Value: score},
		{Name: []byte("value"), Value: value},
	}
}

func LRPush(key []byte, values [][]byte) error {
	lock(LList, key)
	defer unlock(LList, key)

	batch := collection.NewBatch()
	defer batch.Close()

	for _, value := range values {
		score := []byte(fmt.Sprintf("%d", util.GetOrderingKey()))
		id := []byte(string(value) + ":" + string(score))
		if err := listCollection.UpsertRow(&collection.Row{
			Key:     key,
			Id:      id,
			Indexes: newListIndexes(score, value),
			Value:   value,
		}, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func LLPush(key []byte, values [][]byte) error {
	lock(LList, key)
	defer unlock(LList, key)

	batch := collection.NewBatch()
	defer batch.Close()

	for i, value := range values {
		score := []byte(fmt.Sprintf("%d", -util.GetOrderingKey()-int64(i)))
		id := []byte(string(value) + ":" + string(score))
		if err := listCollection.UpsertRow(&collection.Row{
			Key:     key,
			Id:      id,
			Indexes: newListIndexes(score, value),
			Value:   value,
		}, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func LPop(key []byte, values [][]byte) error {
	lock(LList, key)
	defer unlock(LList, key)

	batch := collection.NewBatch()
	defer batch.Close()

	for i := range values {
		rows, err := listCollection.IndexValuePage(key, []byte("value"), values[i], 0, math.MaxUint32)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if err := listCollection.DelRowById(key, row.Id, batch); err != nil {
				return err
			}
		}
	}
	return batch.Commit()
}

func LRange(key []byte, offset int32, limit uint32) ([][]byte, error) {
	rows, err := listCollection.IndexPage(key, []byte("score"), offset, limit)
	if err != nil {
		return nil, err
	}
	res := make([][]byte, len(rows))
	for i := range rows {
		res[i] = rows[i].Value
	}
	return res, nil
}

func LExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i := range values {
		rows, err := listCollection.IndexValuePage(key, []byte("value"), values[i], 0, 1)
		if err != nil {
			return nil, err
		}
		res[i] = len(rows) > 0
	}
	return res, nil
}

func LDel(key []byte) error {
	lock(LList, key)
	defer unlock(LList, key)

	batch := collection.NewBatch()
	defer batch.Close()

	if err := listCollection.DelAll(key, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func LCount(key []byte) (uint32, error) {
	return listCollection.Count(key)
}

func LMembers(key []byte) ([][]byte, error) {
	rows, err := listCollection.IndexPage(key, []byte("score"), 0, math.MaxUint32)
	if err != nil {
		return nil, err
	}
	res := make([][]byte, len(rows))
	for i := range rows {
		res[i] = rows[i].Value
	}
	return res, nil
}
