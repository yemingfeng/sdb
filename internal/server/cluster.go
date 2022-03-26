package server

import (
	"github.com/yemingfeng/sdb/internal/service"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
)

type ClusterServer struct {
	pb.UnimplementedSDBServer
}

func (server *ClusterServer) CInfo(_ context.Context, _ *pb.CInfoRequest) (*pb.CInfoResponse, error) {
	res, err := service.CInfo()
	return &pb.CInfoResponse{Nodes: res}, err
}
