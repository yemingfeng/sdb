package service

import (
	"github.com/tmthrgd/go-bitset"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"math"
)

var bitsetItemSize = uint32(4096)
var bitsetCollection = store.NewCollection(pb.DataType_BITSET)

func BSDel(key []byte) error {
	lock(LBitset, key)
	defer unlock(LBitset, key)

	batch := store.NewBatch()
	defer batch.Close()

	if err := bitsetCollection.DelAll(key, batch); err != nil {
		return err
	}

	if err := PDel(pb.DataType_BITSET, key, batch); err != nil {
		return err
	}

	return batch.Commit()
}

func BSSetRange(key []byte, start uint32, end uint32, value bool) error {
	lock(LBitset, key)
	defer unlock(LBitset, key)

	batch := store.NewBatch()
	defer batch.Close()

	bsMap := make(map[uint32]bitset.Bitset)
	for i := start; i < end; i++ {
		id := i / bitsetItemSize
		if bsMap[id] == nil {
			row, err := bitsetCollection.GetRowById(key, util.UInt32ToBytes(id))
			if err != nil {
				return err
			}
			if row != nil {
				bsMap[id] = row.Value
			} else {
				bsMap[id] = bitset.New(uint(bitsetItemSize))
			}
		}
		bsMap[id].SetTo(uint(i%bitsetItemSize), value)
	}

	for id, bs := range bsMap {
		if err := bitsetCollection.UpsertRow(&store.Row{
			Key:   key,
			Id:    util.UInt32ToBytes(id),
			Value: bs,
		}, batch); err != nil {
			return err
		}
	}

	return batch.Commit()
}

func BSMSet(key []byte, bits []uint32, value bool) error {
	lock(LBitset, key)
	defer unlock(LBitset, key)

	batch := store.NewBatch()
	defer batch.Close()

	bsMap := make(map[uint32]bitset.Bitset)
	for i := 0; i < len(bits); i++ {
		id := bits[i] / bitsetItemSize
		if bsMap[id] == nil {
			row, err := bitsetCollection.GetRowById(key, util.UInt32ToBytes(id))
			if err != nil {
				return err
			}
			if row != nil {
				bsMap[id] = row.Value
			} else {
				bsMap[id] = bitset.New(uint(bitsetItemSize))
			}
		}
		bsMap[id].SetTo(uint(bits[i]%bitsetItemSize), value)
	}

	for id, bs := range bsMap {
		if err := bitsetCollection.UpsertRow(&store.Row{
			Key:   key,
			Id:    util.UInt32ToBytes(id),
			Value: bs,
		}, batch); err != nil {
			return err
		}
	}

	return batch.Commit()
}

func BSGetRange(key []byte, start uint32, end uint32) ([]bool, error) {
	bsMap := make(map[uint32]bitset.Bitset)
	for i := start; i < end; i++ {
		id := i / bitsetItemSize
		if bsMap[id] == nil {
			row, err := bitsetCollection.GetRowById(key, util.UInt32ToBytes(id))
			if err != nil {
				return nil, err
			}
			if row != nil {
				bsMap[id] = row.Value
			} else {
				bsMap[id] = bitset.New(uint(bitsetItemSize))
			}
		}
	}

	res := make([]bool, end-start)
	for i := start; i < end; i++ {
		res[i-start] = bsMap[i/bitsetItemSize].IsSet(uint(i % bitsetItemSize))
	}
	return res, nil
}

func BSMGet(key []byte, bits []uint32) ([]bool, error) {
	bsMap := make(map[uint32]bitset.Bitset)
	for i := 0; i < len(bits); i++ {
		id := bits[i] / bitsetItemSize
		if bsMap[id] == nil {
			row, err := bitsetCollection.GetRowById(key, util.UInt32ToBytes(id))
			if err != nil {
				return nil, err
			}
			if row != nil {
				bsMap[id] = row.Value
			} else {
				bsMap[id] = bitset.New(uint(bitsetItemSize))
			}
		}
	}

	res := make([]bool, len(bits))
	for i := range bits {
		res[i] = bsMap[bits[i]/bitsetItemSize].IsSet(uint(bits[i] % bitsetItemSize))
	}
	return res, nil
}

func BSCount(key []byte) (uint32, error) {
	rows, err := bitsetCollection.Page(key, 0, math.MaxUint32)
	if err != nil {
		return 0, err
	}
	count := uint32(0)
	for i := range rows {
		b := bitset.Bitset(rows[i].Value)
		count += uint32(b.Count())
	}
	return count, nil
}

func BSCountRange(key []byte, start uint32, end uint32) (uint32, error) {
	bsMap := make(map[uint32]bitset.Bitset)
	for i := start; i < end; i++ {
		id := i / bitsetItemSize
		if bsMap[id] == nil {
			row, err := bitsetCollection.GetRowById(key, util.UInt32ToBytes(id))
			if err != nil {
				return 0, err
			}
			if row != nil {
				bsMap[id] = row.Value
			} else {
				bsMap[id] = bitset.New(uint(bitsetItemSize))
			}
		}
	}

	count := uint32(0)
	for i := start; i < end; i++ {
		if bsMap[i/bitsetItemSize].IsSet(uint(i % bitsetItemSize)) {
			count += 1
		}
	}
	return count, nil
}
