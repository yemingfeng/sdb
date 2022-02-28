package collection

import (
	"errors"
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"math"
)

var idEmptyError = errors.New("id is empty")
var keyEmptyError = errors.New("key is empty")
var valueEmptyError = errors.New("value is empty")

// Collection 是对数据结构的抽象，dataType = List/Set/SortedSet
// 一个 Collection 对应包含 row
// 每行 row 以 rowKey 作为唯一值，rowKey = {dataType} + {key} + {id} 联合形成唯一值
// 每行 row 包含 N 个索引
// 每个索引以 indexKey 作为唯一值，indexKey = {dataType} + {key} + idx_{indexName} + {indexValue} + {id}
// 以 ListCollection 为例子，该 List 的 key 为 [l1]，假设该 Collection 有 4 行 Row，每行 Row 都有 value 和 score 的索引
// 那么每行 Row 如下：
// { {key: l1}, {id: 1.1}, {value: aaa}, {score: 1.1}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 1.1} ] }
// { {key: l1}, {id: 2.2}, {value: bbb}, {score: 2.2}, indexes: [ {name: "value", value: bbb}, {name: "score", value: 2.2} ] }
// { {key: l1}, {id: 3.3}, {value: ccc}, {score: 3.3}, indexes: [ {name: "value", value: ccc}, {name: "score", value: 3.3} ] }
// { {key: l1}, {id: 4.4}, {value: aaa}, {score: 4.4}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 4.4} ] }
// 以 id = 1.1 的 Row 为例子，rowKey = 1/l1/1.1, valueIndexKey = 1/l1/idx_value/aaa/1.1, scoreIndexKey = 1/l1/idx_score/1.1/1.1 写入的数据为：
//    rowKey: 1/l1/1.1 -> { {key: l1}, {id: 1.1}, {value: aaa}, {score: 1.1}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 1.1} ] }
//    valueIndexKey: 1/l1/idx_value/aaa/1.1, -> 1/l1/1.1
//    scoreIndexKey: 1/l1/idx_score/1.1/1.1 -> 1/l1/1.1
type Collection struct {
	dataType pb.DataType
}

type Row struct {
	Key     []byte
	Id      []byte
	Value   []byte
	Indexes []Index
}

type Index struct {
	Name  []byte
	Value []byte
}

// NewCollection create collection
func NewCollection(dataType pb.DataType) *Collection {
	return &Collection{dataType: dataType}
}

// DelRowById delete row by id
func (collection *Collection) DelRowById(key []byte, id []byte, batch engine.Batch) error {
	existRow, err := collection.GetRowByIdWithBatch(key, id, batch)
	if err != nil {
		return err
	}
	if existRow != nil {
		// delete exist indexes
		for i := range existRow.Indexes {
			index := existRow.Indexes[i]
			err := batch.Del(indexKey(collection.dataType, key, index.Name, index.Value, id))
			if err != nil {
				return err
			}
		}
	}
	// delete row
	err = batch.Del(rowKey(collection.dataType, key, id))

	return err
}

// UpsertRow update or insert
// batch can be nil, if nil, will auto commit
func (collection *Collection) UpsertRow(row *Row, batch engine.Batch) error {
	if len(row.Key) == 0 {
		return keyEmptyError
	}
	if len(row.Id) == 0 {
		return idEmptyError
	}
	if len(row.Value) == 0 {
		return valueEmptyError
	}

	existRow, err := collection.GetRowByIdWithBatch(row.Key, row.Id, batch)
	if err != nil {
		return err
	}
	if existRow != nil {
		err := collection.DelRowById(row.Key, row.Id, batch)
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
func (collection *Collection) DelAll(key []byte, batch engine.Batch) error {
	return Iterate(rowKeyPrefix(collection.dataType, key),
		0, math.MaxUint32, func(rowKey []byte, rawRow []byte) error {
			row, err := unmarshal(rawRow)
			if err != nil {
				return err
			}
			for i := range row.Indexes {
				index := row.Indexes[i]
				err := batch.Del(indexKey(collection.dataType, row.Key, index.Name, index.Value, row.Id))
				if err != nil {
					return err
				}
			}
			err = batch.Del(rowKey)
			if err != nil {
				return err
			}
			return nil
		})
}

// GetRowByIdWithBatch get row by id
func (collection *Collection) GetRowByIdWithBatch(key []byte, id []byte, batch engine.Batch) (*Row, error) {
	if batch == nil {
		return nil, errors.New("batch is nil")
	}
	value, err := batch.Get(rowKey(collection.dataType, key, id))
	if err != nil {
		return nil, err
	}
	return unmarshal(value)
}

// GetRowById get row by id
func (collection *Collection) GetRowById(key []byte, id []byte) (*Row, error) {
	value, err := Get(rowKey(collection.dataType, key, id))
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
	count := uint32(0)
	if err := Iterate(rowKeyPrefix(collection.dataType, key),
		0, math.MaxUint32, func(_ []byte, _ []byte) error {
			count++
			return nil
		}); err != nil {
		return 0, err
	}
	return count, nil
}

// Page dataType + key
func (collection *Collection) Page(key []byte, offset int32, limit uint32) ([]*Row, error) {
	rows := make([]*Row, 0)
	if err := Iterate(rowKeyPrefix(collection.dataType, key),
		offset, limit, func(_ []byte, rawRow []byte) error {
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
	rows := make([]*Row, 0)
	if err := Iterate(indexKeyPrefix(collection.dataType, key, indexName),
		offset, limit, func(indexKey []byte, rowKey []byte) error {
			rowRaw, err := Get(rowKey)
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
	rows := make([]*Row, 0)
	if err := Iterate(indexKeyValuePrefix(collection.dataType, key, indexName, indexValue),
		offset, limit, func(indexKey []byte, rowKey []byte) error {
			rowRaw, err := Get(rowKey)
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
