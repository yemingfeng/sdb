package collection

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/protobuf/proto"
)

func unmarshal(raw []byte) (*Row, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	pbRow := pb.Row{}
	if err := proto.Unmarshal(raw, &pbRow); err != nil {
		return nil, err
	}
	indexes := make([]Index, len(pbRow.Indexes))
	for i := range pbRow.Indexes {
		indexes[i] = Index{Name: pbRow.Indexes[i].Name, Value: pbRow.Indexes[i].Value}
	}
	return &Row{Key: pbRow.Key, Id: pbRow.Id, Indexes: indexes, Value: pbRow.Value}, nil
}

func marshal(row *Row) ([]byte, error) {
	indexes := make([]*pb.Index, len(row.Indexes))
	for i := range row.Indexes {
		indexes[i] = &pb.Index{Name: row.Indexes[i].Name, Value: row.Indexes[i].Value}
	}
	pbRow := pb.Row{Key: row.Key, Id: row.Id, Indexes: indexes, Value: row.Value}
	return proto.Marshal(&pbRow)
}

func rowKey(dataType pb.DataType, key []byte, id []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/id/%s/", dataType, key, id))
}

func rowKeyPrefix(dataType pb.DataType, key []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/id/", dataType, key))
}

func indexKey(dataType pb.DataType, key []byte, indexName []byte, indexValue []byte, id []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/idx_%s/%s/%s/", dataType, key, indexName, indexValue, id))
}

func indexKeyPrefix(dataType pb.DataType, key []byte, indexName []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/idx_%s/", dataType, key, indexName))
}

func indexKeyValuePrefix(dataType pb.DataType, key []byte, indexName []byte, indexValue []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/idx_%s/%s/", dataType, key, indexName, indexValue))
}
