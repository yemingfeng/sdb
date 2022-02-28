package service

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/protobuf/proto"
	"math"
)

var sortedSetCollection = collection.NewCollection(pb.DataType_SORTED_SET)

func newSortedSetIndexes(score []byte, value []byte) []collection.Index {
	return []collection.Index{
		{Name: []byte("score"), Value: score},
		{Name: []byte("value"), Value: value},
	}
}

func ZPush(key []byte, tuples []*pb.Tuple) error {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	batch := collection.NewBatch()
	defer batch.Close()

	// tuples -> [ {value: a, score: 1.0}, {value:b, score:1.1}, {value: c, score: 0.9} ]
	for _, tuple := range tuples {
		score := []byte(fmt.Sprintf("%32.32f", tuple.Score))
		value, err := proto.Marshal(tuple)
		if err != nil {
			return err
		}
		if err := sortedSetCollection.UpsertRow(&collection.Row{
			Key:     key,
			Id:      tuple.Value,
			Indexes: newSortedSetIndexes(score, tuple.Value),
			Value:   value,
		}, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func ZPop(key []byte, values [][]byte) error {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	batch := collection.NewBatch()
	defer batch.Close()

	for _, value := range values {
		rows, err := sortedSetCollection.IndexValuePage(key, []byte("value"), value, 0, math.MaxUint32)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if err := sortedSetCollection.DelRowById(key, row.Id, batch); err != nil {
				return err
			}
		}
	}
	return batch.Commit()
}

func ZRange(key []byte, offset int32, limit uint32) ([]*pb.Tuple, error) {
	rows, err := sortedSetCollection.IndexPage(key, []byte("score"), offset, limit)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Tuple, len(rows))
	for i := range rows {
		var tuple pb.Tuple
		if err := proto.Unmarshal(rows[i].Value, &tuple); err != nil {
			return nil, err
		}
		res[i] = &tuple
	}
	return res, nil
}

func ZExist(key []byte, values [][]byte) ([]bool, error) {
	res := make([]bool, len(values))
	for i := range values {
		rows, err := sortedSetCollection.IndexValuePage(key, []byte("value"), values[i], 0, math.MaxUint32)
		if err != nil {
			return nil, err
		}
		res[i] = len(rows) > 0
	}
	return res, nil
}

func ZDel(key []byte) error {
	lock(LSortedSet, key)
	defer unlock(LSortedSet, key)

	batch := collection.NewBatch()
	defer batch.Close()

	if err := sortedSetCollection.DelAll(key, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func ZCount(key []byte) (uint32, error) {
	return sortedSetCollection.Count(key)
}

func ZMembers(key []byte) ([]*pb.Tuple, error) {
	rows, err := sortedSetCollection.IndexPage(key, []byte("score"), 0, math.MaxUint32)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Tuple, len(rows))
	for i := range rows {
		var tuple pb.Tuple
		if err := proto.Unmarshal(rows[i].Value, &tuple); err != nil {
			return nil, err
		}
		res[i] = &tuple
	}
	return res, nil
}
