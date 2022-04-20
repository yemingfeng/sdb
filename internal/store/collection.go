package store

import (
	"errors"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"math"
)

var idEmptyError = errors.New("id is empty")
var keyEmptyError = errors.New("key is empty")

// Collection is an abstraction of data structure, dataType = List/Set/SortedSet
// A Collection corresponds to a row containing
// Each row row takes rowKey as a unique value, rowKey = {dataType} + {key} + {id} is combined to form a unique value
// Each row contains N indices
// Each index uses indexKey as a unique value, indexKey = {dataType} + {key} + idx_{indexName} + {indexValue} + {id}
// Take ListCollection as an example, the key of the List is [l1], assuming that the Collection has 4 rows of Row, and each row of Row has the index of value and score
// Then each row of Row is as follows:
// { {key: l1}, {id: 1.1}, {value: aaa}, {score: 1.1}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 1.1} ] }
// { {key: l1}, {id: 2.2}, {value: bbb}, {score: 2.2}, indexes: [ {name: "value", value: bbb}, {name: "score", value: 2.2} ] }
// { {key: l1}, {id: 3.3}, {value: ccc}, {score: 3.3}, indexes: [ {name: "value", value: ccc}, {name: "score", value: 3.3} ] }
// { {key: l1}, {id: 4.4}, {value: aaa}, {score: 4.4}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 4.4} ] }
// Take the Row with id = 1.1 as an example, rowKey = 1/l1/1.1, valueIndexKey = 1/l1/idx_value/aaa/1.1, scoreIndexKey = 1/l1/idx_score/1.1/1.1 The written data is:
//    rowKey: 1/l1/1.1 -> { {key: l1}, {id: 1.1}, {value: aaa}, {score: 1.1}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 1.1} ] }
//    valueIndexKey: 1/l1/idx_value/aaa/1.1, -> 1/l1/1.1
//    scoreIndexKey: 1/l1/idx_score/1.1/1.1 -> 1/l1/1.1
type Collection struct {
	dataType pb.DataType
}

// NewCollection create collection
func NewCollection(dataType pb.DataType) *Collection {
	return &Collection{dataType: dataType}
}

// DelRowById delete row by id
func (collection *Collection) DelRowById(key []byte, id []byte, batch Batch) error {
	existRow, err := collection.GetRowByIdWithBatch(key, id, batch)
	if err != nil {
		return err
	}
	return collection.DelRow(existRow, batch)
}

// DelRow delete row
func (collection *Collection) DelRow(row *Row, batch Batch) error {
	if row == nil {
		return nil
	}
	key := row.Key
	id := row.Id
	// delete exist indexes
	for i := range row.Indexes {
		index := row.Indexes[i]
		err := batch.Del(indexKey(collection.dataType, key, index.Name, index.Value, id))
		if err != nil {
			return err
		}
	}
	// delete row
	return batch.Del(rowKey(collection.dataType, key, id))
}

// UpsertRow update or insert
func (collection *Collection) UpsertRow(row *Row, batch Batch) error {
	if len(row.Key) == 0 {
		return keyEmptyError
	}
	if len(row.Id) == 0 {
		return idEmptyError
	}

	existRow, err := collection.GetRowByIdWithBatch(row.Key, row.Id, batch)
	if err != nil {
		return err
	}
	if existRow != nil {
		err := collection.DelRow(existRow, batch)
		if err != nil {
			return err
		}
	}

	rowKey := rowKey(collection.dataType, row.Key, row.Id)
	for i := range row.Indexes {
		index := row.Indexes[i]
		err := batch.Set(indexKey(collection.dataType, row.Key, index.Name, index.Value, row.Id), rowKey)
		if err != nil {
			return err
		}
	}
	rawRow, err := marshal(row)
	if err != nil {
		return err
	}
	return batch.Set(rowKey, rawRow)
}

// DelAll del all by key
func (collection *Collection) DelAll(key []byte, batch Batch) error {
	return batch.Iterate(&PrefixIteratorOption{Prefix: rowKeyPrefix(collection.dataType, key), Offset: 0, Limit: math.MaxUint32},
		func(rowKey []byte, rawRow []byte) error {
			row, err := unmarshal(rawRow)
			if err != nil {
				return err
			}
			return collection.DelRow(row, batch)
		})
}

// GetRowByIdWithBatch get row by id
func (collection *Collection) GetRowByIdWithBatch(key []byte, id []byte, batch Batch) (*Row, error) {
	value, err := batch.Get(rowKey(collection.dataType, key, id))
	if err != nil {
		return nil, err
	}
	return unmarshal(value)
}

// GetRowById get row by id
func (collection *Collection) GetRowById(key []byte, id []byte) (*Row, error) {
	batch := NewBatch()
	defer batch.Close()

	value, err := batch.Get(rowKey(collection.dataType, key, id))
	if err != nil {
		return nil, err
	}
	return unmarshal(value)
}

// ExistRowById check row exist
func (collection *Collection) ExistRowById(key []byte, id []byte) (bool, error) {
	row, err := collection.GetRowById(key, id)
	if err != nil {
		return false, err
	}
	return row != nil, nil
}

// Count dataType + key
func (collection *Collection) Count(key []byte) (uint32, error) {
	batch := NewBatch()
	defer batch.Close()

	count := uint32(0)
	if err := batch.Iterate(&PrefixIteratorOption{Prefix: rowKeyPrefix(collection.dataType, key), Offset: 0, Limit: math.MaxUint32},
		func(_ []byte, _ []byte) error {
			count++
			return nil
		}); err != nil {
		return 0, err
	}
	return count, nil
}

// Page dataType + key
func (collection *Collection) Page(key []byte, offset int32, limit uint32) ([]*Row, error) {
	batch := NewBatch()
	defer batch.Close()

	rows := make([]*Row, 0)
	if err := batch.Iterate(&PrefixIteratorOption{Prefix: rowKeyPrefix(collection.dataType, key), Offset: offset, Limit: limit},
		func(_ []byte, rawRow []byte) error {
			row, err := unmarshal(rawRow)
			if err != nil {
				return err
			}
			rows = append(rows, row)
			return nil
		}); err != nil {
		return nil, err
	}
	return rows, nil
}

// PageWithBatch dataType + key
func (collection *Collection) PageWithBatch(key []byte, offset int32, limit uint32, batch Batch) ([]*Row, error) {
	rows := make([]*Row, 0)
	if err := batch.Iterate(&PrefixIteratorOption{Prefix: rowKeyPrefix(collection.dataType, key), Offset: offset, Limit: limit},
		func(_ []byte, rawRow []byte) error {
			row, err := unmarshal(rawRow)
			if err != nil {
				return err
			}
			rows = append(rows, row)
			return nil
		}); err != nil {
		return nil, err
	}

	return rows, nil
}

// IndexPage page by index name
func (collection *Collection) IndexPage(key []byte, indexName []byte, offset int32, limit uint32) ([]*Row, error) {
	batch := NewBatch()
	defer batch.Close()

	rows := make([]*Row, 0)

	if err := batch.Iterate(&PrefixIteratorOption{Prefix: indexKeyPrefix(collection.dataType, key, indexName), Offset: offset, Limit: limit},
		func(indexKey []byte, rowKey []byte) error {
			rowRaw, err := batch.Get(rowKey)
			if err != nil {
				return err
			}
			row, err := unmarshal(rowRaw)
			if err != nil {
				return err
			}
			rows = append(rows, row)
			return nil
		}); err != nil {
		return nil, err
	}
	return rows, nil
}

// IndexValuePage page by index value
func (collection *Collection) IndexValuePage(key []byte, indexName []byte, indexValue []byte, offset int32, limit uint32) ([]*Row, error) {
	batch := NewBatch()
	defer batch.Close()

	rows := make([]*Row, 0)

	if err := batch.Iterate(&PrefixIteratorOption{Prefix: indexKeyValuePrefix(collection.dataType, key, indexName, indexValue), Offset: offset, Limit: limit},
		func(indexKey []byte, rowKey []byte) error {
			rowRaw, err := batch.Get(rowKey)
			if err != nil {
				return err
			}
			row, err := unmarshal(rowRaw)
			if err != nil {
				return err
			}
			rows = append(rows, row)
			return nil
		}); err != nil {
		return nil, err
	}
	return rows, nil
}
