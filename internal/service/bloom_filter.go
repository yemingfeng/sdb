package service

import (
	"errors"
	"github.com/devopsfaith/bloomfilter"
	bloomFilter2 "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/pb"
)

var NotFoundBloomFilterError = errors.New("not found bloom filter, please create it")
var BloomFilterExistError = errors.New("bloom filter exist, please delete it or change other")

var bloomFilterCollection = collection.NewCollection(pb.DataType_BLOOM_FILTER)

func BFCreate(key []byte, n uint32, p float64) error {
	lock(LBloomFilter, key)
	defer unlock(LBloomFilter, key)

	batch := collection.NewBatch()
	defer batch.Close()

	exist, err := bloomFilterCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if exist {
		return BloomFilterExistError
	}
	bloomFilter := bloomFilter2.New(
		bloomfilter.Config{N: uint(n), P: p, HashName: bloomfilter.HASHER_DEFAULT})
	value, err := bloomFilter.MarshalBinary()
	if err != nil {
		return err
	}

	if err := bloomFilterCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value}, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func BFDel(key []byte) error {
	lock(LBloomFilter, key)
	defer unlock(LBloomFilter, key)

	batch := collection.NewBatch()
	defer batch.Close()

	if err := bloomFilterCollection.DelRowById(key, key, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func BFAdd(key []byte, values [][]byte) error {
	lock(LBloomFilter, key)
	defer unlock(LBloomFilter, key)

	batch := collection.NewBatch()
	defer batch.Close()

	row, err := bloomFilterCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	if row == nil {
		return NotFoundBloomFilterError
	}

	bloomFilter := &bloomFilter2.Bloomfilter{}
	if err = bloomFilter.UnmarshalBinary(row.Value); err != nil {
		return err
	}

	for _, value := range values {
		bloomFilter.Add(value)
	}

	value, err := bloomFilter.MarshalBinary()
	if err := bloomFilterCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: value}, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func BFExist(key []byte, values [][]byte) ([]bool, error) {
	row, err := bloomFilterCollection.GetRowById(key, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, NotFoundBloomFilterError
	}

	bloomFilter := &bloomFilter2.Bloomfilter{}
	err = bloomFilter.UnmarshalBinary(row.Value)
	if err != nil {
		return nil, err
	}

	res := make([]bool, len(values))
	for i, value := range values {
		res[i] = bloomFilter.Check(value)
	}

	return res, nil
}
