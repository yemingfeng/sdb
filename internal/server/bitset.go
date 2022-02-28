package server

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type BitsetServer struct {
	pb.UnimplementedSDBServer
}

func (server *BitsetServer) BSCreate(_ context.Context, request *pb.BSCreateRequest) (*pb.BSCreateResponse, error) {
	err := service.BSCreate(request.Key, request.Size)
	return &pb.BSCreateResponse{Success: err == nil}, err
}

func (server *BitsetServer) BSDel(_ context.Context, request *pb.BSDelRequest) (*pb.BSDelResponse, error) {
	err := service.BSDel(request.Key)
	return &pb.BSDelResponse{Success: err == nil}, err
}

func (server *BitsetServer) BSSetRange(_ context.Context, request *pb.BSSetRangeRequest) (*pb.BSSetRangeResponse, error) {
	err := service.BSSetRange(request.Key, request.Start, request.End, request.Value)
	return &pb.BSSetRangeResponse{Success: err == nil}, err
}

func (server *BitsetServer) BSMSet(_ context.Context, request *pb.BSMSetRequest) (*pb.BSMSetResponse, error) {
	err := service.BSMSet(request.Key, request.Bits, request.Value)
	return &pb.BSMSetResponse{Success: err == nil}, err
}

func (server *BitsetServer) BSGetRange(_ context.Context, request *pb.BSGetRangeRequest) (*pb.BSGetRangeResponse, error) {
	res, err := service.BSGetRange(request.Key, request.Start, request.End)
	return &pb.BSGetRangeResponse{Values: res}, err
}

func (server *BitsetServer) BSMGet(_ context.Context, request *pb.BSMGetRequest) (*pb.BSMGetResponse, error) {
	res, err := service.BSMGet(request.Key, request.Bits)
	return &pb.BSMGetResponse{Values: res}, err
}

func (server *BitsetServer) BSCount(_ context.Context, request *pb.BSCountRequest) (*pb.BSCountResponse, error) {
	res, err := service.BSCount(request.Key)
	return &pb.BSCountResponse{Count: res}, err
}

func (server *BitsetServer) BSCountRange(_ context.Context, request *pb.BSCountRangeRequest) (*pb.BSCountRangeResponse, error) {
	res, err := service.BSCountRange(request.Key, request.Start, request.End)
	return &pb.BSCountRangeResponse{Count: res}, err
}
