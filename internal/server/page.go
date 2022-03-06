package server

import (
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type PageServer struct {
	pb.UnimplementedSDBServer
}

func (server *PageServer) PList(_ context.Context, request *pb.PListRequest) (*pb.PListResponse, error) {
	res, err := service.PList(request.DataType, request.Key, request.Offset, request.Limit)
	return &pb.PListResponse{Keys: res}, err
}
