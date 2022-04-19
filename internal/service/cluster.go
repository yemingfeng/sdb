package service

import (
	"github.com/yemingfeng/sdb/internal/store"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
)

func CInfo() ([]*pb.Node, error) {
	return store.GetNodes()
}