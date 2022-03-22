package service

import (
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
)

var pageCollection = store.NewCollection(pb.DataType_PAGE)

func PAdd(dataType pb.DataType, key []byte, batch store.Batch) error {
	return pageCollection.UpsertRow(&store.Row{Key: util.ToBytes(int32(dataType)),
		Id: key, Value: nil}, batch)
}

func PDel(dataType pb.DataType, key []byte, batch store.Batch) error {
	return pageCollection.DelRowById(util.ToBytes(int32(dataType)), key, batch)
}

func PList(dataType pb.DataType, key []byte, offset int32, limit uint32) ([][]byte, error) {
	var rows []*store.Row
	var err error

	if len(key) == 0 {
		rows, err = pageCollection.Page(util.ToBytes(int32(dataType)), offset, limit)
	} else {
		var row *store.Row
		row, err = pageCollection.GetRowById(util.ToBytes(int32(dataType)), []byte(key))
		if row != nil {
			rows = []*store.Row{row}
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
