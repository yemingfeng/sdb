package service

import (
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/engine"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"github.com/yemingfeng/sdb/internal/util"
)

var pageCollection = collection.NewCollection(pb.DataType_PAGE)
var emptyValue = []byte{1}

func PAdd(dataType pb.DataType, key []byte, batch engine.Batch) error {
	return pageCollection.UpsertRow(&collection.Row{Key: util.ToBytes(int32(dataType)),
		Id: key, Value: emptyValue}, batch)
}

func PDel(dataType pb.DataType, key []byte, batch engine.Batch) error {
	return pageCollection.DelRowById(util.ToBytes(int32(dataType)), key, batch)
}

func PList(dataType pb.DataType, key []byte, offset int32, limit uint32) ([][]byte, error) {
	var rows []*collection.Row
	var err error

	if len(key) == 0 {
		rows, err = pageCollection.Page(util.ToBytes(int32(dataType)), offset, limit)
	} else {
		var row *collection.Row
		row, err = pageCollection.GetRowById(util.ToBytes(int32(dataType)), []byte(key))
		if row != nil {
			rows = []*collection.Row{row}
		}
	}
	if err != nil {
		return nil, err
	}

	keys := make([][]byte, len(rows))
	for i := range rows {
		keys[i] = rows[i].Id
	}
	return keys, nil
}
