package service

import (
	"errors"
	"github.com/axiomhq/hyperloglog"
	"github.com/yemingfeng/sdb/internal/store"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
)

var NotFoundHyperLogLogError = errors.New("not found hyper log log, please create it")
var HyperLogLogExistError = errors.New("hyper log log exist, please delete it or change other")

var hyperLogLogCollection = store.NewCollection(pb.DataType_HYPER_LOG_LOG)

func HLLCreate(key []byte) error {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	batch := store.NewBatch()
	defer batch.Close()

	exist, err := hyperLogLogCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if exist {
		return HyperLogLogExistError
	}

	h := hyperloglog.New16()
	value, err := h.MarshalBinary()
	if err != nil {
		return err
	}

	if err := hyperLogLogCollection.UpsertRow(&store.Row{
		Key:   key,
		Id:    key,
		Value: value,
	}, batch); err != nil {
		return err
	}

	if err := PAdd(pb.DataType_HYPER_LOG_LOG, key, batch); err != nil {
		return err
	}

	return batch.Commit()
}

func HLLDel(key []byte) error {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	batch := store.NewBatch()
	defer batch.Close()

	if err := hyperLogLogCollection.DelRowById(key, key, batch); err != nil {
		return err
	}

	if err := PDel(pb.DataType_HYPER_LOG_LOG, key, batch); err != nil {
		return err
	}

	return batch.Commit()
}

func HLLAdd(key []byte, values [][]byte) error {
	lock(LHyperLogLog, key)
	defer unlock(LHyperLogLog, key)

	batch := store.NewBatch()
	defer batch.Close()

	row, err := hyperLogLogCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	if row == nil {
		return NotFoundHyperLogLogError
	}
	var hll hyperloglog.Sketch
	err = hll.UnmarshalBinary(row.Value)
	if err != nil {
		return err
	}

	for _, value := range values {
		hll.Insert(value)
	}

	value, err := hll.MarshalBinary()
	if err != nil {
		return err
	}
	if err := hyperLogLogCollection.UpsertRow(&store.Row{
		Key:   key,
		Id:    key,
		Value: value,
	}, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func HLLCount(key []byte) (uint32, error) {
	row, err := hyperLogLogCollection.GetRowById(key, key)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, NotFoundHyperLogLogError
	}
	var hll hyperloglog.Sketch
	err = hll.UnmarshalBinary(row.Value)
	if err != nil {
		return 0, err
	}
	return uint32(hll.Estimate()), nil
}
