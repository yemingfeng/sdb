package server

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type BloomFilterServer struct {
	pb.UnimplementedSDBServer
}

func (server *BloomFilterServer) BFCreate(_ context.Context, request *pb.BFCreateRequest) (*pb.BFCreateResponse, error) {
	err := service.BFCreate(request.Key, request.N, request.P)
	return &pb.BFCreateResponse{Success: err == nil}, err
}

func (server *BloomFilterServer) BFDel(_ context.Context, request *pb.BFDelRequest) (*pb.BFDelResponse, error) {
	err := service.BFDel(request.Key)
	return &pb.BFDelResponse{Success: err == nil}, err
}

func (server *BloomFilterServer) BFAdd(_ context.Context, request *pb.BFAddRequest) (*pb.BFAddResponse, error) {
	err := service.BFAdd(request.Key, request.Values)
	return &pb.BFAddResponse{Success: err == nil}, err
}

func (server *BloomFilterServer) BFExist(_ context.Context, request *pb.BFExistRequest) (*pb.BFExistResponse, error) {
	res, err := service.BFExist(request.Key, request.Values)
	return &pb.BFExistResponse{Exists: res}, err
}
